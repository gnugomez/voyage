package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnugomez/voyage/docker"
	"github.com/gnugomez/voyage/git"
	"github.com/gnugomez/voyage/log"
)

type Syncer interface {
	Sync(subDirs []string) ([]string, error)
}

type Deployer interface {
	DeployCompose(targetPath string, daemonMode bool) error
}

// stringSlice implements flag.Value interface for handling multiple string values
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type DeployCommandParameters struct {
	BaseParameters     `yaml:",inline"`
	Repo               string   `json:"repo" yaml:"repo"`
	Branch             string   `json:"branch" yaml:"branch"`
	OutPath            string   `json:"outPath" yaml:"outPath"`
	RemoteComposePaths []string `json:"remoteComposePaths" yaml:"remoteComposePaths"`
	Force              bool     `json:"force" yaml:"force"`
}

type deployCommand struct {
	params   DeployCommandParameters
	syncer   Syncer
	deployer Deployer
}

func (d *deployCommand) GetBaseParameters() BaseParameters {
	return d.params.BaseParameters
}

func (d *deployCommand) Handle() {
	// Lazy initialization of dependencies. In tests, these will be pre-filled with mocks.
	if d.syncer == nil {
		d.syncer = git.CreateRepository(d.params.Repo, d.params.Branch, d.params.OutPath)
	}
	if d.deployer == nil {
		d.deployer = docker.NewDeployer()
	}

	log.Debug("Running command with parameters", "repo", d.params.Repo, "branch", d.params.Branch, "remoteComposePaths", d.params.RemoteComposePaths, "out-path", d.params.OutPath)

	// Create map of subdirectories to compose paths for change detection
	subDirToComposePaths := make(map[string][]string)
	var subDirs []string

	for _, composePath := range d.params.RemoteComposePaths {
		subDir := filepath.Dir(composePath)
		if subDir == "." {
			subDir = "" // root of repo
		}

		// Add to map
		if _, exists := subDirToComposePaths[subDir]; !exists {
			subDirToComposePaths[subDir] = []string{}
			subDirs = append(subDirs, subDir)
		}
		subDirToComposePaths[subDir] = append(subDirToComposePaths[subDir], composePath)
	}

	updatedSubDirs, err := d.syncer.Sync(subDirs)
	if err != nil {
		log.Error("Error syncing repository", "error", err)
		return
	}

	var composePathsToDeploy []string

	if len(updatedSubDirs) > 0 {
		log.Info("Running docker-compose up for updated subdirectories", "updatedSubDirs", updatedSubDirs)
		// Collect compose files for updated subdirectories
		for _, updatedSubDir := range updatedSubDirs {
			composePaths := subDirToComposePaths[updatedSubDir]
			composePathsToDeploy = append(composePathsToDeploy, composePaths...)
		}
	} else if d.params.Force {
		log.Info("Force flag set, running docker-compose up for all compose files")
		composePathsToDeploy = d.params.RemoteComposePaths
	} else {
		log.Info("No changes detected, skipping docker-compose up")
	}

	// Deploy all collected compose files
	for _, composePath := range composePathsToDeploy {
		log.Info("Deploying compose file", "composePath", composePath)
		err := d.deployer.DeployCompose(filepath.Join(d.params.OutPath, composePath), true)
		if err != nil {
			log.Error("Error running docker-compose up", "error", err, "composePath", composePath)
			return
		}
	}
}

func createDeployCommand() *Command {
	params, printUsage, err := deployCommandParametersParser(os.Args[1:])

	if err != nil {
		var missingParamsErr *missingParamsError
		if errors.As(err, &missingParamsErr) {
			log.Error("Error parsing parameters", "error", err)
			printUsage()
			os.Exit(1)
		} else {
			log.Fatal("Error parsing parameters", "error", err)
		}
	}

	d := &deployCommand{
		params: params,
	}

	return &Command{
		Handle:            d.Handle,
		GetBaseParameters: d.GetBaseParameters,
	}
}
