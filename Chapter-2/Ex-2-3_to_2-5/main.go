package main

import (
	"fmt"
	"math/rand/v2"
	"popcount/clearing"
	"popcount/loop"
	"popcount/shifting"
	"popcount/singleexp"
)

func main() {
	n := rand.Uint64()
	fmt.Printf("%20s : %v -> %v\n", "single expression", n, singleexp.PopCount(n))
	fmt.Printf("%20s : %v -> %v\n", "loop", n, loop.PopCount(n))
	fmt.Printf("%20s : %v -> %v\n", "shifting", n, shifting.PopCount(n))
	fmt.Printf("%20s : %v -> %v\n", "clearing", n, clearing.PopCount(n))
}
