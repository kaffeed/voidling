package models

// EventType represents the type of a trackable or schedulable event
type EventType string

const (
	EventTypeBossOfTheWeek  EventType = "BOSS_OF_THE_WEEK"
	EventTypeSkillOfTheWeek EventType = "SKILL_OF_THE_WEEK"
	EventTypeMass           EventType = "MASS"
	EventTypeWildyWednesday EventType = "WILDY_WEDNESDAY"
)

// HiscoreField represents a skill or boss in OSRS
// This is used for event descriptions and tracking
type HiscoreField string
