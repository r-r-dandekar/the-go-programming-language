// newton_fractal generates a fractal based on the number of iterations needed to find a solution to the function z^4 - 1 = 0 for different complex starting pointsmin(255, n*contrast)
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/big"
	"os"
)

func mainn() {
	const width, height = 256, 256
	prec := uint(256) // Bit precision for big.Float

	// 1. Constants converted to big.Float variables
	xmin := big.NewFloat(-2).SetPrec(prec)
	xmax := big.NewFloat(2).SetPrec(prec)
	ymin := big.NewFloat(-2).SetPrec(prec)
	ymax := big.NewFloat(2).SetPrec(prec)
	fWidth := big.NewFloat(width).SetPrec(prec)
	fHeight := big.NewFloat(height).SetPrec(prec)

	xRange := new(big.Float).Sub(xmax, xmin)
	yRange := new(big.Float).Sub(ymax, ymin)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Reusable scratch variables to reduce allocation overhead inside loops
	pxF := new(big.Float).SetPrec(prec)
	pyF := new(big.Float).SetPrec(prec)
	xBase := new(big.Float).SetPrec(prec)
	yBase := new(big.Float).SetPrec(prec)

	superX, superY := 2, 2

	for py := 0; py < height; py++ {
		pyF.SetFloat64(float64(py))
		// yBase = (py / height) * (ymax - ymin) + ymin
		yBase.Quo(pyF, fHeight).Mul(yBase, yRange).Add(yBase, ymin)

		for px := 0; px < width; px++ {
			pxF.SetFloat64(float64(px))
			// xBase = (px / width) * (xmax - xmin) + xmin
			xBase.Quo(pxF, fWidth).Mul(xBase, xRange).Add(xBase, xmin)

			// Supersampling
			samples := []color.Color{}
			for i := 0; i < superX; i++ {
				for j := 0; j < superY; j++ {
					// i_scaled = (i / superX) / width * (xmax - xmin)
					iScaled := new(big.Float).SetPrec(prec).SetFloat64(float64(i) / float64(superX))
					iScaled.Quo(iScaled, fWidth).Mul(iScaled, xRange)

					// j_scaled = (j / superY) / height * (ymax - ymin)
					jScaled := new(big.Float).SetPrec(prec).SetFloat64(float64(j) / float64(superY))
					jScaled.Quo(jScaled, fHeight).Mul(jScaled, yRange)

					// z = (x + i_scaled) + (y + j_scaled)i
					z := NewBigComplex(prec)
					z.Real.Add(xBase, iScaled)
					z.Imag.Add(yBase, jScaled)

					// Call your high-precision fractal function
					samples = append(samples, newtonFractalBigComplex(z))
				}
			}
			pixel := averageColor(samples)
			img.Set(px, py, pixel)
		}
	}
	png.Encode(os.Stdout, img)
}

func (z *BigComplex) fBigComplex() *BigComplex {
	prec := z.Real.Prec()

	// z^2 = z * z
	z2 := NewBigComplex(prec).Mul(z, z)

	// z^4 = z^2 * z^2
	z4 := NewBigComplex(prec).Mul(z2, z2)

	// z^4 - 1 (subtract 1 from the real part)
	one := big.NewFloat(1.0).SetPrec(prec)
	newReal := new(big.Float).SetPrec(prec).Sub(z4.Real, one)

	return &BigComplex{
		Real: newReal,
		Imag: z4.Imag,
	}
}
func (z *BigComplex) dfBigComplex() *BigComplex {
	prec := z.Real.Prec()

	// 1. z^2 = z * z
	z2 := NewBigComplex(prec).Mul(z, z)

	// 2. z^3 = z^2 * z
	z3 := NewBigComplex(prec).Mul(z2, z)

	// 3. 4 * z^3 (multiply both real and imaginary parts by 4)
	four := big.NewFloat(4.0).SetPrec(prec)

	newReal := new(big.Float).SetPrec(prec).Mul(z3.Real, four)
	newImag := new(big.Float).SetPrec(prec).Mul(z3.Imag, four)

	return &BigComplex{
		Real: newReal,
		Imag: newImag,
	}
}

func newtonFractalBigComplex(z *BigComplex) color.Color {
	const iterations = 200
	const contrast = 15
	prec := z.Real.Prec()

	threshold := big.NewFloat(0.0001).SetPrec(prec)

	// Roots for z^4 - 1 = 0
	root1 := NewFromComplex128(1+0i, prec)  // 1
	root2 := NewFromComplex128(-1+0i, prec) // -1
	root3 := NewFromComplex128(0+1i, prec)  // i
	root4 := NewFromComplex128(0-1i, prec)  // -i

	for n := uint8(0); n < iterations; n++ {
		// f(z) = z^4 - 1
		fz := z.fBigComplex()

		// If |f(z)| < threshold, we've converged close to a root
		if fz.Abs().Cmp(threshold) < 0 {
			nContrast := min(255, uint(n)*contrast)
			var r, g, b uint8

			// Check which root z is closest to: z - root < threshold
			if NewBigComplex(prec).Sub(z, root1).Abs().Cmp(threshold) < 0 {
				r, g, b = uint8(nContrast), 0, 0
			} else if NewBigComplex(prec).Sub(z, root2).Abs().Cmp(threshold) < 0 {
				r, g, b = 0, uint8(nContrast), uint8(nContrast)
			} else if NewBigComplex(prec).Sub(z, root3).Abs().Cmp(threshold) < 0 {
				r, g, b = 0, 0, uint8(nContrast)
			} else if NewBigComplex(prec).Sub(z, root4).Abs().Cmp(threshold) < 0 {
				r, g, b = uint8(nContrast), uint8(nContrast), 0
			}

			return color.RGBA{r, g, b, 255}
		}

		// Newton-Raphson Step: z = z - f(z)/f'(z)
		// f'(z) = 4 * z^3
		dfz := z.dfBigComplex()

		// If the derivative is zero, we hit a critical point trap; exit loop
		if dfz.Abs().Sign() == 0 {
			return color.Black
		}

		quot := NewBigComplex(prec).Div(fz, dfz)
		z = NewBigComplex(prec).Sub(z, quot)
	}

	return color.Black
}
