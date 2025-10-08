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

		// Found a benchmark that was not collected previously.
		if idx == -1 {
			merged = append(merged, s)
			continue
		}

		var unique []Measurement
		previous := previous[idx]

		// Do not add new benchmark measurement with commit hash `h`
		// if it's already presented in collected results from previous runs.
		for _, m := range s.Measurements {
			contains := slices.ContainsFunc(
				previous.Measurements, func(pm Measurement) bool {
					return pm.CommitHash == m.CommitHash
				},
			)

			if contains {
				continue
			}

			unique = append(unique, m)
		}

		previous.Measurements = append(previous.Measurements, unique...)
		merged = append(merged, previous)
	}

	return merged
}
