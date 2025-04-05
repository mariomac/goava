package cache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type counter struct {
	writes     int
	isWriting  int64
	collisions int64
}

func (a *counter) Put(_ string, _ int) {
	a.isWriting = 1
	a.writes += 2
	time.Sleep(time.Nanosecond)
	a.writes--
	a.isWriting = 0
}

func (a *counter) Get(_ string) (value int, ok bool) {
	atomic.AddInt64(&a.collisions, atomic.LoadInt64(&a.isWriting))
	return
}

func (a *counter) Remove(_ string) {
	a.isWriting = 1
	a.writes += 2
	time.Sleep(time.Nanosecond)
	a.writes--
	a.isWriting = 0
}

func TestConcurrentCache(t *testing.T) {
	ic := counter{}
	c := NewConcurrent[string, int](&ic)
	wg := sync.WaitGroup{}
	wg.Add(4)
	for t := 0; t < 4; t++ {
		go func() {
			for i := 0; i < 1000; i++ {
				c.Put("a", 1)
				c.Get("3")
				c.Remove("f")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, 8000, ic.writes)
	assert.Zero(t, ic.collisions)
}
