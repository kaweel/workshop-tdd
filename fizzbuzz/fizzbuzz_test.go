package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{1, "1"},
		{3, "Fizz"},
		{5, "Buzz"},
		{15, "FizzBuzz"},
		{30, "FizzBuzz"},
		{7, "7"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("FizzBuzz(%d)", tt.n), func(t *testing.T) {
			result := FizzBuzz(tt.n)
			assert.Equal(t, tt.expected, result)
		})
	}
}
