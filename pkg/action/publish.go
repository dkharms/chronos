package action

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v6"

	gitops "github.com/dkharms/chronos/pkg/git"
)

const (
	ActionPublishTimeout       = 5 * time.Minute
	ActionPublishCommitMessage = "[chronos] `publish` (%s)"
)

var (
	//go:embed assets/index.ts
	TypescriptFile string
)

func Publish(gctx Context) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ActionSaveTimeout,
	)

	defer cancel()

	return gitops.WithTransient(
		ctx, gctx.Token, gctx.Owner,
		gctx.Repository, gctx.BranchStorage,
		func(
			ctx context.Context,
			repository *git.Repository,
			worktree *git.Worktree,
		) ([]string, string, error) {
			series, err := loadBenchmarksSeries(ChronosMergedFilename)
			if err != nil {
				return nil, "", err
			}

			if err := gitops.Fetch(ctx, repository, gctx.BranchPages); err != nil {
				return nil, "", err
			}

			if err := gitops.Checkout(worktree, gctx.BranchPages); err != nil {
				return nil, "", err
			}

			if err := saveIndexFile(); err != nil {
				return nil, "", err
			}

			return []string{ChronosMergedFilename, "index.ts"},
				fmt.Sprintf(ActionPublishCommitMessage, gctx.CommitHash),
				saveMergedBenchmarks(series)
		},
	)
}

func saveIndexFile() error {
	f, err := os.OpenFile("index.ts", os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(TypescriptFile); err != nil {
		return err
	}

	return nil
}
