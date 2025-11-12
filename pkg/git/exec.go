package git

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func Fetch(ctx context.Context, repo *git.Repository, branch string) error {
	err := repo.FetchContext(ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(branch + ":" + branch)},
	})

	if err == nil || errors.Is(err, git.NoErrAlreadyUpToDate) {
		return nil
	}

	return err
}

func Checkout(worktree *git.Worktree, branch string) error {
	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
}

func Add(worktree *git.Worktree, filename string) error {
	_, err := worktree.Add(filename)
	return err
}

func Commit(worktree *git.Worktree, message string) error {
	author := &object.Signature{
		Name:  "github-actions[bot]",
		Email: "github-actions[bot]@users.noreply.github.com",
		When:  time.Now(),
	}

	_, err := worktree.Commit(message, &git.CommitOptions{
		Author:            author,
		Committer:         author,
		AllowEmptyCommits: true,
	})

	return err
}

func Push(ctx context.Context, repo *git.Repository, branch string) error {
	return repo.PushContext(ctx, &git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(
				fmt.Sprintf(
					"%s:%s",
					plumbing.NewBranchReferenceName(branch),
					plumbing.NewBranchReferenceName(branch),
				),
			),
		},
	})
}
