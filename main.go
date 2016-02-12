package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

//go:generate esc -o static.go -prefix static static
func main() {
	// FS() is created by `esc` and returns a http.Filesystem.
	http.Handle("/", http.FileServer(FS(false)))
	http.HandleFunc("/ws/", wsHandler())
	fmt.Println("serving")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler() http.HandlerFunc {
	H := newHotel()
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got ws req at %s\n", r.URL.Path)
		var name string
		if p := strings.Split(r.URL.Path, "/"); len(p) > 2 && p[2] != "" {
			name = p[2]
		} else {
			// TODO: autogenerate names
			name = fmt.Sprintf("%d", rand.Intn(2))
		}
		c, err := newConnection(w, r)
		if err != nil {
			return
		}
		H.enter <- roomRequest{c, name}
	}
}
