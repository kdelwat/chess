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

// Declare the starting position of a game, in FEN notation.
var startPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Store the current position and current best move of the engine, used globally
// to ensure that search can continue in the background while the engine
// continue to recieve commands.
type globalData struct {
	position position
	bestMove string
}

var engineData globalData

/* Store the global options set via Universal Chess Interface commands for the
engine to follow during runtime.

searchMode can be one of "infinite", "depth", "nodes" and "movetime".

searchMoves is a list of moves to consider when searching, to the exclusion of
others.

ponder places the engine in ponder mode, which searches for the next move during
the opponent's turn.

wtime and btime control the amount of time each player has in the game. winc and
binc determine the time increment added to each player after each move.

depth, nodes, movesToMate, and movetime provide the options for the search
algorithm chosen.
*/
type analysisOptions struct {
	searchMode  string
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

// Start the engine, allowing UCI-compatible programs to communicate.
func startEngine() {
	// Recieve input from stdin.
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	// Continue to recieve each line of input.
	for err == nil {
		// Handle the input.
		result := handleCommand(input)

		// If the program has request to quit, return from the engine loop.
		if !result {
			break
		}

		input, err = reader.ReadString('\n')
	}
}

// Send a command, comprised of the base command and a slice of arguments, to
// stdout.
func sendCommand(command string, args ...string) {
	tokens := append([]string{command}, args...)

	outputCommand := strings.Join(tokens, " ")

	fmt.Println(outputCommand)
}

// Send a debug message, prefixed with "info"
func sendDebug(message string) {
	sendCommand("info", message)
}

// Handle a line of input, representing a UCI command. Return true if the engine
// should continue to recieve input, false otherwise.
func handleCommand(command string) bool {
	// Tokenise the command, splitting on spaces.
	args := strings.Split(strings.TrimSpace(command), " ")

	// Switch based on the first word in the command.
	switch args[0] {
	case "uci":
		handleUCI()
	case "debug":
		// TODO
	case "isready":
		sendCommand("readyok")
	case "ucinewgame":
		handleNewGame()
	case "position":
		setupPosition(args)
	case "go":
		startAnalysis(args)
	case "stop":
		stopAnalysis()
	case "quit":
		return false
	}

	return true
}

// Signal to the user that the engine is in UCI mode.
func handleUCI() {
	sendCommand("id", "name", EngineName)
	sendCommand("id", "author", EngineAuthor)
	sendCommand("uciok")
}

// Establish a new game in the engine data.
func handleNewGame() {
	engineData.position = fromFEN(startPosition)
	sendCommand("isready")
}

// Given a position string, set up the engine position.
func setupPosition(args []string) {
	var fen string

	// The position can be "startpos", meaning a game's initial starting
	// position, or a FEN specified by the interface.
	if args[1] == "startpos" {
		fen = startPosition
	} else {
		fen = strings.TrimSpace(strings.Join(args[2:], " "))
	}

	engineData.position = fromFEN(fen)

	// For each move specified after the initial FEN, apply the move.
	if len(args) > 1 {
		for _, m := range args[:1] {
			engineData.position = applyMove(engineData.position, m)
		}
	}
}

// Start the analysis, returning a best move. Currently, the time-limited and
// depth-limited analysis modes are supported.
func startAnalysis(args []string) {
	var options analysisOptions

	// Determine the analysis mode, based on the command given, and load
	// analysis settings.
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
	} else {
		options.searchMode = "depth"
		options.depth, _ = strconv.Atoi(args[argumentPresent("depth", args)])
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

	// Based on the selected mode, run the analysis.
	switch options.searchMode {
	case "movetime":
		// For the time-limited mode, create a context lasting the length of
		// time specified. This is used to limit searching.
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(options.movetime))

		// Create a channel for moves to be relayed through during iterative
		// deepening of the search.
		ch := make(chan move)

		// Create three goroutines. The first closes the move channel when the
		// context finishes. The second runs the actual search, placing best
		// moves into the move channel as they are found. The third waits for
		// the channel to close, then sends the current best move to the engine.
		go waitToClose(ctx, ch, cancel)
		go runSearch(ctx, engineData.position, 1000, ch)
		go awaitBestMove(engineData.position, ch)

	case "depth":
		// For the depth-limited mode, the search tree is simply searched to the
		// given depth.
		bestMove := search(engineData.position, options.depth, -100000, 100000)
		sendCommand("bestmove", toAlgebraic(engineData.position, bestMove))

	}
}

// When the given context is finished, close the channel and cancel the context.
func waitToClose(ctx context.Context, ch chan move, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			cancel()
			close(ch)
			return
		}
	}

}

// Recieve best moves from a channel until it is closed, then send the current
// best move to the interface.
func awaitBestMove(position position, ch chan move) {
	for move := range ch {
		engineData.bestMove = toAlgebraic(position, move)
	}
	sendCommand("bestmove", engineData.bestMove)
}

// Return the current best move of the engine immediately.
func stopAnalysis() {
	sendCommand("bestmove", engineData.bestMove)
}

// Loop through a slice of arguments, searching for a given string. If found,
// return its index. Otherwise, return -1.
func argumentPresent(arg string, args []string) int {
	for i, a := range args {
		if arg == a {
			return i
		}
	}

	return -1
}

// Apply a move string to the position given.
// This capability is not yet implemented.
func applyMove(position position, move string) position {
	return position
}

// Determine if the move is in algebraic form. This feature is not yet
// implemented.
func isAlgebraic(move string) bool {
	return true
}
