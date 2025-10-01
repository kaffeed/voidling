package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidbound/config"
	"github.com/kaffeed/voidbound/internal/commands"
	"github.com/kaffeed/voidbound/internal/database"
	"github.com/kaffeed/voidbound/internal/models"
	"github.com/kaffeed/voidbound/internal/wiseoldman"
)

// Bot represents the Discord bot
type Bot struct {
	Session         *discordgo.Session
	Config          *config.Config
	DB              *database.Queries
	DBSQL           *sql.DB
	WOMClient       *wiseoldman.Client
	GuildID         string
	commands        []*discordgo.ApplicationCommand
	handlers        map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registerCmds    *commands.RegisterCommands
	trackableCmds   *commands.TrackableCommands
	schedulableCmds *commands.SchedulableCommands
	configCmds      *commands.ConfigCommands
}

// New creates a new Bot instance
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
		handlers:        make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
		registerCmds:    commands.NewRegisterCommands(db, dbSQL, womClient),
		trackableCmds:   commands.NewTrackableCommands(db, dbSQL, womClient),
		schedulableCmds: commands.NewSchedulableCommands(db, dbSQL),
		configCmds:      commands.NewConfigCommands(db, dbSQL),
	}

	// Register interaction handler
	session.AddHandler(bot.interactionHandler)

	return bot, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	b.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

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

// Stop stops the bot
func (b *Bot) Stop() error {
	// Unregister commands
	if err := b.unregisterCommands(); err != nil {
		log.Printf("Error unregistering commands: %v", err)
	}

	return b.Session.Close()
}

// registerCommands registers all slash commands
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
					Description: "What activity (e.g., Corporeal Beast, Nex)",
					Required:    true,
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

// unregisterCommands removes all registered commands
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
	}

	return nil
}

// registerHandler registers a command handler
func (b *Bot) registerHandler(name string, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	b.handlers[name] = handler
}

// handleBOTWCommand routes BOTW subcommands
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

// handleSOTWCommand routes SOTW subcommands
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

// handleConfigCommand routes config subcommands
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
	default:
		log.Printf("Unknown config subcommand: %s", subcommand)
	}
}

// interactionHandler handles all interactions
func (b *Bot) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if handler, ok := b.handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	case discordgo.InteractionMessageComponent:
		b.handleComponentInteraction(s, i)
	case discordgo.InteractionModalSubmit:
		b.handleModalSubmit(s, i)
	}
}

// handleComponentInteraction handles button/select menu interactions
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
	default:
		log.Printf("Unknown component action: %s", action)
	}
}

// handleRegisterForEvent handles registration button clicks
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

	b.trackableCmds.RegisterForEvent(s, i, womCompetitionID, threadID, models.EventType(eventType))
}

// handleListParticipants handles list participants button clicks
func (b *Bot) handleListParticipants(s *discordgo.Session, i *discordgo.InteractionCreate, data string) {
	womCompetitionID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Printf("Invalid WOM competition ID: %s", data)
		return
	}

	b.trackableCmds.ListParticipants(s, i, womCompetitionID)
}

// handleModalSubmit handles modal submissions
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
