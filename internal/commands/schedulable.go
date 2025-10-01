package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidbound/internal/database"
	"github.com/kaffeed/voidbound/internal/embeds"
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
	activity := options[0].StringValue()  // e.g., "Corporeal Beast", "Nex"
	location := options[1].StringValue()  // e.g., "World 444"
	timeStr := options[2].StringValue()   // e.g., "2025-01-15 20:00"
	durationMinutes := options[3].IntValue() // e.g., 60, 120

	// Parse time
	scheduledTime, err := time.Parse("2006-01-02 15:04", timeStr)
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

	// Store event in database
	_, err = sc.DB.CreateSchedulableEvent(ctx, database.CreateSchedulableEventParams{
		Type:           "MASS",
		Activity:       activity,
		Location:       location,
		ScheduledAt:    scheduledTime,
		DiscordEventID: discordEvent.ID,
	})
	if err != nil {
		log.Printf("Error storing event in database: %v", err)
		// Event created in Discord, but failed to store - not critical
	}

	// Send confirmation
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.MassEvent(activity, location, scheduledTime),
		},
	})
}

