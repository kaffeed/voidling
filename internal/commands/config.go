package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidbound/internal/database"
	"github.com/kaffeed/voidbound/internal/embeds"
)

// ConfigCommands handles server configuration commands
type ConfigCommands struct {
	DB    *database.Queries
	DBSQL *sql.DB
}

// NewConfigCommands creates a new ConfigCommands instance
func NewConfigCommands(db *database.Queries, dbSQL *sql.DB) *ConfigCommands {
	return &ConfigCommands{
		DB:    db,
		DBSQL: dbSQL,
	}
}

// HandleSetCoordinatorRole handles /config set-coordinator-role command
func (cc *ConfigCommands) HandleSetCoordinatorRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring response: %v", err)
		return
	}

	// Check if user is server owner or has administrator permission
	if !isServerOwnerOrAdmin(s, i) {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Only the server owner or administrators can configure coordinator roles."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Get role from options (subcommand -> role option)
	options := i.ApplicationCommandData().Options
	if len(options) == 0 || len(options[0].Options) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Missing role parameter. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	roleOption := options[0].Options[0].RoleValue(s, i.GuildID)
	if roleOption == nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to get role information. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Parse guild ID
	guildID, err := strconv.ParseInt(i.GuildID, 10, 64)
	if err != nil {
		log.Printf("Error parsing guild ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to parse guild ID."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Parse role ID
	roleID, err := strconv.ParseInt(roleOption.ID, 10, 64)
	if err != nil {
		log.Printf("Error parsing role ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to parse role ID."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Upsert guild config
	err = cc.DB.UpsertGuildConfig(ctx, database.UpsertGuildConfigParams{
		GuildID:           guildID,
		CoordinatorRoleID: sql.NullInt64{Int64: roleID, Valid: true},
	})
	if err != nil {
		log.Printf("Error upserting guild config: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to save configuration. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Send success message
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.SuccessEmbed(fmt.Sprintf("Coordinator role set to <@&%s>\n\nUsers with this role can now manage BOTW, SOTW, and Mass events.", roleOption.ID)),
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
}

// HandleShowConfig handles /config show command
func (cc *ConfigCommands) HandleShowConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Error deferring response: %v", err)
		return
	}

	// Parse guild ID
	guildID, err := strconv.ParseInt(i.GuildID, 10, 64)
	if err != nil {
		log.Printf("Error parsing guild ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to parse guild ID."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Get guild config
	config, err := cc.DB.GetGuildConfig(ctx, guildID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "**Server Configuration**\n\nCoordinator Role: Not configured\n\nUse `/config set-coordinator-role` to set the coordinator role.",
				Flags:   discordgo.MessageFlagsEphemeral,
			})
			return
		}
		log.Printf("Error getting guild config: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to get configuration."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Build config message
	coordinatorRole := "Not configured"
	if config.CoordinatorRoleID.Valid {
		coordinatorRole = fmt.Sprintf("<@&%d>", config.CoordinatorRoleID.Int64)
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("**Server Configuration**\n\nCoordinator Role: %s", coordinatorRole),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

// isServerOwnerOrAdmin checks if the user is the server owner or has administrator permission
func isServerOwnerOrAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	fmt.Printf("%+v", i.GuildID)
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		return false
	}

	// Check if user is server owner
	if guild.OwnerID == i.Member.User.ID {
		return true
	}

	// Check if user has administrator permission
	permissions, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		log.Printf("Error getting permissions: %v", err)
		return false
	}

	return permissions&discordgo.PermissionAdministrator != 0
}
