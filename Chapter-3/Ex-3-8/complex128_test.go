// using_complex128 uses the complex128 type to generate a fractal based on the number of iterations needed to find a solution to the function z^4 - 1 = 0 for different complex starting pointsmin(255, n*contrast)
package main

import (
	"image"
	"image/color"
	"math/cmplx"
	"testing"
)

func BenchmarkComplex128(b *testing.B) {
	const (
		xmin, ymin, xmax, ymax = -0.5, -0.5, +0.5, +0.5
		width, height          = 256, 256
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < b.N; i++ {
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
						samples = append(samples, newtonFractalComplex128(z))
					}
				}
				pixel := averageColor(samples)

				// Image point (px, py) represents complex value z
				img.Set(px, py, pixel)
			}
		}
	}
}

func fComplex128(z complex128) complex128 {
	return z*z*z*z - 1
}
func dfComplex128(z complex128) complex128 {
	return 4 * z * z * z
}

func newtonFractalComplex128(z complex128) color.Color {
	const iterations = 200
	const contrast = 15
	const threshold = 0.0001

	for n := uint8(0); n < iterations; n++ {
		if cmplx.Abs(fComplex128(z)) < threshold {
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
		z = z - fComplex128(z)/dfComplex128(z) // Newton-Raphson Approximation Method
	}
	return color.Black
}
