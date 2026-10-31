package conversions

// MToF converts a Meter length to Foot.
func MToF(m Meter) Foot { return Foot(m*3.2808) }
// FToM converts a Foot length to Meter.
func FToM(f Foot) Meter { return Meter(f*0.3048) }