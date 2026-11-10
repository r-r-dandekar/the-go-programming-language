// comma inserts commas in a non-negative decimal integer string
package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	s := "-1000000.0000000000000000000"
	s = comma(s)
	fmt.Printf("%v\n", s)
}

func comma(s string) string {
	var buf bytes.Buffer
	var suffix string = ""
	if s[0] == '+' || s[0] == '-' {
		buf.WriteByte(s[0])
		s = s[1:]
	}
	if i := strings.Index(s, "."); i != -1 {
		suffix = s[i:]
		s = s[:i]
	}
	l := len(s)
	commaCount := (l - 1) / 3
	if commaCount > 0 {
		k := (l-1)%3 + 1
		for i := 0; i < l; i++ {
			if k == 0 {
				buf.WriteByte(',')
			}
			k = (k - 1) % 3
			buf.WriteByte(s[i])
		}
		return buf.String() + suffix
	} else {
		return s + suffix
	}
}
