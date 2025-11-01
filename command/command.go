package command

type BaseParameters struct {
	LogLevel string `json:"logLevel" yaml:"logLevel"`
}

type Command struct {
	Handle            func()
	GetBaseParameters func() BaseParameters
}

var Commands = map[string]func() *Command{
	"deploy": createDeployCommand,
}
