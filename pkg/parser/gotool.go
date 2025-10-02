package parser

import (
	"cmp"
	_ "embed"
	"io"

	"golang.org/x/perf/benchfmt"

	"github.com/dkharms/chronos/pkg/benchmark"
)

type goparser struct {
	r io.Reader
}

func NewGoParser(r io.Reader) *goparser {
	return &goparser{r: r}
}

func (p *goparser) Parse() (results []benchmark.Measurement) {
	br := benchfmt.NewReader(p.r, "benchmarks")

	for br.Scan() {
		if v, ok := br.Result().(*benchfmt.Result); ok {
			results = append(results, convert(*v.Clone()))
		}
	}

	return
}

func convert(b benchfmt.Result) benchmark.Measurement {
	r := benchmark.Measurement{
		Name: b.Name.String(),
		Metrics: []benchmark.Metric{{
			Unit:  "iterations",
			Value: float64(b.Iters),
		}},
	}

	for _, v := range b.Values {
		r.Metrics = append(r.Metrics, benchmark.Metric{
			Unit:  cmp.Or(v.OrigUnit, v.Unit),
			Value: cmp.Or(v.OrigValue, v.Value),
		})
	}

	return r
}
