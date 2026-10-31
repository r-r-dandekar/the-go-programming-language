// Package popcount/loop calculates the population count of an int64 by a lookup table and a loop
package loop

// pc[i] is the population count of i
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	count := 0
	for i :=0; i < 8; i++ {
		count += int(pc[byte(x >> (i*8))])
	}
	return count
}