package main

import (
	"math/big"
)

// BigComplex represents a high-precision complex number using big.Float
type BigComplex struct {
	Real, Imag *big.Float
}

// NewBigComplex creates a new BigComplex initialized to zero with a given precision (in bits)
func NewBigComplex(prec uint) *BigComplex {
	return &BigComplex{
		Real: new(big.Float).SetPrec(prec),
		Imag: new(big.Float).SetPrec(prec),
	}
}

// 3. Conversion from complex128 to BigComplex
func NewFromComplex128(c complex128, prec uint) *BigComplex {
	return &BigComplex{
		Real: new(big.Float).SetPrec(prec).SetFloat64(real(c)),
		Imag: new(big.Float).SetPrec(prec).SetFloat64(imag(c)),
	}
}

// 1. Addition: z = x + y
func (z *BigComplex) Add(x, y *BigComplex) *BigComplex {
	if z.Real == nil {
		z.Real = new(big.Float)
	}
	if z.Imag == nil {
		z.Imag = new(big.Float)
	}
	z.Real.Add(x.Real, y.Real)
	z.Imag.Add(x.Imag, y.Imag)
	return z
}

// 1. Subtraction: z = x - y
func (z *BigComplex) Sub(x, y *BigComplex) *BigComplex {
	if z.Real == nil {
		z.Real = new(big.Float)
	}
	if z.Imag == nil {
		z.Imag = new(big.Float)
	}
	z.Real.Sub(x.Real, y.Real)
	z.Imag.Sub(x.Imag, y.Imag)
	return z
}

// 2. Multiplication: z = x * y
// Formula: (ac - bd) + i(ad + bc)
func (z *BigComplex) Mul(x, y *BigComplex) *BigComplex {
	ac := new(big.Float).Mul(x.Real, y.Real)
	bd := new(big.Float).Mul(x.Imag, y.Imag)
	ad := new(big.Float).Mul(x.Real, y.Imag)
	bc := new(big.Float).Mul(x.Imag, y.Real)

	realPart := new(big.Float).Sub(ac, bd)
	imagPart := new(big.Float).Add(ad, bc)

	if z.Real == nil {
		z.Real = new(big.Float)
	}
	if z.Imag == nil {
		z.Imag = new(big.Float)
	}
	z.Real.Set(realPart)
	z.Imag.Set(imagPart)
	return z
}

// 2. Division: z = x / y
// Formula: [(ac + bd) + i(bc - ad)] / (c^2 + d^2)
func (z *BigComplex) Div(x, y *BigComplex) *BigComplex {
	ac := new(big.Float).Mul(x.Real, y.Real)
	bd := new(big.Float).Mul(x.Imag, y.Imag)
	bc := new(big.Float).Mul(x.Imag, y.Real)
	ad := new(big.Float).Mul(x.Real, y.Imag)

	cc := new(big.Float).Mul(y.Real, y.Real)
	dd := new(big.Float).Mul(y.Imag, y.Imag)
	den := new(big.Float).Add(cc, dd)

	realNum := new(big.Float).Add(ac, bd)
	imagNum := new(big.Float).Sub(bc, ad)

	realPart := new(big.Float).Quo(realNum, den)
	imagPart := new(big.Float).Quo(imagNum, den)

	if z.Real == nil {
		z.Real = new(big.Float)
	}
	if z.Imag == nil {
		z.Imag = new(big.Float)
	}
	z.Real.Set(realPart)
	z.Imag.Set(imagPart)
	return z
}

// 4. Abs returns the magnitude as a *big.Float: sqrt(real^2 + imag^2)
func (z *BigComplex) Abs() *big.Float {
	r2 := new(big.Float).Mul(z.Real, z.Real)
	i2 := new(big.Float).Mul(z.Imag, z.Imag)
	sum := new(big.Float).Add(r2, i2)
	return new(big.Float).Sqrt(sum)
}
