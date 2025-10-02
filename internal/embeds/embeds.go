package embeds

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaffeed/voidling/internal/models"
	"github.com/kaffeed/voidling/internal/wiseoldman"
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

	// Get boss info if available
	bossInfo, hasBossInfo := GetBossInfo(string(activity))

	description := fmt.Sprintf("This week's boss challenge: **%s**\n\n", activity)

	// Add boss description if available
	if hasBossInfo {
		description += fmt.Sprintf("*%s*\n\n", bossInfo.Description)
	}

	description += fmt.Sprintf("[View Competition on Wise Old Man](%s)\n\nClick the button below to register and track your kill count!", womURL)

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "How it works",
			Value:  "Register to lock in your starting KC. At the end of the week, we'll check your progress and crown the winner!",
			Inline: false,
		},
	}

	// Add strategy guide link if available
	if hasBossInfo {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "📖 Strategy Guide",
			Value:  fmt.Sprintf("[OSRS Wiki - %s Strategy](%s)", activity, bossInfo.WikiURL),
			Inline: false,
		})
	}

	return &discordgo.MessageEmbed{
		Title:       "🏆 Boss of the Week",
		Description: description,
		Color:       ColorBOTW,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://oldschool.runescape.wiki/images/OSRS_icon.png",
		},
		Fields:    fields,
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

