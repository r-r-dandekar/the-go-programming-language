// popcount calculates the population count of an int64 by a lookup table and a single expression
package main

// pc[i] is the population count of i
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint8) int {
	return int(
		pc[byte(x)])
}
