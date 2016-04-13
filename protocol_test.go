package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONPatch(t *testing.T) {
	s := []byte(`{"1":17,"2":23,"n":"hello\nworld","id":"x","t":"2016-04-13T16:10:18.877Z"}`)
	p := &patch{}
	require.Nil(t, json.Unmarshal(s, p))
	assert.Equal(t, 17, p.Start)
	assert.Equal(t, 23, p.End)
	assert.Equal(t, "hello\nworld", p.NewText)
	b, err := json.Marshal(p)
	require.Nil(t, err)
	p2 := &patch{}
	require.Nil(t, json.Unmarshal(b, p2))
	assert.Exactly(t, p, p2)
}

func BenchmarkPatchUnJSON(b *testing.B) {
	s := []byte(`{"1":17,"2":23,"n":"hello\nworld","id":"x","t":"2016-04-13T16:10:18.877Z"}`)
	p := &patch{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Unmarshal(s, p)
	}
}

func BenchmarkPatchJSON(b *testing.B) {
	p := &patch{17, 23, "hello", "", time.Time{}, ""}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(p)
	}
}
