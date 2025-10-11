package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge_SameBenchmarkSameCommit(t *testing.T) {
	previous := Series{
		Name: "BenchmarkX",
		Measurements: []Measurement{{
			CommitHash: "sha1:something",
			Metrics: []Metric{{
				Unit:  "ns/op",
				Value: 1024,
			}},
		}},
	}

	// Do not change anything basically.
	current := deepCopySeries(previous)

	expected := []Series{{
		Name: "BenchmarkX",
		Measurements: []Measurement{
			{
				CommitHash: "sha1:something",
				Metrics: []Metric{{
					Unit:  "ns/op",
					Value: 1024,
				}},
			},
		},
	}}

	assert.Equal(t, expected, Merge(
		[]Series{previous},
		[]Series{current},
	))
}

func TestMerge_SameBenchmarkDifferentCommit(t *testing.T) {
	previous := Series{
		Name: "BenchmarkX",
		Measurements: []Measurement{{
			CommitHash: "sha1:something",
			Metrics: []Metric{{
				Unit:  "ns/op",
				Value: 1024,
			}},
		}},
	}

	current := deepCopySeries(previous)
	current.Measurements[0].CommitHash = "sha2:something"

	expected := []Series{{
		Name: "BenchmarkX",
		Measurements: []Measurement{
			{
				CommitHash: "sha1:something",
				Metrics: []Metric{{
					Unit:  "ns/op",
					Value: 1024,
				}},
			},
			{
				CommitHash: "sha2:something",
				Metrics: []Metric{{
					Unit:  "ns/op",
					Value: 1024,
				}},
			},
		},
	}}

	assert.Equal(t, expected, Merge(
		[]Series{previous},
		[]Series{current},
	))
}

func deepCopySeries(in Series) Series {
	return Series{
		Name:         in.Name,
		Measurements: deepCopyMeasurements(in.Measurements),
	}
}

func deepCopyMeasurements(in []Measurement) []Measurement {
	var out []Measurement
	for _, m := range in {
		out = append(out, Measurement{
			Name:       m.Name,
			CommitHash: m.CommitHash,
			Metrics:    deepCopyMetrics(m.Metrics),
		})

	}
	return out
}

func deepCopyMetrics(in []Metric) []Metric {
	var out []Metric
	for _, m := range in {
		out = append(out, Metric{
			Unit:  m.Unit,
			Value: m.Value,
		})
	}
	return out
}
