package jobs

import (
	"math/rand"
	"time"

	"gopkg.in/gigablah/dashing-go.v1"
)

type buzzwords struct {
	words []map[string]interface{}
}

func (j *buzzwords) Work(send chan *dashing.Event) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < len(j.words); i++ {
				if 1 < rand.Intn(3) {
					value := j.words[i]["value"].(int)
					j.words[i]["value"] = (value + 1) % 30
				}
			}
			send <- &dashing.Event{"buzzwords", map[string]interface{}{
				"items": j.words,
			}, ""}
		}
	}
}

func init() {
	dashing.Register(&buzzwords{[]map[string]interface{}{
		{"label": "Paradigm shift", "value": 0},
		{"label": "Leverage", "value": 0},
		{"label": "Pivoting", "value": 0},
		{"label": "Turn-key", "value": 0},
		{"label": "Streamlininess", "value": 0},
		{"label": "Exit strategy", "value": 0},
		{"label": "Synergy", "value": 0},
		{"label": "Enterprise", "value": 0},
		{"label": "Web 2.0", "value": 0},
	}})
}
