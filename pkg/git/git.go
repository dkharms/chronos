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
	token, owner, repositoryName string,
	fn func(context.Context, Repository) error,
) error {
	tmp, err := os.MkdirTemp(os.TempDir(), "chronos-*")
	if err != nil {
		return err
	}

	repository, err := git.PlainCloneContext(ctx, tmp, &git.CloneOptions{
		URL: fmt.Sprintf(repoTpl, token, owner, repositoryName),
	})
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer func() { _ = os.Chdir(dir) }()

	if err := os.Chdir(tmp); err != nil {
		return err
	}

	return fn(ctx, Repository{
		r: repository,
		w: worktree,
	})
}

type Repository struct {
	r *git.Repository
	w *git.Worktree
}

func (r Repository) WithBranch(
	ctx context.Context, branch string,
	fn func() ([]string, string, error),
) error {
	cur, err := r.r.Head()
	if err != nil {
		return err
	}
	defer func() { _ = Checkout(r.w, cur.String()) }()

	if err := Fetch(ctx, r.r, branch); err != nil {
		return err
	}

	if err := Checkout(r.w, branch); err != nil {
		return err
	}

	commitable, message, err := fn()
	if err != nil {
		return err
	}

	if len(commitable) == 0 {
		return nil
	}

	for _, c := range commitable {
		if err := Add(r.w, c); err != nil {
			return err
		}
	}

	if err := Commit(r.w, message); err != nil {
		return err
	}

	if err := Push(ctx, r.r, branch); err != nil {
		return err
	}

	return nil
}
