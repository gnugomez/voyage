package command

import (
	"flag"
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
	BaseParameters
	Repo               string
	Branch             string
	OutPath            string
	RemoteComposePaths []string
	Force              bool
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
	if !requiredParamsPresent(d.params) {
		return
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
	params := deployCommandParametersParser()

	d := &deployCommand{
		params: params,
		// Default implementations
		syncer:   git.CreateRepository(params.Repo, params.Branch, params.OutPath),
		deployer: docker.NewDeployer(),
	}

	return &Command{
		Handle:            d.Handle,
		GetBaseParameters: d.GetBaseParameters,
	}
}

func deployCommandParametersParser() DeployCommandParameters {
	params := DeployCommandParameters{}
	var composePaths stringSlice

	flag.StringVar(&params.Repo, "r", "", "repository name")
	flag.Var(&composePaths, "c", "path to docker-compose.yml (can be specified multiple times)")
	flag.StringVar(&params.Branch, "b", "", "branch name")
	flag.StringVar(&params.OutPath, "o", "", "out path")
	flag.BoolVar(&params.Force, "f", false, "force deployment even if no changes detected")
	flag.StringVar(&params.LogLevel, "l", "info", "log level (debug, info, error, fatal)")
	flag.Parse()

	params.RemoteComposePaths = []string(composePaths)
	return params
}

// RequiredParamsPresent checks if all required parameters are present
func requiredParamsPresent(params DeployCommandParameters) bool {
	var missingParams []string

	if params.Repo == "" {
		missingParams = append(missingParams, "r")
	}
	if len(params.RemoteComposePaths) == 0 {
		missingParams = append(missingParams, "c")
	}
	if params.Branch == "" {
		missingParams = append(missingParams, "b")
	}
	if params.OutPath == "" {
		missingParams = append(missingParams, "o")
	}

	if len(missingParams) > 0 {
		log.Error("Missing required parameters", "parameters", missingParams)
		flag.Usage()
		return false
	}

	return true
}
