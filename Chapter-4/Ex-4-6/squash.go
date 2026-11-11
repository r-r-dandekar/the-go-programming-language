package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	s := []byte("Hello, this   is a      string!")
	fmt.Println("Before: ", string(s))
	s = squash(s)
	fmt.Println("After: ", string(s))
}

func squash(str []byte) []byte {
	s := str
	l := len(str)
	ptr := 0
	space := false
	for len(s) > 0 {
		r, size := utf8.DecodeRune(s)
		if unicode.IsSpace(r) && space {
			l -= size
		} else {
			utf8.EncodeRune(str[ptr:], r)
			ptr += size
		}

		if unicode.IsSpace(r) {
			space = true
		} else {
			space = false
		}
		s = s[size:]
	}
	return str[:l]
}
