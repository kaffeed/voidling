package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidbound/internal/database"
	"github.com/kaffeed/voidbound/internal/embeds"
	"github.com/kaffeed/voidbound/internal/timezone"
)

// SchedulableCommands handles Mass and Wildy Wednesday event commands
type SchedulableCommands struct {
	DB    *database.Queries
	DBSQL *sql.DB
}

// NewSchedulableCommands creates a new SchedulableCommands instance
func NewSchedulableCommands(db *database.Queries, dbSQL *sql.DB) *SchedulableCommands {
	return &SchedulableCommands{
		DB:    db,
		DBSQL: dbSQL,
	}
}

// getEffectiveTimezone returns the timezone to use: param > user pref > guild default > UTC
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

// HandleMassEvent handles /mass command
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
	activity := options[0].StringValue()       // e.g., "Corporeal Beast", "Nex"
	location := options[1].StringValue()       // e.g., "World 444"
	timeStr := options[2].StringValue()        // e.g., "2025-01-15 20:00"
	durationMinutes := options[3].IntValue()   // e.g., 60, 120
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
		Type:           "MASS",
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

	// Send confirmation
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.MassEventWithTimezone(activity, location, scheduledTime, tz),
		},
	})
}
