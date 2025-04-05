package rate

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccepter(t *testing.T) {
	accepter := NewAccepter(1, 10*time.Millisecond)
	start := time.Now()
	end := start.Add(100 * time.Millisecond)
	requests := 0
	for time.Now().Before(end) {
		if accepter.Accept() {
			requests++
		}
	}
	// 1 request every millisecond: around 10 effective requests maximum
	assert.InDelta(t, 10, requests, 1)
}

func TestAccepter_NonAccumulate(t *testing.T) {
	accepter := NewAccepter(1, 10*time.Millisecond)

	// a period of inactivity should not accumulate more requests
	// than the maximum allowed
	time.Sleep(50 * time.Millisecond)

	start := time.Now()
	requests := 0
	end := start.Add(100 * time.Millisecond)
	for time.Now().Before(end) {
		if accepter.Accept() {
			requests++
		}
	}
	// 1 request every millisecond: around 10 effective requests maximum
	assert.InDelta(t, 10, requests, 1)
}

func ExampleAccepter() {
	// limits to 10 requests every 5 seconds
	accepter := NewAccepter(10, 5*time.Second)

	start := time.Now()
	end := start.Add(10 * time.Second)
	requests := 0
	for time.Now().Before(end) {
		if accepter.Accept() {
			requests++
		}
	}

	fmt.Println("total requests:", requests)
	//// Output: total requests: 30
	// remove double comment above to run the test.
	// We want to skip it because it's very slow
}
