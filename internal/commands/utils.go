package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// respondToInteraction sends an initial response to an interaction
func respondToInteraction(s *discordgo.Session, i *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	return s.InteractionRespond(i, resp)
}

// sendFollowup sends a followup message to an interaction
func sendFollowup(s *discordgo.Session, i *discordgo.Interaction, wait bool, params *discordgo.WebhookParams) (*discordgo.Message, error) {
	msg, err := s.FollowupMessageCreate(i, wait, params)
	if err != nil {
		log.Printf("Error sending followup message: %v", err)
	}
	return msg, err
}
