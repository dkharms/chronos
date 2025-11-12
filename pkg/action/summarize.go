package action

import (
	"context"
	_ "embed"
	"time"

	"github.com/sethvargo/go-githubactions"

	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
)

const (
	ActionSummarizeTimeout = time.Minute
)

var (
	//go:embed assets/summary.tpl
	summaryTemplate string
)

func Summarize(ctx context.Context, r gitops.Repository, cfg Config, input Input) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ActionSummarizeTimeout,
	)

	defer cancel()

	return r.WithBranch(
		ctx, input.BranchStorage,
		func() ([]string, string, error) {
			incoming, err := loadIncomingBenchmarks(input.CommitHash, input.BenchmarksFilepath)
			if err != nil {
				return nil, "", err
			}

			series, err := loadCollectedBenchmarks(ChronosMergedFilename)
			if err != nil {
				return nil, "", err
			}

			return nil, "", githubactions.AddStepSummaryTemplate(
				summaryTemplate, benchmark.Diff(series, incoming),
			)
		},
	)
}
