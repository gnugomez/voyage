package main

import (
	"gnugomez/voyage/command"
	"gnugomez/voyage/log"
)

func main() {
	log.SetLogger(log.CreateDefaultLogger(log.DebugLevel))
	command.Run()
}
