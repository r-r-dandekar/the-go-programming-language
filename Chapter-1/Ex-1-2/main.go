// Prints its command line arguments
package main

import (
	"fmt"
	"os"
)

func main() {
	var s string
	var sep = "\n"
	for i, arg := range os.Args[0:] {
		s += fmt.Sprintf("%2d : %s%s", i, arg, sep)
	}
	fmt.Println(s)
}