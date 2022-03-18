package signals

import "math"

func BuildSin(Vmax float64, frequency float64, phase float64) []float64 {
	sin := func(t float64) float64 {
		return Vmax * math.Sin(frequency*2*math.Pi*t+phase)
	}

	arr := arange2(math.Pow10(-4), 2*math.Pow10(-2), math.Pow10(-5))

	points := make([]float64, len(arr))

	for i, element := range arr {
		points[i] = sin(element)
	}

	return points
}

func arange2(start, stop, step float64) []float64 {
	N := int(math.Ceil((stop - start) / step))
	rnge := make([]float64, N)
	for x := range rnge {
		rnge[x] = start + step*float64(x)
	}
	return rnge
}
