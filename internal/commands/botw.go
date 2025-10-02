package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/internal/models"
)

// HandleBOTWWildy handles /botw wildy command
func (t *TrackableCommands) HandleBOTWWildy(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Options[0] is the subcommand, Options[0].Options[0] is the boss parameter
	boss := data.Options[0].Options[0].StringValue()

	err := t.StartEvent(s, i, models.EventTypeBossOfTheWeek, boss)
	if err != nil {
		return
	}
}

// HandleBOTWGroup handles /botw group command
func (t *TrackableCommands) HandleBOTWGroup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Options[0] is the subcommand, Options[0].Options[0] is the boss parameter
	boss := data.Options[0].Options[0].StringValue()

	err := t.StartEvent(s, i, models.EventTypeBossOfTheWeek, boss)
	if err != nil {
		return
	}
}

// HandleBOTWQuest handles /botw quest command
func (t *TrackableCommands) HandleBOTWQuest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Options[0] is the subcommand, Options[0].Options[0] is the boss parameter
	boss := data.Options[0].Options[0].StringValue()

	err := t.StartEvent(s, i, models.EventTypeBossOfTheWeek, boss)
	if err != nil {
		return
	}
}

// HandleBOTWSlayer handles /botw slayer command
func (t *TrackableCommands) HandleBOTWSlayer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Options[0] is the subcommand, Options[0].Options[0] is the boss parameter
	boss := data.Options[0].Options[0].StringValue()

	err := t.StartEvent(s, i, models.EventTypeBossOfTheWeek, boss)
	if err != nil {
		return
	}
}

// HandleBOTWWorld handles /botw world command
func (t *TrackableCommands) HandleBOTWWorld(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Options[0] is the subcommand, Options[0].Options[0] is the boss parameter
	boss := data.Options[0].Options[0].StringValue()

	err := t.StartEvent(s, i, models.EventTypeBossOfTheWeek, boss)
	if err != nil {
		return
	}
}

// HandleBOTWFinish handles /botw finish command
func (t *TrackableCommands) HandleBOTWFinish(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := t.FinishEvent(s, i, models.EventTypeBossOfTheWeek)
	if err != nil {
		return
	}
}
