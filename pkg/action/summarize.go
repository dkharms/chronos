package action

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/sethvargo/go-githubactions"

	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
)

var (
	//go:embed assets/summary.tpl
	summaryTemplate string
)

func Summarize(ctx context.Context, r gitops.Repository, cfg Config, input Input) error {
	var calculatedDiff []benchmark.CalculatedDiff

	err := r.WithBranch(
		ctx, input.BranchStorage,
		func() ([]string, string, error) {
			incoming, err := loadIncomingBenchmarks(
				input.LanguageTool,
				input.CommitHash,
				input.BenchmarksFilepath,
			)
			if err != nil {
				return nil, "", err
			}

			series, err := loadCollectedBenchmarks(ChronosMergedFilename)
			if err != nil {
				return nil, "", err
			}

			diff := benchmark.Diff(series, incoming)
			slices.SortFunc(diff, func(x, y benchmark.CalculatedDiff) int {
				return strings.Compare(x.Name, y.Name)
			})

			calculatedDiff = diff
			return nil, "", githubactions.AddStepSummaryTemplate(
				summaryTemplate, diff,
			)
		},
	)

	if err != nil {
		return fmt.Errorf(
			"cannot summarize benchmarks: %w",
			err,
		)
	}

	if input.EventName == "pull_request" {
		return comment(ctx, input, calculatedDiff)
	}

	return nil
}
