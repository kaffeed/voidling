package embeds

import (
	"testing"
	"time"

	"github.com/kaffeed/voidling/internal/models"
	"github.com/kaffeed/voidling/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerInfo(t *testing.T) {
	t.Run("valid player", func(t *testing.T) {
		player := testutil.CreateTestPlayer("TestUser")
		embed := PlayerInfo(player)

		require.NotNil(t, embed)
		assert.Equal(t, "OSRS Player: TestUser", embed.Title)
		assert.Equal(t, ColorInfo, embed.Color)
		assert.NotEmpty(t, embed.Fields, "Should have fields for valid player")
		assert.Contains(t, embed.Thumbnail.URL, "OSRS", "Should have OSRS thumbnail")
	})

	t.Run("nil player", func(t *testing.T) {
		embed := PlayerInfo(nil)

		require.NotNil(t, embed)
		assert.Equal(t, "Player Not Found", embed.Title)
		assert.Equal(t, ColorError, embed.Color)
	})
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"zero", 0, "0"},
		{"small number", 123, "123"},
		{"exactly thousand", 1000, "1,000"},
		{"ten thousand", 10000, "10,000"},
		{"hundred thousand", 100000, "100,000"},
		{"million", 1000000, "1,000,000"},
		{"max int experience", 200000000, "200,000,000"},
		{"single digit", 5, "5"},
		{"two digits", 42, "42"},
		{"three digits", 999, "999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber(tt.input)
			assert.Equal(t, tt.expected, result, "Format number mismatch")
		})
	}
}

