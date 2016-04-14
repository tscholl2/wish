package main

import "errors"

type note struct {
	text    string
	patches []patch
}

func (n *note) update(patches []patch) error {
	for _, p := range patches {
		// TODO sort patches by time and deal with late ones
		if err := n.apply(p); err != nil {
			return err
		} else {
			n.patches = append(n.patches, p)
			for len(n.patches) > 100 {
				n.patches = n.patches[1:]
			}
		}
	}
	return nil
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
