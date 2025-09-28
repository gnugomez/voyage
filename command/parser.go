package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type PrintUsageFunc func()

type missingParamsError struct {
	params []string
}

func (e *missingParamsError) Error() string {
	return fmt.Sprintf("missing required parameters: %s", strings.Join(e.params, ", "))
}

func deployCommandParametersParser(args []string) (DeployCommandParameters, PrintUsageFunc, error) {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nYou can provide parameters either via command-line flags or a JSON configuration file.\n")
		fmt.Fprintf(fs.Output(), "\nExample:\n")
		fmt.Fprintf(fs.Output(), "  voyage deploy -r my-repo -c docker-compose.yml -b main -o /tmp/deploy\n")
		fmt.Fprintf(fs.Output(), "  voyage deploy -config my-config.json\n")
	}

	configPath := fs.String("config", "", "path to a JSON configuration file")
	repo := fs.String("r", "", "repository name")
	var composePaths stringSlice
	fs.Var(&composePaths, "c", "path to docker-compose.yml (can be specified multiple times)")
	branch := fs.String("b", "", "branch name")
	outPath := fs.String("o", "", "out path")
	force := fs.Bool("f", false, "force deployment even if no changes detected")
	logLevel := fs.String("l", "info", "log level (debug, info, error, fatal)")

	if err := fs.Parse(args); err != nil {
		return DeployCommandParameters{}, fs.Usage, err
	}

	params := DeployCommandParameters{}

	if *configPath != "" {
		file, err := os.ReadFile(*configPath)
		if err != nil {
			return DeployCommandParameters{}, fs.Usage, fmt.Errorf("error reading config file %s: %w", *configPath, err)
		}
		if err := json.Unmarshal(file, &params); err != nil {
			return DeployCommandParameters{}, fs.Usage, fmt.Errorf("error parsing config file %s: %w", *configPath, err)
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
		return DeployCommandParameters{}, fs.Usage, &missingParamsError{params: missingParams}
	}

	return params, fs.Usage, nil
}
