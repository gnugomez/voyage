package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GitService defines a set of high-level Git operations.
type GitService interface {
	Clone(path, url, branch string) error
	Fetch(path string) error
	IsBehindRemote(path, branch string) (bool, error)
	Pull(path, branch string) error
	HasChangesInSubdir(path, branch, subDir string) (bool, error)
	IsGitRepository(path string) bool
}

// cliGitService is the implementation of GitService that uses the git command line.
type cliGitService struct{}

func NewCliGitService() GitService {
	return &cliGitService{}
}

func (s *cliGitService) IsGitRepository(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	err := cmd.Run()
	return err == nil
}

func (s *cliGitService) Fetch(path string) error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch: %w, output: %s", err, string(output))
	}
	return nil
}

func (s *cliGitService) IsBehindRemote(path, branch string) (bool, error) {
	cmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..origin/%s", branch, branch))
	cmd.Dir = path
	countOutput, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check if behind remote: %w", err)
	}

	behindCountStr := strings.TrimSpace(string(countOutput))
	behindCount, err := strconv.Atoi(behindCountStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse commit count: %w", err)
	}

	return behindCount > 0, nil
}

func (s *cliGitService) Pull(path, branch string) error {
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull: %w, output: %s", err, string(output))
	}
	return nil
}

func (s *cliGitService) Clone(path, url, branch string) error {
	cmd := exec.Command("git", "clone", "-b", branch, "--single-branch", url, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w, output: %s", err, string(output))
	}
	return nil
}

func (s *cliGitService) HasChangesInSubdir(path, branch, subDir string) (bool, error) {
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("%s..origin/%s", branch, branch), "--", subDir)
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to get diff for subdirectory %s: %w", subDir, err)
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}
