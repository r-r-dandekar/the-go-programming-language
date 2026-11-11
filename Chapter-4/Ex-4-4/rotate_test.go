package main

import (
	"math/rand/v2"
	"testing"
)

func BenchmarkRotate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := rand.Int()%10000 + 1
		s := make([]int, n)
		for i := range n {
			s[i] = i
		}
		for j := 0; j < 1000; j++ {
			k := rand.Int() % n
			rotate(s, k)
		}
	}
}

func BenchmarkRotateInPlace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := rand.Int()%10000 + 1
		s := make([]int, n)
		for i := range n {
			s[i] = i
		}
		for j := 0; j < 1000; j++ {
			k := rand.Int() % n
			rotateInPlace(s, k)
		}
	}
}
