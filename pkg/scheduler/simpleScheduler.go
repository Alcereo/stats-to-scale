package scheduler

import (
	"github.com/robfig/cron/v3"
	"time"
)

func StartRepeatedFunction(
	duration time.Duration,
	task func(),
) chan bool {
	cron := cron.New(cron.WithSeconds())
	cron.AddFunc("@every 3s", task)
	//tick := time.NewTicker(duration)
	//done := make(chan bool)
	//go scheduler(tick, done, task)
	//return done
	cron.Start()
}

func scheduler(
	tick *time.Ticker,
	done chan bool,
	task func(time time.Time),
) {
	task(time.Now())
	for {
		select {
		case t := <-tick.C:
			task(t)
		case <-done:
			return
		}
	}
}
