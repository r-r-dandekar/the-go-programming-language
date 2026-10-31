package main

import (
	"bufio"
	"fmt"
	"os"
	"conversions/conversions"
	"strconv"
)

func main() {
	if (len(os.Args) > 1) {
		for _, arg := range os.Args[1:] {
			printConversions(arg)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			printConversions(word)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "convert: reading standard input:", err)
		}
	}
}

func printConversions(s string) {
		val, err := strconv.ParseFloat(s, 64)
		if (err != nil) {
			fmt.Fprintf(os.Stderr, "covnert: %v\n", err)
			os.Exit(1)
		}
		f := conversions.Fahrenheit(val)
		c := conversions.Celsius(val)
		k := conversions.Kelvin(val)
		m := conversions.Meter(val)
		ft := conversions.Foot(val)
		kg := conversions.Kilogram(val)
		lb := conversions.Pound(val)

		fmt.Printf("%s = %s = %s\n", f, conversions.FToC(f), conversions.FToK(f))
		fmt.Printf("%s = %s = %s\n", c, conversions.CToF(c), conversions.CToK(c))
		fmt.Printf("%s = %s = %s\n", k, conversions.KToF(k), conversions.KToC(k))
		fmt.Printf("\n")
		fmt.Printf("%s = %s\n", m, conversions.MToF(m))
		fmt.Printf("%s = %s\n", ft, conversions.FToM(ft))
		fmt.Printf("\n")
		fmt.Printf("%s = %s\n", kg, conversions.KgToLb(kg))
		fmt.Printf("%s = %s\n", lb, conversions.LbToKg(lb))
}