func TestFormatActivityName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"snake_case single word", "nex", "Nex"},
		{"snake_case two words", "corporeal_beast", "Corporeal Beast"},
		{"snake_case three words", "chambers_of_xeric", "Chambers Of Xeric"},
		{"already capitalized", "Nex", "Nex"},
		{"mixed case", "king_Black_dragon", "King Black Dragon"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatActivityName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase", "attack", "Attack"},
		{"already capitalized", "Attack", "Attack"},
		{"single char", "a", "A"},
		{"empty string", "", ""},
		{"uppercase", "ATTACK", "ATTACK"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := capitalizeFirst(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected rune
	}{
		{"lowercase a", 'a', 'A'},
		{"lowercase z", 'z', 'Z'},
		{"already uppercase", 'A', 'A'},
		{"number", '5', '5'},
		{"special char", '@', '@'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toUpper(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBossOfTheWeek(t *testing.T) {
	activity := models.HiscoreField("corporeal_beast")
	womCompetitionID := int64(12345)

	embed := BossOfTheWeek(activity, womCompetitionID)

	require.NotNil(t, embed)
	assert.Equal(t, "🏆 Boss of the Week", embed.Title)
	assert.Equal(t, ColorBOTW, embed.Color)
	assert.Contains(t, embed.Description, "corporeal_beast")
	assert.Contains(t, embed.Description, "https://wiseoldman.net/competitions/12345")
	assert.NotEmpty(t, embed.Fields, "Should have fields")
}

func TestSkillOfTheWeek(t *testing.T) {
	activity := models.HiscoreField("woodcutting")
	womCompetitionID := int64(67890)

	embed := SkillOfTheWeek(activity, womCompetitionID)

	require.NotNil(t, embed)
	assert.Equal(t, "📚 Skill of the Week", embed.Title)
	assert.Equal(t, ColorSOTW, embed.Color)
	assert.Contains(t, embed.Description, "woodcutting")
	assert.Contains(t, embed.Description, "https://wiseoldman.net/competitions/67890")
	assert.NotEmpty(t, embed.Fields)
}

func TestEventWinners(t *testing.T) {
	winners := []WinnerData{
		{Username: "Player1", StartPoint: 100, EndPoint: 500, Progress: 400, DiscordID: 123},
		{Username: "Player2", StartPoint: 50, EndPoint: 300, Progress: 250, DiscordID: 456},
		{Username: "Player3", StartPoint: 0, EndPoint: 100, Progress: 100, DiscordID: 789},
	}

	t.Run("BOTW winners", func(t *testing.T) {
		embed := EventWinners(models.EventTypeBossOfTheWeek, "corporeal_beast", winners)

		require.NotNil(t, embed)
		assert.Equal(t, "🏆 Boss of the Week - Winners", embed.Title)
		assert.Equal(t, ColorBOTW, embed.Color)
		assert.Len(t, embed.Fields, 3, "Should have 3 winner fields")
		assert.Contains(t, embed.Fields[0].Name, "🥇")
		assert.Contains(t, embed.Fields[1].Name, "🥈")
		assert.Contains(t, embed.Fields[2].Name, "🥉")
	})

	t.Run("SOTW winners", func(t *testing.T) {
		embed := EventWinners(models.EventTypeSkillOfTheWeek, "woodcutting", winners)

		require.NotNil(t, embed)
		assert.Equal(t, "📚 Skill of the Week - Winners", embed.Title)
		assert.Equal(t, ColorSOTW, embed.Color)
	})

	t.Run("fewer than 3 winners", func(t *testing.T) {
		singleWinner := []WinnerData{winners[0]}
		embed := EventWinners(models.EventTypeBossOfTheWeek, "nex", singleWinner)

		require.NotNil(t, embed)
		assert.Len(t, embed.Fields, 1, "Should only have 1 field")
	})
}

func TestMassEvent(t *testing.T) {
	activity := "Nex"
	location := "World 416"
	scheduledAt := time.Now().Add(2 * time.Hour)

	embed := MassEvent(activity, location, scheduledAt)

	require.NotNil(t, embed)
	assert.Equal(t, "⚔️ Mass Event", embed.Title)
	assert.Equal(t, ColorMass, embed.Color)
	assert.Contains(t, embed.Description, activity)
	assert.Len(t, embed.Fields, 3, "Should have Location, Time, and Countdown fields")
}

func TestErrorEmbed(t *testing.T) {
	message := "Something went wrong"
	embed := ErrorEmbed(message)

	require.NotNil(t, embed)
	assert.Equal(t, "❌ Error", embed.Title)
	assert.Equal(t, ColorError, embed.Color)
	assert.Equal(t, message, embed.Description)
}

func TestSuccessEmbed(t *testing.T) {
	message := "Operation completed successfully"
	embed := SuccessEmbed(message)

	require.NotNil(t, embed)
	assert.Equal(t, "✅ Success", embed.Title)
	assert.Equal(t, ColorSuccess, embed.Color)
	assert.Equal(t, message, embed.Description)
}

func TestWelcomeGreeting(t *testing.T) {
	guildName := "Test Guild"
	embed := WelcomeGreeting(guildName)

	require.NotNil(t, embed)
	assert.Contains(t, embed.Title, guildName)
	assert.Equal(t, ColorInfo, embed.Color)
	assert.Contains(t, embed.Description, "RuneScape account")
	assert.NotNil(t, embed.Footer)
}

func TestCompetitionCodeEmbed(t *testing.T) {
	eventName := "Boss of the Week - Nex"
	verificationCode := "test-code-123"
	competitionID := int64(54321)

	embed := CompetitionCodeEmbed(eventName, verificationCode, competitionID)

	require.NotNil(t, embed)
	assert.Equal(t, "🔑 Competition Verification Code", embed.Title)
	assert.Contains(t, embed.Description, eventName)
	assert.Contains(t, embed.Fields[0].Value, verificationCode)
	assert.Contains(t, embed.Fields[1].Value, "https://wiseoldman.net/competitions/54321")
}

func TestGetBossImageURL(t *testing.T) {
	tests := []struct {
		name          string
		activity      string
		expectDefault bool
	}{
		{"known boss", "corporeal_beast", false},
		{"known boss 2", "nex", false},
		{"unknown boss", "unknown_boss", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := getBossImageURL(tt.activity)
			assert.NotEmpty(t, url)

			if tt.expectDefault {
				assert.Contains(t, url, "OSRS_icon.png")
			} else {
				assert.NotContains(t, url, "OSRS_icon.png")
			}
		})
	}
}
