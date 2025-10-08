// Package bot provides the core Discord bot implementation for Voidling.
// It handles command registration, interaction routing, and event processing.
package bot

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/config"
	"github.com/kaffeed/voidling/internal/commands"
	"github.com/kaffeed/voidling/internal/database"
	"github.com/kaffeed/voidling/internal/embeds"
	"github.com/kaffeed/voidling/internal/models"
	"github.com/kaffeed/voidling/internal/timezone"
	"github.com/kaffeed/voidling/internal/wiseoldman"
)

type handlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Bot represents the.
type Bot struct {
	Session         *discordgo.Session
	Config          *config.Config
	DB              *database.Queries
	DBSQL           *sql.DB
	WOMClient       *wiseoldman.Client
	GuildID         string
	commands        []*discordgo.ApplicationCommand
	handlers        map[string]handlerFunc
	registerCmds    *commands.RegisterCommands
	trackableCmds   *commands.TrackableCommands
	schedulableCmds *commands.SchedulableCommands
	configCmds      *commands.ConfigCommands
}

// New creates a new Bot instance.
func New(cfg *config.Config, db *database.Queries, dbSQL *sql.DB) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	womClient := wiseoldman.NewClient()

	bot := &Bot{
		Session:         session,
		Config:          cfg,
		DB:              db,
		DBSQL:           dbSQL,
		WOMClient:       womClient,
		GuildID:         cfg.GuildID,
		handlers:        make(map[string]handlerFunc),
		registerCmds:    commands.NewRegisterCommands(db, dbSQL, womClient),
		trackableCmds:   commands.NewTrackableCommands(db, dbSQL, womClient),
		schedulableCmds: commands.NewSchedulableCommands(db, dbSQL),
		configCmds:      commands.NewConfigCommands(db, dbSQL),
	}

	// Register interaction handler
	session.AddHandler(bot.interactionHandler)

	// Register guild member add handler for auto-greeting
	session.AddHandler(bot.handleGuildMemberAdd)

	return bot, nil
}

// Start starts the bot.
func (b *Bot) Start() error {
	b.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}

	log.Printf("Bot is now running as %s", b.Session.State.User.Username)

	// Register commands
	if err := b.registerCommands(); err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	return nil
}

// Stop stops the bot.
func (b *Bot) Stop() error {
	// Unregister commands
	if err := b.unregisterCommands(); err != nil {
		log.Printf("Error unregistering commands: %v", err)
	}

	return b.Session.Close()
}

// registerCommands registers all slash commands.
func (b *Bot) registerCommands() error {
	// Define commands
	b.commands = []*discordgo.ApplicationCommand{
		{
			Name:        "link-rsn",
			Description: "Link your RuneScape account to Discord",
		},
		{
			Name:        "unlink-rsn",
			Description: "Unlink your RuneScape account from Discord",
		},
		{
			Name:        "botw",
			Description: "Boss of the Week commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "wildy",
					Description: "Start a Wilderness boss of the week",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "boss",
							Description: "Select a wilderness boss",
							Required:    true,
							Choices:     commands.WildyBossChoices(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "group",
					Description: "Start a Group boss of the week",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "boss",
							Description: "Select a group boss",
							Required:    true,
							Choices:     commands.GroupBossChoices(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "quest",
					Description: "Start a Quest boss of the week",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "boss",
							Description: "Select a quest boss",
							Required:    true,
							Choices:     commands.QuestBossChoices(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "slayer",
					Description: "Start a Slayer boss of the week",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "boss",
							Description: "Select a slayer boss",
							Required:    true,
							Choices:     commands.SlayerBossChoices(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "world",
					Description: "Start a World boss of the week",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "boss",
							Description: "Select a world boss",
							Required:    true,
							Choices:     commands.WorldBossChoices(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "finish",
					Description: "Finish the current Boss of the Week event",
				},
			},
		},
		{
			Name:        "sotw",
			Description: "Skill of the Week commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "start",
					Description: "Start a Skill of the Week event",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "skill",
							Description: "Select a skill",
							Required:    true,
							Choices:     commands.SkillChoices(),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "finish",
					Description: "Finish the current Skill of the Week event",
				},
			},
		},
		{
			Name:        "mass",
			Description: "Schedule a mass event",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "activity",
					Description: "Select the boss or activity",
					Required:    true,
					Choices:     commands.MassBossChoices(),
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "location",
					Description: "Where to meet (e.g., World 444)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "When to start (YYYY-MM-DD HH:MM format)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Event duration in minutes (e.g., 60, 120)",
					Required:    true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "timezone",
					Description:  "Timezone for the event time (optional, uses your preference or server default)",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "config",
			Description: "Server configuration commands (Owner/Admin only)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-coordinator-role",
					Description: "Set the role that can manage events",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionRole,
							Name:        "role",
							Description: "The role to assign coordinator permissions",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "show",
					Description: "Show current server configuration",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-competition-code-channel",
					Description: "Set the channel to send competition verification codes to",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "The channel to send competition codes to",
							Required:    true,
							ChannelTypes: []discordgo.ChannelType{
								discordgo.ChannelTypeGuildText,
							},
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-default-timezone",
					Description: "Set the default timezone for this server",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "timezone",
							Description:  "Timezone to use as server default (e.g., America/New_York)",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-my-timezone",
					Description: "Set your personal timezone preference",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "timezone",
							Description:  "Your preferred timezone (e.g., America/New_York)",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-event-notification-channel",
					Description: "Set the channel to post event embeds when events (BOTW, SOTW, Mass) are created",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "Channel for event embeds",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-event-notification-role",
					Description: "Set the role to ping when events (BOTW, SOTW, Mass) are created",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionRole,
							Name:        "role",
							Description: "The role to ping for event notifications",
							Required:    true,
						},
					},
				},
			},
		},
	}

	// Register command handlers
	b.registerHandler("link-rsn", b.registerCmds.HandleLinkRSN)
	b.registerHandler("unlink-rsn", b.registerCmds.HandleUnlinkRSN)
	b.registerHandler("botw", b.handleBOTWCommand)
	b.registerHandler("sotw", b.handleSOTWCommand)
	b.registerHandler("mass", b.schedulableCmds.HandleMassEvent)
	b.registerHandler("config", b.handleConfigCommand)

	// Register commands with Discord
	for _, cmd := range b.commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
		if err != nil {
			return fmt.Errorf("failed to create command %s: %w", cmd.Name, err)
		}
		log.Printf("Registered command: %s", cmd.Name)
	}

	return nil
}

