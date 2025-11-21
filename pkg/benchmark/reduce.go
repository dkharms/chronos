package benchmark

import (
	"math"
	"slices"
)

func reductionFn(function string) func([]float64) float64 {
	switch function {
	case "mean":
		return mean
	case "min":
		return minimal
	case "median":
		return median
	case "max":
		return maximal
	}
	panic("unknown reduce function")
}

func mean(s []float64) float64 {
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}

func minimal(s []float64) float64 {
	return sortedth(s, 0)
}

func maximal(s []float64) float64 {
	return sortedth(s, len(s)-1)
}

func median(s []float64) float64 {
	return sortedth(s, len(s)/2)
}

func sortedth(s []float64, idx int) float64 {
	if len(s) == 0 {
		return math.NaN()
	}

	cp := slices.Clone(s)
	slices.Sort(cp)

	return cp[idx]
}
