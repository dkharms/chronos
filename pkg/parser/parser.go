package parser

import "github.com/dkharms/chronos/pkg/benchmark"

type Parser interface {
	// Parse returns parsed benchmark output.
	//
	// Order of benchmarks is determined by name (ascending).
	// Also all metrics must be sorted by unit name (ascending).
	Parse() []benchmark.Measurement
}
