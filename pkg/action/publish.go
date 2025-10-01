package action

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
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
	htmlTemplate string
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
					if xerr := saveIndexFile(series); xerr != nil {
						return nil, "", xerr
					}

					return []string{"index.html"},
						fmt.Sprintf(ActionPublishCommitMessage, gctx.CommitHash),
						saveMergedBenchmarks(series)
				},
			)
		},
	)
}

func saveIndexFile(series []benchmark.Series) error {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(series)
	if err != nil {
		return err
	}

	f, err := os.OpenFile("index.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
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
