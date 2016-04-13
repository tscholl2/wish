package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	n := note{"hello world", nil}
	arr := []patch{
		patch{0, 1, "a", "h", time.Time{}, "x1"},
		patch{1, 2, "b", "e", time.Time{}, "x2"},
	}
	n.update(arr)
	assert.Equal(t, "abllo world", n.text)
	n.update([]patch{})
	assert.Equal(t, "abllo world", n.text)
	n.unapply(arr[1])
	assert.Equal(t, "aello world", n.text)
}
