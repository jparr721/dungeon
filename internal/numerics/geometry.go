package numerics

import "math"

func RadToDegree(rad float64) float64 {
	return rad * 180 / math.Pi
}

func DegreeToRad(deg float64) float64 {
	return deg * math.Pi / 180
}
