package main

import "fmt"

func main() {
	s := []string{"hello", "world", "world", "this", "is", "a", "a", "this", "test"}
	fmt.Println("Before: ", s)
	s = removeAdjacentDuplicates(s)
	fmt.Println("After: ", s)
}

func removeAdjacentDuplicates(s []string) []string {
	l := len(s)
	if l < 2 {
		return s
	}
	i, j := 0, 1
	for j < l {
		if s[i] != s[j] {
			i++
			s[i] = s[j]
		}
		j++
	}
	return s[:i+1]
}
