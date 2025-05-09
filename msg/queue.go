package msg

type queueConfig struct {
	channelBufferLen int
}

var defaultQueueConfig = queueConfig{
	channelBufferLen: 1,
}

// QueueOpts allow configuring some operation of a queue
type QueueOpts func(*queueConfig)

// ChannelBufferLen sets the length of the channel buffer for the queue.
func ChannelBufferLen(l int) QueueOpts {
	return func(c *queueConfig) {
		c.channelBufferLen = l
	}
}

// Queue is a simple message queue that allows sending messages to multiple subscribers.
// It also allows bypassing messages to other queues, so that a message sent to one queue
// can be received by subscribers of another queue.
// If a message is sent to a queue that has no subscribers, it will not block the sender and the
// message will be lost. This is by design, as the queue is meant to be used for fire-and-forget
type Queue[T any] struct {
	cfg  *queueConfig
	dsts []chan T
	// double-linked list of bypassing queues
	// For simplicity, a Queue instance:
	// - can't bypass to a queue and having other dsts
	// - can only bypass to a single queue, despite multiple queues can bypass to it
	bypassTo *Queue[T]
}

// NewQueue creates a new Queue instance with the given options.
func NewQueue[T any](opts ...QueueOpts) *Queue[T] {
	cfg := defaultQueueConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Queue[T]{cfg: &cfg}
}

func (q *Queue[T]) config() *queueConfig {
	if q.cfg == nil {
		return &defaultQueueConfig
	}
	return q.cfg
}

// Send a message to all subscribers of this queue. If there are no subscribers,
// the message will be lost and the sender will not be blocked.
func (q *Queue[T]) Send(o T) {
	if q.bypassTo != nil {
		q.bypassTo.Send(o)
		return
	}
	for _, d := range q.dsts {
		d <- o
	}
}

// Subscribe to this queue. This will return a channel that will receive messages.
// This operation is not thread-safe.
// You can't subscribe to a queue that is bypassing to another queue.
func (q *Queue[T]) Subscribe() <-chan T {
	if q.bypassTo != nil {
		panic("this queue is already bypassing data to another queue. Can't subscribe to it")
	}
	out := make(chan T, q.config().channelBufferLen)
	q.dsts = append(q.dsts, out)
	return out
}

// Bypass allows this queue to bypass messages to another queue. This means that
// messages sent to this queue will also be sent to the other queue.
// This operation is not thread-safe and does not control for graph cycles.
func (q *Queue[T]) Bypass(to *Queue[T]) {
	if q == to {
		panic("this queue can't bypass to itself")
	}
	if q.bypassTo != nil {
		panic("this queue is already bypassing to another queue")
	}
	q.bypassTo = to
}
