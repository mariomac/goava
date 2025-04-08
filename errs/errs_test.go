package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myError struct{ inside int }

func (myError) Error() string {
	// TODO implement me
	panic("implement me")
}

func TestAs_Matching(t *testing.T) {
	wrappedErr := fmt.Errorf("wrapped: %w", myError{inside: 123})
	err, ok := As[myError](wrappedErr)
	assert.True(t, ok)
	assert.Equal(t, 123, err.inside)
}

func TestAs_NotMatching(t *testing.T) {
	wrappedErr := fmt.Errorf("wrapped: %w", errors.New("foo"))
	_, ok := As[myError](wrappedErr)
	assert.False(t, ok)
}
