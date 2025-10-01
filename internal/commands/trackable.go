package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidbound/internal/database"
	"github.com/kaffeed/voidbound/internal/embeds"
	"github.com/kaffeed/voidbound/internal/models"
	"github.com/kaffeed/voidbound/internal/wiseoldman"
)

// TrackableCommands handles BOTW and SOTW event commands
type TrackableCommands struct {
	DB        *database.Queries
	DBSQL     *sql.DB
	WOMClient *wiseoldman.Client
}

// NewTrackableCommands creates a new TrackableCommands instance
func NewTrackableCommands(db *database.Queries, dbSQL *sql.DB, womClient *wiseoldman.Client) *TrackableCommands {
	return &TrackableCommands{
		DB:        db,
		DBSQL:     dbSQL,
		WOMClient: womClient,
	}
}

// StartEvent creates a new WOM competition with thread and registration buttons
func (t *TrackableCommands) StartEvent(s *discordgo.Session, i *discordgo.InteractionCreate, eventType models.EventType, activity string) error {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return fmt.Errorf("defer response: %w", err)
	}

	// Create thread for event
	eventName := fmt.Sprintf("%s - %s", getEventDisplayName(eventType), FormatActivityName(activity))
	thread, err := s.ThreadStartComplex(i.ChannelID, &discordgo.ThreadStart{
		Name:                eventName,
		AutoArchiveDuration: 10080, // 1 week in minutes
		Type:                discordgo.ChannelTypeGuildPublicThread,
		Invitable:           false,
	})
	if err != nil {
		log.Printf("Error creating thread: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to create event thread. Please try again."),
			},
		})
		return err
	}

	// Create WOM competition
	// Competition runs for 1 week starting 1 minute from now (WOM requires future dates)
	startsAt := time.Now().Add(1 * time.Minute)
	endsAt := startsAt.Add(7 * 24 * time.Hour)

	womResp, err := t.WOMClient.CreateCompetition(ctx, wiseoldman.CreateCompetitionRequest{
		Title:    eventName,
		Metric:   activity,
		StartsAt: startsAt.Format(time.RFC3339),
		EndsAt:   endsAt.Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("Error creating WOM competition: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to create competition on Wise Old Man. Please try again."),
			},
		})
		return err
	}

	// Store competition in database
	eventTypeStr := string(eventType)
	_, err = t.DB.CreateWOMCompetition(ctx, database.CreateWOMCompetitionParams{
		WomCompetitionID: womResp.Competition.ID,
		VerificationCode: womResp.VerificationCode,
		DiscordThreadID:  thread.ID,
		Metric:           activity,
		Type:             eventTypeStr,
	})
	if err != nil {
		log.Printf("Error storing WOM competition: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to save competition. Please try again."),
			},
		})
		return err
	}

	// Send starter message in thread with WOM link
	womURL := fmt.Sprintf("https://wiseoldman.net/competitions/%d", womResp.Competition.ID)
	threadStarterMsg := fmt.Sprintf("**%s** event has started!\n\n🔗 [View on Wise Old Man](%s)\n\nClick the Register button in the channel to join!", eventName, womURL)
	_, err = s.ChannelMessageSend(thread.ID, threadStarterMsg)
	if err != nil {
		log.Printf("Error sending thread starter message: %v", err)
	}

	// Create embed and buttons
	var embed *discordgo.MessageEmbed
	if eventType == models.EventTypeBossOfTheWeek {
		embed = embeds.BossOfTheWeek(models.HiscoreField(activity), womResp.Competition.ID)
	} else {
		embed = embeds.SkillOfTheWeek(models.HiscoreField(activity), womResp.Competition.ID)
	}

	eventAbbrev := "botw"
	if eventType == models.EventTypeSkillOfTheWeek {
		eventAbbrev = "sotw"
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Register",
					Style:    discordgo.SuccessButton,
					CustomID: fmt.Sprintf("register-for-%s:%d,%s", eventAbbrev, womResp.Competition.ID, thread.ID),
				},
				discordgo.Button{
					Label:    "List Participants",
					Style:    discordgo.SecondaryButton,
					CustomID: fmt.Sprintf("list-participants-%s:%d", eventAbbrev, womResp.Competition.ID),
				},
			},
		},
	}

	// Send event announcement
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Printf("Error sending event announcement: %v", err)
		return err
	}

	return nil
}

