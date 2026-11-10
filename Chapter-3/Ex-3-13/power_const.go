// power_const has constant declarations for KB, MB upto YB
package main

import "fmt"

const (
	KB = 1000
	MB = KB * KB
	GB = MB * KB
	TB = GB * KB
	PB = TB * KB
	EB = PB * KB
	ZB = EB * KB
	YB = ZB * KB
)

func main() {
	fmt.Printf("KB : %v\n", KB)
	fmt.Printf("MB : %v\n", MB)
	fmt.Printf("GB : %v\n", GB)
	fmt.Printf("EB : %v\n", EB)
}
