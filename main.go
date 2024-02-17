package main

import (
	"fmt"
	"io"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cmdMap = map[string]cliCommand{}

func commandHelp() error {
	for _, cmd := range cmdMap {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
func commandExit() error {
	return io.EOF
}
func init() {
	cmdMap["help"] = cliCommand{
		name:        "help",
		description: "Show help message",
		callback:    commandHelp,
	}
	cmdMap["exit"] = cliCommand{
		name:        "exit",
		description: "exit CLI",
		callback:    commandExit,
	}
}

func main() {

	startRepl()
}
