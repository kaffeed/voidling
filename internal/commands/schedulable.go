package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/internal/database"
	"github.com/kaffeed/voidling/internal/embeds"
	"github.com/kaffeed/voidling/internal/timezone"
)

// SchedulableCommands handles Mass and Wildy Wednesday event commands.
type SchedulableCommands struct {
	DB    *database.Queries
	DBSQL *sql.DB
}

// NewSchedulableCommands creates a new SchedulableCommands instance.
func NewSchedulableCommands(db *database.Queries, dbSQL *sql.DB) *SchedulableCommands {
	return &SchedulableCommands{
		DB:    db,
		DBSQL: dbSQL,
	}
}

// getEffectiveTimezone returns the timezone to use: param > user pref > guild default > UTC.
func (sc *SchedulableCommands) getEffectiveTimezone(ctx context.Context, guildID, userID int64, paramTZ string) string {
	// 1. If timezone parameter provided, use it
	if paramTZ != "" {
		if err := timezone.ValidateTimezone(paramTZ); err == nil {
			return paramTZ
		}
	}

	// 2. Try user preference
	userPref, err := sc.DB.GetUserTimezone(ctx, userID)
	if err == nil {
		return userPref.Timezone
	}

	// 3. Try guild default
	guildConfig, err := sc.DB.GetGuildConfig(ctx, guildID)
	if err == nil && guildConfig.DefaultTimezone.Valid {
		return guildConfig.DefaultTimezone.String
	}

	// 4. Fallback to UTC
	return "UTC"
}

// HandleMassEvent handles /mass command.
func (sc *SchedulableCommands) HandleMassEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Error deferring response: %v", err)
		return
	}

	// Get options
	options := i.ApplicationCommandData().Options
	activity := options[0].StringValue()     // e.g., "Corporeal Beast", "Nex"
	location := options[1].StringValue()     // e.g., "World 444"
	timeStr := options[2].StringValue()      // e.g., "2025-01-15 20:00"
	durationMinutes := options[3].IntValue() // e.g., 60, 120
	var timezoneParam string
	if len(options) > 4 {
		timezoneParam = options[4].StringValue() // Optional timezone
	}

	// Parse IDs
	guildID, _ := strconv.ParseInt(i.GuildID, 10, 64)
	userID, _ := strconv.ParseInt(i.Member.User.ID, 10, 64)

	// Get effective timezone
	tz := sc.getEffectiveTimezone(ctx, guildID, userID, timezoneParam)

	// Parse time in the specified timezone
	scheduledTime, err := timezone.ParseInTimezone(timeStr, tz)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Invalid time format. Please use YYYY-MM-DD HH:MM (e.g., 2025-01-15 20:00)"),
			},
		})
		return
	}

	// Check if time is in the future
	if scheduledTime.Before(time.Now()) {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Event time must be in the future!"),
			},
		})
		return
	}

	// Calculate end time
	endTime := scheduledTime.Add(time.Duration(durationMinutes) * time.Minute)

	// Create Discord scheduled event
	eventName := fmt.Sprintf("Mass: %s", activity)
	eventDescription := fmt.Sprintf("Join us for a mass event at %s!\n\nClick 'Interested' to RSVP and get a reminder before the event starts.", location)

	discordEvent, err := s.GuildScheduledEventCreate(i.GuildID, &discordgo.GuildScheduledEventParams{
		Name:               eventName,
		Description:        eventDescription,
		ScheduledStartTime: &scheduledTime,
		ScheduledEndTime:   &endTime,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: location,
		},
		PrivacyLevel: discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	})
	if err != nil {
		log.Printf("Error creating Discord event: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to create Discord event. Please try again."),
			},
		})
		return
	}

	// Store event in database with timezone
	_, err = sc.DB.CreateSchedulableEvent(ctx, database.CreateSchedulableEventParams{
		Type:           "Mass",
		Activity:       activity,
		Location:       location,
		ScheduledAt:    scheduledTime,
		DiscordEventID: discordEvent.ID,
		Timezone:       sql.NullString{String: tz, Valid: true},
	})
	if err != nil {
		log.Printf("Error storing event in database: %v", err)
		// Event created in Discord, but failed to store - not critical
	}

	// Get notification role if configured
	content := ""
	guildConfig, err := sc.DB.GetGuildConfig(ctx, guildID)
	if err == nil && guildConfig.EventNotificationRoleID.Valid {
		content = fmt.Sprintf("<@&%d>", guildConfig.EventNotificationRoleID.Int64)
	}

	embed := embeds.MassEventWithTimezone(activity, location, scheduledTime, tz)

	// Create participation button
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "I'll Participate",
					Style:    discordgo.SuccessButton,
					CustomID: fmt.Sprintf("participate-mass:%s", discordEvent.ID),
				},
				discordgo.Button{
					Label:    "List Participants",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("list-participants-mass:%s", discordEvent.ID),
				},
			},
		},
	}

	// If notification channel is configured, post there. Otherwise post in command channel
	if err == nil && guildConfig.EventNotificationChannelID.Valid {
		// Post to event notification channel
		notificationChannelID := strconv.FormatInt(guildConfig.EventNotificationChannelID.Int64, 10)
		_, err = s.ChannelMessageSendComplex(notificationChannelID, &discordgo.MessageSend{
			Content:    content,
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		})
		if err != nil {
			log.Printf("Error posting to event notification channel: %v", err)
			// Fallback to command channel on error
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content:    content,
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
			})
		} else {
			// Success - send confirmation in command channel
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					embeds.SuccessEmbed(fmt.Sprintf("Event created! Check <#%s> for details.", notificationChannelID)),
				},
			})
		}
	} else {
		// No notification channel configured - post in command channel
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content:    content,
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		})
	}
}

