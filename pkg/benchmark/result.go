package benchmark

import (
	"slices"

	"golang.org/x/perf/benchunit"
)

type Series struct {
	Name   string
	Points []Result
}

type Result struct {
	Name       string
	CommitHash string
	Metrics    []Metric
}

type Metric struct {
	Unit  string
	Value float64
}

func (m Metric) UnitClass() benchunit.Class {
	return benchunit.ClassOf(m.Unit)
}

func Merge(collected []Series, incoming []Series) []Series {
	var merged []Series

	for _, s := range incoming {
		idx := slices.IndexFunc(collected, func(v Series) bool {
			return v.Name == s.Name
		})

		if idx == -1 {
			merged = append(merged, s)
			continue
		}

		previous := collected[idx]
		previous.Points = append(previous.Points, s.Points...)

		merged = append(merged, previous)
	}

	return merged
}
