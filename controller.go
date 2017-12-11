package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type globalOptions struct {
	debug bool
}

var engineOptions = globalOptions{
	debug: false,
}

func startEngine() {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	for err == nil {
		handleCommand(input)
		input, err = reader.ReadString('\n')
	}
}

func handleCommand(command string) {
	args := strings.Split(strings.TrimSpace(command), " ")

	fmt.Printf("Args: %v\n", args)
	switch args[0] {
	case "uci":
		handleUCI()
	case "debug":
		toggleDebug(args[1])
	}
}

func toggleDebug(setting string) {
	if setting == "on" {
		engineOptions.debug = true
	} else {
		engineOptions.debug = false
	}
}
func handleUCI() {
	sendCommand("id", "name", "Ultimate Engine")
	sendCommand("id", "author", "Cadel Watson")
}

func sendCommand(command string, args ...string) {
	tokens := append([]string{command}, args...)

	outputCommand := strings.Join(tokens, " ")

	fmt.Println(outputCommand)
}

func sendDebug(message string) {
	sendCommand("info", message)
}
