package wiseoldman

import "time"

// Player represents a Wise Old Man player.
type Player struct {
	ID             int64      `json:"id"`
	Username       string     `json:"username"`
	DisplayName    string     `json:"displayName"`
	Type           string     `json:"type"`
	Build          string     `json:"build"`
	Country        *string    `json:"country"`
	Status         string     `json:"status"`
	Patron         bool       `json:"patron"`
	Exp            int64      `json:"exp"`
	EHP            float64    `json:"ehp"`
	EHB            float64    `json:"ehb"`
	TTM            float64    `json:"ttm"`
	TT200M         float64    `json:"tt200m"`
	RegisteredAt   time.Time  `json:"registeredAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	LastChangedAt  *time.Time `json:"lastChangedAt"`
	LastImportedAt *time.Time `json:"lastImportedAt"`
	CombatLevel    int        `json:"combatLevel"`
	Archive        *Archive   `json:"archive"`
	LatestSnapshot *Snapshot  `json:"latestSnapshot"`
}

// Archive represents archived player data.
type Archive struct {
	PreviousUsername string    `json:"previousUsername"`
	ArchivedAt       time.Time `json:"archivedAt"`
	RestoredAt       time.Time `json:"restoredAt"`
}

// Snapshot represents a player snapshot.
type Snapshot struct {
	ID         int64        `json:"id"`
	PlayerID   int64        `json:"playerId"`
	CreatedAt  time.Time    `json:"createdAt"`
	ImportedAt *time.Time   `json:"importedAt"`
	Data       SnapshotData `json:"data"`
}

// SnapshotData contains all player metrics.
type SnapshotData struct {
	Skills     map[string]SkillData    `json:"skills"`
	Bosses     map[string]BossData     `json:"bosses"`
	Activities map[string]ActivityData `json:"activities"`
	Computed   map[string]ComputedData `json:"computed"`
}

// SkillData represents skill metrics.
type SkillData struct {
	Metric     string  `json:"metric"`
	Experience int64   `json:"experience"`
	Rank       int     `json:"rank"`
	Level      int     `json:"level"`
	EHP        float64 `json:"ehp"`
}

// BossData represents boss kill count metrics.
type BossData struct {
	Metric string  `json:"metric"`
	Kills  int     `json:"kills"`
	Rank   int     `json:"rank"`
	EHB    float64 `json:"ehb"`
}

// ActivityData represents activity metrics.
type ActivityData struct {
	Metric string `json:"metric"`
	Score  int    `json:"score"`
	Rank   int    `json:"rank"`
}

// ComputedData represents computed metrics (EHP, EHB, etc).
type ComputedData struct {
	Metric string  `json:"metric"`
	Value  float64 `json:"value"`
	Rank   int     `json:"rank"`
}

// GetSkill returns skill data for a given skill name, or nil if not found.
func (p *Player) GetSkill(skillName string) *SkillData {
	if p.LatestSnapshot == nil {
		return nil
	}
	skill, ok := p.LatestSnapshot.Data.Skills[skillName]
	if !ok {
		return nil
	}
	return &skill
}

// GetBoss returns boss data for a given boss name, or nil if not found.
func (p *Player) GetBoss(bossName string) *BossData {
	if p.LatestSnapshot == nil {
		return nil
	}
	boss, ok := p.LatestSnapshot.Data.Bosses[bossName]
	if !ok {
		return nil
	}
	return &boss
}

// Competition represents a WOM competition.
type Competition struct {
	ID               int64                      `json:"id"`
	Title            string                     `json:"title"`
	Metric           string                     `json:"metric"`
	Type             string                     `json:"type"`
	StartsAt         time.Time                  `json:"startsAt"`
	EndsAt           time.Time                  `json:"endsAt"`
	GroupID          *int64                     `json:"groupId"`
	Score            int                        `json:"score"`
	CreatedAt        time.Time                  `json:"createdAt"`
	UpdatedAt        time.Time                  `json:"updatedAt"`
	ParticipantCount int                        `json:"participantCount"`
	Participations   []CompetitionParticipation `json:"participations"`
}

// CompetitionParticipation represents a player's participation in a competition.
type CompetitionParticipation struct {
	PlayerID      int64                  `json:"playerId"`
	CompetitionID int64                  `json:"competitionId"`
	TeamName      *string                `json:"teamName"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	Player        Player                 `json:"player"`
	Progress      *ParticipationProgress `json:"progress,omitempty"`
}

// ParticipationProgress represents progress gained during a competition.
type ParticipationProgress struct {
	Start  int64 `json:"start"`
	End    int64 `json:"end"`
	Gained int64 `json:"gained"`
}

// CreateCompetitionRequest is the request body for creating a competition.
type CreateCompetitionRequest struct {
	Title        string   `json:"title"`
	Metric       string   `json:"metric"`
	StartsAt     string   `json:"startsAt"`
	EndsAt       string   `json:"endsAt"`
	Participants []string `json:"participants,omitempty"`
}

// CreateCompetitionResponse is the response from creating a competition.
type CreateCompetitionResponse struct {
	Competition      Competition `json:"competition"`
	VerificationCode string      `json:"verificationCode"`
}

// AddParticipantsRequest is the request body for adding participants.
type AddParticipantsRequest struct {
	VerificationCode string   `json:"verificationCode"`
	Participants     []string `json:"participants"`
}

// AddParticipantsResponse is the response from adding participants.
type AddParticipantsResponse struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}
