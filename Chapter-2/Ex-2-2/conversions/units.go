// Package conversions performs conversions between various units
package conversions

import "fmt"

// Temperature
type Celsius float64
type Fahrenheit float64
type Kelvin float64

// Length and Weight
// MKS 
type Meter float64
type Kilogram float64
// Imperial
type Foot float64
type Pound float64

const (
	AbsoluteZeroC 	Celsius = -273.15
	FreezingC		Celsius = 0
	BoilingC		Celsius = 100
)

func (c Celsius) String() string { return fmt.Sprintf("%.2f°C", c) }
func (f Fahrenheit) String() string { return fmt.Sprintf("%.2f°F", f) }
func (k Kelvin) String() string { return fmt.Sprintf("%.2fK", k) }

func (m Meter) String() string { return fmt.Sprintf("%.2fm", m) }
func (k Kilogram) String() string { return fmt.Sprintf("%.2fkg", k) }
func (f Foot) String() string { return fmt.Sprintf("%.2fft", f) }
func (p Pound) String() string { return fmt.Sprintf("%.2flb", p) }