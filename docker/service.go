package docker

import (
	"fmt"
	"io"
	"os/exec"
)

// DockerService defines a set of high-level Docker operations.
type DockerService interface {
	IsDaemonRunning() (bool, error)
	IsComposeInstalled() (bool, error)
	ComposeUp(composeFilePath string, daemonMode bool, stdout, stderr io.Writer) error
}

// cliDockerService is the implementation of DockerService that uses the docker command line.
type cliDockerService struct{}

func NewCliDockerService() DockerService {
	return &cliDockerService{}
}

func (s *cliDockerService) IsDaemonRunning() (bool, error) {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *cliDockerService) IsComposeInstalled() (bool, error) {
	cmd := exec.Command("docker", "compose", "--version")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *cliDockerService) ComposeUp(composeFilePath string, daemonMode bool, stdout, stderr io.Writer) error {
	args := []string{"compose", "-f", composeFilePath, "up"}
	if daemonMode {
		args = append(args, "-d")
	}
	cmd := exec.Command("docker", args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run docker compose: %w", err)
	}

	return nil
}
