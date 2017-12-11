package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type globalData struct {
	position position
	bestMove string
}

type globalOptions struct {
	debug bool
}

type analysisOptions struct {
	searchMode  string // one of "infinite", "depth", "nodes"
	searchMoves []string
	ponder      bool
	wtime       int
	btime       int
	winc        int
	binc        int
	movestogo   int
	depth       int
	nodes       int
	movesToMate int
	movetime    int
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
	case "go":
		startAnalysis(args)
	case "stop":
		stopAnalysis()
	case "ponderhit":
		ponderHit()
	case "quit":
		os.Exit(0)
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
	// can send options here
	sendCommand("uciok")
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

	if args[1] == "startpos" {
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

func startAnalysis(args []string) {
	var options analysisOptions

	// add ponder but no idea what it means

	if argumentPresent("infinite", args) != -1 {
		options.searchMode = "infinite"
	} else if argumentPresent("depth", args) != -1 {
		options.searchMode = "depth"
		options.depth, _ = strconv.Atoi(args[argumentPresent("depth", args)+1])
	} else if argumentPresent("movetime", args) != -1 {
		options.searchMode = "movetime"
		options.movetime, _ = strconv.Atoi(args[argumentPresent("depth", args)+1])
	} else if argumentPresent("nodes", args) != -1 {
		options.searchMode = "nodes"
		options.movetime, _ = strconv.Atoi(args[argumentPresent("nodes", args)+1])
	} else if argumentPresent("mate", args) != -1 {
		options.searchMode = "mate"
		options.movesToMate, _ = strconv.Atoi(args[argumentPresent("mate", args)+1])

	}

	if argumentPresent("searchmoves", args) != -1 {
		var moves []string
		for i := argumentPresent("searchmoves", args) + 1; isAlgebraic(args[i]); i++ {
			moves = append(moves, args[i])
		}
		options.searchMoves = moves
	}

	if argumentPresent("wtime", args) != -1 {
		options.wtime, _ = strconv.Atoi(args[argumentPresent("wtime", args)+1])
	}

	if argumentPresent("btime", args) != -1 {
		options.btime, _ = strconv.Atoi(args[argumentPresent("btime", args)+1])
	}

	if argumentPresent("winc", args) != -1 {
		options.winc, _ = strconv.Atoi(args[argumentPresent("winc", args)+1])
	}

	if argumentPresent("binc", args) != -1 {
		options.binc, _ = strconv.Atoi(args[argumentPresent("binc", args)+1])
	}

	if argumentPresent("movestogo", args) != -1 {
		options.movestogo, _ = strconv.Atoi(args[argumentPresent("movestogo", args)+1])
	}

	engineData.bestMove = getBestMove(engineData.position)
}

func stopAnalysis() {
	sendCommand("bestmove", engineData.bestMove)
}

func ponderHit() {
	// impl
	return
}

func isAlgebraic(move string) bool {
	return true // implement this
}

func argumentPresent(arg string, args []string) int {
	for i, a := range args {
		if arg == a {
			return i
		}
	}

	return -1
}

// can add info and options later
