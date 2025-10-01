package embeds

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidbound/internal/models"
	"github.com/kaffeed/voidbound/internal/wiseoldman"
)

// Color constants for embeds
const (
	ColorSuccess = 0x2ecc71 // Green
	ColorError   = 0xe74c3c // Red
	ColorInfo    = 0x3498db // Blue
	ColorWarning = 0xf39c12 // Orange
	ColorBOTW    = 0x9b59b6 // Purple
	ColorSOTW    = 0x1abc9c // Turquoise
	ColorMass    = 0xe67e22 // Orange
	ColorWildy   = 0xc0392b // Dark Red
)

// PlayerInfo creates an embed showing player information from Wise Old Man
func PlayerInfo(player *wiseoldman.Player) *discordgo.MessageEmbed {
	if player == nil {
		return &discordgo.MessageEmbed{
			Title:       "Player Not Found",
			Description: "Unable to fetch player data from Wise Old Man.",
			Color:       ColorError,
			Timestamp:   time.Now().Format(time.RFC3339),
		}
	}

	// Get overall stats
	overall := player.GetSkill("overall")
	if overall == nil {
		return &discordgo.MessageEmbed{
			Title:       "Data Unavailable",
			Description: "Player snapshot data is not available.",
			Color:       ColorError,
			Timestamp:   time.Now().Format(time.RFC3339),
		}
	}

	// Create fields for key stats
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Total Level",
			Value:  fmt.Sprintf("%d", overall.Level),
			Inline: true,
		},
		{
			Name:   "Total XP",
			Value:  fmt.Sprintf("%s", formatNumber(overall.Experience)),
			Inline: true,
		},
		{
			Name:   "Rank",
			Value:  fmt.Sprintf("#%s", formatNumber(int64(overall.Rank))),
			Inline: true,
		},
		{
			Name:   "Combat Level",
			Value:  fmt.Sprintf("%d", player.CombatLevel),
			Inline: true,
		},
		{
			Name:   "EHP",
			Value:  fmt.Sprintf("%.1f", player.EHP),
			Inline: true,
		},
		{
			Name:   "EHB",
			Value:  fmt.Sprintf("%.1f", player.EHB),
			Inline: true,
		},
	}

	// Add some notable skill levels (combat stats)
	notableSkills := []string{
		"attack",
		"strength",
		"defence",
		"hitpoints",
		"ranged",
		"magic",
	}

	skillValues := ""
	for _, skillName := range notableSkills {
		if skill := player.GetSkill(skillName); skill != nil {
			// Capitalize first letter
			displayName := skillName
			if len(displayName) > 0 {
				displayName = string(displayName[0]-32) + displayName[1:]
			}
			skillValues += fmt.Sprintf("**%s**: %d\n", displayName, skill.Level)
		}
	}

	if skillValues != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Combat Stats",
			Value:  skillValues,
			Inline: false,
		})
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("OSRS Player: %s", player.DisplayName),
		Description: "Is this your account?",
		Color:       ColorInfo,
		Fields:      fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://oldschool.runescape.wiki/images/thumb/OSRS_icon.png/200px-OSRS_icon.png",
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Data from Wise Old Man",
		},
	}
}

