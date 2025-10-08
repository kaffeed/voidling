package bot

import (
	"context"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// PermissionLevel represents the required permission level for a command.
type PermissionLevel int

const (
	// PermissionEveryone - anyone can use.
	PermissionEveryone PermissionLevel = iota
	// PermissionCoordinator - requires Coordinator role.
	PermissionCoordinator
	// PermissionAdmin - requires administrator permission.
	PermissionAdmin
)

// HasPermission checks if a user has the required permission level.
func (b *Bot) HasPermission(s *discordgo.Session, i *discordgo.InteractionCreate, level PermissionLevel) bool {
	// Everyone level always passes
	if level == PermissionEveryone {
		return true
	}

	// Get guild member
	member := i.Member
	if member == nil {
		log.Printf("No member data in interaction")
		return false
	}

	// Check for administrator permission
	if level == PermissionAdmin {
		perms, err := s.UserChannelPermissions(member.User.ID, i.ChannelID)
		if err != nil {
			log.Printf("Error checking permissions: %v", err)
			return false
		}
		return perms&discordgo.PermissionAdministrator != 0
	}

	// Check for Coordinator role
	if level == PermissionCoordinator {
		// First, check guild config from database
		ctx := context.Background()
		guildID, err := strconv.ParseInt(i.GuildID, 10, 64)
		if err == nil {
			config, err := b.DB.GetGuildConfig(ctx, guildID)
			if err == nil && config.CoordinatorRoleID.Valid {
				roleIDStr := strconv.FormatInt(config.CoordinatorRoleID.Int64, 10)
				return hasRole(member, roleIDStr)
			}
		}

		// Fallback: If coordinator role ID is configured in env, use it
		if b.Config.CoordinatorRoleID != "" {
			return hasRole(member, b.Config.CoordinatorRoleID)
		}

		// Otherwise, check for role named "Coordinator" (case-insensitive)
		guild, err := s.Guild(i.GuildID)
		if err != nil {
			log.Printf("Error fetching guild: %v", err)
			return false
		}

		// Find Coordinator role
		var coordinatorRoleID string
		for _, role := range guild.Roles {
			if role.Name == "Coordinator" || role.Name == "coordinator" {
				coordinatorRoleID = role.ID
				break
			}
		}

		if coordinatorRoleID == "" {
			log.Printf("Coordinator role not found in guild")
			return false
		}

		return hasRole(member, coordinatorRoleID)
	}

	return false
}

// hasRole checks if a member has a specific role.
func hasRole(member *discordgo.Member, roleID string) bool {
	for _, r := range member.Roles {
		if r == roleID {
			return true
		}
	}
	return false
}

// RequirePermission wraps a handler with permission checking.
func (b *Bot) RequirePermission(level PermissionLevel, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if !b.HasPermission(s, i, level) {
			// Send permission denied message
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ You don't have permission to use this command. Coordinator role required.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		// Permission granted, call the handler
		handler(s, i)
	}
}
