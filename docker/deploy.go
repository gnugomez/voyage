package docker

import (
	"fmt"
	"os"
	"os/exec"
)

func DeployCompose(targetPath string, daemonMode bool) error {
	if err := checkForDockerEnv(); err != nil {
		return err
	}
	if err := checkTargetPath(targetPath); err != nil {
		return err
	}

	args := []string{"compose", "-f", targetPath, "up"}
	if daemonMode {
		args = append(args, "-d")
	}
	cmd := exec.Command("docker", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run docker compose: %w", err)
	}

	return nil
}

func checkTargetPath(targetPath string) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("target path does not exist: %v", err)
	}
	return nil
}

func checkForDockerEnv() error {
	if sockRunning, err := isDockerSocketRunning(); !sockRunning && err != nil {
		return fmt.Errorf("docker socket is not running: %v", err)
	}
	if installed, err := isDockerComposeInstalled(); !installed && err != nil {
		return fmt.Errorf("docker compose is not installed: %v", err)
	}
	return nil
}

func isDockerComposeInstalled() (bool, error) {
	cmd := exec.Command("docker", "compose", "--version")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

func isDockerSocketRunning() (bool, error) {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
