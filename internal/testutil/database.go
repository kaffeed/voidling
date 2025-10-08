// Package testutil provides testing utilities and helpers for the Voidling bot tests.
package testutil

import (
	"database/sql"
	"testing"

	"github.com/kaffeed/voidling/internal/database"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

// SetupTestDB creates an in-memory SQLite database for testing with migrations applied.
func SetupTestDB(t *testing.T) (*sql.DB, *database.Queries) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:?cache=shared")
	require.NoError(t, err, "Failed to create in-memory database")

	// Run migrations
	err = goose.SetDialect("sqlite3")
	require.NoError(t, err, "Failed to set goose dialect")
	err = goose.Up(db, "../../migrations")
	if err != nil {
		// Try alternative path for tests in different directories
		err = goose.Up(db, "./migrations")
		require.NoError(t, err, "Failed to run migrations")
	}

	queries := database.New(db)
	return db, queries
}

// CleanupTestDB closes the test database connection.
func CleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	if db != nil {
		_ = db.Close()
	}
}

// CreateTestAccountLink creates a test account link for testing purposes.
func CreateTestAccountLink(t *testing.T, q *database.Queries, discordID int64, rsn string, active bool) database.AccountLink {
	t.Helper()

	link, err := q.CreateAccountLink(t.Context(), database.CreateAccountLinkParams{
		DiscordMemberID: discordID,
		RunescapeName:   rsn,
		IsActive:        active,
	})
	require.NoError(t, err, "Failed to create test account link")
	return link
}