// BossOfTheWeek creates an embed for Boss of the Week events
func BossOfTheWeek(activity models.HiscoreField, womCompetitionID int64) *discordgo.MessageEmbed {
	womURL := fmt.Sprintf("https://wiseoldman.net/competitions/%d", womCompetitionID)
	return &discordgo.MessageEmbed{
		Title:       "🏆 Boss of the Week",
		Description: fmt.Sprintf("This week's boss challenge: **%s**\n\n[View Competition on Wise Old Man](%s)\n\nClick the button below to register and track your kill count!", activity, womURL),
		Color:       ColorBOTW,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://oldschool.runescape.wiki/images/OSRS_icon.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "How it works",
				Value:  "Register to lock in your starting KC. At the end of the week, we'll check your progress and crown the winner!",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// SkillOfTheWeek creates an embed for Skill of the Week events
func SkillOfTheWeek(activity models.HiscoreField, womCompetitionID int64) *discordgo.MessageEmbed {
	womURL := fmt.Sprintf("https://wiseoldman.net/competitions/%d", womCompetitionID)
	return &discordgo.MessageEmbed{
		Title:       "📚 Skill of the Week",
		Description: fmt.Sprintf("This week's skill challenge: **%s**\n\n[View Competition on Wise Old Man](%s)\n\nClick the button below to register and track your experience gains!", activity, womURL),
		Color:       ColorSOTW,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://oldschool.runescape.wiki/images/OSRS_icon.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "How it works",
				Value:  "Register to lock in your starting XP. At the end of the week, we'll check your progress and crown the winner!",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// EventWinners creates an embed showing event winners
func EventWinners(eventType models.EventType, activity models.HiscoreField, winners []WinnerData) *discordgo.MessageEmbed {
	title := ""
	switch eventType {
	case models.EventTypeBossOfTheWeek:
		title = "🏆 Boss of the Week - Winners"
	case models.EventTypeSkillOfTheWeek:
		title = "📚 Skill of the Week - Winners"
	}

	description := fmt.Sprintf("**%s** has concluded!\n\nHere are the top performers:", activity)

	fields := []*discordgo.MessageEmbedField{}
	medals := []string{"🥇", "🥈", "🥉"}

	for i, winner := range winners {
		if i >= 3 {
			break
		}
		medal := medals[i]

		unit := "KC"
		if eventType == models.EventTypeSkillOfTheWeek {
			unit = "XP"
		}

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s %s", medal, winner.Username),
			Value: fmt.Sprintf("Progress: **%s %s**\nStart: %s | End: %s",
				formatNumber(winner.Progress),
				unit,
				formatNumber(winner.StartPoint),
				formatNumber(winner.EndPoint)),
			Inline: false,
		})
	}

	color := ColorBOTW
	if eventType == models.EventTypeSkillOfTheWeek {
		color = ColorSOTW
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Fields:      fields,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}

// MassEvent creates an embed for mass events
func MassEvent(activity string, location string, scheduledAt time.Time) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "⚔️ Mass Event",
		Description: fmt.Sprintf("Join us for **%s**!", activity),
		Color:       ColorMass,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Location",
				Value:  location,
				Inline: true,
			},
			{
				Name:   "Time",
				Value:  fmt.Sprintf("<t:%d:F>", scheduledAt.Unix()),
				Inline: true,
			},
			{
				Name:   "Countdown",
				Value:  fmt.Sprintf("<t:%d:R>", scheduledAt.Unix()),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://oldschool.runescape.wiki/images/OSRS_icon.png",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// ScheduledEventReminder creates an embed for event reminders
func ScheduledEventReminder(eventType models.EventType, activity models.HiscoreField, location string, scheduledAt time.Time) *discordgo.MessageEmbed {
	title := "📅 Event Reminder"
	color := ColorInfo

	if eventType == models.EventTypeWildyWednesday {
		title = "💀 Wildy Wednesday Reminder"
		color = ColorWildy
	} else if eventType == models.EventTypeMass {
		title = "⚔️ Mass Event Reminder"
		color = ColorMass
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: fmt.Sprintf("**%s** is starting soon at **%s**!", activity, location),
		Color:       color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Time",
				Value:  fmt.Sprintf("<t:%d:F>", scheduledAt.Unix()),
				Inline: false,
			},
			{
				Name:   "Starts in",
				Value:  fmt.Sprintf("<t:%d:R>", scheduledAt.Unix()),
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Don't forget to show up!",
		},
	}
}

// ErrorEmbed creates a generic error embed
func ErrorEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "❌ Error",
		Description: message,
		Color:       ColorError,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}

// SuccessEmbed creates a generic success embed
func SuccessEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "✅ Success",
		Description: message,
		Color:       ColorSuccess,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}

// WinnerData holds winner information for display
type WinnerData struct {
	Username   string
	StartPoint int64
	EndPoint   int64
	Progress   int64
	DiscordID  uint64
}

// formatNumber formats numbers with commas for readability
func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}
