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
			Name: "FindSequence_Random/tiny-16", Metrics: []benchmark.Metric{
				{"iterations", 102575256}, {"ns/op", 11.72}, {"MB/s", 5461.91}, {"B/op", 0}, {"allocs/op", 0},
			},
		},
		{
			Name: "FindSequence_Random/small-16", Metrics: []benchmark.Metric{
				{"iterations", 76243033}, {"ns/op", 15.72}, {"MB/s", 16286.37}, {"B/op", 0}, {"allocs/op", 0},
			},
		},
		{
			Name: "FindSequence_Random/medium-16", Metrics: []benchmark.Metric{
				{"iterations", 18174138}, {"ns/op", 55.84}, {"MB/s", 18339.66}, {"B/op", 0}, {"allocs/op", 0},
			},
		},
		{
			Name: "FindSequence_Random/large-16", Metrics: []benchmark.Metric{
				{"iterations", 2262572}, {"ns/op", 531.7}, {"MB/s", 30814.31}, {"B/op", 0}, {"allocs/op", 0},
			},
		},
		{
			Name: "FindSequence_Random/extra-large-16", Metrics: []benchmark.Metric{
				{"iterations", 20184}, {"ns/op", 59050}, {"MB/s", 17757.45}, {"B/op", 0}, {"allocs/op", 0},
			},
		},
	}, actual)
}
