package main

import "github.com/sergi/go-diff/diffmatchpatch"

type text struct {
	dmp      *diffmatchpatch.DiffMatchPatch
	snapshot string
	patches  []diffmatchpatch.Patch
}

func (t *text) update(patches []diffmatchpatch.Patch) error {
	if t.dmp == nil {
		t.dmp = diffmatchpatch.New()
	}
	s, _ := t.dmp.PatchApply(patches, t.snapshot) // TODO: error handling?
	t.patches = append(t.patches, patches...)
	if len(t.patches) > 10 {
		t.patches = t.patches[len(t.patches)-10 : len(t.patches)]
	}
	t.snapshot = s
	return nil
}
