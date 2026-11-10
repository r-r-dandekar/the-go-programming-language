package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {

	bits := flag.Int("bits", 256, "Number of bits in the SHA hash")
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		os.Exit(1)
	}

	if *bits == 256 {
		fmt.Printf("%x\n", sha256.Sum256([]byte(input)))
	} else if *bits == 384 {
		fmt.Printf("%x\n", sha512.Sum384([]byte(input)))
	} else if *bits == 512 {
		fmt.Printf("%x\n", sha512.Sum512([]byte(input)))
	} else {
		fmt.Fprintln(os.Stderr, "Error: bits must be either 256, 384 or 512!")
	}
}
