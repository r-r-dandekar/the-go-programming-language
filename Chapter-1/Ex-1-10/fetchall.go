// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "fetchall: needs at least 2 args\n")
		os.Exit(1)
	}
	filename := os.Args[1]
	outfile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetchall: writing to %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer outfile.Close()

	for _, url := range os.Args[2:] {
		prefix := "http://"
		if !strings.HasPrefix(url, prefix) {
			url = prefix + url
		}
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[2:] {
		fmt.Fprintf(outfile, "%v\n", <-ch) // receive from channel ch
	}
	fmt.Printf("%2.fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	var buf bytes.Buffer
	start := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(&buf, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	sec := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%s\n\n%.2fs  %7d  %s", buf.String(), sec, nbytes, url)
}
