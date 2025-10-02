package benchmark

import (
	"slices"
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
