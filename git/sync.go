package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnugomez/voyage/log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repository represents a Git repository with sync capabilities
type Repository struct {
	URL     string
	Branch  string
	OutPath string
	repo    *git.Repository
}

// CreateRepository creates a new Repository instance
func CreateRepository(url, branch, outPath string) *Repository {
	return &Repository{
		URL:     url,
		Branch:  branch,
		OutPath: outPath,
	}
}

// Sync synchronizes the repository, checking multiple subdirectories for changes
// Returns a slice of subdirectories that had updates and any error encountered
func (r *Repository) Sync(subDirs []string) ([]string, error) {
	log.Info("Trying to sync repository", "repo", r.URL, "branch", r.Branch, "subDirs", subDirs)
	defer log.Info("Syncing repository finished without errors")

	var updatedSubDirs []string

	if !directoryExists(r.OutPath) {
		err := r.clone()
		if err != nil {
			return updatedSubDirs, err
		}
		// If we cloned, all subdirectories are considered "updated"
		return subDirs, nil
	}

	gitRepo, err := git.PlainOpen(r.OutPath)
	if err != nil {
		return updatedSubDirs, fmt.Errorf("failed to open repository: %w", err)
	}
	r.repo = gitRepo

	// Fetch latest changes from remote
	err = r.repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Force:      true,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return updatedSubDirs, fmt.Errorf("failed to fetch: %w", err)
	}

	// Check which subdirectories have changes between local and remote
	for _, subDir := range subDirs {
		changed, err := r.HasRemoteChanges(subDir)
		if err != nil {
			return updatedSubDirs, err
		}
		if changed {
			updatedSubDirs = append(updatedSubDirs, subDir)
		}
	}

	if len(updatedSubDirs) == 0 {
		log.Debug("No changes in any subdirectory")
		return updatedSubDirs, nil
	}

	// Pull changes if any subDir changed
	pulled, err := r.pullNewDiff()
	if err != nil {
		return updatedSubDirs, err
	}

	// If pull failed or no changes were actually pulled, return empty slice
	if !pulled {
		return []string{}, nil
	}

	return updatedSubDirs, nil
}

func (r *Repository) pullNewDiff() (bool, error) {
	log.Debug("Checking for changes", "repo", r.URL, "branch", r.Branch)

	if r.repo == nil {
		gitRepo, err := git.PlainOpen(r.OutPath)
		if err != nil {
			return false, fmt.Errorf("failed to open repository: %w", err)
		}
		r.repo = gitRepo
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
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

func (r *Repository) clone() error {
	_, err := git.PlainClone(r.OutPath, false, &git.CloneOptions{
		URL:           r.URL,
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
		SingleBranch:  true,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// HasRemoteChanges checks if subDir has changes between local and remote branch.
func (r *Repository) HasRemoteChanges(subDir string) (bool, error) {
	if r.repo == nil {
		gitRepo, err := git.PlainOpen(r.OutPath)
		if err != nil {
			return false, fmt.Errorf("failed to open repository: %w", err)
		}
		r.repo = gitRepo
	}

	localRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(r.Branch), true)
	if err != nil {
		return false, fmt.Errorf("failed to get local branch ref: %w", err)
	}
	remoteRef, err := r.repo.Reference(plumbing.NewRemoteReferenceName("origin", r.Branch), true)
	if err != nil {
		return false, fmt.Errorf("failed to get remote branch ref: %w", err)
	}

	localCommit, err := r.repo.CommitObject(localRef.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get local commit: %w", err)
	}
	remoteCommit, err := r.repo.CommitObject(remoteRef.Hash())
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

func directoryExists(outPath string) bool {
	if _, err := os.Stat(outPath); err == nil {
		log.Debug("Directory already exists", "path", outPath)
		return true
	}
	return false
}
