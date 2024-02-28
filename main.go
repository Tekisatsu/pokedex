package main

import (
	"fmt"
	"io"
	"time"
	pokeapi "github.com/tekisatsu/pokedex/pokeApi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.CliContext) error
	context     *pokeapi.CliContext
}

var cmdMap = map[string]cliCommand{}

func commandHelp(_ *pokeapi.CliContext) error {
	for _, cmd := range cmdMap {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}
func commandExit(_ *pokeapi.CliContext) error {
	return io.EOF
}

func init() {
	shareContext := &pokeapi.CliContext{
		State: &pokeapi.JsonConfig{},
		Cache: pokeapi.NewCache(15 * time.Minute),
		Pokedex: make(map[string][]byte),
	}
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
		callback:    pokeapi.GetMapUrl,
		context:     shareContext,
	}
	cmdMap["mapb"] = cliCommand{
		name:        "mapb",
		description: "show the previous 20 area locations",
		callback:    pokeapi.GetPrevMapUrl,
		context:     shareContext,
	}
	cmdMap["encounter"] = cliCommand{
		name: "encounter",
		description: "show encounters in an area",
		callback: pokeapi.Encounter,
		context: shareContext,
	}
	cmdMap["pokemon"] = cliCommand{
		name: "pokemon",
		description: "add pokemon to your Pokedex",
		callback: pokeapi.Catch,
		context: shareContext,
	}
	cmdMap["inspect"] = cliCommand{
		name: "inspect",
		description: "show caught Pokemon stats",
		callback: pokeapi.Inspect,
		context: shareContext,
	}
	cmdMap["pokedex"] = cliCommand{
		name: "pokedex",
		description: "show all Pokemon in Pokedex",
		callback: pokeapi.Pokedex,
		context: shareContext,
	}
}

func main() {
	startRepl()
}
