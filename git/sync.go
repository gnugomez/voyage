package git

import (
	"fmt"
	"gnugomez/voyage/log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func SyncRepo(repo string, branch string, outPath string) (bool, error) {
	log.Debug("Trying to clone repository", "repo", repo, "branch", branch)

	if directoryExists(outPath) {
		return pullNewDiff(repo, branch, outPath)
	}

	err := cloneRepo(repo, branch, outPath)
	return true, err
}

func pullNewDiff(repo string, branch string, executionPath string) (bool, error) {
	log.Debug("Checking for changes", "repo", repo, "branch", branch)

	r, err := git.PlainOpen(executionPath)
	if err != nil {
		return false, fmt.Errorf("failed to open repository: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Force:         true,
	})

	if err == git.NoErrAlreadyUpToDate {
		log.Debug("Already up to date")
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to pull: %w", err)
	}

	log.Debug("Changes pulled successfully")
	return true, nil
}

func cloneRepo(repo string, branch string, outPath string) error {
	_, err := git.PlainClone(outPath, false, &git.CloneOptions{
		URL:           repo,
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

func directoryExists(outPath string) bool {
	if _, err := os.Stat(outPath); err == nil {
		log.Debug("Directory already exists", "path", outPath)
		return true
	}
	return false
}
