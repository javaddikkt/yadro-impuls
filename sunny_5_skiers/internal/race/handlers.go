package race

import (
	"fmt"
	"strconv"
	"time"

	"sunny_5_skiers/internal/model"
)

func (r *Race) initHandlers() {
	r.handlers = map[int]func(*model.Skier, *model.Event) error{
		model.EventRegistration: r.handleRegister,
		model.EventTimeDraw:     r.handleDraw,
		model.EventOnLine:       r.handleOnStart,
		model.EventOnLap:        r.handleOnLap,
		model.EventOnRange:      r.handleOnRange,
		model.EventHit:          r.handleHit,
		model.EventOffRange:     r.handleOffRange,
		model.EventOnPenalty:    r.handleOnPenalty,
		model.EventOffPenalty:   r.handleOffPenalty,
		model.EventOffLap:       r.handleOffLap,
		model.EventFail:         r.handleFail,
	}
}

func (r *Race) handleEvent(e *model.Event) error {
	if e.ID == model.EventRegistration {
		return r.handlers[e.ID](nil, e)
	}
	skier, ok := r.skiers[e.SkierID]
	if !ok {
		return fmt.Errorf("skier %d not registered", e.SkierID)
	}
	h := r.handlers[e.ID]
	if h == nil {
		return fmt.Errorf("unknown event ID %d", e.ID)
	}
	return h(skier, e)
}

func (r *Race) handleRegister(_ *model.Skier, e *model.Event) error {
	if _, dup := r.skiers[e.SkierID]; dup {
		return fmt.Errorf("duplicate registration for %d", e.SkierID)
	}
	r.skiers[e.SkierID] = &model.Skier{
		ID:            e.SkierID,
		State:         model.StateRegistered,
		LapTimes:      make([]time.Duration, r.Config.Laps),
		CurrentHits:   make(map[int]bool, 5),
		CurrentRanges: make(map[int]bool, r.Config.FiringLines),
	}
	return nil
}

func (r *Race) handleDraw(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateRegistered {
		return fmt.Errorf("skier %d wrong state for handleDraw", sk.ID)
	}
	t, err := time.Parse("15:04:05", e.Param)
	if err != nil {
		return fmt.Errorf("handleDraw parse Time: %w", err)
	}
	sk.PlannedStart = t.Sub(r.Config.Start)
	sk.State = model.StateTimeDraw
	return nil
}

func (r *Race) handleOnStart(sk *model.Skier, _ *model.Event) error {
	if sk.State != model.StateTimeDraw {
		return fmt.Errorf("skier %d not drawn", sk.ID)
	}
	sk.State = model.StateOnLine
	return nil
}

func (r *Race) handleOnLap(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateOnLine {
		return fmt.Errorf("skier %d not at line", sk.ID)
	}
	if e.Time > sk.PlannedStart+r.Config.StartDelta {
		sk.State = model.StateDisqualified
		return nil
	}
	sk.State = model.StateOnLap
	sk.ActualStart = e.Time
	sk.CurrentRanges = make(map[int]bool, r.Config.FiringLines)
	return nil
}

func (r *Race) handleOnRange(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateOnLap {
		return fmt.Errorf("skier %d not on lap", sk.ID)
	}
	n, err := strconv.Atoi(e.Param)
	if err != nil || n < 1 || n > r.Config.FiringLines {
		return fmt.Errorf("invalid range num %q", e.Param)
	}
	if sk.CurrentRanges[n] {
		return fmt.Errorf("skier %d repeated range %d", sk.ID, n)
	}
	sk.CurrentRanges[n] = true
	sk.CurrentHits = make(map[int]bool, 5)
	sk.Hits = 0
	sk.State = model.StateOnRange
	return nil
}

func (r *Race) handleHit(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateOnRange && sk.State != model.StateHit {
		return fmt.Errorf("skier %d not at range", sk.ID)
	}
	target, err := strconv.Atoi(e.Param)
	if err != nil || target < 1 || target > 5 {
		return fmt.Errorf("bad target %q", e.Param)
	}
	if sk.CurrentHits[target] {
		return fmt.Errorf("skier %d duplicate target %d", sk.ID, target)
	}
	sk.CurrentHits[target] = true
	sk.Hits++
	if sk.Hits > 5 {
		return fmt.Errorf("skier %d too many shots", sk.ID)
	}
	sk.State = model.StateHit
	return nil
}

func (r *Race) handleOffRange(sk *model.Skier, _ *model.Event) error {
	if sk.State != model.StateHit && sk.State != model.StateOnRange {
		return fmt.Errorf("skier %d not finished shooting", sk.ID)
	}
	misses := 5 - sk.Hits
	sk.PenaltyLaps += misses
	sk.State = model.StateOffRange
	return nil
}

func (r *Race) handleOnPenalty(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateOffRange {
		return fmt.Errorf("skier %d cannot handleOnLap penalty", sk.ID)
	}
	sk.State = model.StateOnPenalty
	sk.PenaltyStart = e.Time
	return nil
}

func (r *Race) handleOffPenalty(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateOnPenalty {
		return fmt.Errorf("skier %d not on penalty", sk.ID)
	}
	sk.State = model.StateOffPenalty
	sk.PenaltyTime += e.Time - sk.PenaltyStart
	return nil
}

func (r *Race) handleOffLap(sk *model.Skier, e *model.Event) error {
	if sk.State != model.StateOffPenalty && !(sk.State == model.StateOffRange && sk.Hits == 5) {
		return fmt.Errorf("skier %d wrong state to end lap", sk.ID)
	}
	dur := e.Time - sk.ActualStart
	sk.LapTimes[sk.LapsCompleted] = dur
	sk.TotalMain += dur
	sk.LapsCompleted++

	if sk.LapsCompleted == r.Config.Laps {
		sk.State = model.StateOffLap
		return nil
	}

	sk.State = model.StateOnLap
	sk.ActualStart = e.Time
	sk.CurrentRanges = make(map[int]bool, r.Config.FiringLines)
	sk.CurrentHits = make(map[int]bool, 5)
	sk.Hits = 0
	return nil
}

func (r *Race) handleFail(sk *model.Skier, _ *model.Event) error {
	sk.State = model.StateFail
	return nil
}
