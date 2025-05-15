package main

import (
	"gnugomez/voyage/log"
	"os"
	"os/exec"
)

func SyncRepo(repo string, branch string, outPath string) (bool, error) {
	log.Debug("Trying to clone repository", "repo", repo, "branch", branch)

	if directoryExists(outPath) {
		return pullNewDiff(repo, branch, outPath)
	}

	err := cloneRepo(repo, outPath)
	return true, err
}

func pullNewDiff(repo string, branch string, executionPath string) (bool, error) {
	log.Debug("Checking for changes", "repo", repo, "branch", branch)

	err := os.Chdir(executionPath)
	if err != nil {
		return false, err
	}

	output, err := execCommandWithLogging(exec.Command("git", "diff", "--name-only", "origin/"+branch))
	if err != nil {
		return false, err
	}

	if len(output) > 0 {
		log.Debug("New changes found")
		output, err = execCommandWithLogging(exec.Command("git", "pull", "origin", branch))
		if err != nil {
			return false, err
		}
		log.Debug("Changes pulled")
		return true, nil
	}
	return false, nil
}
func cloneRepo(repo string, outPath string) error {
	_, err := execCommandWithLogging(exec.Command("git", "clone", repo, outPath))
	return err
}

func directoryExists(outPath string) bool {
	if _, err := os.Stat(outPath); err == nil {
		log.Debug("Directory already exists", "path", outPath)
		return true
	}
	return false
}

func execCommandWithLogging(cmd *exec.Cmd) ([]byte, error) {
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Debug("Command error", "error", err, "output", string(output))
	}
	log.Debug("Command output", "output", string(output))
	return output, err
}