// unregisterCommands removes all registered commands.
func (b *Bot) unregisterCommands() error {
	commands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, b.GuildID)
	if err != nil {
		return fmt.Errorf("failed to fetch commands: %w", err)
	}

	for _, cmd := range commands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, b.GuildID, cmd.ID)
		if err != nil {
			log.Printf("Failed to delete command %s: %v", cmd.Name, err)
		}
		log.Printf("Removed cmd %s", cmd.Name)
	}

	return nil
}

// registerHandler registers a command handler.
func (b *Bot) registerHandler(name string, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	b.handlers[name] = handler
}

// handleBOTWCommand routes BOTW subcommands.
func (b *Bot) handleBOTWCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return
	}

	subcommand := data.Options[0].Name

	// All BOTW commands require Coordinator permission
	if !b.HasPermission(s, i, PermissionCoordinator) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ You don't have permission to use this command. Coordinator role required.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	switch subcommand {
	case "wildy":
		b.trackableCmds.HandleBOTWWildy(s, i)
	case "group":
		b.trackableCmds.HandleBOTWGroup(s, i)
	case "quest":
		b.trackableCmds.HandleBOTWQuest(s, i)
	case "slayer":
		b.trackableCmds.HandleBOTWSlayer(s, i)
	case "world":
		b.trackableCmds.HandleBOTWWorld(s, i)
	case "finish":
		b.trackableCmds.HandleBOTWFinish(s, i)
	default:
		log.Printf("Unknown BOTW subcommand: %s", subcommand)
	}
}

// handleSOTWCommand routes SOTW subcommands.
func (b *Bot) handleSOTWCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return
	}

	subcommand := data.Options[0].Name

	// All SOTW commands require Coordinator permission
	if !b.HasPermission(s, i, PermissionCoordinator) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ You don't have permission to use this command. Coordinator role required.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	switch subcommand {
	case "start":
		b.trackableCmds.HandleSOTWStart(s, i)
	case "finish":
		b.trackableCmds.HandleSOTWFinish(s, i)
	default:
		log.Printf("Unknown SOTW subcommand: %s", subcommand)
	}
}

// handleConfigCommand routes config subcommands.
func (b *Bot) handleConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return
	}

	subcommand := data.Options[0].Name

	switch subcommand {
	case "set-coordinator-role":
		b.configCmds.HandleSetCoordinatorRole(s, i)
	case "show":
		b.configCmds.HandleShowConfig(s, i)
	case "set-competition-code-channel":
		b.configCmds.HandleSetCompetitionCodeChannel(s, i)
	case "set-default-timezone":
		b.configCmds.HandleSetDefaultTimezone(s, i)
	case "set-my-timezone":
		b.configCmds.HandleSetMyTimezone(s, i)
	case "set-event-notification-channel":
		b.configCmds.HandleSetEventNotificationChannel(s, i)
	case "set-event-notification-role":
		b.configCmds.HandleSetEventNotificationRole(s, i)
	default:
		log.Printf("Unknown config subcommand: %s", subcommand)
	}
}

// interactionHandler handles all interactions.
func (b *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if handler, ok := b.handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		b.handleTimezoneAutocomplete(s, i)
	case discordgo.InteractionMessageComponent:
		b.handleComponentInteraction(s, i)
	case discordgo.InteractionModalSubmit:
		b.handleModalSubmit(s, i)
	}
}

