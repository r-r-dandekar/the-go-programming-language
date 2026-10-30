package main

import (
	"tempconv/tempconv"
	"fmt"
)

func main() {
	fmt.Printf("Hi!\n")
	var c tempconv.Celsius = 20
	var f tempconv.Fahrenheit = tempconv.CToF(c)
	var k tempconv.Kelvin = tempconv.FToK(f)
	c = tempconv.KToC(k)
 	fmt.Printf("%s = %s = %s\n", c.String(), f.String(), k.String())
}