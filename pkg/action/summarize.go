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

func Summarize(gctx Context) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ActionSummarizeTimeout,
	)

	defer cancel()

	return gitops.WithTransient(
		ctx, gctx.Token, gctx.Owner, gctx.Repository,
		func(ctx context.Context, r gitops.Repository) error {
			return r.WithBranch(
				ctx, gctx.BranchStorage,
				func() ([]string, string, error) {
					incoming, err := parseBenchmarkSeries(
						gctx.CommitHash, gctx.InputFilepath,
					)
					if err != nil {
						return nil, "", err
					}

					series, err := loadBenchmarksSeries(ChronosMergedFilename)
					if err != nil {
						return nil, "", err
					}

					return nil, "", githubactions.AddStepSummaryTemplate(
						summaryTemplate, benchmark.Diff(series, incoming),
					)
				},
			)
		},
	)
}
