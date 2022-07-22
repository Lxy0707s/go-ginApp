package number

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// FixFloat returns float value with fix bits
func FixFloat(f float64, fix int) float64 {
	pow := math.Pow10(fix)
	i := int64(f * pow)
	return float64(i) / pow
}

// Float64 converts value to float64
func Float64(v interface{}) (float64, error) {
	if v == nil {
		return 0, errors.New("nil value")
	}
	var (
		floatValue float64
		err        error
	)
	isValid := true
	switch value := v.(type) {
	case string:
		if floatValue, err = strconv.ParseFloat(value, 64); err != nil {
			isValid = false
		}
	case float64:
		floatValue = value
	case float32:
		floatValue = float64(value)
	case int:
		floatValue = float64(value)
	case int32:
		floatValue = float64(value)
	case int64:
		floatValue = float64(value)
	case uint:
		floatValue = float64(value)
	case uint32:
		floatValue = float64(value)
	case uint64:
		floatValue = float64(value)
	default:
		isValid = false
	}
	if !isValid {
		return 0, errors.New("wrong value type")
	}
	return floatValue, nil
}

func FriendlyFloat(raw float64) string {
	val := strconv.FormatFloat(raw, 'f', 5, 64)
	if strings.Contains(val, ".") {
		val = strings.TrimRight(val, "0")
		val = strings.TrimRight(val, ".")
	}
	return val
}

func EtoFloat64String(value string) (rv string) {
	var rf float64
	_, err := fmt.Sscanf(value, "%e", &rf)
	if err != nil {
		rv = "0.0"
	} else {
		rv = fmt.Sprintf("%f", rf)
	}
	if strings.Contains(rv, ".") {
		rv = strings.TrimRight(rv, "0")
		rv = strings.TrimRight(rv, ".")
	}
	return rv
}
