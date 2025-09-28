package docker

import (
	"fmt"
	"io"
	"os"
)

// Deployer handles the logic for deploying a docker-compose application.
type Deployer struct {
	dockerService DockerService
	fileExists    func(path string) bool
	stdout        io.Writer
	stderr        io.Writer
}

// NewDeployer creates a new Deployer with default dependencies.
func NewDeployer() *Deployer {
	return &Deployer{
		dockerService: NewCliDockerService(),
		fileExists:    osFileExists,
		stdout:        os.Stdout,
		stderr:        os.Stderr,
	}
}

// DeployCompose checks the environment and runs 'docker compose up'.
func (d *Deployer) DeployCompose(targetPath string, daemonMode bool) error {
	// Check docker availability
	if err := d.isDockerAvailable(); err != nil {
		return err
	}

	// Check target file
	if !d.fileExists(targetPath) {
		return fmt.Errorf("target path does not exist: %s", targetPath)
	}

	// Run compose
	return d.dockerService.ComposeUp(targetPath, daemonMode, d.stdout, d.stderr)
}

func (d *Deployer) isDockerAvailable() error {
	daemonRunning, err := d.dockerService.IsDaemonRunning()
	if err != nil || !daemonRunning {
		return fmt.Errorf("docker daemon is not running: %w", err)
	}

	composeInstalled, err := d.dockerService.IsComposeInstalled()
	if err != nil || !composeInstalled {
		return fmt.Errorf("docker compose is not installed: %w", err)
	}

	return nil
}

func osFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
