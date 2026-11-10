package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	c1, c2 := "x", "X"
	h1 := sha256.Sum256([]byte(c1))
	h2 := sha256.Sum256([]byte(c2))
	d := shaDiff(h1, h2)
	fmt.Printf("%v vs %v : %v bytes differ\n", c1, c2, d)
}

func shaDiff(h1, h2 [32]uint8) int {
	count := 0
	for i := range 32 {
		count += PopCount(h1[i] ^ h2[i])
	}
	return count
}
