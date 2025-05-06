package model

import "time"

const (
	EventRegistration = 1
	EventTimeDraw     = 2
	EventOnLine       = 3
	EventOnLap        = 4
	EventOnRange      = 5
	EventHit          = 6
	EventOffRange     = 7
	EventOnPenalty    = 8
	EventOffPenalty   = 9
	EventOffLap       = 10
	EventFail         = 11
)

type Event struct {
	Time    time.Duration
	ID      int
	SkierID int64
	Param   string
}
