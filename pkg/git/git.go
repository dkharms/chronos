package git

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
)

const (
	repoTpl = "https://x-access-token:%s@github.com/%s/%s.git"
)

func WithTransient(
	ctx context.Context,
	token, owner, repositoryName, branch string,
	do func(context.Context, *git.Repository, *git.Worktree) ([]string, string, error),
) error {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	repository, err := git.PlainCloneContext(ctx, tmp, &git.CloneOptions{
		URL: fmt.Sprintf(repoTpl, token, owner, repositoryName),
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(dir)

	if err := os.Chdir(tmp); err != nil {
		return err
	}

	if err := Fetch(ctx, repository, branch); err != nil {
		return err
	}

	if err := Checkout(worktree, branch); err != nil {
		return err
	}

	commitable, message, err := do(ctx, repository, worktree)
	if err != nil {
		return err
	}

	if len(commitable) == 0 {
		return nil
	}

	for _, c := range commitable {
		if err := Add(worktree, c); err != nil {
			return err
		}
	}

	if err := Commit(worktree, message); err != nil {
		return err
	}

	if err := Push(ctx, repository, branch); err != nil {
		return err
	}

	return nil
}
