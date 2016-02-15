package dashing

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"gopkg.in/husobee/vestigo.v1"
	"gopkg.in/karlseguin/gerb.v0"
)

// A Server contains webservice parameters and handlers.
type Server struct {
	dev    bool
	port   int
	broker *Broker
	Router http.Handler
}

// IndexHandler redirects to the default dashboard.
func (h *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("dashboards/*.gerb")

	for _, file := range files {
		dashboard := file[11 : len(file)-5]
		if dashboard != "layout" {
			http.Redirect(w, r, "/"+dashboard, http.StatusTemporaryRedirect)
			return
		}
	}

	http.NotFound(w, r)
}

// EventsHandler opens a keepalive connection and pushes events to the client.
func (h *Server) EventsHandler(w http.ResponseWriter, r *http.Request) {
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
	// send this client events.
	events := make(chan *Event)

	// Add this client to the map of those that should
	// receive updates
	h.broker.newClients <- events

	// Remove this client from the map of attached clients
	// when the handler exits.
	defer func() {
		h.broker.defunctClients <- events
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	closer := c.CloseNotify()

	for {
		select {
		case event := <-events:
			data := event.Body
			data["id"] = event.ID
			data["updatedAt"] = int32(time.Now().Unix())
			json, err := json.Marshal(data)
			if err != nil {
				continue
			}
			if event.Target != "" {
				fmt.Fprintf(w, "event: %s\n", event.Target)
			}
			fmt.Fprintf(w, "data: %s\n\n", json)
			f.Flush()
		case <-closer:
			log.Println("Closing connection")
			return
		}
	}
}

// DashboardHandler serves the dashboard layout template.
func (h *Server) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	template, err := gerb.ParseFile(true, "dashboards/"+vestigo.Param(r, "dashboard")+".gerb", "dashboards/layout.gerb")

	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	template.Render(w, map[string]interface{}{
		"dashboard":   vestigo.Param(r, "dashboard"),
		"development": h.dev,
		"request":     r,
	})
}

// DashboardEventHandler accepts dashboard events.
func (h *Server) DashboardEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	var data map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	h.broker.events <- &Event{vestigo.Param(r, "id"), data, "dashboards"}

	w.WriteHeader(http.StatusNoContent)
}

// WidgetHandler serves widget templates.
func (h *Server) WidgetHandler(w http.ResponseWriter, r *http.Request) {
	template, err := gerb.ParseFile(true, "widgets/"+vestigo.Param(r, "widget")+"/"+vestigo.Param(r, "widget")+".html")

	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	template.Render(w, nil)
}

// WidgetEventHandler accepts widget data.
func (h *Server) WidgetEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	var data map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	h.broker.events <- &Event{vestigo.Param(r, "id"), data, ""}

	w.WriteHeader(http.StatusNoContent)
}

// Start listening to requests.
func (h *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", h.port), h.Router))
}

// NewServer creates a Server instance.
func NewServer(p int, b *Broker, h http.Handler) *Server {
	return &Server{
		dev:    false,
		port:   p,
		broker: b,
		Router: h,
	}
}
