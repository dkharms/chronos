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
	benchByName := make(map[string][]benchmark.Measurement)
	benchReader := benchfmt.NewReader(p.r, "benchmarks")

	for benchReader.Scan() {
		v, ok := benchReader.Result().(*benchfmt.Result)
		if !ok {
			continue
		}

		name := v.Name.String()
		benchByName[name] = append(
			benchByName[name],
			convert(*v.Clone()),
		)
	}

	for benchName, benchResults := range benchByName {
		// metrics maps unit to values
		metrics := make(map[string][]float64)

		for _, benchResult := range benchResults {
			for _, metric := range benchResult.Metrics {
				metrics[metric.Unit] = append(
					metrics[metric.Unit],
					metric.Values...,
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
			Name:    benchName,
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
	}

	for _, v := range b.Values {
		r.Metrics = append(r.Metrics, benchmark.Metric{
			Unit:   cmp.Or(v.OrigUnit, v.Unit),
			Values: []float64{cmp.Or(v.OrigValue, v.Value)},
		})
	}

	return r
}
