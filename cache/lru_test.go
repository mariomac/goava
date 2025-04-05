package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type value string

func (v value) SizeBytes() int {
	return len(v)
}

func TestPutGetRemove(t *testing.T) {
	lru := NewLRU[string, value](1000)

	lru.Put("foo", "bar")
	v, ok := lru.Get("foo")
	assert.True(t, ok)
	assert.EqualValues(t, "bar", string(v))

	_, ok = lru.Get("baz")
	assert.False(t, ok)

	lru.Put("baz", "bae")
	v, ok = lru.Get("baz")
	assert.True(t, ok)
	assert.EqualValues(t, "bae", string(v))

	lru.Remove("foo")
	_, ok = lru.Get("foo")
	assert.False(t, ok)
}

func TestLRUOldestRemoval(t *testing.T) {
	// GIVEN a cache with preexisting values
	lru := NewLRU[string, value](10)
	lru.Put("foo", "bar")
	lru.Put("baz", "bae")
	lru.Put("tri", "tra")
	v, ok := lru.Get("foo")
	require.True(t, ok)
	require.EqualValues(t, "bar", string(v))
	v, ok = lru.Get("baz")
	require.True(t, ok)
	require.EqualValues(t, "bae", string(v))
	v, ok = lru.Get("tri")
	require.True(t, ok)
	require.EqualValues(t, "tra", string(v))

	// WHEN we add a value that overflows the max size
	lru.Put("toma", "ya!")

	// THEN the LRU value is removed
	_, ok = lru.Get("foo")
	assert.False(t, ok)
	v, ok = lru.Get("baz")
	assert.True(t, ok)
	assert.EqualValues(t, "bae", string(v))
	v, ok = lru.Get("tri")
	assert.True(t, ok)
	assert.EqualValues(t, "tra", string(v))
	v, ok = lru.Get("toma")
	assert.True(t, ok)
	assert.EqualValues(t, "ya!", string(v))

	// AND when we access a given element
	lru.Get("baz")
	// THEN it is not the LRU anymore
	lru.Put("chin", "cha")
	v, ok = lru.Get("baz")
	assert.True(t, ok)
	assert.EqualValues(t, "bae", string(v))
	_, ok = lru.Get("tri")
	assert.False(t, ok)

	// AND when the overridable value is bigger than multiple elements
	lru.Put("toma", "1234567")
	// THEN multiple elements are evicted from the cache
	_, ok = lru.Get("foo")
	assert.False(t, ok)
	_, ok = lru.Get("tri")
	assert.False(t, ok)
	_, ok = lru.Get("chin")
	assert.False(t, ok)

	v, ok = lru.Get("toma")
	assert.True(t, ok)
	assert.EqualValues(t, "1234567", string(v))

	v, ok = lru.Get("baz")
	assert.True(t, ok)
	assert.EqualValues(t, "bae", string(v))
}
