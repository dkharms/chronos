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
	"github.com/dkharms/chronos/pkg/parser"
)

const (
	TimeoutActionSave     = time.Minute
	ChronosMergedFilename = ".chronos"
	CommitMessage         = "chore: new benchmark series just dropped"
)

func Save(gctx Context) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		TimeoutActionSave,
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
			return []string{ChronosMergedFilename},
				CommitMessage,
				ProcessBenchmarks(gctx)
		},
	)
}

func ProcessBenchmarks(gctx Context) error {
	collected, err := loadBenchmarksSeries(ChronosMergedFilename)
	if err != nil {
		return err
	}

	incoming, err := parseBenchmarkSeries(gctx.CommitHash, gctx.InputFilepath)
	if err != nil {
		return err
	}

	if err := saveMergedBenchmarks(benchmark.Merge(collected, incoming)); err != nil {
		return err
	}

	return nil
}

func loadBenchmarksSeries(path string) ([]benchmark.Series, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0o644)
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

func parseBenchmarkSeries(commitHash, path string) ([]benchmark.Series, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	raw := parser.NewGoParser(f).Parse()

	var series []benchmark.Series
	for _, r := range raw {
		r.CommitHash = commitHash
		series = append(series, benchmark.Series{
			Name:   r.Name,
			Points: []benchmark.Result{r},
		})
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

	if err := json.NewEncoder(out).Encode(merged); err != nil {
		return err
	}

	return os.Rename(out.Name(), ChronosMergedFilename)
}
