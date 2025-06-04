package command

import (
	"flag"
	"path/filepath"

	"github.com/gnugomez/voyage/docker"
	"github.com/gnugomez/voyage/git"
	"github.com/gnugomez/voyage/log"
)

type DeployCommandParameters struct {
	BaseParameters
	Repo        string
	ComposePath string
	Branch      string
	OutPath     string
	Force       bool
	Daemon      bool
}

type deployCommand struct {
	params DeployCommandParameters
}

func (d *deployCommand) GetBaseParameters() BaseParameters {
	return d.params.BaseParameters
}

func (d *deployCommand) Handle() {
	if !requiredParamsPresent(d.params) {
		return
	}

	log.Debug("Running command with parameters", "repo", d.params.Repo, "branch", d.params.Branch, "compose-path", d.params.ComposePath, "out-path", d.params.OutPath)

	subDir := filepath.Dir(d.params.ComposePath)
	if subDir == "." {
		subDir = "" // root of repo
	}

	hasChanges, err := git.SyncRepo(d.params.Repo, d.params.Branch, d.params.OutPath, subDir)

	if err != nil {
		log.Fatal("Error syncing repository", "error", err)
		return
	}

	if hasChanges || d.params.Force {
		if hasChanges {
			log.Info("Running docker-compose up")
		} else if d.params.Force {
			log.Info("Force flag set, running docker-compose up")
		}

		err := docker.DeployCompose(filepath.Join(d.params.OutPath, d.params.ComposePath), d.params.Daemon)

		if err != nil {
			log.Fatal("Error running docker-compose up", "error", err)
			return
		}
	} else {
		log.Info("No changes detected, skipping docker-compose up")
	}
}

func createDeployCommand() *Command {
	params := deployCommandParametersParser()

	d := &deployCommand{
		params: params,
	}

	return &Command{
		Handle:            d.Handle,
		GetBaseParameters: d.GetBaseParameters,
	}
}

func deployCommandParametersParser() DeployCommandParameters {
	params := DeployCommandParameters{}
	flag.StringVar(&params.Repo, "r", "", "repository name")
	flag.StringVar(&params.ComposePath, "c", "", "path to docker-compose.yml")
	flag.StringVar(&params.Branch, "b", "", "branch name")
	flag.StringVar(&params.OutPath, "o", "", "out path")
	flag.BoolVar(&params.Force, "f", false, "force deployment even if no changes detected")
	flag.BoolVar(&params.Daemon, "d", true, "run docker compose in daemon mode")
	flag.StringVar(&params.LogLevel, "l", "info", "log level (debug, info, error, fatal)")
	flag.Parse()
	return params
}

// RequiredParamsPresent checks if all required parameters are present
func requiredParamsPresent(params DeployCommandParameters) bool {
	requiredParams := map[string]string{
		"r": params.Repo,
		"c": params.ComposePath,
		"b": params.Branch,
		"o": params.OutPath,
	}

	hasErrors := false

	for name, value := range requiredParams {
		if value == "" {
			log.Error("Missing required parameter", "parameter", name)
			hasErrors = true
		}
	}

	if hasErrors {
		flag.Usage()
		return false
	}

	return true
}
