package helper

import "math"

func RoundTwoDigits(x float64) float64 {
	return math.Round(x*100) / 100
}
