package main

import "testing"

func BenchmarkInefficient(b *testing.B) {
	const n_args = 10000
	arr := makeArr(n_args)
	for i := 0; i < b.N; i++ {
		InefficientEcho(arr);
	}
}

func BenchmarkEfficient(b *testing.B) {
	const n_args = 10000
	arr := makeArr(n_args)
	for i := 0; i < b.N; i++ {
		EfficientEcho(arr);
	}
}

func makeArr(n int) []string {
	var arr []string
	for i := 0; i < n; i++ {
		arr = append(arr, "benchmarking!")
	}
	return arr
}