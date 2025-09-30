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
		panic(err)
	}

	owner, repo := ctx.Repo()
	newVar := action.Context{
		Token:      act.GetInput("github-token"),
		Owner:      owner,
		Repository: repo,
		CommitHash: ctx.Ref,

		InputFilepath: act.GetInput("benchmarks-file-path"),
		BranchStorage: act.GetInput("branch-storage"),
		BranchPages:   act.GetInput("branch-github-pages"),
	}

	var actErr error
	switch act.GetInput("action-to-perform") {
	case "save":
		actErr = action.Save(newVar)
	default:
		actErr = errors.New("unknown 'action-to-perform'")
	}

	if actErr != nil {
		act.Fatalf("action failed: %s", actErr)
	}
}
