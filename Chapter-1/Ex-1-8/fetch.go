// Fetch prints the content found at a URL.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		protocol := "http://"
		if !strings.HasPrefix(url, protocol) {
			url = protocol + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		var count int64
		count = 1
		for count > 0 {
			count, err = io.Copy(os.Stdout, resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
				resp.Body.Close()
				os.Exit(1)
			}
		}
		resp.Body.Close()
	}
}
