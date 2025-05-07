package race

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/model"
)

func makeRace(t *testing.T) *Race {
	t.Helper()
	js := `{"laps":1,"lapLen":100,"penaltyLen":10,"firingLines":1,"start":"09:00:00","startDelta":"00:00:10"}`
	var cfg config.Config
	assert.NoError(t, cfg.Decode(strings.NewReader(js)))
	buf := &bytes.Buffer{}
	return NewRace(&cfg, bufio.NewWriter(buf))
}

func TestHandleRegister_Duplicate(t *testing.T) {
	r := makeRace(t)
	e := &model.Event{ID: model.EventRegistration, SkierID: 7}
	assert.NoError(t, r.handleRegister(nil, e))
	err := r.handleRegister(nil, e)
	assert.Error(t, err)
}

func TestHandleDraw_WrongState(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateTimeDraw}
	e := &model.Event{ID: model.EventTimeDraw, SkierID: sk.ID, Param: "09:00:05"}
	err := r.handleDraw(sk, e)
	assert.Error(t, err)
}

func TestHandleOnLine_NoDraw(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateRegistered}
	e := &model.Event{ID: model.EventOnLine, SkierID: sk.ID}
	err := r.handleOnStart(sk, e)
	assert.Error(t, err)
}

func TestHandleOnLap_Disqualify(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOnLine, PlannedStart: 5 * time.Second}
	e := &model.Event{ID: model.EventOnLap, SkierID: sk.ID, Time: 20 * time.Second}
	assert.NoError(t, r.handleOnLap(sk, e))
	assert.Equal(t, model.StateDisqualified, sk.State)
}

func TestHandleOnRange_Duplicate(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOnLap, CurrentRanges: map[int]bool{1: true}}
	e := &model.Event{ID: model.EventOnRange, SkierID: sk.ID, Param: "1"}
	err := r.handleOnRange(sk, e)
	assert.Error(t, err)
}

func TestHandleHit_InvalidTarget(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOnRange}
	e := &model.Event{ID: model.EventHit, SkierID: sk.ID, Param: "9"}
	err := r.handleHit(sk, e)
	assert.Error(t, err)
}

func TestHandleOffRange_WrongState(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOnLap}
	e := &model.Event{ID: model.EventOffRange, SkierID: sk.ID}
	err := r.handleOffRange(sk, e)
	assert.Error(t, err)
}

func TestHandleOnPenalty_WrongState(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOnRange}
	e := &model.Event{ID: model.EventOnPenalty, SkierID: sk.ID}
	err := r.handleOnPenalty(sk, e)
	assert.Error(t, err)
}

func TestHandleOffPenalty_WrongState(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOffRange}
	e := &model.Event{ID: model.EventOffPenalty, SkierID: sk.ID}
	err := r.handleOffPenalty(sk, e)
	assert.Error(t, err)
}

func TestHandleOffLap_WrongState(t *testing.T) {
	r := makeRace(t)
	sk := &model.Skier{State: model.StateOnLap}
	e := &model.Event{ID: model.EventOffLap, SkierID: sk.ID, Time: 1 * time.Second}
	err := r.handleOffLap(sk, e)
	assert.Error(t, err)
}
