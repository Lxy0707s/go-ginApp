package prof

import (
	"math"
)

// fixFloat returns float value with fix bits
func fixFloat(f float64, fix int) float64 {
	pow := math.Pow10(fix)
	i := int64(f * pow)
	return float64(i) / pow
}
