package main

import "fmt"

const l = 7

func main() {
	arr := [l]int{1, 2, 3, 4, 5, 6, 7}
	reverse(&arr)
	fmt.Println(arr)
}

func reverse(arr *[l]int) {
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
