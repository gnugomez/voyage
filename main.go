package main

import (
	"flag"
	"gnugomez/voyage/log"
)

type Parameters struct {
	repo        string
	composePath string
	branch      string
	outPath     string
}

func main() {
	log.SetLogger(log.CreateDefaultLogger(log.DebugLevel))

	params := Parameters{}

	flag.StringVar(&params.repo, "repo", "", "repository name")
	flag.StringVar(&params.composePath, "compose-path", "", "path to docker-compose.yml")
	flag.StringVar(&params.branch, "branch", "", "branch name")
	flag.StringVar(&params.outPath, "out-path", "", "out path")
	flag.Parse()

	if !areRequiredParamsPresent(params) {
		return
	}

	log.Debug("Running command with parameters", "repo", params.repo, "branch", params.branch, "compose-path", params.composePath, "out-path", params.outPath)

	SyncRepo(params.repo, params.branch, params.outPath)
}

func areRequiredParamsPresent(params Parameters) bool {
	requiredParams := map[string]string{
		"repo":         params.repo,
		"compose-path": params.composePath,
		"branch":       params.branch,
		"out-path":     params.outPath,
	}

	for name, value := range requiredParams {
		if value == "" {
			log.Error("Missing required parameter", "parameter", name)
			flag.Usage()
			return false
		}
	}
	return true
}
