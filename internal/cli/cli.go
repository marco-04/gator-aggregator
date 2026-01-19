package cli

import (
	"fmt"
	"os"

	"github.com/marco-04/gator-aggregator/internal/state"
)

type cmd struct {
	description string
	argNum      int
	usageStr    string
	callback    func(s *state.State, args []string) error
}

var availableCommands map[string]cmd

func validateArgs(args []string, argNum int) error {
	if len(args) < argNum {
		return fmt.Errorf("not enough arguments: got %d, expected %d", len(args), argNum)
	}
	return nil
}

func Dispatch() {
	args := os.Args

	if len(args) == 1 {
		fmt.Println("error: no command specified")
		printHelp(nil, nil)
		os.Exit(1)
	}
	
	cmdName := args[1]
	
	// Initialize arguments to callbacks
	var subArgs []string
	if len(args[1:]) == 1 {
		subArgs = nil
	} else {
		subArgs = args[2:]
	}

	cmd, ok := availableCommands[cmdName]
	if !ok {
		fmt.Printf("error: command %s does not exist\n", cmdName)
		printHelp(nil, nil)
		os.Exit(1)
	}

	// Validate number of arguments before passing them to callbacks
	if err := validateArgs(subArgs, cmd.argNum); err != nil {
		fmt.Printf("error: %v\n", err)
		printHelp(nil, nil)
		os.Exit(1)
	}

	// Initialize application state
	s, err := state.Init()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		printHelp(nil, nil)
	}

	if err := availableCommands[cmdName].callback(&s, subArgs); err != nil {
		fmt.Printf("error: %s: %v\n", cmdName, err)
		os.Exit(1)
	}
}

func printHelp(s *state.State, args []string) error {
	progName := os.Args[0]
	fmt.Printf("Usage: %s <cmd> [args]\nCommands:\n", progName)

	for k, v := range availableCommands {
		fmt.Printf(" - %s: %s\n   usage: %s %s\n\n", k, v.description, k, v.usageStr)
	}

	return nil
}

func init() {
	availableCommands = cmds
}
