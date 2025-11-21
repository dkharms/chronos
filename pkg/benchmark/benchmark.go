package benchmark

import (
	"math"
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
	Unit   string    `json:"unit"`
	Values []float64 `json:"values"`
}

type CompareVerdict int

const (
	CompareVerdictSame CompareVerdict = iota
	CompareVerdictBetter
	CompareVerdictWorse
)

type MetricDescriptor struct {
	Unit       string
	Reduce     func([]float64) float64
	Comparator func(float64, float64) CompareVerdict
}

func NewMetricDescriptor(
	unit string, better string,
	threshold float64, reduce string,
) MetricDescriptor {
	m := MetricDescriptor{
		Unit:   unit,
		Reduce: reductionFn(reduce),
	}

	switch better {
	case "lower":
		m.Comparator = comparator(
			func(previous, current float64) CompareVerdict {
				if current < previous {
					return CompareVerdictBetter
				}
				return CompareVerdictWorse
			},
			threshold,
		)
	case "higher":
		m.Comparator = comparator(
			func(previous, current float64) CompareVerdict {
				if current > previous {
					return CompareVerdictBetter
				}
				return CompareVerdictWorse
			},
			threshold,
		)
	}

	return m
}

func comparator(
	cmp func(float64, float64) CompareVerdict,
	threshold float64,
) func(float64, float64) CompareVerdict {
	return func(previous, current float64) CompareVerdict {
		within := previous == current ||
			withinThreshold(previous, current, threshold)

		if within {
			return CompareVerdictSame
		}

		return cmp(previous, current)
	}
}

func withinThreshold(previous, current, threshold float64) bool {
	if previous == 0 || current == 0 {
		return false
	}
	return math.Abs(1-(previous/current)) <= threshold
}

var (
	metricDescriptorRegistry = MetricDescriptorRegistry{}
)

type MetricDescriptorRegistry map[string]MetricDescriptor

func AddMetricDescriptor(m MetricDescriptor) {
	metricDescriptorRegistry[m.Unit] = m
}

func GetMetricDescriptor(unit string) (MetricDescriptor, bool) {
	m, ok := metricDescriptorRegistry[unit]
	return m, ok
}
