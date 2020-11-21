package broker

type Status int

const (
	Success Status = iota
	Retry
)

type Queue struct {
	Name       string
	RouteKey   string
	RetryQueue []int64
	Handle     func([]byte) Status
}

type Broker interface {
	Consume(queue *Queue) error
	Publish(key string, body []byte) error
}