// RegisterForEvent handles registration button clicks
func (t *TrackableCommands) RegisterForEvent(s *discordgo.Session, i *discordgo.InteractionCreate, womCompetitionID int64, threadID string, eventType models.EventType) error {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("defer response: %w", err)
	}

	// Get WOM competition from database
	comp, err := t.DB.GetWOMCompetitionByWOMID(ctx, womCompetitionID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Competition not found. It may have been deleted.",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return err
	}

	// Get user's linked account
	discordID, err := strconv.ParseInt(i.Member.User.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("parse discord id: %w", err)
	}

	link, err := t.DB.GetAccountLinkByDiscordID(ctx, discordID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "You need to link your RuneScape account first using `/link-rsn`",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return err
	}

	// Update player on WOM (ensures fresh data)
	_, err = t.WOMClient.UpdatePlayer(ctx, link.RunescapeName)
	if err != nil {
		log.Printf("Warning: failed to update player %s: %v", link.RunescapeName, err)
	}

	// Add participant to WOM competition
	resp, err := t.WOMClient.AddParticipants(ctx, womCompetitionID, []string{link.RunescapeName}, comp.VerificationCode)
	if err != nil {
		log.Printf("Error adding participant to WOM: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to register for competition. Please try again."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return err
	}

	// Send confirmation
	message := fmt.Sprintf("Registered **%s** for **%s**! %s", link.RunescapeName, FormatActivityName(comp.Metric), resp.Message)

	// Post in thread
	_, err = s.ChannelMessageSend(threadID, message)
	if err != nil {
		log.Printf("Error sending message to thread: %v", err)
	}

	// Send ephemeral confirmation
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: message,
		Flags:   discordgo.MessageFlagsEphemeral,
	})

	return nil
}

// ListParticipants shows the list of participants for a WOM competition
func (t *TrackableCommands) ListParticipants(s *discordgo.Session, i *discordgo.InteractionCreate, womCompetitionID int64) error {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		return fmt.Errorf("defer response: %w", err)
	}

	// Get competition from WOM
	competition, err := t.WOMClient.GetCompetition(ctx, womCompetitionID)
	if err != nil {
		log.Printf("Error fetching competition: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to fetch competition details."),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return err
	}

	if len(competition.Participations) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "No participants yet! Be the first to register!",
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return nil
	}

	// Build participant list with current standings
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("**Participants for %s:**\n\n", competition.Title))
	for i, p := range competition.Participations {
		gained := int64(0)
		if p.Progress != nil {
			gained = p.Progress.Gained
		}
		msg.WriteString(fmt.Sprintf("%d. %s - %d gained\n", i+1, p.Player.DisplayName, gained))
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg.String(),
		Flags:   discordgo.MessageFlagsEphemeral,
	})

	return nil
}

// FinishEvent ends a WOM competition and announces winners
func (t *TrackableCommands) FinishEvent(s *discordgo.Session, i *discordgo.InteractionCreate, eventType models.EventType) error {
	ctx := context.Background()

	// Defer the response
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return fmt.Errorf("defer response: %w", err)
	}

	// Get latest competition of this type
	eventTypeStr := string(eventType)
	comp, err := t.DB.GetLatestWOMCompetitionByType(ctx, eventTypeStr)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("There's no active %s competition ongoing!", getEventDisplayName(eventType)),
		})
		return err
	}

	// Fetch competition details from WOM
	competition, err := t.WOMClient.GetCompetition(ctx, comp.WomCompetitionID)
	if err != nil {
		log.Printf("Error fetching competition: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{
				embeds.ErrorEmbed("Failed to fetch competition details. Please try again."),
			},
		})
		return err
	}

	if len(competition.Participations) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Sadly there were no participants this time! :(",
		})
		return nil
	}

	// Sort participations by progress (WOM API should return them sorted)
	// Get top 3 winners
	winnersData := make([]embeds.WinnerData, 0, 3)
	for i, p := range competition.Participations {
		if i >= 3 {
			break
		}

		gained := int64(0)
		if p.Progress != nil {
			gained = p.Progress.Gained
		}

		if gained > 0 {
			// Find Discord ID for this player
			link, err := t.DB.GetAccountLinkByUsername(ctx, p.Player.Username)
			discordID := uint64(0)
			if err == nil {
				discordID = uint64(link.DiscordMemberID)
			}

			winnersData = append(winnersData, embeds.WinnerData{
				Username:  p.Player.DisplayName,
				DiscordID: discordID,
				Progress:  gained,
			})
		}
	}

	if len(winnersData) == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "No one made any progress during this competition!",
		})
		return nil
	}

	// Announce winner
	firstPlace := winnersData[0]
	unit := "KC"
	if eventType == models.EventTypeSkillOfTheWeek {
		unit = "XP"
	}

	content := fmt.Sprintf("Winner of this week's %s is **%s** with **%d %s**! Congratulations!",
		getEventDisplayName(eventType),
		firstPlace.Username,
		firstPlace.Progress,
		unit)

	if firstPlace.DiscordID > 0 {
		content = fmt.Sprintf("Winner of this week's %s is <@%d> with **%d %s**! Congratulations!",
			getEventDisplayName(eventType),
			firstPlace.DiscordID,
			firstPlace.Progress,
			unit)
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		log.Printf("Error sending winner announcement: %v", err)
	}

	// Send winners embed
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embeds.EventWinners(eventType, models.HiscoreField(comp.Metric), winnersData),
		},
	})

	return err
}

func getEventDisplayName(eventType models.EventType) string {
	switch eventType {
	case models.EventTypeBossOfTheWeek:
		return "Boss of the Week"
	case models.EventTypeSkillOfTheWeek:
		return "Skill of the Week"
	default:
		return string(eventType)
	}
}
