package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gnugomez/voyage/log"
)

// Repository represents a Git repository with sync capabilities
type Repository struct {
	URL     string
	Branch  string
	OutPath string
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

	// Check if it's a valid git repository
	if !r.isGitRepository() {
		return updatedSubDirs, fmt.Errorf("directory %s is not a git repository", r.OutPath)
	}

	// Fetch latest changes from remote
	err := r.fetch()
	if err != nil {
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

// isGitRepository checks if the directory is a valid git repository
func (r *Repository) isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = r.OutPath
	err := cmd.Run()
	return err == nil
}

// fetch retrieves the latest changes from the remote repository
func (r *Repository) fetch() error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = r.OutPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch: %w, output: %s", err, string(output))
	}
	return nil
}

func (r *Repository) pullNewDiff() (bool, error) {
	log.Debug("Checking for changes", "repo", r.URL, "branch", r.Branch)

	// Check if local branch is behind remote
	cmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..origin/%s", r.Branch, r.Branch))
	cmd.Dir = r.OutPath
	countOutput, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check if behind remote: %w", err)
	}

	behindCount := strings.TrimSpace(string(countOutput))
	if behindCount == "0" {
		log.Debug("Already up to date")
		return false, nil
	}

	// Pull the changes
	cmd = exec.Command("git", "pull", "origin", r.Branch)
	cmd.Dir = r.OutPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to pull: %w, output: %s", err, string(output))
	}

	log.Debug("Changes pulled successfully")
	return true, nil
}

func (r *Repository) clone() error {
	cmd := exec.Command("git", "clone", "-b", r.Branch, "--single-branch", r.URL, r.OutPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w, output: %s", err, string(output))
	}
	return nil
}

// HasRemoteChanges checks if subDir has changes between local and remote branch.
func (r *Repository) HasRemoteChanges(subDir string) (bool, error) {
	// Get the diff between local and remote branch for the specific subdirectory
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("%s..origin/%s", r.Branch, r.Branch), "--", subDir)
	cmd.Dir = r.OutPath
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to get diff for subdirectory %s: %w", subDir, err)
	}

	// If there's any output, it means there are changes in the subdirectory
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func directoryExists(outPath string) bool {
	if _, err := os.Stat(outPath); err == nil {
		log.Debug("Directory already exists", "path", outPath)
		return true
	}
	return false
}
