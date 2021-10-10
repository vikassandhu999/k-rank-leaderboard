package main

import (
	"bufio"
	"os"
	"strings"
)

type TerminalI interface {
	run(commander *CommanderI)
}

type Resolver func(...interface{}) interface{}

type Terminal struct {
	currentLDB string
}

func createLDB(cmd *CommanderI, args []string) {
	name := strings.Trim(args[1], " ")
	(*cmd).CreateLDB(name)
}

func (trmnl Terminal) run(cmd *CommanderI) {

	reader := bufio.NewReader(os.Stdin)

	for {
		textInput, _ := reader.ReadString('\n')
		args := strings.Split(textInput, " ")
		switch args[0] {
		case "CREATE_LDB":
			createLDB(cmd, args)
		case "EXIT":
			break
		}

	}

}
