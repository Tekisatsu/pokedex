package main

import (
	"fmt"
	"io"

	pokeapi "github.com/tekisatsu/pokedex/pokeApi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.JsonConfig) error
	state       *pokeapi.JsonConfig
}

var cmdMap = map[string]cliCommand{}

func commandHelp(_ *pokeapi.JsonConfig) error {
	for _, cmd := range cmdMap {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
func commandExit(_ *pokeapi.JsonConfig) error {
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
	cmdMap["map"] = cliCommand{
		name:        "map",
		description: "show the next 20 area locations",
		callback:    pokeapi.CommandMap,
		state:       &pokeapi.JsonConfig{},
	}
}

func main() {

	startRepl()
}
