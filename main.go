package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

//go:generate esc -o static.go -prefix static static
// remember to run `go get github.com/mjibson/esc`
func main() {
	// FS() is created by `esc` and returns a http.Filesystem.
	http.Handle("/", http.FileServer(FS(true)))
	http.HandleFunc("/ws/", wsHandler())
	fmt.Println("serving")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler() http.HandlerFunc {
	H := new(hotel).run()
	return func(w http.ResponseWriter, r *http.Request) {
		var name string
		if p := strings.Split(r.URL.Path, "/"); len(p) > 2 && p[2] != "" {
			name = p[2]
		} else {
			// TODO: autogenerate names
			name = fmt.Sprintf("%d", rand.Intn(1))
		}
		c, err := newConnection(w, r)
		if err != nil {
			return
		}
		H.enter <- roomRequest{c, name}
	}
}
