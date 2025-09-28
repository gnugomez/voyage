package command

import (
	"errors"
	"testing"
)

// --- Mocks ---

type mockSyncer struct {
	SyncFunc func(subDirs []string) ([]string, error)
}

func (m *mockSyncer) Sync(subDirs []string) ([]string, error) {
	if m.SyncFunc != nil {
		return m.SyncFunc(subDirs)
	}
	return nil, nil
}

type mockDeployer struct {
	DeployComposeFunc func(targetPath string, daemonMode bool) error
}

func (m *mockDeployer) DeployCompose(targetPath string, daemonMode bool) error {
	if m.DeployComposeFunc != nil {
		return m.DeployComposeFunc(targetPath, daemonMode)
	}
	return nil
}

func TestDeployCommand_Handle(t *testing.T) {
	t.Run("Changes detected, should deploy", func(t *testing.T) {
		syncer := &mockSyncer{}
		deployer := &mockDeployer{}

		dc := &deployCommand{
			params: DeployCommandParameters{
				Repo:               "repo",
				Branch:             "main",
				OutPath:            "/tmp",
				RemoteComposePaths: []string{"app1/docker-compose.yml"},
			},
			syncer:   syncer,
			deployer: deployer,
		}

		syncer.SyncFunc = func(subDirs []string) ([]string, error) {
			return []string{"app1"}, nil
		}

		deployerCalled := false
		deployer.DeployComposeFunc = func(targetPath string, daemonMode bool) error {
			deployerCalled = true
			return nil
		}

		dc.Handle()

		if !deployerCalled {
			t.Error("Deployer.DeployCompose should have been called, but it wasn't")
		}
	})

	t.Run("No changes detected, should not deploy", func(t *testing.T) {
		syncer := &mockSyncer{}
		deployer := &mockDeployer{}

		dc := &deployCommand{
			params: DeployCommandParameters{
				Repo:               "repo",
				Branch:             "main",
				OutPath:            "/tmp",
				RemoteComposePaths: []string{"app1/docker-compose.yml"},
				Force:              false, // Explicitly false
			},
			syncer:   syncer,
			deployer: deployer,
		}

		syncer.SyncFunc = func(subDirs []string) ([]string, error) {
			return []string{}, nil // No changes
		}

		deployerCalled := false
		deployer.DeployComposeFunc = func(targetPath string, daemonMode bool) error {
			deployerCalled = true
			return nil
		}

		dc.Handle()

		if deployerCalled {
			t.Error("Deployer.DeployCompose should not have been called, but it was")
		}
	})

	t.Run("No changes, but force flag is set, should deploy", func(t *testing.T) {
		syncer := &mockSyncer{}
		deployer := &mockDeployer{}

		dc := &deployCommand{
			params: DeployCommandParameters{
				Repo:               "repo",
				Branch:             "main",
				OutPath:            "/tmp",
				RemoteComposePaths: []string{"app1/docker-compose.yml"},
				Force:              true, // Force is true
			},
			syncer:   syncer,
			deployer: deployer,
		}

		syncer.SyncFunc = func(subDirs []string) ([]string, error) {
			return []string{}, nil // No changes
		}

		deployerCalled := false
		deployer.DeployComposeFunc = func(targetPath string, daemonMode bool) error {
			deployerCalled = true
			return nil
		}

		dc.Handle()

		if !deployerCalled {
			t.Error("Deployer.DeployCompose should have been called with force flag, but it wasn't")
		}
	})

	t.Run("Sync fails, should not deploy", func(t *testing.T) {
		syncer := &mockSyncer{}
		deployer := &mockDeployer{}

		dc := &deployCommand{
			params: DeployCommandParameters{
				Repo:               "repo",
				Branch:             "main",
				OutPath:            "/tmp",
				RemoteComposePaths: []string{"app1/docker-compose.yml"},
			},
			syncer:   syncer,
			deployer: deployer,
		}

		syncer.SyncFunc = func(subDirs []string) ([]string, error) {
			return nil, errors.New("sync failed")
		}

		deployerCalled := false
		deployer.DeployComposeFunc = func(targetPath string, daemonMode bool) error {
			deployerCalled = true
			return nil
		}

		dc.Handle()

		if deployerCalled {
			t.Error("Deployer.DeployCompose should not have been called when sync fails, but it was")
		}
	})
}
