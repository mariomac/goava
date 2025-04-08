package svc

import (
	"context"
	"sync"
	"sync/atomic"
)

// Swarm coordinates a set of service runners that start and end at the same time.
type Swarm struct {
	started  atomic.Bool
	creators []Creator
}

// Creator is a function that creates a service Runner.
type Creator func(context.Context) (Runner, error)

// Runner is a function that runs a service. The service should be stopped via the passed context Done function.
type Runner func(context.Context)

// Register a service creator to the swarm. The creator will be called when the swarm starts, and must return
// a Runner instance that will execute the actual operation of the service.
func (s *Swarm) Register(c Creator) {
	if s.started.Load() {
		panic("service.Swarm already started")
	}
	s.creators = append(s.creators, c)
}

// Start the Swarm. It calls all registered service creators and, if all succeed,
// it runs the returned Runner instances.
// If any of the creators return an error, the swarm will stop and return the error. The context
// that is passed to the rest of creators will be cancelled.
// No service Runner internal instance is started until all the Creators are successfully
// created. This means that if any of the creators fail, no service Runner is started.
func (s *Swarm) Start(ctx context.Context) error {
	if s.started.Swap(true) {
		panic("service.Swarm already started")
	}
	buildCtx, cancel := context.WithCancel(ctx)
	runners := make([]Runner, 0, len(s.creators))
	for _, creator := range s.creators {
		runner, err := creator(buildCtx)
		if err != nil {
			cancel()
			return err
		}
		runners = append(runners, runner)
	}
	// call the cancel function to avoid leaking the context
	wg := sync.WaitGroup{}
	wg.Add(len(runners))
	go func() {
		wg.Wait()
		cancel()
	}()
	for i := range runners {
		runner := runners[i]
		go func() {
			runner(ctx)
			wg.Done()
		}()
	}
	return nil
}
