package dashing

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/husobee/vestigo.v1"
)

// An Event contains the widget ID, a body of data,
// and an optional target (only "dashboard" for now).
type Event struct {
	ID     string
	Body   map[string]interface{}
	Target string
}

// A Job does periodic work and sends events to a channel.
type Job interface {
	Work(send chan *Event)
}

var reg []Job

// A Dashing instance contains an event broker, a webservice
// and a collection of registered jobs.
type Dashing struct {
	Broker   *Broker
	Server   *Server
	registry []Job
}

// Start all jobs and listen to requests.
func (d *Dashing) Start() {
	d.Broker.Start()

	// Start the jobs
	for _, j := range d.registry {
		go j.Work(d.Broker.events)
	}

	// Start the webservice
	d.Server.Start()
}

// NewDashing sets up the event broker, router and webservice.
func NewDashing() *Dashing {
	var port int
	var err error
	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatalf("Invalid port: %s", os.Getenv("PORT"))
		}
	} else {
		port = 8080
	}

	broker := NewBroker()
	router := vestigo.NewRouter()
	server := NewServer(port, broker, router)
	if os.Getenv("DEV") != "" {
		server.dev = true
	}

	router.Get("/", server.IndexHandler)
	router.Get("/events", server.EventsHandler)
	router.Get("/:dashboard", server.DashboardHandler)
	router.Post("/dashboards/:id", server.DashboardEventHandler)
	router.Get("/views/:widget.html", server.WidgetHandler)
	router.Post("/widgets/:id", server.WidgetEventHandler)

	return &Dashing{
		Broker:   broker,
		Server:   server,
		registry: reg,
	}
}

// Register a job to be kicked off upon starting the server.
func Register(j Job) {
	if j == nil {
		panic("Can't register nil job")
	}
	reg = append(reg, j)
}
