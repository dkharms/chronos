package main

import (
	"context"
	"errors"
	"time"

	"github.com/kkyr/fig"
	"github.com/sethvargo/go-githubactions"

	"github.com/dkharms/chronos/pkg/action"
	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
)

const (
	ChronosConfigFilename = ".config.yaml"
	Timeout               = time.Minute * 5
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		Timeout,
	)
	defer cancel()

	act := githubactions.New()
	gctx, err := act.Context()
	if err != nil {
		exit(err)
	}

	owner, repo := gctx.Repo()
	input := action.Input{
		Token: act.GetInput("github-token"),

		Owner:      owner,
		Repository: repo,
		CommitHash: gctx.SHA,

		LanguageTool: act.GetInput("language-tool"),

		BenchmarksFilepath: act.GetInput("benchmarks-file-path"),
		BranchStorage:      act.GetInput("branch-storage"),
	}

	r, err := gitops.WithRepository(
		context.Background(), input.Token,
		input.Owner, input.Repository,
	)
	if err != nil {
		exit(err)
	}

	cfg, err := loadConfig(ctx, r, input.BranchStorage)
	if err != nil {
		exit(err)
	}

	var actErr error
	switch act.GetInput("action-to-perform") {
	case "save":
		actErr = action.Save(ctx, r, cfg, input)
	case "publish":
		actErr = action.Publish(ctx, r, cfg, input)
	case "summarize":
		actErr = action.Summarize(ctx, r, cfg, input)
	default:
		actErr = errors.New("unknown 'action-to-perform'")
	}

	if actErr != nil {
		exit(actErr)
	}
}

func loadConfig(ctx context.Context, r gitops.Repository, branch string) (action.Config, error) {
	var cfg action.Config

	err := r.WithBranch(ctx, branch, func() ([]string, string, error) {
		err := fig.Load(&cfg,
			fig.AllowNoFile(),
			fig.File(ChronosConfigFilename),
		)
		if err != nil {
			return nil, "", err
		}

		for _, unit := range cfg.Units {
			benchmark.AddMetricDescriptor(
				benchmark.NewMetricDescriptor(
					unit.Name, unit.Better,
					unit.Threshold, unit.Reduction,
				),
			)
		}

		return nil, "", nil
	})

	return cfg, err
}

func exit(err error) {
	githubactions.Fatalf("failed: %s", err)
}
