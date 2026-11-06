package main

import (
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
)

func main() {
	const width, height = 16, 16

	xmin := big.NewRat(-2, 1)
	xmax := big.NewRat(2, 1)
	ymin := big.NewRat(-2, 1)
	ymax := big.NewRat(2, 1)
	fWidth := big.NewRat(width, 1)
	fHeight := big.NewRat(height, 1)

	xRange := new(big.Rat).Sub(xmax, xmin)
	yRange := new(big.Rat).Sub(ymax, ymin)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	pxRat := new(big.Rat)
	pyRat := new(big.Rat)
	xBase := new(big.Rat)
	yBase := new(big.Rat)

	superX, superY := 2, 2

	for py := 0; py < height; py++ {
		pyRat.SetInt64(int64(py))
		// yBase = (py / height) * (ymax - ymin) + ymin
		yBase.Quo(pyRat, fHeight).Mul(yBase, yRange).Add(yBase, ymin)

		for px := 0; px < width; px++ {
			pxRat.SetInt64(int64(px))
			// xBase = (px / width) * (xmax - xmin) + xmin
			xBase.Quo(pxRat, fWidth).Mul(xBase, xRange).Add(xBase, xmin)

			// Supersampling
			samples := []color.Color{}
			for i := 0; i < superX; i++ {
				for j := 0; j < superY; j++ {
					iScaled := big.NewRat(int64(i), int64(superX))
					iScaled.Quo(iScaled, fWidth).Mul(iScaled, xRange)

					jScaled := big.NewRat(int64(j), int64(superY))
					jScaled.Quo(jScaled, fHeight).Mul(jScaled, yRange)

					// z = (x + i_scaled) + (y + j_scaled)i
					z := NewBigComplexUnbounded()
					z.Real.Add(xBase, iScaled)
					z.Imag.Add(yBase, jScaled)

					// Call your high-precision fractal function
					samples = append(samples, newtonFractalBigComplexUnbounded(z))
				}
			}
			pixel := averageColor(samples)
			img.Set(px, py, pixel)
		}
	}
	png.Encode(os.Stdout, img)
}

func (z *BigComplexUnbounded) fBigComplexUnbounded() *BigComplexUnbounded {
	// z^2 = z * z
	z2 := NewBigComplexUnbounded().Mul(z, z)

	// z^4 = z^2 * z^2
	z4 := NewBigComplexUnbounded().Mul(z2, z2)

	// z^4 - 1 (subtract 1 from the real part)
	one := big.NewRat(1, 1)
	newReal := new(big.Rat).Sub(z4.Real, one)

	return &BigComplexUnbounded{
		Real: newReal,
		Imag: z4.Imag,
	}
}

func (z *BigComplexUnbounded) dfBigComplexUnbounded() *BigComplexUnbounded {
	// 1. z^2 = z * z
	z2 := NewBigComplexUnbounded().Mul(z, z)

	// 2. z^3 = z^2 * z
	z3 := NewBigComplexUnbounded().Mul(z2, z)

	// 3. 4 * z^3 (multiply both real and imaginary parts by 4)
	four := big.NewRat(4, 1)

	newReal := new(big.Rat).Mul(z3.Real, four)
	newImag := new(big.Rat).Mul(z3.Imag, four)

	return &BigComplexUnbounded{
		Real: newReal,
		Imag: newImag,
	}
}

func newtonFractalBigComplexUnbounded(z *BigComplexUnbounded) color.Color {
	const iterations = 200
	const contrast = 15

	thresholdFloat := big.NewFloat(0.0001)
	outOfBoundsThreshold := big.NewFloat(4.0)

	// Roots for z^4 - 1 = 0
	root1 := NewFromComplex128Unbounded(1 + 0i)  // 1
	root2 := NewFromComplex128Unbounded(-1 + 0i) // -1
	root3 := NewFromComplex128Unbounded(0 + 1i)  // i
	root4 := NewFromComplex128Unbounded(0 - 1i)  // -i

	for n := uint8(0); n < iterations; n++ {
		// f(z) = z^4 - 1
		fz := z.fBigComplexUnbounded()

		if fz.AbsFloat(53).Cmp(outOfBoundsThreshold) > 0 {
			return color.Black
		}

		// If |f(z)| < threshold, we've converged close to a root
		if fz.AbsFloat(53).Cmp(thresholdFloat) < 0 {
			nContrast := min(255, uint(n)*contrast)
			var r, g, b uint8

			// Check which root z is closest to: z - root < threshold
			if NewBigComplexUnbounded().Sub(z, root1).AbsFloat(53).Cmp(thresholdFloat) < 0 {
				r, g, b = uint8(nContrast), 0, 0
			} else if NewBigComplexUnbounded().Sub(z, root2).AbsFloat(53).Cmp(thresholdFloat) < 0 {
				r, g, b = 0, uint8(nContrast), uint8(nContrast)
			} else if NewBigComplexUnbounded().Sub(z, root3).AbsFloat(53).Cmp(thresholdFloat) < 0 {
				r, g, b = 0, 0, uint8(nContrast)
			} else if NewBigComplexUnbounded().Sub(z, root4).AbsFloat(53).Cmp(thresholdFloat) < 0 {
				r, g, b = uint8(nContrast), uint8(nContrast), 0
			}

			return color.RGBA{r, g, b, 255}
		}

		// Newton-Raphson Step: z = z - f(z)/f'(z)
		// f'(z) = 4 * z^3
		dfz := z.dfBigComplexUnbounded()

		// If the derivative is zero, we hit a critical point trap; exit loop
		if dfz.Real.Sign() == 0 && dfz.Imag.Sign() == 0 {
			return color.Black
		}

		quot := NewBigComplexUnbounded().Div(fz, dfz)
		z = NewBigComplexUnbounded().Sub(z, quot)
	}

	return color.Black
}
