package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)

	wordcount := make(map[string]int)

	for input.Scan() {
		t := input.Text()
		wordcount[t]++
	}

	if err := input.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}

	for key, value := range wordcount {
		fmt.Printf("%v : %v\n", key, value)
	}
}
