package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type globalData struct {
	position position
	bestMove string
}

type globalOptions struct {
	debug bool
	log   *os.File
}

type analysisOptions struct {
	searchMode  string // one of "infinite", "depth", "nodes", "movetime"
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
	errorLog := initialiseLog("/tmp/cadelChessLog.txt")

	defer errorLog.Close()

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	for err == nil {
		_, _ = errorLog.WriteString(input)
		errorLog.Sync()
		handleCommand(input)
		input, err = reader.ReadString('\n')
	}
}

func initialiseLog(filename string) *os.File {
	errorLog, err := os.Create("/tmp/cadelChessLog.txt")
	if err != nil {
		fmt.Printf("Couldn't create file %v\n", err)
	}

	engineOptions.log = errorLog

	return errorLog
}

func sendCommand(command string, args ...string) {
	tokens := append([]string{command}, args...)

	outputCommand := strings.Join(tokens, " ")

	_, _ = engineOptions.log.WriteString("> " + outputCommand + "\n")
	fmt.Println(outputCommand)
}

func sendDebug(message string) {
	sendCommand("info", message)
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
	sendCommand("id", "name", EngineName)
	sendCommand("id", "author", EngineAuthor)
	// can send options here
	sendCommand("uciok")
}

func handleNewGame() {
	// do cleanup
	sendCommand("isready")
}

func setupPosition(args []string) {
	var fen string

	if args[1] == "startpos" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	} else {
		fen = strings.TrimSpace(strings.Join(args[2:], " "))
	}

	engineData.position = fromFEN(fen)

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
		options.movetime, _ = strconv.Atoi(args[argumentPresent("movetime", args)+1])
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

	// run the actual analysis
	switch options.searchMode {
	case "movetime":
		fmt.Printf("Running with movetime %v\n", options.movetime)

		//ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(options.movetime))
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(options.movetime))

		ch := make(chan move)
		go runSearch(ctx, cancel, engineData.position, 1000, ch)
		go awaitBestMove(engineData.position, ch)

		//close(ch)

	}
	//engineData.bestMove = getBestMove(engineData.position)
	//sendCommand("bestmove", engineData.bestMove)

}

func awaitBestMove(position position, ch chan move) {
	for move := range ch {
		engineData.bestMove = toAlgebraic(position, move)
	}
	sendCommand("bestmove", engineData.bestMove)
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
