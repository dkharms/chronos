package benchmark

import (
	"math"
	"slices"
)

type CalculatedDiff struct {
	Name       string
	MetricDiff []MetricDiff
}

type MetricDiff struct {
	Unit string

	PreviousCommit string
	CurrentCommit  string

	PreviousValue float64
	CurrentValue  float64
}

func (m MetricDiff) Ratio() float64 {
	if math.IsNaN(m.PreviousValue) || m.PreviousValue == 0 || math.IsNaN(m.CurrentValue) {
		return math.NaN()
	}
	return m.CurrentValue / m.PreviousValue
}

func Diff(previous, current []Series) []CalculatedDiff {
	var calculated []CalculatedDiff

	for _, s := range current {
		idx := slices.IndexFunc(previous, func(v Series) bool {
			return v.Name == s.Name
		})

		var md []MetricDiff
		if idx == -1 {
			md = metricDiff(
				Measurement{},
				s.Measurements[len(s.Measurements)-1],
			)
		} else {
			md = metricDiff(
				previous[idx].Measurements[len(previous[idx].Measurements)-1],
				s.Measurements[len(s.Measurements)-1],
			)
		}

		calculated = append(calculated, CalculatedDiff{
			Name:       s.Name,
			MetricDiff: md,
		})
	}

	return calculated
}

func metricDiff(previous, current Measurement) []MetricDiff {
	var (
		prevCommit = "------"
		prevValue  = math.NaN()
		diff       []MetricDiff
	)

	for _, cm := range current.Metrics {
		idx := slices.IndexFunc(previous.Metrics, func(v Metric) bool {
			return v.Unit == cm.Unit
		})

		if idx >= 0 {
			prevCommit = previous.CommitHash
			prevValue = mean(previous.Metrics[idx].Values)
		}

		diff = append(diff, MetricDiff{
			Unit: cm.Unit,

			PreviousCommit: prevCommit,
			CurrentCommit:  current.CommitHash,

			PreviousValue: prevValue,
			CurrentValue:  mean(cm.Values),
		})
	}

	return diff
}

func mean(s []float64) float64 {
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}
