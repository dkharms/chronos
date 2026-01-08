package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
	"github.com/dkharms/chronos/pkg/parser"
)

const (
	ChronosMergedFilename   = ".chronos"
	ActionSaveCommitMessage = "[chronos] `save` (%s)"
)

func Save(ctx context.Context, r gitops.Repository, cfg Config, input Input) error {
	err := r.WithBranch(
		ctx, input.BranchStorage,
		func() ([]string, string, error) {
			return []string{ChronosMergedFilename},
				fmt.Sprintf(ActionSaveCommitMessage, input.CommitHash),
				processBenchmarks(cfg, input)
		},
	)

	if err != nil {
		return fmt.Errorf(
			"cannot save benchmarks: %w",
			err,
		)
	}

	return nil
}

func processBenchmarks(cfg Config, input Input) error {
	collected, err := loadCollectedBenchmarks(ChronosMergedFilename)
	if err != nil {
		return err
	}

	incoming, err := loadIncomingBenchmarks(
		input.LanguageTool,
		input.CommitHash,
		input.BenchmarksFilepath,
	)
	if err != nil {
		return err
	}

	merged := benchmark.Merge(collected, incoming)
	for i := range len(merged) {
		sub := max(0, len(merged[i].Measurements)-cfg.Storage.Capacity)
		merged[i].Measurements = merged[i].Measurements[sub:]
	}

	return saveMergedBenchmarks(merged)
}

func loadCollectedBenchmarks(path string) ([]benchmark.Series, error) {
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

func loadIncomingBenchmarks(tool, commitHash, path string) ([]benchmark.Series, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p, err := parser.New(tool, f)
	if err != nil {
		return nil, err
	}
	raw := p.Parse()

	var series []benchmark.Series
	for _, r := range raw {
		r.CommitHash = commitHash
		series = append(series, benchmark.Series{
			Name:         r.Name,
			Measurements: []benchmark.Measurement{r},
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
