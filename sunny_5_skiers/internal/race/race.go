package race

import (
	"bufio"
	"fmt"
	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/model"
)

type Race struct {
	skiers    map[int64]*model.Skier
	handlers  map[int]func(*model.Skier, *model.Event) error
	Config    *config.Config
	LogWriter *bufio.Writer
}

func NewRace(cfg *config.Config, writer *bufio.Writer) *Race {
	r := &Race{
		skiers:    make(map[int64]*model.Skier),
		Config:    cfg,
		LogWriter: writer,
	}
	r.initHandlers()
	return r
}

func (r *Race) HandleEvent(e *model.Event) error {
	if e.ID == model.EventRegistration {
		return r.handlers[e.ID](nil, e)
	}
	sk, ok := r.skiers[e.SkierID]
	if !ok {
		return fmt.Errorf("skier %d doesn't exist", e.SkierID)
	}
	if sk.State == model.StateFail || sk.State == model.StateDisqualified {
		return nil
	}
	return r.handlers[e.ID](sk, e)
}

func (r *Race) Results() { printResults(r) }
