// Prints its command line arguments
package main

import (
	"fmt"
	"strings"
)
func EfficientEcho(args []string) {
	fmt.Println(strings.Join(args, " "));
}
func InefficientEcho(args []string) {
	var s string
	var sep = ""
	for _, arg := range args {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}