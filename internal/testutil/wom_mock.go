package testutil

import (
	"context"

	"github.com/kaffeed/voidling/internal/wiseoldman"
	"github.com/stretchr/testify/mock"
)

// MockWOMClient is a mock implementation of the Wise Old Man API client
type MockWOMClient struct {
	mock.Mock
}

// GetPlayer mocks fetching a player from WOM API
func (m *MockWOMClient) GetPlayer(ctx context.Context, username string) (*wiseoldman.Player, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wiseoldman.Player), args.Error(1)
}

// UpdatePlayer mocks updating a player in WOM
func (m *MockWOMClient) UpdatePlayer(ctx context.Context, username string) (*wiseoldman.Player, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wiseoldman.Player), args.Error(1)
}

// CreateCompetition mocks creating a WOM competition
func (m *MockWOMClient) CreateCompetition(ctx context.Context, req wiseoldman.CreateCompetitionRequest) (*wiseoldman.CreateCompetitionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*wiseoldman.CreateCompetitionResponse), args.Error(1)
}

// AddParticipantsToCompetition mocks adding participants to a competition
func (m *MockWOMClient) AddParticipantsToCompetition(ctx context.Context, competitionID int64, usernames []string, verificationCode string) error {
	args := m.Called(ctx, competitionID, usernames, verificationCode)
	return args.Error(0)
}

// CreateTestPlayer creates a test WOM player with realistic data
func CreateTestPlayer(username string) *wiseoldman.Player {
	return &wiseoldman.Player{
		ID:          12345,
		Username:    username,
		DisplayName: username,
		Type:        "regular",
		Build:       "main",
		CombatLevel: 126,
		EHP:         500.5,
		EHB:         250.3,
		LatestSnapshot: &wiseoldman.Snapshot{
			Data: wiseoldman.SnapshotData{
				Skills: map[string]wiseoldman.SkillData{
					"overall": {
						Metric:     "overall",
						Rank:       10000,
						Level:      2277,
						Experience: 4600000000,
						EHP:        500.5,
					},
					"attack": {
						Metric:     "attack",
						Rank:       5000,
						Level:      99,
						Experience: 200000000,
						EHP:        50.0,
					},
					"strength": {
						Metric:     "strength",
						Rank:       5001,
						Level:      99,
						Experience: 200000000,
						EHP:        50.0,
					},
					"defence": {
						Metric:     "defence",
						Rank:       5002,
						Level:      99,
						Experience: 200000000,
						EHP:        50.0,
					},
					"hitpoints": {
						Metric:     "hitpoints",
						Rank:       5003,
						Level:      99,
						Experience: 200000000,
						EHP:        50.0,
					},
					"ranged": {
						Metric:     "ranged",
						Rank:       5004,
						Level:      99,
						Experience: 200000000,
						EHP:        50.0,
					},
					"magic": {
						Metric:     "magic",
						Rank:       5005,
						Level:      99,
						Experience: 200000000,
						EHP:        50.0,
					},
					"prayer": {
						Metric:     "prayer",
						Rank:       5006,
						Level:      99,
						Experience: 13034431,
						EHP:        25.0,
					},
					"woodcutting": {
						Metric:     "woodcutting",
						Rank:       10000,
						Level:      99,
						Experience: 15000000,
						EHP:        30.0,
					},
				},
				Bosses: map[string]wiseoldman.BossData{
					"chambers_of_xeric": {
						Metric: "chambers_of_xeric",
						Rank:   1000,
						Kills:  500,
						EHB:    25.5,
					},
					"theatre_of_blood": {
						Metric: "theatre_of_blood",
						Rank:   2000,
						Kills:  300,
						EHB:    30.2,
					},
					"corporeal_beast": {
						Metric: "corporeal_beast",
						Rank:   1500,
						Kills:  250,
						EHB:    20.1,
					},
					"nex": {
						Metric: "nex",
						Rank:   3000,
						Kills:  100,
						EHB:    15.3,
					},
				},
			},
		},
	}
}

// CreateTestCompetitionResponse creates a test competition creation response
func CreateTestCompetitionResponse(title string, metric string, competitionID int64) *wiseoldman.CreateCompetitionResponse {
	return &wiseoldman.CreateCompetitionResponse{
		Competition: wiseoldman.Competition{
			ID:     competitionID,
			Title:  title,
			Metric: metric,
		},
		VerificationCode: "test-verification-code-123",
	}
}
