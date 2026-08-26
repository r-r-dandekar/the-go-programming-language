// Prints its command line arguments
package main

import (
	"fmt"
	"os"
)

func main() {
	var s string
	var sep = "\n"
	for _, arg := range os.Args[0:] {
		s += arg + sep
	}
	fmt.Println(s)
}