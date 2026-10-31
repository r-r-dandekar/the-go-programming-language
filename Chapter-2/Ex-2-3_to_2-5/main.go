package main

import (
	"popcount/singleexp"
	"popcount/loop"
	"popcount/shifting"
	"popcount/clearing"
	"math/rand/v2"
	"fmt"
)

func main() {
	n := rand.Uint64()
	fmt.Printf("%20s : %v -> %v\n", "single expression", n, singleexp.PopCount(n))
	fmt.Printf("%20s : %v -> %v\n", "loop", n, loop.PopCount(n))
	fmt.Printf("%20s : %v -> %v\n", "shifting", n, shifting.PopCount(n))
	fmt.Printf("%20s : %v -> %v\n", "clearing", n, clearing.PopCount(n))
}