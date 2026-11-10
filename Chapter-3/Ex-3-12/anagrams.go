// anagrams reports whether two strings are anagrams of each other
package main

import "fmt"

func main() {
	var s1, s2 string
	s1, s2 = "hello", "ollhe"
	fmt.Printf("anagrams(%s, %s) = %v\n", s1, s2, anagrams(s1, s2))
	s1, s2 = "er", "rt"
	fmt.Printf("anagrams(%s, %s) = %v\n", s1, s2, anagrams(s1, s2))
	s1, s2 = "ee", "ollhe"
	fmt.Printf("anagrams(%s, %s) = %v\n", s1, s2, anagrams(s1, s2))
	s1, s2 = "341", "341"
	fmt.Printf("anagrams(%s, %s) = %v\n", s1, s2, anagrams(s1, s2))
}

func anagrams(s1, s2 string) bool {
	m1 := make(map[rune]int)
	m2 := make(map[rune]int)

	for _, r := range s1 {
		if _, exists := m1[r]; exists {
			m1[r]++
		} else {
			m1[r] = 1
		}
	}
	for _, r := range s2 {
		if _, exists := m2[r]; exists {
			m2[r]++
		} else {
			m2[r] = 0
		}
	}

	if len(m1) != len(m2) {
		return false
	}
	for key, val1 := range m1 {
		if val2, exists := m2[key]; exists {
			if val1 != val2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
