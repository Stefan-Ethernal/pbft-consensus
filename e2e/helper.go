package e2e

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func generateNodeNames(from int, count int, prefix string) []string {
	var names []string
	for j := from; j < count; j++ {
		names = append(names, prefix+strconv.Itoa(j))
	}
	return names
}

func executeInTimerAndWait(tickTime time.Duration, duration time.Duration, fn func(time.Duration)) {
	end := executeInTimer(tickTime, duration, fn)
	<-end
}

func executeInTimer(tickTime time.Duration, duration time.Duration, fn func(time.Duration)) chan struct{} {
	tick := time.NewTicker(tickTime)
	tickerDone := make(chan struct{})
	end := make(chan struct{})
	startTime := time.Now()
	go func() {
		for {
			select {
			case v := <-tick.C:
				elapsedTime := v.Sub(startTime)
				fn(elapsedTime)
			case <-tickerDone:
				close(end)
				return
			}
		}
	}()

	after := time.After(duration)
	go func() {
		<-after
		tick.Stop()
		close(tickerDone)
	}()
	return end
}

func isFuzzEnabled(t *testing.T) {
	if os.Getenv("FUZZ") != "true" {
		t.Skip("Fuzz tests are disabled.")
	}
}
