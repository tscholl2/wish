package main

import "errors"

type note struct {
	text    string
	patches []patch
}

func (n *note) update(patches []patch) error {
	for _, p := range patches {

		if err := n.apply(p); err != nil {
			return err
		}
	}
	return nil
}

func (n *note) addPatch(p patch) error {
	if len(n.patches) > 0 && p.Time.Before(n.patches[0].Time) {
		return errors.New("too old")
	}
	// find insertion point
	i := len(n.patches) - 1
	if i < 0 {
		i = 0
	}
	for i > 0 && n.patches[i].Time.Before(p.Time) {
		i--
	}
	// unapply post patches
	for j := len(n.patches) - 1; j > i; j-- {
		n.unapply(n.patches[j])
	}
	// insert and apply
	newPatches := append(p, n.patches[:i])
	n.patches = append(append(n.patches[:i], p), n.patches[i:]...)
	for j := i; i < len(n.patches); j++ {
		n.apply(p)
	}
}

func (n *note) apply(p patch) error {
	if p.Start < 0 || p.Start > len(n.text) || p.End < 0 || p.End > len(n.text) {
		return errors.New("index out of bounds")
	}
	n.text = n.text[:p.Start] + p.NewText + n.text[p.End:]
	return nil
}

func (n *note) unapply(p patch) error {
	if p.Start < 0 || p.Start > len(n.text) || p.Start+len(p.NewText) > len(n.text) {
		return errors.New("patch index out of bounds")
	}
	n.text = n.text[:p.Start] + p.OldText + n.text[p.Start+len(p.NewText):]
	return nil
}
