package msg

// Queue is a simple message queue that allows sending messages to multiple subscribers.
// It also allows bypassing messages to other queues, so that a message sent to one queue
// can be received by subscribers of another queue.
// If a message is sent to a queue that has no subscribers, it will not block the sender and the
// message will be lost. This is by design, as the queue is meant to be used for fire-and-forget
type Queue[T any] struct {
	dsts []chan T
	// double-linked list of bypassing queues
	// For simplicity, a Queue instance:
	// - can't bypass to a queue and having other dsts
	// - can only bypass to a single queue, despite multiple queues can bypass to it
	bypassTo *Queue[T]
}

// Send a message to all subscribers of this queue. If there are no subscribers,
// the message will be lost and the sender will not be blocked.
func (h *Queue[T]) Send(o T) {
	if h.bypassTo != nil {
		h.bypassTo.Send(o)
		return
	}
	for _, d := range h.dsts {
		d <- o
	}
}

// Subscribe to this queue. This will return a channel that will receive messages.
// This operation is not thread-safe.
// You can't subscribe to a queue that is bypassing to another queue.
func (h *Queue[T]) Subscribe() <-chan T {
	if h.bypassTo != nil {
		panic("this queue is already bypassing data to another queue. Can't subscribe to it")
	}
	out := make(chan T, 1)
	h.dsts = append(h.dsts, out)
	return out
}

// Bypass allows this queue to bypass messages to another queue. This means that
// messages sent to this queue will also be sent to the other queue.
// This operation is not thread-safe and does not control for graph cycles.
func (h *Queue[T]) Bypass(to *Queue[T]) {
	if h == to {
		panic("this queue can't bypass to itself")
	}
	if h.bypassTo != nil {
		panic("this queue is already bypassing to another queue")
	}
	h.bypassTo = to
}