// handleComponentInteraction handles button/select menu interactions.
func (b *Bot) handleComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	// Parse custom ID: "action:data"
	parts := strings.SplitN(customID, ":", 2)
	if len(parts) < 1 {
		log.Printf("Invalid component custom ID: %s", customID)
		return
	}

	action := parts[0]
	data := ""
	if len(parts) == 2 {
		data = parts[1]
	}

	// Route based on action
	switch action {
	case "dm-link-rsn":
		// Handle DM link button - show modal (reuse existing handler)
		b.registerCmds.HandleLinkRSN(s, i)
	case "confirm-rsn":
		b.registerCmds.HandleConfirmRSN(s, i, data)
	case "cancel-rsn":
		b.registerCmds.HandleCancelRSN(s, i, data)
	case "register-for-botw":
		b.handleRegisterForEvent(s, i, data, "botw")
	case "register-for-sotw":
		b.handleRegisterForEvent(s, i, data, "sotw")
	case "list-participants-botw":
		b.handleListParticipants(s, i, data)
	case "list-participants-sotw":
		b.handleListParticipants(s, i, data)
	case "participate-mass":
		b.schedulableCmds.HandleParticipateInMass(s, i, data)
	case "list-participants-mass":
		b.schedulableCmds.HandleListParticipantsMass(s, i, data)
	default:
		log.Printf("Unknown component action: %s", action)
	}
}

// handleRegisterForEvent handles registration button clicks.
func (b *Bot) handleRegisterForEvent(s *discordgo.Session, i *discordgo.InteractionCreate, data string, eventTypeStr string) {
	// data format: "womCompetitionID,threadID"
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		log.Printf("Invalid register data format: %s", data)
		return
	}

	womCompetitionID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.Printf("Invalid WOM competition ID: %s", parts[0])
		return
	}

	threadID := parts[1]

	var eventType string
	if eventTypeStr == "botw" {
		eventType = string(models.EventTypeBossOfTheWeek)
	} else {
		eventType = string(models.EventTypeSkillOfTheWeek)
	}

	if err := b.trackableCmds.RegisterForEvent(s, i, womCompetitionID, threadID, models.EventType(eventType)); err != nil {
		slog.Error("failed to register user for trackable event",
			"error", err,
			"competition_id", womCompetitionID,
			"event_type", eventType,
		)
	}
}

// handleListParticipants handles list participants button clicks.
func (b *Bot) handleListParticipants(s *discordgo.Session, i *discordgo.InteractionCreate, data string) {
	womCompetitionID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Printf("Invalid WOM competition ID: %s", data)
		return
	}

	if err := b.trackableCmds.ListParticipants(s, i, womCompetitionID); err != nil {
		slog.Error("failed to list participants for trackable event",
			"error", err,
			"competition_id", womCompetitionID,
		)
	}
}

// handleModalSubmit handles modal submissions.
func (b *Bot) handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.ModalSubmitData().CustomID

	// Route modal submissions
	switch customID {
	case "link-rsn-modal":
		b.registerCmds.HandleLinkRSNModal(s, i)
	default:
		log.Printf("Unknown modal submit: %s", customID)
	}
}

// handleTimezoneAutocomplete handles timezone autocomplete.
func (b *Bot) handleTimezoneAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var focusedOption *discordgo.ApplicationCommandInteractionDataOption

	// Find the focused option
	for _, opt := range data.Options {
		if opt.Focused {
			focusedOption = opt
			break
		}
		// Check subcommand options
		for _, subOpt := range opt.Options {
			if subOpt.Focused {
				focusedOption = subOpt
				break
			}
		}
	}

	if focusedOption == nil {
		return
	}

	// Search timezones based on user input
	query := focusedOption.StringValue()
	matches := timezone.SearchTimezones(query)

	// Convert to Discord choices
	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(matches))
	for _, tz := range matches {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  tz,
			Value: tz,
		})
	}

	// Respond with filtered choices
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		log.Printf("Error responding to autocomplete: %v", err)
	}
}

// handleGuildMemberAdd sends a greeting DM to new members.
func (b *Bot) handleGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Get guild information for greeting
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Printf("Error fetching guild info for greeting: %v", err)
		return
	}

	// Create DM channel with the new member
	dmChannel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		log.Printf("Error creating DM channel for user %s: %v", m.User.Username, err)
		return
	}

	// Create greeting embed
	embed := embeds.WelcomeGreeting(guild.Name)

	// Create "Link My Account" button
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "🔗 Link My RuneScape Account",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("dm-link-rsn:%s", m.GuildID), // Include guild ID for context
				},
			},
		},
	}

	// Send greeting message in DM
	_, err = s.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Printf("Error sending greeting DM to user %s: %v", m.User.Username, err)
		// User likely has DMs disabled, fail silently
		return
	}

	log.Printf("Sent greeting DM to new member: %s (ID: %s) in guild: %s", m.User.Username, m.User.ID, guild.Name)
}
