package benchmark

import (
	"math"
	"slices"

	"golang.org/x/perf/benchunit"
)

type Series struct {
	Name         string        `json:"name"`
	Measurements []Measurement `json:"measurements"`
}

type Measurement struct {
	Name       string   `json:"-"`
	CommitHash string   `json:"commit_hash"`
	Metrics    []Metric `json:"metrics"`
}

type Metric struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

func (m Metric) UnitClass() benchunit.Class {
	return benchunit.ClassOf(m.Unit)
}

func Merge(previous, current []Series) []Series {
	var merged []Series

	for _, s := range current {
		idx := slices.IndexFunc(previous, func(v Series) bool {
			return v.Name == s.Name
		})

		if idx == -1 {
			merged = append(merged, s)
			continue
		}

		previous := previous[idx]
		previous.Measurements = append(previous.Measurements, s.Measurements...)

		merged = append(merged, previous)
	}

	return merged
}

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
		if idx != -1 {
			md = metricDiff(Measurement{}, s.Measurements[len(s.Measurements)-1])
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
		prevCommit = "-"
		prevValue  = math.NaN()
		diff       []MetricDiff
	)

	for _, cm := range current.Metrics {
		idx := slices.IndexFunc(previous.Metrics, func(v Metric) bool {
			return v.Unit == cm.Unit
		})

		if idx > 0 {
			prevCommit = previous.CommitHash
			prevValue = previous.Metrics[idx].Value
		}

		diff = append(diff, MetricDiff{
			Unit: cm.Unit,

			PreviousCommit: prevCommit,
			CurrentCommit:  current.CommitHash,

			PreviousValue: prevValue,
			CurrentValue:  cm.Value,
		})
	}

	return diff
}
