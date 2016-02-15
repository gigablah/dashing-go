package jobs

import (
	"math/rand"
	"time"

	"gopkg.in/gigablah/dashing-go.v1"
)

type sample struct{}

func (j *sample) Work(send chan *dashing.Event) {
	ticker := time.NewTicker(1 * time.Second)
	var lastValuation, lastKarma, currentValuation, currentKarma int
	for {
		select {
		case <-ticker.C:
			lastValuation, currentValuation = currentValuation, rand.Intn(100)
			lastKarma, currentKarma = currentKarma, rand.Intn(200000)
			send <- &dashing.Event{"valuation", map[string]interface{}{
				"current": currentValuation,
				"last":    lastValuation,
			}, ""}
			send <- &dashing.Event{"karma", map[string]interface{}{
				"current": currentKarma,
				"last":    lastKarma,
			}, ""}
			send <- &dashing.Event{"synergy", map[string]interface{}{
				"value": rand.Intn(100),
			}, ""}
		}
	}
}

func init() {
	dashing.Register(&sample{})
}
