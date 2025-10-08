// Package commands provides Discord slash command handlers for the Voidling bot.
// This file contains registration-related commands for linking/unlinking RuneScape accounts.
package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/internal/database"
	"github.com/kaffeed/voidling/internal/embeds"
	"github.com/kaffeed/voidling/internal/wiseoldman"
)

// ErrNoGuildContext is returned when a guild context is required but not available.
var ErrNoGuildContext = errors.New("no guild context")

// RegisterCommands holds the handlers for account registration commands.
type RegisterCommands struct {
	DB        *database.Queries
	DBSQL     *sql.DB
	WOMClient *wiseoldman.Client
}

// NewRegisterCommands creates a new RegisterCommands instance.
func NewRegisterCommands(db *database.Queries, dbSQL *sql.DB, womClient *wiseoldman.Client) *RegisterCommands {
	return &RegisterCommands{
		DB:        db,
		DBSQL:     dbSQL,
		WOMClient: womClient,
	}
}

// HandleLinkRSN shows the modal for linking a RuneScape account.
func (r *RegisterCommands) HandleLinkRSN(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "link-rsn-modal",
			Title:    "Link RuneScape Account",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "rsn-input",
							Label:       "RuneScape Username",
							Style:       discordgo.TextInputShort,
							Placeholder: "Enter your RSN",
							Required:    true,
							MaxLength:   12,
							MinLength:   1,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error showing link-rsn modal: %v", err)
	}
}

// HandleLinkRSNModal processes the modal submission for linking.
func (r *RegisterCommands) HandleLinkRSNModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := r.deferEphemeralResponse(s, i); err != nil {
		return
	}

	username := r.extractUsernameFromModal(i)
	if username == "" {
		r.sendErrorFollowup(s, i, "Username cannot be empty.")
		return
	}

	log.Printf("User %s wants to link RSN: %s", i.Member.User.Username, username)

	// Fetch player from Wise Old Man API
	ctx := context.Background()
	player, err := r.WOMClient.GetPlayer(ctx, username)
	if err != nil {
		log.Printf("Error fetching player %s: %v", username, err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed(fmt.Sprintf("Failed to fetch player data for '%s'. Make sure the username is correct and try again.", username)),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Create confirmation buttons
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "That's me!",
					Style:    discordgo.SuccessButton,
					CustomID: fmt.Sprintf("confirm-rsn:%s", username),
				},
				discordgo.Button{
					Label:    "Not me",
					Style:    discordgo.DangerButton,
					CustomID: fmt.Sprintf("cancel-rsn:%s", username),
				},
			},
		},
	}

	// Show player info with confirmation buttons
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embeds.PlayerInfo(player)},
		Components: components,
		Flags:      discordgo.MessageFlagsEphemeral,
	})
}

// HandleConfirmRSN handles the confirmation button for linking.
func (r *RegisterCommands) HandleConfirmRSN(s *discordgo.Session, i *discordgo.InteractionCreate, username string) {
	if err := r.deferEphemeralResponse(s, i); err != nil {
		return
	}

	ctx := context.Background()
	userID, guildID := r.getUserAndGuildIDs(i)

	discordID, err := r.parseDiscordID(userID)
	if err != nil {
		log.Printf("Error parsing Discord ID: %v", err)
		r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Invalid Discord ID."))
		return
	}

	log.Printf("Confirming RSN link for Discord user %s (%d) with RSN: %s", userID, discordID, username)

	// Start a transaction
	tx, err := r.DBSQL.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Database error. Please try again later."))
		return
	}
	defer func() { _ = tx.Rollback() }() // Rollback is safe to call even after commit

	qtx := r.DB.WithTx(tx)

	// Check if this exact account link already exists and is active
	existingLink, err := qtx.GetExistingAccountLink(ctx, database.GetExistingAccountLinkParams{
		DiscordMemberID: discordID,
		LOWER:           strings.ToLower(username),
	})
	linkExists := (err == nil)

	if linkExists && existingLink.IsActive {
		log.Printf("Account link already exists and is active for user %d", discordID)
		_ = tx.Rollback()
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "This account is already linked and active!",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Deactivate all existing links for this user
	if err = qtx.DeactivateAllAccountLinksForUser(ctx, discordID); err != nil {
		log.Printf("Error deactivating existing links: %v", err)
		r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Failed to update account links. Please try again."))
		return
	}

	// If the link exists but was inactive, reactivate it; otherwise create new
	if linkExists {
		log.Printf("Reactivating existing account link: %d", existingLink.ID)
		if err = qtx.ActivateAccountLink(ctx, existingLink.ID); err != nil {
			log.Printf("Error reactivating account link: %v", err)
			r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Failed to activate account link. Please try again."))
			return
		}
	} else {
		log.Printf("Creating new account link for user %d with RSN %s", discordID, username)
		if _, err = qtx.CreateAccountLink(ctx, database.CreateAccountLinkParams{
			DiscordMemberID: discordID,
			RunescapeName:   username,
			IsActive:        true,
		}); err != nil {
			log.Printf("Error creating account link: %v", err)
			r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Failed to link account. Please try again."))
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Failed to save changes. Please try again."))
		return
	}

	log.Printf("Successfully linked RSN %s to Discord user %d", username, discordID)

	// Update the original message to remove buttons
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &[]discordgo.MessageComponent{},
	}) // Ignore error - success message sent below

	// Send appropriate success message based on nickname update result
	successMsg := r.buildSuccessMessage(s, guildID, userID, username)
	r.sendEmbedFollowup(s, i, embeds.SuccessEmbed(successMsg))
}

