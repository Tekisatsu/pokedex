package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		command := strings.Fields(text)
		t, ok := cmdMap[command[0]]
		if ok && len(command) == 1{
			tCall := t.callback(t.context)
			if tCall == io.EOF {
				break
			}
		}
		if ok && len(command) == 2 {
			t.context.Args = []string{command[1]}
			tCall := t.callback(t.context)
			if tCall != nil {
				fmt.Println(tCall)
			}
		}
		if !ok {
			err := errors.New("not a valid command")
			fmt.Println(err)
		}
	}
}
