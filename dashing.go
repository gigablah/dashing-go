package dashing

import (
    "log"
    "fmt"
    "net/http"
    "github.com/codegangsta/martini"
    "github.com/codegangsta/martini-contrib/encoder"
)

// The Martini instance
var m *martini.Martini

// Message broker
var b *Broker

func init() {
    m = martini.New()

    // Setup middleware
    m.Use(martini.Recovery())
    m.Use(martini.Logger())
    m.Use(martini.Static("public"))

    // Setup encoder
    m.Use(func(c martini.Context, w http.ResponseWriter) {
        c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
    })

    // Setup and inject message broker
    b = NewBroker()
    m.Map(b)

    // Setup routes
    r := martini.NewRouter()

    r.Get("/events", func(w http.ResponseWriter, r *http.Request, e encoder.Encoder, b *Broker) {

        f, ok := w.(http.Flusher)
        if !ok {
            http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
            return
        }

        c, ok := w.(http.CloseNotifier)
        if !ok {
            http.Error(w, "Close notification unsupported!", http.StatusInternalServerError)
            return
        }

        // Create a new channel, over which the broker can
        // send this client messages.
        messageChan := make(chan *Message)

        // Add this client to the map of those that should
        // receive updates
        b.newClients <- messageChan

        // Remove this client from the map of attached clients
        // when the handler exits.
        defer func() {
            b.defunctClients <- messageChan
        }()

        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")
        w.Header().Set("X-Accel-Buffering", "no")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        closer := c.CloseNotify()

        for {
            select {
            case msg := <-messageChan:
                fmt.Fprintf(w, "data: %s\n\n", encoder.Must(e.Encode(msg.Data)))
                f.Flush()
            case <-closer:
                log.Println("Closing connection")
                return
            }
        }

    })

    // Add the router action
    m.Action(r.Handle)
}

func Start() {
    // Start the message broker
    b.Start()

    // Start the jobs
    for _, j := range registry {
        go j.Work(b.messages)
    }

    // Start Martini
    m.Run()
}
