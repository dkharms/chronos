package git

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

const (
	repoPath = "/tmp/chronos"
)

func WithTransient(
	ctx context.Context,
	token, owner, repository string,
	do func(context.Context, *git.Repository, *git.Worktree) error,
) error {
	repo, err := git.PlainCloneContext(ctx, repoPath, &git.CloneOptions{
		URL: fmt.Sprintf("https://github.com/%s/%s.git", owner, repository),
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: token,
		},
	})
	if err != nil {
		return err
	}
	defer os.RemoveAll(repoPath)

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(dir)

	if err := os.Chdir(repoPath); err != nil {
		return err
	}

	return do(ctx, repo, worktree)
}
