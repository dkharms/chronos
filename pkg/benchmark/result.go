package benchmark

import (
	"slices"

	"golang.org/x/perf/benchunit"
)

type Series struct {
	Name   string   `json:"name"`
	Points []Result `json:"measurements"`
}

type Result struct {
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

func Merge(collected, incoming []Series) []Series {
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
