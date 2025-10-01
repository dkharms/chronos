package main

import (
	"errors"

	"github.com/sethvargo/go-githubactions"

	"github.com/dkharms/chronos/pkg/action"
)

func main() {
	act := githubactions.New()

	ctx, err := act.Context()
	if err != nil {
		exit(err)
	}

	owner, repo := ctx.Repo()
	gctx := action.Context{
		Token:      act.GetInput("github-token"),
		Owner:      owner,
		Repository: repo,
		CommitHash: ctx.SHA,

		InputFilepath: act.GetInput("benchmarks-file-path"),
		BranchStorage: act.GetInput("branch-storage"),
		BranchPages:   act.GetInput("branch-github-pages"),
	}

	var actErr error
	switch act.GetInput("action-to-perform") {
	case "save":
		actErr = action.Save(gctx)
	case "publish":
		actErr = action.Publish(gctx)
	default:
		actErr = errors.New("unknown 'action-to-perform'")
	}

	if actErr != nil {
		exit(actErr)
	}
}

func exit(err error) {
	githubactions.Fatalf("failed: %s", err)
}
