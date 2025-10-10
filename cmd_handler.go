package krap

// Takes in list of parameters
type CmdHandler = func([]string)

type CmdConfig struct {
	Command   string
	MinParams int
	Docs      string
	Handler   CmdHandler
}

// Creates a new CmdConfig
func NewCommand(command string, minParams int, docs string, handler CmdHandler) *CmdConfig {
	return &CmdConfig{command, minParams, docs, handler}
}

// Creates a new map of command => CmdConfigs
func NewCommandMap(cfgs ...*CmdConfig) map[string]*CmdConfig {
	commands := make(map[string]*CmdConfig)
	for _, cfg := range cfgs {
		commands[cfg.Command] = cfg
	}
	return commands
}
