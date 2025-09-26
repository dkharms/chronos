package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/sethvargo/go-githubactions"

	gitops "github.com/dkharms/chronos/pkg/git"
	"github.com/dkharms/chronos/pkg/parser"
)

func main() {
	action := githubactions.New()

	ctx, err := action.Context()
	if err != nil {
		panic(err)
	}

	branch := "chronos-storage"
	owner, repository := ctx.Repo()

	err = gitops.WithTransient(
		context.Background(), action.GetInput("github-token"), owner, repository,
		func(ctx context.Context, repo *git.Repository, worktree *git.Worktree) error {
			if err := gitops.Fetch(ctx, repo, branch); err != nil {
				return err
			}

			if err := gitops.Checkout(worktree, branch); err != nil {
				return err
			}

			results := parser.NewGoParser(strings.NewReader(parser.GotoolOutput)).Parse()

			f, err := os.OpenFile(".chronos", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0o644)
			if err != nil {
				return err
			}
			defer f.Close()

			enc := json.NewEncoder(f)
			if err := enc.Encode(results); err != nil {
				return err
			}

			if err := gitops.Add(worktree, ".chronos"); err != nil {
				return err
			}

			if err := gitops.Commit(worktree, "some message", ".chronos"); err != nil {
				return err
			}

			if err := gitops.Push(ctx, repo, branch); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		panic(err)
	}
}
