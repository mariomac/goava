package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myErr struct{ inside int }

func (myErr) Error() string {
	//TODO implement me
	panic("implement me")
}

func TestAs_Matching(t *testing.T) {
	wrappedErr := fmt.Errorf("wrapped: %w", myErr{inside: 123})
	err, ok := As[myErr](wrappedErr)
	assert.True(t, ok)
	assert.Equal(t, 123, err.inside)
}

func TestAs_NotMatching(t *testing.T) {
	wrappedErr := fmt.Errorf("wrapped: %w", errors.New("foo"))
	_, ok := As[myErr](wrappedErr)
	assert.False(t, ok)
}
