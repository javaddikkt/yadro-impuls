package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/model"
	"sunny_5_skiers/internal/race"
)

func ParseEvents(path string, cfg *config.Config) error {
	f, _ := os.Open(path)
	defer f.Close()

	r := race.NewRace(cfg)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		e, err := parseEvent(sc.Text(), cfg)
		if err != nil {
			return err
		}
		if err := r.HandleEvent(e); err != nil {
			return err
		}
	}
	if sc.Err() != nil {
		return sc.Err()
	}

	r.Results()
	return nil
}

func parseEvent(line string, cfg *config.Config) (*model.Event, error) {
	p := strings.Fields(line)
	ts, _ := time.Parse("15:04:05.000", strings.Trim(p[0], "[]"))
	id, _ := strconv.Atoi(p[1])
	sk, _ := strconv.ParseInt(p[2], 10, 64)
	param := strings.Join(p[3:], " ")
	return &model.Event{Time: ts.Sub(cfg.Start), ID: id, SkierID: sk, Param: param}, nil
}
