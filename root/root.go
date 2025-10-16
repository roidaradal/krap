package root

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/str"
	"golang.org/x/term"
)

const (
	allCommands string = "*"
	cmdHelp     string = "help"
	cmdExit     string = "exit"
	cmdSearch   string = "cmd"
	cmdGlue     string = "/"
)

var cmdMap = map[string]*CmdConfig{}

var (
	errInvalidCommand    = errors.New("invalid command")
	errInvalidParamCount = errors.New("invalid param count")
	getHelp              = fmt.Sprintf("Type `%s` for list of commands, `%s <keyword>` to search for command", cmdHelp, cmdSearch)
	helpSkipCommands     = []string{cmdHelp, cmdExit, cmdSearch}
)

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

// Sets the command map
func SetCommandMap(commands map[string]*CmdConfig) {
	cmdMap = commands
}

// Authenticate Root account in command-line app
func Authenticate(authFn func(string) error) {
	fmt.Print("Enter password: ")
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Println()
	password := strings.TrimSpace(string(pwd))
	err = authFn(password)
	if err != nil {
		log.Fatal("Root authentication failed")
	}
}

// Root's main loop
func MainLoop(onExit func()) {
	var err error
	var line, command string
	var params []string

	fmt.Println("Commands:", len(cmdMap))
	fmt.Printf("Root: type `%s` for list of commands, `%s` to close\n", cmdHelp, cmdExit)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		line, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		command, params = getCommandParams(line)
		if command == "" {
			continue
		}
		switch command {
		case cmdExit:
			onExit()
			return
		case cmdHelp:
			if len(params) == 0 {
				command = allCommands
			} else {
				command = params[0]
			}
			displayHelp(command)
		case cmdSearch:
			keyword := params[0]
			searchCommand(keyword)
		default:
			c := cmdMap[command]
			c.Handler(params)
		}
	}
}

// Check if the command exists and if it meets the min parameter count
func validateCommandParams(command string, params []string) error {
	cfg, ok := cmdMap[command]
	if !ok {
		return errInvalidCommand
	}
	if len(params) < cfg.MinParams {
		return errInvalidParamCount
	}
	return nil
}

// Get command and params from line
func getCommandParams(line string) (string, []string) {
	if strings.TrimSpace(line) == "" {
		fmt.Println(getHelp)
		return "", nil
	}
	args := str.SpaceSplit(line)
	command, params := args[0], args[1:]
	command = strings.ToLower(command)
	err := validateCommandParams(command, params)
	if err != nil {
		fmt.Println("Error:", err)
		if errors.Is(err, errInvalidCommand) {
			fmt.Println(getHelp)
		} else if errors.Is(err, errInvalidParamCount) {
			displayHelp(command)
		}
		return "", nil
	}
	return command, params
}

// Display help list
func displayHelp(targetCommand string) {
	targetCommand = strings.ToLower(targetCommand)
	if _, ok := cmdMap[targetCommand]; !ok && targetCommand != allCommands {
		fmt.Println("Error: unknown command: ", targetCommand)
		fmt.Println(getHelp)
		return
	}
	fmt.Println("Usage: <command> <params>")
	fmt.Println("\nCommands and params:")

	commands := dict.Keys(cmdMap)
	sort.Strings(commands)
	for _, command := range commands {
		if slices.Contains(helpSkipCommands, command) {
			continue
		}
		cfg := cmdMap[command]
		if targetCommand == allCommands || targetCommand == command {
			fmt.Printf("%-30s\t%s\n", command, cfg.Docs)
		}
	}
}

// Search for command keyword
func searchCommand(keyword string) {
	keyword = strings.ToLower(keyword)
	commands := dict.Keys(cmdMap)
	slices.Sort(commands)
	if keyword == allCommands {
		stems := ds.NewSet[string]()
		for _, command := range commands {
			if slices.Contains(helpSkipCommands, command) {
				continue
			}
			stem := str.CleanSplit(command, cmdGlue)[0]
			stems.Add(stem)
		}
		heads := stems.Items()
		slices.Sort(heads)
		for _, head := range heads {
			fmt.Println(head)
		}
	} else {
		for _, command := range commands {
			if slices.Contains(helpSkipCommands, command) {
				continue
			}
			if strings.Contains(command, keyword) {
				cfg := cmdMap[command]
				fmt.Printf("%-30s\t%s\n", command, cfg.Docs)
			}
		}
	}
}
