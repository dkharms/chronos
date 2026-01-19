package action

import (
	"context"
	_ "embed"
	"strings"
	"text/template"

	"github.com/google/go-github/v81/github"
	"golang.org/x/oauth2"

	"github.com/dkharms/chronos/pkg/benchmark"
)

var (
	//go:embed assets/pr-comment-degraded.tpl
	degradedTpl string
)

func comment(
	ctx context.Context, input Input, diffs []benchmark.CalculatedDiff,
) error {
	client, err := ghclient(ctx, input.Token)
	if err != nil {
		return err
	}

	degradedDiffs := filterDegraded(diffs)
	if len(degradedDiffs) == 0 {
		return nil
	}

	body, err := renderTemplate(degradedDiffs)
	if err != nil {
		return err
	}

	_, _, err = client.Issues.CreateComment(
		ctx, input.Owner, input.Repository,
		input.PRNumber, &github.IssueComment{Body: github.Ptr(body)},
	)

	return err
}

func ghclient(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

func filterDegraded(diffs []benchmark.CalculatedDiff) []benchmark.CalculatedDiff {
	var degraded []benchmark.CalculatedDiff

	for _, diff := range diffs {
		var degradedMetrics []benchmark.MetricDiff

		for _, metric := range diff.MetricDiff {
			if metric.Emoji() == "🔴" {
				degradedMetrics = append(degradedMetrics, metric)
			}
		}

		if len(degradedMetrics) > 0 {
			degraded = append(degraded, benchmark.CalculatedDiff{
				Name:           diff.Name,
				PreviousCommit: diff.PreviousCommit,
				CurrentCommit:  diff.CurrentCommit,
				MetricDiff:     degradedMetrics,
			})
		}
	}

	return degraded
}

func renderTemplate(diffs []benchmark.CalculatedDiff) (string, error) {
	var b strings.Builder

	t, err := template.New("comment").Parse(degradedTpl)
	if err != nil {
		return "", err
	}

	if err := t.Execute(&b, diffs); err != nil {
		return "", err
	}

	return b.String(), nil
}
