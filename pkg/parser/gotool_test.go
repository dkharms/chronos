package parser

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/dkharms/chronos/pkg/benchmark"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/gotool.txt
	gotoolOutput string
)

func TestParseGoTool(t *testing.T) {
	actual := NewGoParser(strings.NewReader(gotoolOutput)).Parse()
	require.Equal(t, []benchmark.Measurement{
		{
			Name: "FindSequence_Random/extra-large-16", Metrics: []benchmark.Metric{
				{"B/op", 0, []float64{0}},
				{"MB/s", 0, []float64{17757.45}},
				{"allocs/op", 0, []float64{0}},
				{"iterations", 0, []float64{20184}},
				{"ns/op", 0, []float64{59050}},
			},
		},
		{
			Name: "FindSequence_Random/large-16", Metrics: []benchmark.Metric{
				{"B/op", 0, []float64{0}},
				{"MB/s", 0, []float64{30814.31}},
				{"allocs/op", 0, []float64{0}},
				{"iterations", 0, []float64{2262572}},
				{"ns/op", 0, []float64{531.7}},
			},
		},
		{
			Name: "FindSequence_Random/medium-16", Metrics: []benchmark.Metric{
				{"B/op", 0, []float64{0}},
				{"MB/s", 0, []float64{18339.66}},
				{"allocs/op", 0, []float64{0}},
				{"iterations", 0, []float64{18174138}},
				{"ns/op", 0, []float64{55.84}},
			},
		},
		{
			Name: "FindSequence_Random/small-16", Metrics: []benchmark.Metric{
				{"B/op", 0, []float64{0}},
				{"MB/s", 0, []float64{16286.37}},
				{"allocs/op", 0, []float64{0}},
				{"iterations", 0, []float64{76243033}},
				{"ns/op", 0, []float64{15.72}},
			},
		},
		{
			Name: "FindSequence_Random/tiny-16", Metrics: []benchmark.Metric{
				{"B/op", 0, []float64{0}},
				{"MB/s", 0, []float64{5461.91}},
				{"allocs/op", 0, []float64{0}},
				{"iterations", 0, []float64{102575256}},
				{"ns/op", 0, []float64{11.72}},
			},
		},
	}, actual)
}
