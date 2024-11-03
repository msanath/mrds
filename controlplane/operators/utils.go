package operators

import "time"

func newImmediatelyFiringTicker(d time.Duration) (ticks <-chan time.Time, stop func()) {
	tick := make(chan time.Time)
	ticker := time.NewTicker(d)

	go func() {
		tick <- time.Now()
		for range ticker.C {
			tick <- time.Now()
		}
	}()

	return tick, ticker.Stop
}
