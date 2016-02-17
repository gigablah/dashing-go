package dashing

import (
	"net/http"
	"os"
	"path/filepath"
)

// An Event contains the widget ID, a body of data,
// and an optional target (only "dashboard" for now).
type Event struct {
	ID     string
	Body   map[string]interface{}
	Target string
}

// Dashing struct definition.
type Dashing struct {
	started bool
	Broker  *Broker
	Worker  *Worker
	Server  *Server
	Router  http.Handler
}

// ServeHTTP implements the HTTP Handler.
func (d *Dashing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !d.started {
		panic("dashing.Start() has not been called")
	}
	d.Router.ServeHTTP(w, r)
}

// Start actives the broker and workers.
func (d *Dashing) Start() *Dashing {
	if !d.started {
		if d.Router == nil {
			d.Router = d.Server.NewRouter()
		}
		d.Broker.Start()
		d.Worker.Start()
		d.started = true
	}
	return d
}

// NewDashing sets up the event broker, workers and webservice.
func NewDashing() *Dashing {
	broker := NewBroker()
	worker := NewWorker(broker)
	server := NewServer(broker)

	if os.Getenv("WEBROOT") != "" {
		server.webroot = filepath.Clean(os.Getenv("WEBROOT")) + "/"
	}
	if os.Getenv("DEV") != "" {
		server.dev = true
	}

	return &Dashing{
		started: false,
		Broker:  broker,
		Worker:  worker,
		Server:  server,
	}
}
