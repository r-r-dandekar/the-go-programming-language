package main

import (
	"math/big"
)

// BigComplexUnbounded represents an exact, arbitrary-precision complex number using big.Rat.
type BigComplexUnbounded struct {
	Real, Imag *big.Rat
}

// NewBigComplexUnbounded creates a new BigComplexUnbounded initialized to 0/1 + 0/1 i.
func NewBigComplexUnbounded() *BigComplexUnbounded {
	return &BigComplexUnbounded{
		Real: new(big.Rat),
		Imag: new(big.Rat),
	}
}

// NewFromComplex128Unbounded converts a complex128 into a BigComplexUnbounded.
func NewFromComplex128Unbounded(c complex128) *BigComplexUnbounded {
	z := NewBigComplexUnbounded()
	z.Real.SetFloat64(real(c))
	z.Imag.SetFloat64(imag(c))
	return z
}

// Addition: z = x + y
func (z *BigComplexUnbounded) Add(x, y *BigComplexUnbounded) *BigComplexUnbounded {
	if z.Real == nil {
		z.Real = new(big.Rat)
	}
	if z.Imag == nil {
		z.Imag = new(big.Rat)
	}

	z.Real.Add(x.Real, y.Real)
	z.Imag.Add(x.Imag, y.Imag)
	return z
}

// Subtraction: z = x - y
func (z *BigComplexUnbounded) Sub(x, y *BigComplexUnbounded) *BigComplexUnbounded {
	if z.Real == nil {
		z.Real = new(big.Rat)
	}
	if z.Imag == nil {
		z.Imag = new(big.Rat)
	}

	z.Real.Sub(x.Real, y.Real)
	z.Imag.Sub(x.Imag, y.Imag)
	return z
}

// Multiplication: z = x * y
// Formula: (ac - bd) + i(ad + bc)
func (z *BigComplexUnbounded) Mul(x, y *BigComplexUnbounded) *BigComplexUnbounded {
	ac := new(big.Rat).Mul(x.Real, y.Real)
	bd := new(big.Rat).Mul(x.Imag, y.Imag)
	ad := new(big.Rat).Mul(x.Real, y.Imag)
	bc := new(big.Rat).Mul(x.Imag, y.Real)

	realPart := new(big.Rat).Sub(ac, bd)
	imagPart := new(big.Rat).Add(ad, bc)

	if z.Real == nil {
		z.Real = new(big.Rat)
	}
	if z.Imag == nil {
		z.Imag = new(big.Rat)
	}

	z.Real.Set(realPart)
	z.Imag.Set(imagPart)
	return z
}

// Division: z = x / y
// Formula: [(ac + bd) + i(bc - ad)] / (c^2 + d^2)
func (z *BigComplexUnbounded) Div(x, y *BigComplexUnbounded) *BigComplexUnbounded {
	ac := new(big.Rat).Mul(x.Real, y.Real)
	bd := new(big.Rat).Mul(x.Imag, y.Imag)
	bc := new(big.Rat).Mul(x.Imag, y.Real)
	ad := new(big.Rat).Mul(x.Real, y.Imag)

	cc := new(big.Rat).Mul(y.Real, y.Real)
	dd := new(big.Rat).Mul(y.Imag, y.Imag)
	den := new(big.Rat).Add(cc, dd)

	if z.Real == nil {
		z.Real = new(big.Rat)
	}
	if z.Imag == nil {
		z.Imag = new(big.Rat)
	}

	// Safety check for division by zero denominator (0 + 0i)
	if den.Num().Sign() == 0 {
		z.Real.SetInt64(0)
		z.Imag.SetInt64(0)
		return z
	}

	realNum := new(big.Rat).Add(ac, bd)
	imagNum := new(big.Rat).Sub(bc, ad)

	// Using invDen = 1 / den to avoid calling .Inv() directly on den
	invDen := new(big.Rat).Inv(den)

	z.Real.Mul(realNum, invDen)
	z.Imag.Mul(imagNum, invDen)
	return z
}

// Abs2 returns the squared magnitude |z|^2 = real^2 + imag^2 as an exact *big.Rat.
// Avoids square root non-rational conversions.
func (z *BigComplexUnbounded) Abs2() *big.Rat {
	r2 := new(big.Rat).Mul(z.Real, z.Real)
	i2 := new(big.Rat).Mul(z.Imag, z.Imag)
	return new(big.Rat).Add(r2, i2)
}

// AbsFloat returns the magnitude |z| as a *big.Float converted to a specific bit precision.
func (z *BigComplexUnbounded) AbsFloat(prec uint) *big.Float {
	r2 := new(big.Float).SetPrec(prec).SetRat(z.Real)
	i2 := new(big.Float).SetPrec(prec).SetRat(z.Imag)

	r2.Mul(r2, r2)
	i2.Mul(i2, i2)

	sum := new(big.Float).SetPrec(prec).Add(r2, i2)
	return new(big.Float).SetPrec(prec).Sqrt(sum)
}
