package main

import "fmt"

func main() {
	s := []byte("Hello, ひらがな, this is a string!")

	fmt.Println("Before: ", string(s))
	s = reverseUTF8(s)
	fmt.Println("Reversed: ", string(s))
	s = reverseUTF8(s)
	fmt.Println("Re-reversed: ", string(s))
}

func reverseUTF8(s []byte) []byte {
	n := len(s)

	// First, reverse the slice
	reverse(s)

	// Then traverse backwards, and for each byte that indicates the
	// start of a multi-byte utf-8 character, reverse those many bytes
	for i := n - 1; i >= 0; i-- {
		if s[i]&0xE0 == 0xA0 { // 110xxxxx indicates 2 bytes
			reverse(s[i-1 : i+1])
			i -= 1
		} else if s[i]&0xF0 == 0xE0 { // 1110xxxx indicates 3 bytes
			reverse(s[i-2 : i+1])
			i -= 2
		} else if s[i]&0xF8 == 0xF0 { // 11110xxx indicates 4 bytes
			reverse(s[i-3 : i+1])
			i -= 3
		}
	}

	return s
}

// reverse reverses a slice of bytes in place.
func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