// HandleParticipateInMass handles mass event participation button clicks.
func (sc *SchedulableCommands) HandleParticipateInMass(s *discordgo.Session, i *discordgo.InteractionCreate, discordEventID string) {
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

	// Check if user has linked RSN
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

	accountLink, err := sc.DB.GetAccountLinkByDiscordID(ctx, userID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("You must link your RuneScape account first! Use `/link-rsn` to get started."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Get event from database using Discord event ID
	event, err := sc.DB.GetSchedulableEventByDiscordID(ctx, discordEventID)
	if err != nil {
		log.Printf("Error getting event: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Event not found. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Check if already registered
	_, err = sc.DB.GetSchedulableParticipation(ctx, database.GetSchedulableParticipationParams{
		EventID:       event.ID,
		AccountLinkID: accountLink.ID,
	})
	if err == nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("You're already registered for this event!"),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Register for event
	_, err = sc.DB.CreateSchedulableParticipation(ctx, database.CreateSchedulableParticipationParams{
		EventID:       event.ID,
		AccountLinkID: accountLink.ID,
		Notified:      false,
	})
	if err != nil {
		log.Printf("Error registering for event: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to register for event. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.SuccessEmbed(fmt.Sprintf("You're registered for **%s**!\n\nYou'll receive a reminder before the event starts.", event.Activity)),
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
}

// HandleListParticipantsMass handles listing mass event participants.
func (sc *SchedulableCommands) HandleListParticipantsMass(s *discordgo.Session, i *discordgo.InteractionCreate, discordEventID string) {
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

	// Get event from database
	event, err := sc.DB.GetSchedulableEventByDiscordID(ctx, discordEventID)
	if err != nil {
		log.Printf("Error getting event: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Event not found."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Get participants
	participants, err := sc.DB.GetSchedulableParticipationsByEvent(ctx, event.ID)
	if err != nil {
		log.Printf("Error getting participants: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to get participant list."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	if len(participants) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("**%s**\n\nNo participants yet. Be the first to register!", event.Activity),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// Build participant list
	participantList := ""
	for i, p := range participants {
		participantList += fmt.Sprintf("%d. <@%d> - %s\n", i+1, p.DiscordMemberID, p.RunescapeName)
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("**%s - Participants (%d)**\n\n%s", event.Activity, len(participants), participantList),
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}
