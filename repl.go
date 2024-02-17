package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		t, ok := cmdMap[text]

		if ok {
			tCall := t.callback()
			if tCall == io.EOF {
				break
			}
		}
		if !ok {
			err := errors.New("not a valid command")
			fmt.Println(err)
		}
	}
}
