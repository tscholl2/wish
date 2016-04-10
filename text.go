package main

import "fmt"

type note struct {
	text string
}

func (n *note) update(patches []patch) (err error) {
	for _, p := range patches {
		if p.Start < 0 || p.End > len(n.text) {
			return fmt.Errorf("note: out of bounds update")
		}
		n.text = n.text[:p.Start] + p.Text + n.text[p.End:]
	}
	return
}
