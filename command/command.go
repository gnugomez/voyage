package command

type BaseParameters struct {
	LogLevel string
}

type Command struct {
	Handle            func()
	GetBaseParameters func() BaseParameters
}

var Commands = map[string]func() *Command{
	"deploy": createDeployCommand,
}