// CompetitionCodeEmbed creates an embed for WOM competition verification codes
func CompetitionCodeEmbed(eventName string, verificationCode string, competitionID int64) *discordgo.MessageEmbed {
	womURL := fmt.Sprintf("https://wiseoldman.net/competitions/%d", competitionID)

	return &discordgo.MessageEmbed{
		Title:       "🔑 Competition Verification Code",
		Description: fmt.Sprintf("**%s**", eventName),
		Color:       ColorBOTW,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Verification Code",
				Value:  fmt.Sprintf("`%s`", verificationCode),
				Inline: false,
			},
			{
				Name:   "Wise Old Man Link",
				Value:  fmt.Sprintf("[View Competition](%s)", womURL),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use this code to manage the competition on Wise Old Man",
		},
		Timestamp: time.Now().Format(time.RFC3339),
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

// getBossImageURL returns the OSRS Wiki image URL for a boss/activity
func getBossImageURL(activity string) string {
	// Map of activity names (snake_case) to OSRS Wiki image URLs
	bossImages := map[string]string{
		// PvM Bosses
		"corporeal_beast":      "https://oldschool.runescape.wiki/images/thumb/Corporeal_Beast.png/270px-Corporeal_Beast.png",
		"nex":                  "https://oldschool.runescape.wiki/images/thumb/Nex.png/270px-Nex.png",
		"nightmare":            "https://oldschool.runescape.wiki/images/thumb/The_Nightmare.png/250px-The_Nightmare.png",
		"phosanis_nightmare":   "https://oldschool.runescape.wiki/images/thumb/Phosani%27s_Nightmare.png/250px-Phosani%27s_Nightmare.png",
		"commander_zilyana":    "https://oldschool.runescape.wiki/images/thumb/Commander_Zilyana.png/250px-Commander_Zilyana.png",
		"kril_tsutsaroth":      "https://oldschool.runescape.wiki/images/thumb/K%27ril_Tsutsaroth.png/250px-K%27ril_Tsutsaroth.png",
		"general_graardor":     "https://oldschool.runescape.wiki/images/thumb/General_Graardor.png/250px-General_Graardor.png",
		"kreearra":             "https://oldschool.runescape.wiki/images/thumb/Kree%27arra.png/250px-Kree%27arra.png",
		"godwars":              "https://oldschool.runescape.wiki/images/thumb/God_Wars_Dungeon.png/300px-God_Wars_Dungeon.png",
		// Raids
		"theatre_of_blood":     "https://oldschool.runescape.wiki/images/thumb/Theatre_of_Blood_logo.png/250px-Theatre_of_Blood_logo.png",
		"chambers_of_xeric":    "https://oldschool.runescape.wiki/images/thumb/Chambers_of_Xeric_logo.png/250px-Chambers_of_Xeric_logo.png",
		"tombs_of_amascut":     "https://oldschool.runescape.wiki/images/thumb/Tombs_of_Amascut.png/300px-Tombs_of_Amascut.png",
		// Wilderness Bosses
		"king_black_dragon":    "https://oldschool.runescape.wiki/images/thumb/King_Black_Dragon.png/280px-King_Black_Dragon.png",
		"scorpia":              "https://oldschool.runescape.wiki/images/thumb/Scorpia.png/300px-Scorpia.png",
		"artio":                "https://oldschool.runescape.wiki/images/thumb/Artio.png/250px-Artio.png",
		"callisto":             "https://oldschool.runescape.wiki/images/thumb/Callisto.png/300px-Callisto.png",
		"calvarion":            "https://oldschool.runescape.wiki/images/thumb/Calvarion.png/300px-Calvarion.png",
		"chaos_elemental":      "https://oldschool.runescape.wiki/images/thumb/Chaos_Elemental.png/280px-Chaos_Elemental.png",
		"chaos_fanatic":        "https://oldschool.runescape.wiki/images/thumb/Chaos_Fanatic.png/200px-Chaos_Fanatic.png",
		"crazy_archaeologist":  "https://oldschool.runescape.wiki/images/thumb/Crazy_Archaeologist.png/200px-Crazy_Archaeologist.png",
		"spindel":              "https://oldschool.runescape.wiki/images/thumb/Spindel.png/300px-Spindel.png",
		"venenatis":            "https://oldschool.runescape.wiki/images/thumb/Venenatis.png/300px-Venenatis.png",
		"vetion":               "https://oldschool.runescape.wiki/images/thumb/Vet%27ion.png/250px-Vet%27ion.png",
		// Skilling Bosses
		"tempoross":            "https://oldschool.runescape.wiki/images/thumb/Tempoross.png/280px-Tempoross.png",
		"wintertodt":           "https://oldschool.runescape.wiki/images/thumb/Wintertodt.png/300px-Wintertodt.png",
		"guardians_of_the_rift": "https://oldschool.runescape.wiki/images/thumb/The_Great_Guardian.png/250px-The_Great_Guardian.png",
	}

	// Return boss-specific image or default OSRS icon
	if url, exists := bossImages[activity]; exists {
		return url
	}
	return "https://oldschool.runescape.wiki/images/OSRS_icon.png"
}

// MassEventWithTimezone creates an embed for mass events with timezone information
func MassEventWithTimezone(activity, location string, scheduledTime time.Time, timezone string) *discordgo.MessageEmbed {
	// Get timezone abbreviation
	tzAbbrev := scheduledTime.Format("MST")

	// Create Discord timestamp (auto-converts to user's local time)
	discordTimestamp := fmt.Sprintf("<t:%d:F>", scheduledTime.Unix())
	relativeTime := fmt.Sprintf("<t:%d:R>", scheduledTime.Unix())

	// Format activity name for display (convert snake_case to Title Case)
	displayName := formatActivityName(activity)

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("⚔️ Mass Event: %s", displayName),
		Description: fmt.Sprintf("**Location:** %s\n\n**Time:** %s\n**Starts:** %s\n\n**📢 Important:** Click the \"I'll Participate\" button below to register for this event! This helps us plan and you'll get a reminder before the event starts.", location, discordTimestamp, relativeTime),
		Color:       ColorMass,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: getBossImageURL(activity),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Scheduled in %s (%s)", timezone, tzAbbrev),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// formatActivityName converts snake_case to Title Case for display
func formatActivityName(activity string) string {
	result := ""
	capitalize := true
	for _, c := range activity {
		if c == '_' {
			result += " "
			capitalize = true
		} else if capitalize {
			if c >= 'a' && c <= 'z' {
				result += string(c - 32) // Convert to uppercase
			} else {
				result += string(c)
			}
			capitalize = false
		} else {
			result += string(c)
		}
	}
	return result
}
