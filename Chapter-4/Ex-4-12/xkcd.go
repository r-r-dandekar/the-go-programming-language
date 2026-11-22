package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"xkcd/xkcd"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "index" {
			buildIndex()
		} else {
			strs := os.Args[1:]
			for _, s := range strs {
				i, err := strconv.Atoi(s)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Invalid comic number:", err)
					continue
				}
				printComic(i)
			}
		}
	} else {
		for true {
			fmt.Print("xkcd > ")
			var num int
			if _, err := fmt.Scanln(&num); err != nil {
				fmt.Fprintln(os.Stderr, "Invalid comic number:", err)
				continue
			}
			printComic(num)
		}
	}
}

func printComic(num int) {
	c, err := xkcd.Get(num)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching comic:", err)
		return
	}

	prettyPrint(c)
}

func prettyPrint(c *xkcd.Comic) {
	transcript, extra, found := strings.Cut(c.Transcript, "{{")
	if found {
		extra = strings.TrimSuffix(extra, "}}")
		transcript = strings.TrimSpace(transcript)
		extra = strings.TrimSpace(extra)
	}

	fmt.Printf("\n%s\n======================================================================\n"+
		"%s\n----------------------------------------------------------------------\n\n"+
		"%s\n\nImg: %s\n\n", c.Title, transcript, extra, c.Img)
}

func buildIndex() {
	i := 1
	for true {
		present, err := xkcd.Exists(i)
		if err == nil && present {
			i++
			continue
		}
		time.Sleep(500 * time.Microsecond) // Limit request rate
		_, err = xkcd.FetchAndPut(i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while building index: %v\n", err)
			return
		}
		fmt.Println("Indexed comic", i)
		i++
	}
}
