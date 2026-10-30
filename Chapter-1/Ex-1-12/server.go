package main

import (
	"log"
	"net/http"
	"strconv"
)

// Web server to server Lissajous pattern GIFs

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	cycles := 5
	s := r.URL.Query().Get("cycles")
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Print(err)
		} else {
			cycles = i
		}
	}
	lissajous(w, cycles)
}
