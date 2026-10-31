package main

import (
	"testing"
	"popcount/singleexp"
	"popcount/loop"
	"popcount/shifting"
	"popcount/clearing"
	"math/rand/v2"
)

var k = 100000

func BenchmarkSingleExpression(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := rand.Uint64()
		for j := 0; j < k; j++ {
			singleexp.PopCount(n)
		}
	}
}

func BenchmarkLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := rand.Uint64()
		for j := 0; j < k; j++ {
			loop.PopCount(n)
		}
	}
}

func BenchmarkShifting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := rand.Uint64()
		for j := 0; j < k; j++ {
			shifting.PopCount(n)
		}
	}
}

func BenchmarkClearing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := rand.Uint64()
		for j := 0; j < k; j++ {
			clearing.PopCount(n)
		}
	}
}
