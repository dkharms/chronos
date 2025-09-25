package parser

import (
	"github.com/dkharms/chronos/pkg/benchmark"
)

type Parser interface {
	Parse() []benchmark.Result
}
