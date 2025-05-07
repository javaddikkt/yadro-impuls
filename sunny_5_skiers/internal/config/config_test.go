package config

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecode_Valid(t *testing.T) {
	js := `{
        "laps":2,
        "lapLen":3651,
        "penaltyLen":50,
        "firingLines":1,
        "start":"09:30:00",
        "startDelta":"00:00:30"
    }`
	var cfg Config
	err := cfg.Decode(strings.NewReader(js))
	assert.NoError(t, err)
	assert.Equal(t, 2, cfg.Laps)
	assert.Equal(t, 3651, cfg.LapLen)
	assert.Equal(t, 50, cfg.PenaltyLen)
	assert.Equal(t, 1, cfg.FiringLines)
	wantStart, _ := time.Parse("15:04:05", "09:30:00")
	assert.True(t, cfg.Start.Equal(wantStart))
	assert.Equal(t, 30*time.Second, cfg.StartDelta)
}

func TestDecode_BadJSON(t *testing.T) {
	var cfg Config
	err := cfg.Decode(strings.NewReader("{bad json"))
	assert.Error(t, err)
}

func TestUnmarshalJSON_BadStart(t *testing.T) {
	bad := `{"laps":1,"lapLen":1,"penaltyLen":1,"firingLines":1,"start":"xx","startDelta":"00:00:10"}`
	var cfg Config
	err := json.Unmarshal([]byte(bad), &cfg)
	assert.Error(t, err, "invalid start must error")
}

func TestUnmarshalJSON_BadDelta(t *testing.T) {
	bad := `{"laps":1,"lapLen":1,"penaltyLen":1,"firingLines":1,"start":"00:00:00","startDelta":"zz"}`
	var cfg Config
	err := json.Unmarshal([]byte(bad), &cfg)
	assert.Error(t, err, "invalid delta must error")
}
