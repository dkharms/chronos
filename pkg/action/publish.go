package action

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
)

const (
	ActionPublishTimeout       = 5 * time.Minute
	ActionPublishCommitMessage = "[chronos] `publish` (%s)"
)

var (
	//go:embed assets/index.html
	TypescriptFile string
)

func Publish(gctx Context) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ActionSaveTimeout,
	)

	defer cancel()

	return gitops.WithTransient(
		ctx, gctx.Token, gctx.Owner, gctx.Repository,
		func(ctx context.Context, r gitops.Repository) error {
			var series []benchmark.Series

			err := r.WithBranch(
				ctx, gctx.BranchStorage,
				func() ([]string, string, error) {
					xseries, err := loadBenchmarksSeries(ChronosMergedFilename)
					if err != nil {
						return nil, "", err
					}
					series = xseries
					return nil, "", nil
				},
			)

			if err != nil {
				return err
			}

			return r.WithBranch(
				ctx, gctx.BranchPages,
				func() ([]string, string, error) {
					if xerr := saveIndexFile(); xerr != nil {
						return nil, "", xerr
					}

					return []string{ChronosMergedFilename, "index.html"},
						fmt.Sprintf(ActionPublishCommitMessage, gctx.CommitHash),
						saveMergedBenchmarks(series)
				},
			)
		},
	)
}

func saveIndexFile() error {
	f, err := os.OpenFile("index.html", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(TypescriptFile); err != nil {
		return err
	}

	return nil
}
