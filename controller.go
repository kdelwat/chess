package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
