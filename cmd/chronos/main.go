package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/object"
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

	_, repository := ctx.Repo()
	err = gitops.WithRepositoryAndBranch(
		context.Background(), repository, "chronos-storage",
		func(r *git.Repository, w *git.Worktree, path string) error {
			seriesPath := filepath.Join(path, ".chronos")
			bpath := action.GetInput("benchmarks-file-path")

			in, err := os.Open(bpath)
			if err != nil {
				return err
			}
			defer in.Close()
			incoming := parser.NewGoParser(in).Parse()

			out, err := os.OpenFile(
				seriesPath,
				os.O_RDWR|os.O_CREATE, 0o644,
			)
			if err != nil {
				return err
			}
			defer out.Close()

			enc := json.NewEncoder(out)
			if err := enc.Encode(incoming); err != nil {
				return err
			}

			if _, err := w.Add(seriesPath); err != nil {
				return err
			}

			author := object.Signature{
				Name:  "chronos",
				Email: "chronos@noreply.com",
				When:  time.Now(),
			}

			if _, err := w.Commit("[chronos] save new measurements", &git.CommitOptions{
				All:               true,
				AllowEmptyCommits: false,
				Author:            &author,
				Committer:         &author,
			}); err != nil {
				return err
			}

			return r.PushContext(context.TODO(), &git.PushOptions{
				RemoteName: "origin",
				RefSpecs:   []config.RefSpec{"chronos-storage:chronos-storage"},
			})
		},
	)

	if err != nil {
		panic(err)
	}
}
