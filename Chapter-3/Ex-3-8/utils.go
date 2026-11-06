package main

import "image/color"

// averageColor takes a slice of color.Color and returns an averaged color.RGBA
func averageColor(colors []color.Color) color.RGBA {
	if len(colors) == 0 {
		return color.RGBA{0, 0, 0, 255}
	}

	var sumR, sumG, sumB, sumA uint64

	for _, c := range colors {
		// RGBA() returns values in the range [0, 65535]
		r, g, b, a := c.RGBA()
		sumR += uint64(r)
		sumG += uint64(g)
		sumB += uint64(b)
		sumA += uint64(a)
	}

	count := uint64(len(colors))

	// Average the 16-bit components and shift right by 8
	// to convert them down to standard 8-bit values (0 - 255)
	avgR := uint8((sumR / count) >> 8)
	avgG := uint8((sumG / count) >> 8)
	avgB := uint8((sumB / count) >> 8)
	avgA := uint8((sumA / count) >> 8)

	return color.RGBA{R: avgR, G: avgG, B: avgB, A: avgA}
}
