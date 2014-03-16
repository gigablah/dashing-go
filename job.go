package dashing

type Job interface {
    Work(send chan *Message)
}

var registry = []Job{}

func Register(j Job) {
    if j == nil {
        panic("Can't register nil job")
    }
    registry = append(registry, j)
}
