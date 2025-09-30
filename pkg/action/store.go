package action

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/go-git/go-git/v6"

	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
)

const (
	TimeoutActionStore    = time.Minute
	ChronosMergedFilename = ".chronos"
	CommitMessage         = "chore: new benchmark series just dropped"
)

func Save(gctx Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutActionStore)
	defer cancel()

	gitops.WithTransient(
		ctx, gctx.Token, gctx.Owner, gctx.Repository,
		func(ctx context.Context, r *git.Repository, w *git.Worktree) error {
			if err := gitops.Fetch(ctx, r, gctx.BranchStorage); err != nil {
				return err
			}

			if err := gitops.Checkout(w, gctx.BranchStorage); err != nil {
				return err
			}

			collected, err := loadBenchmarksSeries(ChronosMergedFilename)
			if err != nil {
				return err
			}

			incoming, err := loadBenchmarksSeries(gctx.InputFilepath)
			if err != nil {
				return err
			}

			if err := saveMergedBenchmarks(benchmark.Merge(collected, incoming)); err != nil {
				return err
			}

			if err := gitops.Add(w, ChronosMergedFilename); err != nil {
				return err
			}

			if err := gitops.Commit(w, CommitMessage, ChronosMergedFilename); err != nil {
				return err
			}

			if err := gitops.Push(ctx, r, gctx.BranchStorage); err != nil {
				return err
			}

			return nil
		},
	)

	return nil
}

func loadBenchmarksSeries(path string) ([]benchmark.Series, error) {
	f, err := os.OpenFile(".chronos", os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var series []benchmark.Series
	if err := json.NewDecoder(f).Decode(&series); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	return series, nil
}

func saveMergedBenchmarks(merged []benchmark.Series) error {
	out, err := os.CreateTemp(".", "*.chronos")
	if err != nil {
		return err
	}

	defer func(path string) {
		out.Close()
		os.Remove(path)
	}(out.Name())

	if json.NewEncoder(out).Encode(merged); err != nil {
		return err
	}

	return os.Rename(out.Name(), ChronosMergedFilename)
}
