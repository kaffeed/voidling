package testutil

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

// MockDiscordSession is a mock implementation of discordgo.Session for testing
type MockDiscordSession struct {
	mock.Mock
}

// InteractionRespond mocks the Discord interaction response
func (m *MockDiscordSession) InteractionRespond(i *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	args := m.Called(i, resp)
	return args.Error(0)
}

// FollowupMessageCreate mocks creating a followup message
func (m *MockDiscordSession) FollowupMessageCreate(i *discordgo.Interaction, wait bool, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	args := m.Called(i, wait, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

// InteractionResponseEdit mocks editing an interaction response
func (m *MockDiscordSession) InteractionResponseEdit(i *discordgo.Interaction, edit *discordgo.WebhookEdit) (*discordgo.Message, error) {
	args := m.Called(i, edit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

// GuildMemberNickname mocks updating a guild member's nickname
func (m *MockDiscordSession) GuildMemberNickname(guildID, userID, nickname string) error {
	args := m.Called(guildID, userID, nickname)
	return args.Error(0)
}

// ThreadStartComplex mocks creating a Discord thread
func (m *MockDiscordSession) ThreadStartComplex(channelID string, params *discordgo.ThreadStart) (*discordgo.Channel, error) {
	args := m.Called(channelID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Channel), args.Error(1)
}

// ChannelMessageSendComplex mocks sending a complex message
func (m *MockDiscordSession) ChannelMessageSendComplex(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	args := m.Called(channelID, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

// UserChannelCreate mocks creating a DM channel
func (m *MockDiscordSession) UserChannelCreate(userID string) (*discordgo.Channel, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Channel), args.Error(1)
}

// Guild mocks fetching guild information
func (m *MockDiscordSession) Guild(guildID string) (*discordgo.Guild, error) {
	args := m.Called(guildID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Guild), args.Error(1)
}

// CreateTestInteraction creates a test Discord interaction for command testing
func CreateTestInteraction(commandName string, userID string, guildID string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				ID:   "test-command-id",
				Name: commandName,
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID:       userID,
					Username: "testuser",
				},
			},
			GuildID:   guildID,
			ChannelID: "test-channel-123",
		},
	}
}

// CreateTestModalSubmit creates a test modal submission interaction
func CreateTestModalSubmit(customID string, userID string, componentValues map[string]string) *discordgo.InteractionCreate {
	components := make([]discordgo.MessageComponent, 0)

	for id, value := range componentValues {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID: id,
					Value:    value,
				},
			},
		})
	}

	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionModalSubmit,
			Data: discordgo.ModalSubmitInteractionData{
				CustomID:   customID,
				Components: components,
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID:       userID,
					Username: "testuser",
				},
			},
		},
	}
}

// CreateTestButtonInteraction creates a test button interaction
func CreateTestButtonInteraction(customID string, userID string, guildID string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: customID,
			},
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID:       userID,
					Username: "testuser",
				},
			},
			GuildID:   guildID,
			ChannelID: "test-channel-123",
			Message: &discordgo.Message{
				ID:        "test-message-123",
				ChannelID: "test-channel-123",
			},
		},
	}
}
