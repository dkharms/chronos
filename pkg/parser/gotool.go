package parser

import (
	"cmp"
	"io"
	"maps"
	"slices"
	"strings"

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
	bname := make(map[string][]benchmark.Measurement)
	br := benchfmt.NewReader(p.r, "benchmarks")

	for br.Scan() {
		v, ok := br.Result().(*benchfmt.Result)
		if !ok {
			continue
		}

		name := v.Name.String()
		bname[name] = append(
			bname[name],
			convert(*v.Clone()),
		)
	}

	for name, bresults := range bname {
		// metrics maps unit to values
		metrics := make(map[string][]float64)

		for _, r := range bresults {
			for _, m := range r.Metrics {
				metrics[m.Unit] = append(
					metrics[m.Unit],
					m.Value,
				)
			}
		}

		var smetrics []benchmark.Metric
		for unit, values := range maps.All(metrics) {
			smetrics = append(smetrics, benchmark.Metric{
				Unit:   unit,
				Values: values,
			})
		}

		slices.SortFunc(smetrics, func(x, y benchmark.Metric) int {
			return strings.Compare(x.Unit, y.Unit)
		})

		results = append(results, benchmark.Measurement{
			Name:    name,
			Metrics: smetrics,
		})
	}

	slices.SortFunc(results, func(x, y benchmark.Measurement) int {
		return strings.Compare(x.Name, y.Name)
	})

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
