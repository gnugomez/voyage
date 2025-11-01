package command

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDeployCommandParametersParser(t *testing.T) {
	t.Run("Loads from valid JSON config file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		configContent := `{
			"repo": "my-repo",
			"branch": "main",
			"outPath": "/tmp/voyage",
			"remoteComposePaths": ["docker-compose.yml"],
			"force": true,
			"logLevel": "debug"
		}`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		args := []string{"-config", configPath}
		params, _, err := deployCommandParametersParser(args)

		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		expected := DeployCommandParameters{
			Repo:               "my-repo",
			Branch:             "main",
			OutPath:            "/tmp/voyage",
			RemoteComposePaths: []string{"docker-compose.yml"},
			Force:              true,
			BaseParameters:     BaseParameters{LogLevel: "debug"},
		}

		// Can't directly compare RemoteComposePaths because one is stringSlice and other is []string
		if params.Repo != expected.Repo ||
			params.Branch != expected.Branch ||
			params.OutPath != expected.OutPath ||
			!reflect.DeepEqual(params.RemoteComposePaths, expected.RemoteComposePaths) ||
			params.Force != expected.Force ||
			params.LogLevel != expected.LogLevel {
			t.Errorf("Parsed params do not match expected.\nGot:      %+v\nExpected: %+v", params, expected)
		}
	})

	t.Run("Loads valid YAML config file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.yaml")
		configContent := `
repo: my-repo
branch: main
outPath: /tmp/voyage
remoteComposePaths:
  - docker-compose.yml
force: true
logLevel: debug
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		args := []string{"-config", configPath}
		params, _, err := deployCommandParametersParser(args)

		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		expected := DeployCommandParameters{
			Repo:               "my-repo",
			Branch:             "main",
			OutPath:            "/tmp/voyage",
			RemoteComposePaths: []string{"docker-compose.yml"},
			Force:              true,
			BaseParameters:     BaseParameters{LogLevel: "debug"},
		}

		// Can't directly compare RemoteComposePaths because one is stringSlice and other is []string
		if params.Repo != expected.Repo ||
			params.Branch != expected.Branch ||
			params.OutPath != expected.OutPath ||
			!reflect.DeepEqual(params.RemoteComposePaths, expected.RemoteComposePaths) ||
			params.Force != expected.Force ||
			params.LogLevel != expected.LogLevel {
			t.Errorf("Parsed params do not match expected.\nGot:      %+v\nExpected: %+v", params, expected)
		}
	})

	t.Run("Returns error if required flags are missing without config", func(t *testing.T) {
		args := []string{"-r", "my-repo"} // Missing other required flags
		_, _, err := deployCommandParametersParser(args)
		if err == nil {
			t.Fatal("Expected an error for missing flags, but got nil")
		}
	})

	t.Run("Parses flags correctly without config", func(t *testing.T) {
		args := []string{
			"-r", "my-repo-from-flags",
			"-b", "develop",
			"-o", "/tmp/flags",
			"-c", "service1/docker-compose.yml",
		}

		params, _, err := deployCommandParametersParser(args)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if params.Repo != "my-repo-from-flags" {
			t.Errorf("Expected repo to be 'my-repo-from-flags', got '%s'", params.Repo)
		}
		if params.Branch != "develop" {
			t.Errorf("Expected branch to be 'develop', got '%s'", params.Branch)
		}
	})

	t.Run(("Merges config file and flags, with flags taking precedence"), func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		configContent := `{
			"repo": "my-repo",
			"branch": "main",
			"outPath": "/tmp/voyage",
			"remoteComposePaths": ["docker-compose.yml"],
			"force": false,
			"logLevel": "info"
		}`

		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		args := []string{
			"-config", configPath,
			"-b", "feature-branch", // Override branch
			"-f", // Override force to true
		}

		params, _, err := deployCommandParametersParser(args)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if params.Repo != "my-repo" {
			t.Errorf("Expected repo to be 'my-repo', got '%s'", params.Repo)
		}
		if params.Branch != "feature-branch" {
			t.Errorf("Expected branch to be 'feature-branch', got '%s'", params.Branch)
		}
		if !params.Force {
			t.Errorf("Expected force to be true, got false")
		}
		if params.LogLevel != "info" {
			t.Errorf("Expected logLevel to be 'info', got '%s'", params.LogLevel)
		}
	})

	t.Run("Returns error for non-existent config file", func(t *testing.T) {
		args := []string{"-config", "/path/to/non-existent-config.json"}
		_, _, err := deployCommandParametersParser(args)
		if err == nil {
			t.Fatal("Expected an error for non-existent config file, but got nil")
		}
	})
}
