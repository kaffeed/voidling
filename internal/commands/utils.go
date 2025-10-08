package commands

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// respondToInteraction sends an initial response to an interaction.
func respondToInteraction(s *discordgo.Session, i *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	return s.InteractionRespond(i, resp)
}

// sendFollowup sends a followup message to an interaction.
// Errors are logged internally, so callers can safely ignore the return value.
func sendFollowup(s *discordgo.Session, i *discordgo.Interaction, params *discordgo.WebhookParams) (*discordgo.Message, error) {
	msg, err := s.FollowupMessageCreate(i, true, params)
	if err != nil {
		slog.Error("failed to send Discord followup message",
			"error", err,
			"interaction_id", i.ID,
			"interaction_type", i.Type,
		)
	}
	return msg, err
}
