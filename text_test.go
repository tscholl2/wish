package main

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	txt := new(text)
	for i := 0; i < 100; i++ {
		txt.update([]diffmatchpatch.Patch{diffmatchpatch.Patch{}})
	}
	assert.Equal(t, 10, len(txt.patches))
	assert.Equal(t, "", txt.snapshot)
}
