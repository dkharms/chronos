package parser

import (
	"fmt"
	"io"

	"github.com/dkharms/chronos/pkg/benchmark"
)

type Parser interface {
	// Parse returns parsed benchmark output.
	//
	// Order of benchmarks is determined by name (ascending).
	// Also all metrics must be sorted by unit name (ascending).
	Parse() []benchmark.Measurement
}

func New(tool string, r io.Reader) (Parser, error) {
	if tool == "gotool" {
		return NewGoParser(r), nil
	}
	return nil, fmt.Errorf("unknown language-tool %s", tool)
}
