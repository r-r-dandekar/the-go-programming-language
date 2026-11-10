// comma inserts commas in a non-negative decimal integer string
package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "1000000000"
	s = comma(s)
	fmt.Printf("%v\n", s)
}

func comma(s string) string {
	l := utf8.RuneCountInString(s)
	commaCount := (l - 1) / 3
	var buf bytes.Buffer
	if commaCount > 0 {
		k := (l-1)%3 + 1
		for _, r := range s {
			if k == 0 {
				buf.WriteRune(',')
			}
			k = (k - 1) % 3
			buf.WriteRune(r)
		}
		return buf.String()
	} else {
		return s
	}
}
