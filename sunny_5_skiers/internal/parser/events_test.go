package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/model"
)

func mustCfg(t *testing.T) *config.Config {
	t.Helper()
	js := `{"laps":1,"lapLen":100,"penaltyLen":10,"firingLines":1,"start":"09:00:00","startDelta":"00:00:30"}`
	var c config.Config
	if err := c.Decode(strings.NewReader(js)); err != nil {
		t.Fatalf("config decode failed: %v", err)
	}
	return &c
}

func TestParseEventLine_Table(t *testing.T) {
	cfg := mustCfg(t)

	tests := []struct {
		name      string
		line      string
		wantErr   bool
		wantID    int
		wantSkier int64
		wantParam string
	}{
		{
			name:      "valid registration",
			line:      "[09:00:10.000] 1 42",
			wantErr:   false,
			wantID:    model.EventRegistration,
			wantSkier: 42,
			wantParam: "",
		},
		{
			name:      "valid with extra param",
			line:      "[09:01:02.500] 4 7 extra",
			wantErr:   false,
			wantID:    model.EventOnLap,
			wantSkier: 7,
			wantParam: "extra",
		},
		{
			name:    "too few fields",
			line:    "[09:00:00.000] 1",
			wantErr: true,
		},
		{
			name:    "bad timestamp",
			line:    "[xx:yy:zz.www] 2 3",
			wantErr: true,
		},
		{
			name:    "non-integer id",
			line:    "[09:00:00.000] foo 3",
			wantErr: true,
		},
		{
			name:    "non-integer skierID",
			line:    "[09:00:00.000] 1 bar",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, err := parseEvent(tc.line, cfg)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantID, e.ID)
				assert.Equal(t, tc.wantSkier, e.SkierID)
				assert.Equal(t, tc.wantParam, e.Param)
			}
		})
	}
}
