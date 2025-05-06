package config

import (
	"encoding/json"
	"io"
	"time"
)

type Config struct {
	Laps        int
	LapLen      int
	PenaltyLen  int
	FiringLines int
	Start       time.Time
	StartDelta  time.Duration
}

func (c *Config) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(c)
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var aux struct {
		Laps        int    `json:"laps"`
		LapLen      int    `json:"lapLen"`
		PenaltyLen  int    `json:"penaltyLen"`
		FiringLines int    `json:"firingLines"`
		Start       string `json:"start"`
		StartDelta  string `json:"startDelta"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.Laps, c.LapLen, c.PenaltyLen, c.FiringLines =
		aux.Laps, aux.LapLen, aux.PenaltyLen, aux.FiringLines

	t, err := time.Parse("15:04:05", aux.Start)
	if err != nil {
		return err
	}
	c.Start = t

	zero, _ := time.Parse("15:04:05", "00:00:00")
	delta, _ := time.Parse("15:04:05", aux.StartDelta)
	c.StartDelta = delta.Sub(zero)
	return nil
}
