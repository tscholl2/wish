package main

import "github.com/boltdb/bolt"

type store struct {
	db *bolt.DB
}

func newStore(filename string) (*store, error) {
	db, err := bolt.Open(filename, 0600, nil)
	return &store{db: db}, err
}
