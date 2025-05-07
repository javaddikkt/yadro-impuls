package race

import (
	"fmt"
	"sort"
	"time"

	"sunny_5_skiers/internal/model"
)

func printResults(r *Race) {
	var list []*model.Skier
	for _, sk := range r.skiers {
		list = append(list, sk)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].TotalMain < list[j].TotalMain
	})

	for _, sk := range list {
		switch sk.State {
		case model.StateFail:
			fmt.Print("[NotFinished] ")
		case model.StateDisqualified:
			fmt.Print("[NotStarted] ")
		default:
			fmt.Printf("[%s] ", formatDuration(sk.TotalMain))
		}

		fmt.Print(sk.ID, " [")
		for i := 0; i < r.Config.Laps; i++ {
			dur := sk.LapTimes[i]
			if dur == 0 {
				fmt.Print("{,}, ")
				continue
			}
			speed := float64(r.Config.LapLen) / (float64(dur) / float64(time.Second))
			fmt.Printf("{%s, %.3f}, ", formatDuration(dur), speed)
		}
		fmt.Print("] ")

		if sk.PenaltyTime == 0 {
			fmt.Print("{,}, ")
		} else {
			penSp := float64(r.Config.PenaltyLen*sk.PenaltyLaps) / (float64(sk.PenaltyTime) / float64(time.Second))
			fmt.Printf("{%s, %.3f} ", formatDuration(sk.PenaltyTime), penSp)
		}
		fmt.Printf("%d/%d\n", r.Config.FiringLines*5-sk.PenaltyLaps, r.Config.FiringLines*5)
	}
}

func printLog(r *Race, log string, e *model.Event) error {
	_, err := fmt.Fprintf(r.LogWriter, "[%s] %s\n", r.Config.Start.Add(e.Time).Format("15:04:05.000"), log)
	if err != nil {
		return fmt.Errorf("error writing to log: %v", err)
	}
	return nil
}

func formatDuration(d time.Duration) string {
	return time.Unix(0, d.Nanoseconds()).UTC().Format("15:04:05.000")
}
