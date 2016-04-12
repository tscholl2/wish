package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONPatch(t *testing.T) {
	s := []byte(`{"1":17,"2":23,"s":"hello\nworld"}`)
	p := &patch{}
	require.Nil(t, json.Unmarshal(s, p))
	assert.Equal(t, 17, p.Start)
	assert.Equal(t, 23, p.End)
	assert.Equal(t, "hello\nworld", p.Text)
	b, err := json.Marshal(p)
	require.Nil(t, err)
	p2 := &patch{}
	require.Nil(t, json.Unmarshal(b, p2))
	assert.Exactly(t, p, p2)
}

func BenchmarkPatchUnJSON(b *testing.B) {
	s := []byte(`{"1":17,"2":23,"s":"hello\nworld"}`)
	p := &patch{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Unmarshal(s, p)
	}
}

func BenchmarkPatchJSON(b *testing.B) {
	p := &patch{17, 23, "hello"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(p)
	}
}
