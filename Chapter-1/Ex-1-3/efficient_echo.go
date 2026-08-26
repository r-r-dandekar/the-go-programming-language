// Prints its command line arguments
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func eff_old() {
	start := time.Now()
	fmt.Println(strings.Join(os.Args[1:], " "))
	secs := time.Since(start).Seconds()
	fmt.Println("Time: ", secs)
}