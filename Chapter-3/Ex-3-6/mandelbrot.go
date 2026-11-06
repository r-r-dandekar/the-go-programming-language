// Mandelbrot emits a PNG image of the Mandelbrot fractal
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin

			// Supersampling
			samples := []color.Color{}
			superX := 2
			superY := 2
			for i := 0; i < superX; i++ {
				for j := 0; j < superY; j++ {
					i_scaled := float64(i) / float64(superX) / width * (xmax - xmin)
					j_scaled := float64(j) / float64(superY) / height * (ymax - ymin)
					z := complex(x+i_scaled, y+j_scaled)
					samples = append(samples, mandelbrot(z))
				}
			}
			pixel := averageColor(samples)

			// Image point (px, py) represents complex value z
			img.Set(px, py, pixel)
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			var red uint8 = min(255, n*contrast)
			var green uint8 = min(255, n*contrast)
			var blue uint8 = 255 - red
			return color.RGBA{red, green, blue, 255}
		}
	}
	return color.Black
}

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