// buildSuccessMessage creates a success message and attempts nickname update.
func (r *RegisterCommands) buildSuccessMessage(s *discordgo.Session, guildID, userID, username string) string {
	baseMsg := fmt.Sprintf("Successfully linked your account to **%s**!", username)

	if guildID == "" {
		return baseMsg
	}

	// Attempt to update nickname
	if err := r.updateMemberNickname(s, guildID, userID, username); err != nil {
		log.Printf("Failed to update nickname for user %s in guild %s: %v", userID, guildID, err)
		return baseMsg + "\n\n*Note: I couldn't update your server nickname automatically. Please ask a server admin to update it.*"
	}

	log.Printf("Updated nickname for user %s to %s in guild %s", userID, username, guildID)
	return baseMsg + " Your server nickname has been updated too!"
}

// HandleCancelRSN handles the cancel button for linking.
func (r *RegisterCommands) HandleCancelRSN(s *discordgo.Session, i *discordgo.InteractionCreate, username string) {
	log.Printf("User %s cancelled linking RSN: %s", i.Member.User.Username, username)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf("Account linking cancelled for '%s'.", username),
			Components: []discordgo.MessageComponent{}, // Remove buttons
			Embeds:     []*discordgo.MessageEmbed{},    // Remove embed
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

// HandleUnlinkRSN handles unlinking a RuneScape account.
func (r *RegisterCommands) HandleUnlinkRSN(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := r.deferEphemeralResponse(s, i); err != nil {
		return
	}

	ctx := context.Background()
	discordID, err := r.parseDiscordID(i.Member.User.ID)
	if err != nil {
		log.Printf("Error parsing Discord ID: %v", err)
		r.sendErrorFollowup(s, i, "Invalid Discord ID.")
		return
	}

	log.Printf("Unlinking RSN for Discord user %s (%d)", i.Member.User.Username, discordID)

	// Get active account link
	activeLink, err := r.DB.GetAccountLinkByDiscordID(ctx, discordID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("No active account found for user %d", discordID)
		r.sendErrorFollowup(s, i, "You don't have any linked account. Use `/link-rsn` to link your RuneScape account.")
		return
	}
	if err != nil {
		log.Printf("Error fetching account link: %v", err)
		r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Database error. Please try again later."))
		return
	}

	// Deactivate the link
	if err = r.DB.DeactivateAccountLink(ctx, activeLink.ID); err != nil {
		log.Printf("Error deactivating account link: %v", err)
		r.sendEmbedFollowup(s, i, embeds.ErrorEmbed("Failed to unlink account. Please try again."))
		return
	}

	log.Printf("Successfully unlinked RSN %s from Discord user %d", activeLink.RunescapeName, discordID)
	r.sendEmbedFollowup(s, i, embeds.SuccessEmbed(fmt.Sprintf("Successfully unlinked your account from **%s**.", activeLink.RunescapeName)))
}

// Helper methods for cleaner code

// deferEphemeralResponse defers an ephemeral response for the interaction.
func (r *RegisterCommands) deferEphemeralResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring response: %v", err)
	}
	return err
}

// extractUsernameFromModal extracts and trims the username from modal submission.
func (r *RegisterCommands) extractUsernameFromModal(i *discordgo.InteractionCreate) string {
	data := i.ModalSubmitData()
	actionRow, ok := data.Components[0].(*discordgo.ActionsRow)
	if !ok {
		return ""
	}
	textInput, ok := actionRow.Components[0].(*discordgo.TextInput)
	if !ok {
		return ""
	}
	return strings.TrimSpace(textInput.Value)
}

// sendErrorFollowup sends an error message as a followup.
func (r *RegisterCommands) sendErrorFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: message,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

// sendEmbedFollowup sends an embed as a followup message.
func (r *RegisterCommands) sendEmbedFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}

// getUserAndGuildIDs extracts user ID and guild ID from interaction context.
func (r *RegisterCommands) getUserAndGuildIDs(i *discordgo.InteractionCreate) (userID, guildID string) {
	if i.Member != nil {
		// Guild context (slash command)
		return i.Member.User.ID, i.GuildID
	}

	if i.User != nil {
		// DM context (button from greeting)
		userID = i.User.ID
		guildID = r.extractGuildIDFromMessage(i.Message)
		return userID, guildID
	}

	return "", ""
}

// extractGuildIDFromMessage extracts guild ID from message components.
func (r *RegisterCommands) extractGuildIDFromMessage(msg *discordgo.Message) string {
	if msg == nil || len(msg.Components) == 0 {
		return ""
	}

	for _, component := range msg.Components {
		actionRow, ok := component.(*discordgo.ActionsRow)
		if !ok {
			continue
		}

		for _, comp := range actionRow.Components {
			button, ok := comp.(*discordgo.Button)
			if !ok {
				continue
			}

			// Parse "dm-link-rsn:GUILD_ID" from button custom ID
			parts := strings.SplitN(button.CustomID, ":", 2)
			if len(parts) == 2 && parts[0] == "dm-link-rsn" {
				return parts[1]
			}
		}
	}

	return ""
}

// parseDiscordID safely parses a Discord ID string to int64.
func (r *RegisterCommands) parseDiscordID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}

// updateMemberNickname attempts to update a guild member's nickname.
func (r *RegisterCommands) updateMemberNickname(s *discordgo.Session, guildID, userID, nickname string) error {
	if guildID == "" {
		return ErrNoGuildContext
	}
	return s.GuildMemberNickname(guildID, userID, nickname)
}
