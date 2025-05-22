package main

import (
	"flag"
	"path/filepath"

	"github.com/gnugomez/voyage/docker"
	"github.com/gnugomez/voyage/git"
	"github.com/gnugomez/voyage/log"
)

type Parameters struct {
	repo        string
	composePath string
	branch      string
	outPath     string
	force       bool
	daemon      bool
	logLevel    string
}

func main() {
	params := Parameters{}

	flag.StringVar(&params.repo, "r", "", "repository name")
	flag.StringVar(&params.composePath, "c", "", "path to docker-compose.yml")
	flag.StringVar(&params.branch, "b", "", "branch name")
	flag.StringVar(&params.outPath, "o", "", "out path")
	flag.BoolVar(&params.force, "f", false, "force deployment even if no changes detected")
	flag.BoolVar(&params.daemon, "d", true, "run docker compose in daemon mode")
	flag.StringVar(&params.logLevel, "l", "info", "log level (debug, info, error, fatal)")
	flag.Parse()

	log.SetLogger(log.CreateDefaultLogger(log.ParseLogLevel(params.logLevel)))

	if !requiredParamsPresent(params) {
		return
	}

	log.Debug("Running command with parameters", "repo", params.repo, "branch", params.branch, "compose-path", params.composePath, "out-path", params.outPath)

	subDir := filepath.Dir(params.composePath)
	if subDir == "." {
		subDir = "" // root of repo
	}

	hasChanges, err := git.SyncRepo(params.repo, params.branch, params.outPath, subDir)

	if err != nil {
		log.Fatal("Error syncing repository", "error", err)
		return
	}

	if hasChanges || params.force {
		if hasChanges {
			log.Info("Running docker-compose up")
		} else if params.force {
			log.Info("Force flag set, running docker-compose up")
		}

		err := docker.DeployCompose(filepath.Join(params.outPath, params.composePath), params.daemon)

		if err != nil {
			log.Fatal("Error running docker-compose up", "error", err)
			return
		}
	} else {
		log.Info("No changes detected, skipping docker-compose up")
	}
}

func requiredParamsPresent(params Parameters) bool {
	requiredParams := map[string]string{
		"r": params.repo,
		"c": params.composePath,
		"b": params.branch,
		"o": params.outPath,
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
