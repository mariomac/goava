package test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventually_Error(t *testing.T) {
	innerTest := &testing.T{}
	fatalEventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
		t.Errorf("fooo")
	})
	assert.True(t, innerTest.Failed())
}

func TestEventually_Fatal(t *testing.T) {
	innerTest := &testing.T{}
	fatalEventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
		t.FailNow()
	})
	assert.True(t, innerTest.Failed())
}

func TestEventually_Timeout(t *testing.T) {
	innerTest := &testing.T{}
	Eventually(innerTest, 10*time.Millisecond, func(_ require.TestingT) {
		time.Sleep(5 * time.Second)
	})
	assert.True(t, innerTest.Failed())
}

func TestEventually_Success(t *testing.T) {
	num := 3
	Eventually(t, 5*time.Second, func(t require.TestingT) {
		num--
		require.Equal(t, 0, num)
	})
}

func TestEventually_Interval(t *testing.T) {
	innerTest := &testing.T{}
	executions := 0
	fatalEventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
		executions++
		t.FailNow()
	}, Interval(20*time.Second))
	assert.True(t, innerTest.Failed())
	assert.Equal(t, 1, executions)
}

func TestEventually_InternallyContinueAfterAssertFail(t *testing.T) {
	continued := atomic.Bool{}
	innerTest := &testing.T{}
	Eventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
		assert.True(t, false)
		continued.Store(true)
		assert.True(t, true) // does not matter if a later assertion succeeds
	})
	assert.True(t, innerTest.Failed())
	assert.True(t, continued.Load())
}

func TestEventually_InternallyDoNotContinueAfterRequireFail(t *testing.T) {
	continued := atomic.Bool{}
	innerTest := &testing.T{}
	fatalEventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
		require.True(t, false)
		continued.Store(true)
		assert.True(t, true) // does not matter if a later assertion succeeds
	})
	assert.True(t, innerTest.Failed())
	assert.False(t, continued.Load())
}

func TestEventually_ErrorIfClauseOnlyContainsAssertions(t *testing.T) {
	// Wrapping the tests inside another Eventually clause
	continued := atomic.Bool{}
	innerTest := &testing.T{}
	Eventually(innerTest, 100*time.Millisecond, func(_ require.TestingT) {
		Eventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
			assert.True(t, false)
		})
		// should continue
		continued.Store(true)
		// inner test should not pass because previous eventually test failed
		Eventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
			assert.True(t, true)
		})
	})
	assert.True(t, innerTest.Failed())
	assert.True(t, continued.Load())
}

func TestEventually_FatalIfClauseContainsRequires(t *testing.T) {
	// Wrapping the tests inside another Eventually clause
	continued := atomic.Bool{}
	innerTest := &testing.T{}
	Eventually(innerTest, 100*time.Millisecond, func(_ require.TestingT) {
		Eventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
			require.True(t, false)
		})
		// should continue
		continued.Store(true)
		// inner test should not pass because previous eventually test failed
		Eventually(innerTest, 10*time.Millisecond, func(t require.TestingT) {
			assert.True(t, true)
		})
	})
	assert.True(t, innerTest.Failed())
	assert.False(t, continued.Load())
}

// fatalEventually wraps the execution of an eventually function that should end with fatal,
// avoiding that any internal "fatal" invocation panics because of the goruntime.Goexit invocation
func fatalEventually(t *testing.T, timeout time.Duration, testFunc func(_ require.TestingT), options ...EventuallyOption) {
	returnCh := make(chan struct{})
	go func() {
		defer close(returnCh)
		Eventually(t, timeout, testFunc, options...)
	}()

	<-returnCh
}

func TestEventually_AllGoroutinesEndBeforeEventuallyEnds(t *testing.T) {
	count := atomic.Int64{}
	count.Store(int64(0))
	for i := 0; i < 100; i++ {
		add := atomic.Int64{}
		add.Store(0)
		Eventually(t, 100*time.Millisecond, func(t require.TestingT) {
			defer count.Add(1)
			add.Add(1)
			// this test is expected to fail until number 33
			assert.EqualValues(t, 10, add.Load())
		})
		require.EqualValues(t, 10*(i+1), count.Load())
	}
}
