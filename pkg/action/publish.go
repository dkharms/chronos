package action

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/dkharms/chronos/pkg/benchmark"
	gitops "github.com/dkharms/chronos/pkg/git"
)

const (
	ActionPublishCommitMessage = "[chronos] `publish` (%s)"
)

var (
	//go:embed assets/index.html
	htmlTemplate string
)

func Publish(ctx context.Context, r gitops.Repository, cfg Config, input Input) error {
	var series []benchmark.Series

	err := r.WithBranch(
		ctx, input.BranchStorage,
		func() ([]string, string, error) {
			xseries, xerr := loadCollectedBenchmarks(ChronosMergedFilename)
			if xerr != nil {
				return nil, "", xerr
			}
			series = xseries
			return nil, "", nil
		},
	)

	if err != nil {
		return fmt.Errorf(
			"cannot load collected benchmarks: %w",
			err,
		)
	}

	err = r.WithBranch(
		ctx, cfg.GithubPages.Branch,
		func() ([]string, string, error) {
			p := filepath.Join(cfg.GithubPages.Directory, "index.html")
			return []string{p},
				fmt.Sprintf(ActionPublishCommitMessage, input.CommitHash),
				saveIndexFile(p, series)
		},
	)

	if err != nil {
		return fmt.Errorf(
			"cannot save rendered HTML page: %w",
			err,
		)
	}

	return nil
}

func saveIndexFile(p string, series []benchmark.Series) error {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(series)
	if err != nil {
		return err
	}

	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		BenchmarkData template.JS
	}{
		BenchmarkData: template.JS(jsonData),
	}

	return tmpl.Execute(f, data)
}
