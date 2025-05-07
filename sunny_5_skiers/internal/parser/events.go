package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sunny_5_skiers/internal/config"
	"sunny_5_skiers/internal/model"
	"sunny_5_skiers/internal/race"
)

func ParseEvents(eventsPath string, outputPath string, cfg *config.Config) error {
	events, err := os.Open(eventsPath)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}
	defer events.Close()

	out, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}
	w := bufio.NewWriter(out)
	defer out.Close()
	defer w.Flush()

	r := race.NewRace(cfg, w)
	sc := bufio.NewScanner(events)
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
		return fmt.Errorf("error scanning events: %v", sc.Err())
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
