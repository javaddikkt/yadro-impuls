package race

import (
	"bufio"
	"os"
	"strings"
	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/model"
	"testing"
	"time"
)

func FuzzHandleEvent(f *testing.F) {
	const js = `{
        "laps": 1,
        "lapLen": 100,
        "penaltyLen": 10,
        "firingLines": 1,
        "start": "09:00:00",
        "startDelta": "00:00:10"
    }`
	var cfg config.Config
	if err := cfg.Decode(strings.NewReader(js)); err != nil {
		f.Fatalf("config decode failed: %v", err)
	}

	out, _ := os.OpenFile("race/temp/output", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	w := bufio.NewWriter(out)
	defer out.Close()
	defer os.Remove("race/temp/output")
	defer w.Flush()

	f.Add(model.EventRegistration, int64(1), "", int64(0))
	f.Add(model.EventTimeDraw, int64(1), "09:00:05", int64(5*time.Second))
	f.Add(model.EventOnLine, int64(1), "", int64(0))
	f.Add(model.EventOnLap, int64(1), "", int64(12*time.Second))
	f.Fuzz(func(t *testing.T, id int, skierID int64, param string, timeNs int64) {
		r := NewRace(&cfg, w)

		e := &model.Event{
			ID:      id,
			SkierID: skierID,
			Param:   param,
			Time:    time.Duration(timeNs),
		}

		_ = r.HandleEvent(e)
	})
}
