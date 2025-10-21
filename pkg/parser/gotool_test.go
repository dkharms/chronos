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
				{"B/op", []float64{0, 0, 0}},
				{"MB/s", []float64{17757.45, 17757.45, 17757.45}},
				{"allocs/op", []float64{0, 0, 0}},
				{"iterations", []float64{20184, 20184, 20184}},
				{"ns/op", []float64{59050, 59050, 59050}},
			},
		},
		{
			Name: "FindSequence_Random/large-16", Metrics: []benchmark.Metric{
				{"B/op", []float64{0, 0, 0}},
				{"MB/s", []float64{30814.31, 30814.31, 30814.31}},
				{"allocs/op", []float64{0, 0, 0}},
				{"iterations", []float64{2262572, 2262572, 2262572}},
				{"ns/op", []float64{531.7, 531.7, 531.7}},
			},
		},
		{
			Name: "FindSequence_Random/medium-16", Metrics: []benchmark.Metric{
				{"B/op", []float64{0, 0, 0}},
				{"MB/s", []float64{18339.66, 18339.66, 18339.66}},
				{"allocs/op", []float64{0, 0, 0}},
				{"iterations", []float64{18174138, 18174138, 18174138}},
				{"ns/op", []float64{55.84, 55.84, 55.84}},
			},
		},
		{
			Name: "FindSequence_Random/small-16", Metrics: []benchmark.Metric{
				{"B/op", []float64{0, 0, 0}},
				{"MB/s", []float64{16286.37, 16286.37, 16286.37}},
				{"allocs/op", []float64{0, 0, 0}},
				{"iterations", []float64{76243033, 76243033, 76243033}},
				{"ns/op", []float64{15.72, 15.72, 15.72}},
			},
		},
		{
			Name: "FindSequence_Random/tiny-16", Metrics: []benchmark.Metric{
				{"B/op", []float64{0, 0, 0}},
				{"MB/s", []float64{5461.91, 5461.91, 5461.91}},
				{"allocs/op", []float64{0, 0, 0}},
				{"iterations", []float64{102575256, 102575256, 102575256}},
				{"ns/op", []float64{11.72, 11.72, 11.72}},
			},
		},
	}, actual)
}
