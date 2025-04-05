package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFreeTCPPort(t *testing.T) {
	tp, err := FreeTCPPort()
	require.NoError(t, err)
	assert.NotZero(t, tp)
}

func TestFreeUDPPort(t *testing.T) {
	tp, err := FreeUDPPort()
	require.NoError(t, err)
	assert.NotZero(t, tp)
}
