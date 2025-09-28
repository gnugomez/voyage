package git

import (
	"errors"
	"reflect"
	"testing"
)

// mockGitService is a mock implementation of the GitService interface for testing.
type mockGitService struct {
	IsGitRepositoryFunc    func(path string) bool
	FetchFunc              func(path string) error
	IsBehindRemoteFunc     func(path, branch string) (bool, error)
	PullFunc               func(path, branch string) error
	CloneFunc              func(path, url, branch string) error
	HasChangesInSubdirFunc func(path, branch, subDir string) (bool, error)
}

func (m *mockGitService) IsGitRepository(path string) bool {
	if m.IsGitRepositoryFunc != nil {
		return m.IsGitRepositoryFunc(path)
	}
	return false
}

func (m *mockGitService) Fetch(path string) error {
	if m.FetchFunc != nil {
		return m.FetchFunc(path)
	}
	return nil
}

func (m *mockGitService) IsBehindRemote(path, branch string) (bool, error) {
	if m.IsBehindRemoteFunc != nil {
		return m.IsBehindRemoteFunc(path, branch)
	}
	return false, nil
}

func (m *mockGitService) Pull(path, branch string) error {
	if m.PullFunc != nil {
		return m.PullFunc(path, branch)
	}
	return nil
}

func (m *mockGitService) Clone(path, url, branch string) error {
	if m.CloneFunc != nil {
		return m.CloneFunc(path, url, branch)
	}
	return nil
}

func (m *mockGitService) HasChangesInSubdir(path, branch, subDir string) (bool, error) {
	if m.HasChangesInSubdirFunc != nil {
		return m.HasChangesInSubdirFunc(path, branch, subDir)
	}
	return false, nil
}

func TestSync(t *testing.T) {
	subDirs := []string{"app1", "app2"}

	t.Run("Clone flow", func(t *testing.T) {
		mock := &mockGitService{}
		repo := &Repository{
			URL:             "url",
			Branch:          "branch",
			OutPath:         "path",
			gitService:      mock,
			directoryExists: func(s string) bool { return false },
		}

		cloneCalled := false
		mock.CloneFunc = func(path, url, branch string) error {
			cloneCalled = true
			return nil
		}

		updated, err := repo.Sync(subDirs)
		if err != nil {
			t.Fatalf("Sync() returned an unexpected error: %v", err)
		}

		if !cloneCalled {
			t.Error("Expected Clone to be called, but it wasn't")
		}

		if !reflect.DeepEqual(updated, subDirs) {
			t.Errorf("Expected all subdirs to be updated on clone, got %v", updated)
		}
	})

	t.Run("No changes flow", func(t *testing.T) {
		mock := &mockGitService{}
		repo := &Repository{
			gitService:      mock,
			directoryExists: func(s string) bool { return true },
		}

		mock.IsGitRepositoryFunc = func(path string) bool { return true }
		mock.FetchFunc = func(path string) error { return nil }
		mock.HasChangesInSubdirFunc = func(path, branch, subDir string) (bool, error) { return false, nil }

		updated, err := repo.Sync(subDirs)
		if err != nil {
			t.Fatalf("Sync() returned an unexpected error: %v", err)
		}

		if len(updated) != 0 {
			t.Errorf("Expected no updated dirs, but got %v", updated)
		}
	})

	t.Run("Pull changes flow", func(t *testing.T) {
		mock := &mockGitService{}
		repo := &Repository{
			URL:             "url",
			Branch:          "branch",
			OutPath:         "path",
			gitService:      mock,
			directoryExists: func(s string) bool { return true },
		}

		mock.IsGitRepositoryFunc = func(path string) bool { return true }
		mock.FetchFunc = func(path string) error { return nil }
		mock.HasChangesInSubdirFunc = func(path, branch, subDir string) (bool, error) {
			// Only app1 has changes
			return subDir == "app1", nil
		}
		mock.IsBehindRemoteFunc = func(path, branch string) (bool, error) { return true, nil }

		pullCalled := false
		mock.PullFunc = func(path, branch string) error {
			pullCalled = true
			return nil
		}

		updated, err := repo.Sync(subDirs)
		if err != nil {
			t.Fatalf("Sync() returned an unexpected error: %v", err)
		}

		if !pullCalled {
			t.Error("Expected Pull to be called, but it wasn't")
		}

		if !reflect.DeepEqual(updated, []string{"app1"}) {
			t.Errorf("Expected updated dirs to be [app1], but got %v", updated)
		}
	})

	t.Run("Error on fetch", func(t *testing.T) {
		mock := &mockGitService{}
		repo := &Repository{
			gitService:      mock,
			directoryExists: func(s string) bool { return true },
		}

		mock.IsGitRepositoryFunc = func(path string) bool { return true }
		mock.FetchFunc = func(path string) error { return errors.New("fetch failed") }

		_, err := repo.Sync(subDirs)
		if err == nil {
			t.Fatal("Expected an error on fetch, but got nil")
		}
	})
}
