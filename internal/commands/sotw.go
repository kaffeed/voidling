package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/internal/models"
)

// HandleSOTWStart handles /sotw start command.
func (t *TrackableCommands) HandleSOTWStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Options[0] is the subcommand, Options[0].Options[0] is the skill parameter
	skill := data.Options[0].Options[0].StringValue()

	err := t.StartEvent(s, i, models.EventTypeSkillOfTheWeek, skill)
	if err != nil {
		return
	}
}

// HandleSOTWFinish handles /sotw finish command.
func (t *TrackableCommands) HandleSOTWFinish(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := t.FinishEvent(s, i, models.EventTypeSkillOfTheWeek)
	if err != nil {
		return
	}
}
