package docker

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// mockDockerService is a mock implementation of the DockerService interface for testing.
type mockDockerService struct {
	IsDaemonRunningFunc    func() (bool, error)
	IsComposeInstalledFunc func() (bool, error)
	ComposeUpFunc          func(composeFilePath string, daemonMode bool, stdout, stderr io.Writer) error
}

func (m *mockDockerService) IsDaemonRunning() (bool, error) {
	if m.IsDaemonRunningFunc != nil {
		return m.IsDaemonRunningFunc()
	}
	return false, nil
}

func (m *mockDockerService) IsComposeInstalled() (bool, error) {
	if m.IsComposeInstalledFunc != nil {
		return m.IsComposeInstalledFunc()
	}
	return false, nil
}

func (m *mockDockerService) ComposeUp(composeFilePath string, daemonMode bool, stdout, stderr io.Writer) error {
	if m.ComposeUpFunc != nil {
		return m.ComposeUpFunc(composeFilePath, daemonMode, stdout, stderr)
	}
	return nil
}

func TestDeployer_DeployCompose(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		mock := &mockDockerService{}
		d := &Deployer{
			dockerService: mock,
			fileExists:    func(path string) bool { return true },
			stdout:        io.Discard,
			stderr:        io.Discard,
		}

		mock.IsDaemonRunningFunc = func() (bool, error) { return true, nil }
		mock.IsComposeInstalledFunc = func() (bool, error) { return true, nil }

		composeUpCalled := false
		mock.ComposeUpFunc = func(path string, daemon bool, stdout, stderr io.Writer) error {
			composeUpCalled = true
			return nil
		}

		err := d.DeployCompose("docker-compose.yml", false)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if !composeUpCalled {
			t.Error("Expected ComposeUp to be called, but it wasn't")
		}
	})

	t.Run("Docker daemon not running", func(t *testing.T) {
		mock := &mockDockerService{}
		d := &Deployer{dockerService: mock}

		mock.IsDaemonRunningFunc = func() (bool, error) { return false, errors.New("daemon error") }

		if err := d.DeployCompose("path", false); err == nil {
			t.Fatal("Expected an error, but got nil")
		}
	})

	t.Run("Compose file not found", func(t *testing.T) {
		mock := &mockDockerService{}
		d := &Deployer{
			dockerService: mock,
			fileExists:    func(path string) bool { return false },
		}

		mock.IsDaemonRunningFunc = func() (bool, error) { return true, nil }
		mock.IsComposeInstalledFunc = func() (bool, error) { return true, nil }

		if err := d.DeployCompose("path", false); err == nil {
			t.Fatal("Expected an error, but got nil")
		}
	})

	t.Run("ComposeUp fails", func(t *testing.T) {
		mock := &mockDockerService{}
		d := &Deployer{
			dockerService: mock,
			fileExists:    func(path string) bool { return true },
			stdout:        &bytes.Buffer{},
			stderr:        &bytes.Buffer{},
		}

		mock.IsDaemonRunningFunc = func() (bool, error) { return true, nil }
		mock.IsComposeInstalledFunc = func() (bool, error) { return true, nil }
		mock.ComposeUpFunc = func(path string, daemon bool, stdout, stderr io.Writer) error {
			return errors.New("compose failed")
		}

		if err := d.DeployCompose("path", false); err == nil {
			t.Fatal("Expected an error, but got nil")
		}
	})
}
