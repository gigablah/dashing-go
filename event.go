package dashing

import (
    _ "log"
)

type Message struct {
    Data interface{}
}

type Broker struct {
    // Create a map of clients, the keys of the map are the channels
    // over which we can push messages to attached clients. (The values
    // are just booleans and are meaningless)
    clients map[chan *Message]bool

    // Channel into which new clients can be pushed
    newClients chan chan *Message

    // Channel into which disconnected clients should be pushed
    defunctClients chan chan *Message

    // Channel into which messages are pushed to be broadcast out
    // to attached clients
    messages chan *Message
}

func (b *Broker) Start() {
    go func() {
        // Loop endlessly
        for {
            // Block until we receive from one of the
            // three following channels.
            select {
            case s := <-b.newClients:
                // There is a new client attached and we
                // want to start sending them messages.
                b.clients[s] = true
                // log.Println("Added new client")
            case s := <-b.defunctClients:
                // A client has detached and we want to
                // stop sending them messages.
                delete(b.clients, s)
                // log.Println("Removed client")
            case msg := <-b.messages:
                // There is a new message to send. For each
                // attached client, push the new message
                // into the client's message channel.
                for s, _ := range b.clients {
                    s <- msg
                }
                // log.Printf("Broadcast message to %d clients", len(b.clients))
            }
        }
    }()
}

func NewBroker() *Broker {
    return &Broker{
        make(map[chan *Message]bool),
        make(chan (chan *Message)),
        make(chan (chan *Message)),
        make(chan *Message),
    }
}
