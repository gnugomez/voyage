package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func deployCommandParametersParser(args []string) (DeployCommandParameters, error) {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)

	configPath := fs.String("config", "", "path to a JSON configuration file")
	repo := fs.String("r", "", "repository name")
	var composePaths stringSlice
	fs.Var(&composePaths, "c", "path to docker-compose.yml (can be specified multiple times)")
	branch := fs.String("b", "", "branch name")
	outPath := fs.String("o", "", "out path")
	force := fs.Bool("f", false, "force deployment even if no changes detected")
	logLevel := fs.String("l", "info", "log level (debug, info, error, fatal)")

	if err := fs.Parse(args); err != nil {
		return DeployCommandParameters{}, err
	}

	params := DeployCommandParameters{}

	if *configPath != "" {
		file, err := os.ReadFile(*configPath)
		if err != nil {
			return DeployCommandParameters{}, fmt.Errorf("error reading config file %s: %w", *configPath, err)
		}
		if err := json.Unmarshal(file, &params); err != nil {
			return DeployCommandParameters{}, fmt.Errorf("error parsing config file %s: %w", *configPath, err)
		}
	} else {
		params.Repo = *repo
		params.RemoteComposePaths = composePaths
		params.Branch = *branch
		params.OutPath = *outPath
		params.Force = *force
		params.LogLevel = *logLevel
	}

	// Validate the final parameters
	var missingParams []string
	if params.Repo == "" {
		missingParams = append(missingParams, "-r (repository)")
	}
	if len(params.RemoteComposePaths) == 0 {
		missingParams = append(missingParams, "-c (compose path)")
	}
	if params.Branch == "" {
		missingParams = append(missingParams, "-b (branch)")
	}
	if params.OutPath == "" {
		missingParams = append(missingParams, "-o (out path)")
	}

	if len(missingParams) > 0 {
		return DeployCommandParameters{}, fmt.Errorf("missing required parameters: %s", strings.Join(missingParams, ", "))
	}

	return params, nil
}
