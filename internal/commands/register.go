package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/internal/database"
	"github.com/kaffeed/voidling/internal/embeds"
	"github.com/kaffeed/voidling/internal/wiseoldman"
)

// RegisterCommands holds the handlers for account registration commands
type RegisterCommands struct {
	DB        *database.Queries
	DBSQL     *sql.DB
	WOMClient *wiseoldman.Client
}

// NewRegisterCommands creates a new RegisterCommands instance
func NewRegisterCommands(db *database.Queries, dbSQL *sql.DB, womClient *wiseoldman.Client) *RegisterCommands {
	return &RegisterCommands{
		DB:        db,
		DBSQL:     dbSQL,
		WOMClient: womClient,
	}
}

// HandleLinkRSN shows the modal for linking a RuneScape account
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

// HandleLinkRSNModal processes the modal submission for linking
func (r *RegisterCommands) HandleLinkRSNModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the response to give us time to fetch player data
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring link-rsn modal response: %v", err)
		return
	}

	// Get the username from the modal
	data := i.ModalSubmitData()
	username := strings.TrimSpace(data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)

	if username == "" {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Username cannot be empty.",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
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

// HandleConfirmRSN handles the confirmation button for linking
func (r *RegisterCommands) HandleConfirmRSN(s *discordgo.Session, i *discordgo.InteractionCreate, username string) {
	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring confirm-rsn response: %v", err)
		return
	}

	ctx := context.Background()
	discordIDStr := i.Member.User.ID
	discordID, err := strconv.ParseInt(discordIDStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing Discord ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Invalid Discord ID.")},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
		return
	}

	log.Printf("Confirming RSN link for Discord user %s (%d) with RSN: %s", i.Member.User.Username, discordID, username)

	// Start a transaction
	tx, err := r.DBSQL.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Database error. Please try again later.")},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
		return
	}
	defer tx.Rollback()

	qtx := r.DB.WithTx(tx)

	// Check if this exact account link already exists and is active
	existingLink, err := qtx.GetExistingAccountLink(ctx, database.GetExistingAccountLinkParams{
		DiscordMemberID: discordID,
		LOWER:           strings.ToLower(username),
	})
	linkExists := (err == nil)

	if linkExists && existingLink.IsActive {
		log.Printf("Account link already exists and is active for user %d", discordID)
		tx.Rollback()
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "This account is already linked and active!",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Deactivate all existing links for this user
	err = qtx.DeactivateAllAccountLinksForUser(ctx, discordID)
	if err != nil {
		log.Printf("Error deactivating existing links: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Failed to update account links. Please try again.")},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// If the link exists but was inactive, reactivate it
	if linkExists {
		log.Printf("Reactivating existing account link: %d", existingLink.ID)
		err = qtx.ActivateAccountLink(ctx, existingLink.ID)
		if err != nil {
			log.Printf("Error reactivating account link: %v", err)
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Failed to activate account link. Please try again.")},
				Flags:  discordgo.MessageFlagsEphemeral,
			})
			return
		}
	} else {
		// Create new link
		log.Printf("Creating new account link for user %d with RSN %s", discordID, username)
		_, err = qtx.CreateAccountLink(ctx, database.CreateAccountLinkParams{
			DiscordMemberID: discordID,
			RunescapeName:   username,
			IsActive:        true,
		})
		if err != nil {
			log.Printf("Error creating account link: %v", err)
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Failed to link account. Please try again.")},
				Flags:  discordgo.MessageFlagsEphemeral,
			})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Failed to save changes. Please try again.")},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
		return
	}

	log.Printf("Successfully linked RSN %s to Discord user %d", username, discordID)

	// Update the original message to remove buttons and show success
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &[]discordgo.MessageComponent{},
	})

	// Send success message
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embeds.SuccessEmbed(fmt.Sprintf("Successfully linked your account to **%s**!", username))},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}

// HandleCancelRSN handles the cancel button for linking
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

// HandleUnlinkRSN handles unlinking a RuneScape account
func (r *RegisterCommands) HandleUnlinkRSN(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring unlink-rsn response: %v", err)
		return
	}

	ctx := context.Background()
	discordIDStr := i.Member.User.ID
	discordID, err := strconv.ParseInt(discordIDStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing Discord ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Invalid Discord ID.",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	log.Printf("Unlinking RSN for Discord user %s (%d)", i.Member.User.Username, discordID)

	// Get active account link
	activeLink, err := r.DB.GetAccountLinkByDiscordID(ctx, discordID)
	if err == sql.ErrNoRows {
		log.Printf("No active account found for user %d", discordID)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "You don't have any linked account. Use `/link-rsn` to link your RuneScape account.",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}
	if err != nil {
		log.Printf("Error fetching account link: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Database error. Please try again later.")},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Deactivate the link
	err = r.DB.DeactivateAccountLink(ctx, activeLink.ID)
	if err != nil {
		log.Printf("Error deactivating account link: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embeds.ErrorEmbed("Failed to unlink account. Please try again.")},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
		return
	}

	log.Printf("Successfully unlinked RSN %s from Discord user %d", activeLink.RunescapeName, discordID)

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embeds.SuccessEmbed(fmt.Sprintf("Successfully unlinked your account from **%s**.", activeLink.RunescapeName))},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}
