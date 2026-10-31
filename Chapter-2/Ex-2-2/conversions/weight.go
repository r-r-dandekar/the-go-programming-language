package conversions

// KgToLb converts a Kilogram length to Pound.
func KgToLb(k Kilogram) Pound { return Pound(k*2.2046) }
// LbToKg converts a Pound length to Kilogram.
func LbToKg(p Pound) Kilogram { return Kilogram(p*0.3048) }