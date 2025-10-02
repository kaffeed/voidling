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
	"github.com/kaffeed/voidbound/internal/timezone"
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

	// Get existing config to preserve other settings
	existingConfig, err := cc.DB.GetGuildConfig(ctx, guildID)
	competitionChannel := sql.NullInt64{Valid: false}
	defaultTZ := sql.NullString{String: "UTC", Valid: true}
	if err == nil {
		if existingConfig.CompetitionCodeChannelID.Valid {
			competitionChannel = existingConfig.CompetitionCodeChannelID
		}
		if existingConfig.DefaultTimezone.Valid {
			defaultTZ = existingConfig.DefaultTimezone
		}
	}

	// Upsert guild config
	err = cc.DB.UpsertGuildConfig(ctx, database.UpsertGuildConfigParams{
		GuildID:                  guildID,
		CoordinatorRoleID:        sql.NullInt64{Int64: roleID, Valid: true},
		CompetitionCodeChannelID: competitionChannel,
		DefaultTimezone:          defaultTZ,
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
				Content: "**Server Configuration**\n\nCoordinator Role: Not configured\nCompetition Code Channel: Not configured\n\nUse `/config set-coordinator-role` to set the coordinator role.\nUse `/config set-competition-code-channel` to set the competition code channel.",
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

	competitionCodeChannel := "Not configured"
	if config.CompetitionCodeChannelID.Valid {
		competitionCodeChannel = fmt.Sprintf("<#%d>", config.CompetitionCodeChannelID.Int64)
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("**Server Configuration**\n\nCoordinator Role: %s\nCompetition Code Channel: %s", coordinatorRole, competitionCodeChannel),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

// HandleSetCompetitionCodeChannel handles /config set-competition-code-channel command
func (cc *ConfigCommands) HandleSetCompetitionCodeChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				embeds.ErrorEmbed("Only the server owner or administrators can configure the competition code channel."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Get channel from options (subcommand -> channel option)
	options := i.ApplicationCommandData().Options
	if len(options) == 0 || len(options[0].Options) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Missing channel parameter. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	channelOption := options[0].Options[0].ChannelValue(s)
	if channelOption == nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to get channel information. Please try again."),
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

	// Parse channel ID
	channelID, err := strconv.ParseInt(channelOption.ID, 10, 64)
	if err != nil {
		log.Printf("Error parsing channel ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to parse channel ID."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Update guild config
	err = cc.DB.UpdateCompetitionCodeChannel(ctx, database.UpdateCompetitionCodeChannelParams{
		CompetitionCodeChannelID: sql.NullInt64{Int64: channelID, Valid: true},
		GuildID:                  guildID,
	})
	if err != nil {
		log.Printf("Error updating competition code channel: %v", err)
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
			embeds.SuccessEmbed(fmt.Sprintf("Competition code channel set to <#%s>\n\nWhen BOTW or SOTW events are created, the verification code will be sent to this channel.", channelOption.ID)),
		},
		Flags: discordgo.MessageFlagsEphemeral,
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

// HandleSetDefaultTimezone handles /config set-default-timezone command
func (cc *ConfigCommands) HandleSetDefaultTimezone(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

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

	if !isServerOwnerOrAdmin(s, i) {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Only the server owner or administrators can configure the default timezone."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 || len(options[0].Options) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Missing timezone parameter."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	timezoneStr := options[0].Options[0].StringValue()

	// Validate timezone
	if err := timezone.ValidateTimezone(timezoneStr); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed(fmt.Sprintf("Invalid timezone: %s", timezoneStr)),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

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

	err = cc.DB.UpdateDefaultTimezone(ctx, database.UpdateDefaultTimezoneParams{
		DefaultTimezone: sql.NullString{String: timezoneStr, Valid: true},
		GuildID:         guildID,
	})
	if err != nil {
		log.Printf("Error updating default timezone: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to save timezone configuration."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.SuccessEmbed(fmt.Sprintf("Default server timezone set to **%s**", timezoneStr)),
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
}

// HandleSetMyTimezone handles /config set-my-timezone command
func (cc *ConfigCommands) HandleSetMyTimezone(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

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

	options := i.ApplicationCommandData().Options
	if len(options) == 0 || len(options[0].Options) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Missing timezone parameter."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	timezoneStr := options[0].Options[0].StringValue()
	log.Printf("Received timezone value for set-my-timezone: '%s' (length: %d)", timezoneStr, len(timezoneStr))

	// Validate timezone
	if err := timezone.ValidateTimezone(timezoneStr); err != nil {
		log.Printf("Timezone validation error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed(fmt.Sprintf("Invalid timezone: %s", timezoneStr)),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	userID, err := strconv.ParseInt(i.Member.User.ID, 10, 64)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to parse user ID."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	err = cc.DB.UpsertUserTimezone(ctx, database.UpsertUserTimezoneParams{
		DiscordUserID: userID,
		Timezone:      timezoneStr,
	})
	if err != nil {
		log.Printf("Error saving user timezone: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to save timezone preference."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.SuccessEmbed(fmt.Sprintf("Your timezone set to **%s**\n\nThis will be used as the default when creating events.", timezoneStr)),
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
}
