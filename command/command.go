package command

import (
	"flag"
	"gnugomez/voyage/log"
)

type commandParameters struct {
	repo        string
	composePath string
	branch      string
}

func Run() {
	params := commandParameters{}

	flag.StringVar(&params.repo, "repo", "", "repository name")
	flag.StringVar(&params.composePath, "compose-path", "", "path to docker-compose.yml")
	flag.StringVar(&params.branch, "branch", "", "branch name")
	flag.Parse()

	log.Debug("Running command with parameters:", params)
}
