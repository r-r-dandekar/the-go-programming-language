// using_complex64 uses the complex64 type to generate a fractal based on the number of iterations needed to find a solution to the function z^4 - 1 = 0 for different complex starting pointsmin(255, n*contrast)
package main

import (
	"image"
	"testing"

	"image/color"
	"math/cmplx"
)

func BenchmarkComplex64(b *testing.B) {
	const (
		xmin, ymin, xmax, ymax = -0.5, -0.5, +0.5, +0.5
		width, height          = 256, 256
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < b.N; i++ {
		for py := 0; py < height; py++ {
			y := float32(py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {
				x := float32(px)/width*(xmax-xmin) + xmin

				// Supersampling
				samples := []color.Color{}
				superX := 2
				superY := 2
				for i := 0; i < superX; i++ {
					for j := 0; j < superY; j++ {
						i_scaled := float32(i) / float32(superX) / width * (xmax - xmin)
						j_scaled := float32(j) / float32(superY) / height * (ymax - ymin)
						z := complex(x+i_scaled, y+j_scaled)
						samples = append(samples, newtonFractalComplex64(z))
					}
				}
				pixel := averageColor(samples)

				// Image point (px, py) represents complex value z
				img.Set(px, py, pixel)
			}
		}
	}
}

func fComplex64(z complex64) complex64 {
	return z*z*z*z - 1
}
func dfComplex64(z complex64) complex64 {
	return 4 * z * z * z
}

func newtonFractalComplex64(z complex64) color.Color {
	const iterations = 200
	const contrast = 15
	const threshold = 0.0001

	for n := uint8(0); n < iterations; n++ {
		if cmplx.Abs(complex128(fComplex64(z))) < threshold {
			var r, g, b uint8
			if cmplx.Abs(complex128(z-1)) < threshold {
				r, g, b = min(255, n*contrast), 0, 0
			} else if cmplx.Abs(complex128(z+1)) < threshold {
				r, g, b = 0, min(255, n*contrast), min(255, n*contrast)
			} else if cmplx.Abs(complex128(z-1i)) < threshold {
				r, g, b = 0, 0, min(255, n*contrast)
			} else if cmplx.Abs(complex128(z+1i)) < threshold {
				r, g, b = min(255, n*contrast), min(255, n*contrast), 0
			}
			return color.RGBA{r, g, b, 255}
		}
		z = z - fComplex64(z)/dfComplex64(z) // Newton-Raphson Approximation Method
	}
	return color.Black
}
