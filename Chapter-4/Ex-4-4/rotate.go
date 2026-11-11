package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	s = rotate(s, 4)
	fmt.Println(s)
	s = []int{1, 2, 3, 4, 5, 6, 7}
	s = rotateInPlace(s, 4)
	fmt.Println(s)
}

// This is not as efficient as the in-place method
func rotate(s []int, k int) []int {
	if k == 0 {
		return s
	}
	l := len(s)
	k = k % l
	q := make([]int, k) // circular queue
	for i := 0; i < k; i++ {
		q[i] = s[l-k+i]
	}
	head := 0
	tail := k - 1

	for i := 0; i < l; i++ {
		temp := s[i]
		s[i] = q[head]
		head = (head + 1) % k
		tail = (tail + 1) % k
		q[tail] = temp
	}

	return s
}

func rotateInPlace(s []int, k int) []int {
	if k == 0 {
		return s
	}
	l := len(s)
	k = k % l

	reverse(s)
	reverse(s[:k])
	reverse(s[k:])

	return s
}

// reverse reverses a slice of ints in place.
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
