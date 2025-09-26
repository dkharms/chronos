package parser

import (
	_ "embed"
	"io"

	"golang.org/x/perf/benchfmt"

	"github.com/dkharms/chronos/pkg/benchmark"
)

var (
	//go:embed testdata/gotool.txt
	GotoolOutput string
)

type goparser struct {
	r io.Reader
}

func NewGoParser(r io.Reader) *goparser {
	return &goparser{r: r}
}

func (p *goparser) Parse() (results []benchmark.Result) {
	br := benchfmt.NewReader(p.r, "benchmarks")

	for br.Scan() {
		switch v := br.Result().(type) {
		case *benchfmt.Result:
			results = append(results, convert(*v.Clone()))
		}
	}

	return
}

func convert(b benchfmt.Result) benchmark.Result {
	r := benchmark.Result{
		Name: b.Name.String(),
		Metrics: []benchmark.Metric{{
			Unit:  "iterations",
			Value: float64(b.Iters),
		}},
	}

	for _, v := range b.Values {
		r.Metrics = append(r.Metrics, benchmark.Metric{
			Unit:  v.Unit,
			Value: v.Value,
		})
	}

	return r
}
