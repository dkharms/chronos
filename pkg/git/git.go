package git

import (
	"context"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

const (
	repoPath = "/tmp/chronos"
)

func WithRepositoryAndBranch(
	ctx context.Context,
	repoURL, branch string,
	do func(*git.Repository, *git.Worktree, string) error,
) error {
	repo, err := git.PlainCloneContext(ctx, repoPath, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(repoPath)

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	}); err != nil {
		return err
	}

	if err := do(repo, worktree, repoPath); err != nil {
		return err
	}

	return nil
}
