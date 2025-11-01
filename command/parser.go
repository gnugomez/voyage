package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Constants for default values and configuration
const (
	defaultLogLevel = "info"
)

// PrintUsageFunc represents a function that prints command usage information
type PrintUsageFunc func()

// missingParamsError represents an error when required parameters are missing
type missingParamsError struct {
	params []string
}

// Error implements the error interface for missingParamsError
func (e *missingParamsError) Error() string {
	if len(e.params) == 1 {
		return fmt.Sprintf("missing required parameter: %s", e.params[0])
	}
	return fmt.Sprintf("missing required parameters: %s", strings.Join(e.params, ", "))
}

// deployCommandParametersParser parses command line arguments and configuration file
// to create DeployCommandParameters for the deploy command.
//
// It supports both command-line flags and JSON configuration files, with command-line
// flags taking precedence over configuration file values.
//
// Returns the parsed parameters, a usage function, and any error encountered.
func deployCommandParametersParser(args []string) (DeployCommandParameters, PrintUsageFunc, error) {
	fs := setupDeployFlags()

	if err := fs.Parse(args); err != nil {
		return DeployCommandParameters{}, fs.Usage, err
	}

	params, err := loadConfigFromFile(fs)
	if err != nil {
		return DeployCommandParameters{}, fs.Usage, err
	}

	params = overrideWithFlags(fs, params)

	if err := validateParameters(params); err != nil {
		return DeployCommandParameters{}, fs.Usage, err
	}

	return params, fs.Usage, nil
}

// setupDeployFlags creates and configures the flag set for deploy command
func setupDeployFlags() *flag.FlagSet {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nYou can provide parameters either via command-line flags or a JSON configuration file.\n")
		fmt.Fprintf(fs.Output(), "\nExample:\n")
		fmt.Fprintf(fs.Output(), "  voyage deploy -r my-repo -c docker-compose.yml -b main -o /tmp/deploy\n")
		fmt.Fprintf(fs.Output(), "  voyage deploy -config my-config.json\n")
	}

	fs.String("config", "", "path to a JSON configuration file")
	fs.String("r", "", "repository name")
	fs.Var(&stringSlice{}, "c", "path to docker-compose.yml (can be specified multiple times)")
	fs.String("b", "", "branch name")
	fs.String("o", "", "out path")
	fs.Bool("f", false, "force deployment even if no changes detected")
	fs.String("l", defaultLogLevel, "log level (debug, info, error, fatal)")

	return fs
}

// loadConfigFromFile loads parameters from configuration file if provided
func loadConfigFromFile(fs *flag.FlagSet) (DeployCommandParameters, error) {
	params := DeployCommandParameters{}
	configPath := fs.Lookup("config").Value.String()

	if configPath != "" {
		file, err := os.ReadFile(configPath)
		if err != nil {
			return params, fmt.Errorf("error reading config file %s: %w", configPath, err)
		}

		if strings.HasSuffix(configPath, ".json") {
			if err := json.Unmarshal(file, &params); err != nil {
				return params, fmt.Errorf("error parsing JSON config file %s: %w", configPath, err)
			}
		} else if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
			if err := yaml.Unmarshal(file, &params); err != nil {
				return params, fmt.Errorf("error parsing YAML config file %s: %w", configPath, err)
			}
		} else {
			return params, fmt.Errorf("unsupported config file format: %s", configPath)
		}

	}

	return params, nil
}

// overrideWithFlags applies command-line flag values to override config file values
func overrideWithFlags(fs *flag.FlagSet, params DeployCommandParameters) DeployCommandParameters {
	// Override with command-line flags if they were provided
	if repo := fs.Lookup("r").Value.String(); repo != "" {
		params.Repo = repo
	}

	if composePathsFlag := fs.Lookup("c"); composePathsFlag != nil {
		if composePaths, ok := composePathsFlag.Value.(*stringSlice); ok && len(*composePaths) > 0 {
			params.RemoteComposePaths = []string(*composePaths)
		}
	}

	if branch := fs.Lookup("b").Value.String(); branch != "" {
		params.Branch = branch
	}

	if outPath := fs.Lookup("o").Value.String(); outPath != "" {
		params.OutPath = outPath
	}

	// Handle boolean flag - check if it was explicitly set
	if forceFlag := fs.Lookup("f"); forceFlag.Value.String() == "true" {
		params.Force = true
	}

	// Handle log level - always override if different from default
	if logLevel := fs.Lookup("l").Value.String(); logLevel != defaultLogLevel {
		params.LogLevel = logLevel
	} else if params.LogLevel == "" {
		params.LogLevel = defaultLogLevel
	}

	return params
}

// validateParameters validates that all required parameters are present
func validateParameters(params DeployCommandParameters) error {
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
		return &missingParamsError{params: missingParams}
	}

	return nil
}
