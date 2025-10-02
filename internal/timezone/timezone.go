package timezone

import (
	"fmt"
	"strings"
	"time"
)

// CommonTimezones returns a list of commonly used IANA timezone names
func CommonTimezones() []string {
	return []string{
		// UTC
		"UTC",
		// Americas
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/Toronto",
		"America/Vancouver",
		"America/Phoenix",
		"America/Anchorage",
		"America/Sao_Paulo",
		"America/Argentina/Buenos_Aires",
		"America/Mexico_City",
		// Europe
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Europe/Amsterdam",
		"Europe/Madrid",
		"Europe/Rome",
		"Europe/Stockholm",
		"Europe/Moscow",
		"Europe/Athens",
		"Europe/Istanbul",
		// Asia
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Hong_Kong",
		"Asia/Singapore",
		"Asia/Seoul",
		"Asia/Dubai",
		"Asia/Kolkata",
		"Asia/Bangkok",
		"Asia/Jakarta",
		"Asia/Manila",
		// Australia/Pacific
		"Australia/Sydney",
		"Australia/Melbourne",
		"Australia/Brisbane",
		"Australia/Perth",
		"Pacific/Auckland",
		"Pacific/Fiji",
		"Pacific/Honolulu",
	}
}

// ValidateTimezone checks if a timezone string is valid
func ValidateTimezone(tz string) error {
	if tz == "" {
		return fmt.Errorf("timezone cannot be empty")
	}

	_, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone '%s': %w", tz, err)
	}

	return nil
}

// ParseInTimezone parses a time string in a specific timezone
// timeStr format: "2006-01-02 15:04"
// Returns the time in UTC
func ParseInTimezone(timeStr, tz string) (time.Time, error) {
	// Validate timezone
	if err := ValidateTimezone(tz); err != nil {
		return time.Time{}, err
	}

	// Load the timezone location
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone: %w", err)
	}

	// Parse the time in the specified timezone
	t, err := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return t, nil
}

// ConvertToTimezone converts a UTC time to a specific timezone
func ConvertToTimezone(t time.Time, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone: %w", err)
	}

	return t.In(loc), nil
}

// FormatForDiscord returns a Discord timestamp string
// Format: <t:UNIX:F> for full date/time
func FormatForDiscord(t time.Time) string {
	return fmt.Sprintf("<t:%d:F>", t.Unix())
}

// FormatForDiscordRelative returns a Discord relative timestamp string
// Format: <t:UNIX:R> for "in X hours/days"
func FormatForDiscordRelative(t time.Time) string {
	return fmt.Sprintf("<t:%d:R>", t.Unix())
}

// GetTimezoneAbbreviation returns a friendly abbreviation for display
// e.g., "America/New_York" -> "EST" or "EDT" depending on date
func GetTimezoneAbbreviation(t time.Time, tz string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return tz
	}

	return t.In(loc).Format("MST")
}

// FormatTimeWithTimezone returns a formatted string with timezone info
// e.g., "January 15, 2025 8:00 PM EST"
func FormatTimeWithTimezone(t time.Time, tz string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return t.Format("January 2, 2006 3:04 PM")
	}

	localTime := t.In(loc)
	abbrev := localTime.Format("MST")
	formatted := localTime.Format("January 2, 2006 3:04 PM")

	return fmt.Sprintf("%s %s", formatted, abbrev)
}

// SearchTimezones returns timezones that match the search query
// Used for autocomplete functionality
func SearchTimezones(query string) []string {
	query = strings.ToLower(query)
	results := []string{}

	for _, tz := range CommonTimezones() {
		if strings.Contains(strings.ToLower(tz), query) {
			results = append(results, tz)
			if len(results) >= 25 { // Discord autocomplete limit
				break
			}
		}
	}

	return results
}
