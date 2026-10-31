// Package popcount/shifting calculates the population count of an int64 by shifting the argument through 64 bit positions
package shifting

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	count := 0
	for i :=0; i < 64; i++ {
		count += int(x&1)
		x >>= 1
	}
	return count
}