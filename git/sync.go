package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnugomez/voyage/log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// SyncRepo now only pulls if subDir has changes in the remote branch.
func SyncRepo(repo string, branch string, outPath string, subDir string) (bool, error) {
	log.Info("Trying to sync repository", "repo", repo, "branch", branch, "subDir", subDir)
	defer log.Info("Syncing repository finished without errors")

	if !directoryExists(outPath) {
		err := cloneRepo(repo, branch, outPath)
		return true, err
	}

	r, err := git.PlainOpen(outPath)
	if err != nil {
		return false, fmt.Errorf("failed to open repository: %w", err)
	}

	// Fetch latest changes from remote
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return false, fmt.Errorf("failed to fetch: %w", err)
	}

	// Check if subDir has changes between local and remote
	changed, err := subDirHasRemoteChanges(r, branch, subDir)
	if err != nil {
		return false, err
	}
	if !changed {
		log.Debug("No changes in subDir", "subDir", subDir)
		return false, nil
	}

	// Pull changes if subDir changed
	return pullNewDiff(repo, branch, outPath)
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

// subDirHasRemoteChanges checks if subDir has changes between local and remote branch.
func subDirHasRemoteChanges(r *git.Repository, branch, subDir string) (bool, error) {
	localRef, err := r.Reference(plumbing.NewBranchReferenceName(branch), true)
	if err != nil {
		return false, fmt.Errorf("failed to get local branch ref: %w", err)
	}
	remoteRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
	if err != nil {
		return false, fmt.Errorf("failed to get remote branch ref: %w", err)
	}

	localCommit, err := r.CommitObject(localRef.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get local commit: %w", err)
	}
	remoteCommit, err := r.CommitObject(remoteRef.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get remote commit: %w", err)
	}

	patch, err := remoteCommit.Patch(localCommit)
	if err != nil {
		return false, fmt.Errorf("failed to get patch: %w", err)
	}

	for _, fileStat := range patch.FilePatches() {
		from, to := fileStat.Files()
		if (from != nil && isInSubDir(from.Path(), subDir)) || (to != nil && isInSubDir(to.Path(), subDir)) {
			return true, nil
		}
	}
	return false, nil
}

// isInSubDir checks if the file path is within the subDir.
func isInSubDir(filePath, subDir string) bool {
	cleanSubDir := filepath.Clean(subDir) + string(os.PathSeparator)
	cleanFile := filepath.Clean(filePath)
	return len(cleanFile) >= len(cleanSubDir) && cleanFile[:len(cleanSubDir)] == cleanSubDir
}
