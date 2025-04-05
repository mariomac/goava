package casing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelToDots(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "camel.case"},
		{"dromedaryCase", "dromedary.case"},
		{"SimpleTest", "simple.test"},
		{"AnotherExampleTest", "another.example.test"},
	}

	for _, test := range tests {
		result := CamelToDots(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "camel_case"},
		{"dromedaryCase", "dromedary_case"},
		{"SimpleTest", "simple_test"},
		{"AnotherExampleTest", "another_example_test"},
	}

	for _, test := range tests {
		result := CamelToSnake(test.input)
		assert.Equal(t, test.expected, result)
	}
}
