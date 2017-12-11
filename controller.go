package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type globalData struct {
	position position
}

type globalOptions struct {
	debug bool
}

var engineData globalData
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
	case "isready":
		sendCommand("readyok")
	case "setoption":
		// TODO
	case "register":
		// TODO
	case "ucinewgame":
		handleNewGame()
	case "position":
		setupPosition(args)
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

func handleNewGame() {
	// do cleanup
	sendCommand("isready")
}

func sendCommand(command string, args ...string) {
	tokens := append([]string{command}, args...)

	outputCommand := strings.Join(tokens, " ")

	fmt.Println(outputCommand)
}

func sendDebug(message string) {
	sendCommand("info", message)
}

func setupPosition(args []string) {
	var fen string

	if args[0] == "startpos" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" //from wikipedia
	} else {
		fen = args[0]
	}

	engineData.position = fromFen(fen)

	if len(args) > 1 {
		for _, m := range args[:1] {
			engineData.position = applyMove(engineData.position, m)
		}
	}
}

func applyMove(position position, move string) position {
	// todo apply algebraic move
	return position
}
