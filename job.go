package dashing

// A Job does periodic work and sends events to a channel.
type Job interface {
    Work(send chan *Event)
}

var registry = []Job{}

// Register a job to be kicked off upon starting the server.
func Register(j Job) {
    if j == nil {
        panic("Can't register nil job")
    }
    registry = append(registry, j)
}
