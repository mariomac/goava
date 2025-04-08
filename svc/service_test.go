package svc

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mariomac/guara/test"
)

func TestSwarm_RegisterAfterStart(t *testing.T) {
	var s Swarm
	s.started.Store(true)
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	s.Register(func(_ context.Context) (Runner, error) {
		return func(_ context.Context) {}, nil
	})
}
func TestSwarm_StartWithError(t *testing.T) {
	var s Swarm
	s.Register(func(_ context.Context) (Runner, error) {
		return nil, errors.New("creation error")
	})
	assert.Error(t, s.Start(context.Background()))
}

func TestSwarm_StartTwice(t *testing.T) {
	var s Swarm
	s.Register(func(_ context.Context) (Runner, error) {
		return func(_ context.Context) {}, nil
	})
	err := s.Start(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	_ = s.Start(context.Background())
}

func TestSwarm_RunnerExecution(t *testing.T) {
	var s Swarm
	runnerExecuted := atomic.Bool{}
	s.Register(func(_ context.Context) (Runner, error) {
		return func(_ context.Context) {
			runnerExecuted.Store(true)
		}, nil
	})
	require.NoError(t, s.Start(context.Background()))
	test.Eventually(t, 5*time.Second, func(t require.TestingT) {
		assert.True(t, runnerExecuted.Load(), "runner was not executed")
	})
}

func TestSwarm_CreatorFailure(t *testing.T) {
	var s Swarm
	runnerStarted := atomic.Bool{}
	c1cancel := atomic.Bool{}
	c3exec := atomic.Bool{}
	s.Register(func(ctx context.Context) (Runner, error) {
		go func() {
			<-ctx.Done()
			c1cancel.Store(true)
		}()
		return func(_ context.Context) {
			runnerStarted.Store(true)
		}, nil
	})
	s.Register(func(_ context.Context) (Runner, error) {
		return nil, errors.New("creation error")
	})
	s.Register(func(_ context.Context) (Runner, error) {
		c3exec.Store(true)
		return func(_ context.Context) {}, nil
	})

	// second creator fails, so the first one should be cancelled and the third one should not be executed
	require.Error(t, s.Start(context.Background()))
	test.Eventually(t, 5*time.Second, func(t require.TestingT) {
		assert.True(t, c1cancel.Load(), "c1 was not cancelled")
	})
	assert.False(t, c3exec.Load(), "c3 was executed")
	assert.False(t, runnerStarted.Load(), "runner was started")
}

func TestSwarm_ContextPassed(t *testing.T) {
	startWg := sync.WaitGroup{}
	startWg.Add(3)
	doneWg := sync.WaitGroup{}
	doneWg.Add(3)
	s := Swarm{}
	innerRunner := func(ctx context.Context) {
		startWg.Done()
		<-ctx.Done()
		doneWg.Done()
	}
	s.Register(func(_ context.Context) (Runner, error) { return innerRunner, nil })
	s.Register(func(_ context.Context) (Runner, error) { return innerRunner, nil })
	s.Register(func(_ context.Context) (Runner, error) { return innerRunner, nil })
	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, s.Start(ctx))
	test.Eventually(t, 5*time.Second, func(_ require.TestingT) {
		startWg.Wait()
	})
	cancel()
	test.Eventually(t, 5*time.Second, func(_ require.TestingT) {
		doneWg.Wait()
	})

}
