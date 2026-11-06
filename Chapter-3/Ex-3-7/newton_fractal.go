// newton_fractal generates a fractal based on the number of iterations needed to find a solution to the function z^4 - 1 = 0 for different complex starting pointsmin(255, n*contrast)
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
					samples = append(samples, newtonFractal(z))
				}
			}
			pixel := averageColor(samples)

			// Image point (px, py) represents complex value z
			img.Set(px, py, pixel)
		}
	}
	png.Encode(os.Stdout, img)
}

func f(z complex128) complex128 {
	return z*z*z*z - 1
}
func df(z complex128) complex128 {
	return 4 * z * z * z
}

func newtonFractal(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	const threshold = 0.0001

	for n := uint8(0); n < iterations; n++ {
		if cmplx.Abs(f(z)) < threshold {
			var r, g, b uint8
			if cmplx.Abs(z-1) < threshold {
				r, g, b = min(255, n*contrast), 0, 0
			} else if cmplx.Abs(z+1) < threshold {
				r, g, b = 0, min(255, n*contrast), min(255, n*contrast)
			} else if cmplx.Abs(z-1i) < threshold {
				r, g, b = 0, 0, min(255, n*contrast)
			} else if cmplx.Abs(z+1i) < threshold {
				r, g, b = min(255, n*contrast), min(255, n*contrast), 0
			}
			return color.RGBA{r, g, b, 255}
		}
		z = z - f(z)/df(z) // Newton-Raphson Approximation Method
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
