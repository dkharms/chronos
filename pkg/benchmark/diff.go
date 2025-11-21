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

	PreviousValues []float64
	CurrentValues  []float64
}

func (m MetricDiff) PreviousValue() float64 {
	return m.reduce(m.PreviousValues)
}
func (m MetricDiff) CurrentValue() float64 {
	return m.reduce(m.CurrentValues)
}

func (m MetricDiff) reduce(s []float64) float64 {
	if d, ok := GetMetricDescriptor(m.Unit); ok {
		return d.Reduce(s)
	}
	return mean(s)
}

func (m MetricDiff) Ratio() float64 {
	prev, curr := m.PreviousValue(), m.CurrentValue()

	nan := prev == 0 ||
		math.IsNaN(prev) ||
		math.IsNaN(curr)

	if nan {
		return math.NaN()
	}

	return curr / prev
}

func (m MetricDiff) Emoji() string {
	if d, ok := GetMetricDescriptor(m.Unit); ok {
		prev, curr := m.PreviousValue(), m.CurrentValue()
		switch d.Comparator(prev, curr) {
		case CompareVerdictBetter:
			return "🟢"
		case CompareVerdictWorse:
			return "🔴"
		case CompareVerdictSame:
			return "🟰"
		default:
			panic("unknown compare verdict")
		}
	}
	return "❓"
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
		prevCommit = "------" // Exactly 6 symbols
		prevValues []float64
		diff       []MetricDiff
	)

	for _, cm := range current.Metrics {
		idx := slices.IndexFunc(previous.Metrics, func(v Metric) bool {
			return v.Unit == cm.Unit
		})

		if idx >= 0 {
			prevCommit = previous.CommitHash
			prevValues = previous.Metrics[idx].Values
		}

		diff = append(diff, MetricDiff{
			Unit: cm.Unit,

			PreviousCommit: prevCommit,
			CurrentCommit:  current.CommitHash,

			PreviousValues: prevValues,
			CurrentValues:  cm.Values,
		})
	}

	return diff
}
