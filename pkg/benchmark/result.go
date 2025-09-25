package benchmark

import "golang.org/x/perf/benchunit"

type Result struct {
	Name    string
	Metrics []Metric
}

type Metric struct {
	Unit  string
	Value float64
}

func (m Metric) UnitClass() benchunit.Class {
	return benchunit.ClassOf(m.Unit)
}
