package model

import "time"

type SkierState int

const (
	StateRegistered SkierState = iota
	StateTimeDraw
	StateOnLine
	StateOnLap
	StateOnRange
	StateHit
	StateOffRange
	StateOnPenalty
	StateOffPenalty
	StateOffLap
	StateFail
	StateDisqualified
)

type Skier struct {
	ID int64

	State SkierState

	PlannedStart time.Duration
	ActualStart  time.Duration

	LapTimes  []time.Duration
	TotalMain time.Duration

	CurrentHits   map[int]bool
	CurrentRanges map[int]bool
	Hits          int

	PenaltyStart time.Duration
	PenaltyTime  time.Duration
	PenaltyLaps  int

	LapsCompleted int
}
