package git

import (
	"fmt"
	"os"

	"github.com/gnugomez/voyage/log"
)

type Repository struct {
	URL             string
	Branch          string
	OutPath         string
	gitService      GitService
	directoryExists func(string) bool
}

// CreateRepository creates a new Repository instance
func CreateRepository(url, branch, outPath string) *Repository {
	return &Repository{
		URL:             url,
		Branch:          branch,
		OutPath:         outPath,
		gitService:      NewCliGitService(),
		directoryExists: osDirectoryExists,
	}
}

// Sync synchronizes the repository, checking multiple subdirectories for changes
// Returns a slice of subdirectories that had updates and any error encountered
func (r *Repository) Sync(subDirs []string) ([]string, error) {
	log.Info("Trying to sync repository", "repo", r.URL, "branch", r.Branch, "subDirs", subDirs)

	if !r.directoryExists(r.OutPath) {
		err := r.gitService.Clone(r.OutPath, r.URL, r.Branch)
		if err != nil {
			return nil, err
		}
		// If we cloned, all subdirectories are considered "updated"
		return subDirs, nil
	}

	if !r.gitService.IsGitRepository(r.OutPath) {
		return nil, fmt.Errorf("directory %s is not a git repository", r.OutPath)
	}

	err := r.gitService.Fetch(r.OutPath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}

	var updatedSubDirs []string
	for _, subDir := range subDirs {
		changed, err := r.gitService.HasChangesInSubdir(r.OutPath, r.Branch, subDir)
		if err != nil {
			return nil, err
		}
		if changed {
			updatedSubDirs = append(updatedSubDirs, subDir)
		}
	}

	if len(updatedSubDirs) == 0 {
		log.Debug("No changes in any subdirectory")
		return nil, nil
	}

	// If we detected changes in subdirectories, pull the code.
	log.Debug("Changes detected in subdirectories, pulling changes.", "subDirs", updatedSubDirs)
	isBehind, err := r.gitService.IsBehindRemote(r.OutPath, r.Branch)
	if err != nil {
		return nil, err
	}

	if isBehind {
		err = r.gitService.Pull(r.OutPath, r.Branch)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug("Subdirectory changes detected, but remote is not ahead. No pull needed.")
	}

	return updatedSubDirs, nil
}

func osDirectoryExists(outPath string) bool {
	if _, err := os.Stat(outPath); err == nil {
		log.Debug("Directory already exists", "path", outPath)
		return true
	}
	return false
}


