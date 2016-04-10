package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	n := note{"hello world"}
	arr := []patch{
		patch{0, 1, "a"},
		patch{1, 2, "b"},
	}
	n.update(arr)
	assert.Equal(t, "abllo world", n.text)
	n.update([]patch{})
	assert.Equal(t, "abllo world", n.text)
}
