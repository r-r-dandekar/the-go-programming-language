// Mandelbrot hosts a server that serves a PNG image of the Mandelbrot fractal
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler)
	address := "localhost:8000"
	fmt.Printf("Starting server at: http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	x, ok := getFloatParam(r, w, "x", 0)
	if !ok {
		return
	}

	y, ok := getFloatParam(r, w, "y", 0)
	if !ok {
		return
	}

	zoom, ok := getFloatParam(r, w, "zoom", 1)
	if !ok {
		return
	}

	mandelbrotImage(w, x, y, zoom)
}

// Helper to parse a required float query parameter
func getFloatParam(r *http.Request, w http.ResponseWriter, key string, defaultVal float64) (float64, bool) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal, true
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid parameter '%s': %v", key, err), http.StatusBadRequest)
		return 0, false
	}
	return val, true
}
