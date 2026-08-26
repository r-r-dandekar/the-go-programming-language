// Prints its command line arguments
package main

import (
	"fmt"
	"os"
	"time"
)

func ineff_old() {
	start := time.Now()
	var s string
	var sep = ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
	secs := time.Since(start).Seconds()
	fmt.Println("Time: ", secs)
}