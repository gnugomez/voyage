package main

import (
	"os"

	"github.com/gnugomez/voyage/command"
	"github.com/gnugomez/voyage/log"
)

func main() {
	if len(os.Args) < 2 {
		log.Error("No command provided. Usage: voyage <command> [options]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)

	createCommand, ok := command.Commands[cmd]

	if !ok {
		log.Error("Unknown command", "command", cmd)
		os.Exit(1)
	}

	command := createCommand()
	baseParams := command.GetBaseParameters()

	log.SetLogger(log.CreateDefaultLogger(log.ParseLogLevel(baseParams.LogLevel)))

	command.Handle()
}
