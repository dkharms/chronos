package action

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"
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
	htmlTemplate string
)

func Publish(ctx context.Context, r gitops.Repository, cfg Config, input Input) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		ActionSaveTimeout,
	)

	defer cancel()

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
		return err
	}

	return r.WithBranch(
		ctx, cfg.GithubPages.Branch,
		func() ([]string, string, error) {
			filepath := path.Join(cfg.GithubPages.Path, "index.html")
			return []string{filepath},
				fmt.Sprintf(ActionPublishCommitMessage, input.CommitHash),
				saveIndexFile(filepath, series)
		},
	)
}

func saveIndexFile(filepath string, series []benchmark.Series) error {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(series)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
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
