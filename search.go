package main

import (
	"context"
	"fmt"
)

var cutoffs int
var cutoffPositions = make(map[int]int)

func runSearch(ctx context.Context, position position, depth int, ch chan move) {
	cutoffs = 0

	alpha := -100000
	beta := 100000

	for i := 1; i <= depth; i++ {
		result := search(position, i, alpha, beta)

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Searching to depth %v\n", i)
			ch <- result
		}
	}
}

func search(position position, depth int, alpha int, beta int) move {
	moves := generateLegalMoves(position)

	bestScore := -100000
	var bestMove move

	for _, move := range moves {
		artifacts := makeMove(&position, move)
		negamaxScore := -alphaBeta(&position, alpha, beta, depth)

		if negamaxScore >= bestScore {
			bestScore = negamaxScore
			bestMove = move
		}

		unmakeMove(&position, move, artifacts)
	}

	return bestMove
}

// alpha beta algorithm from pseudocode on
// https://chessprogramming.wikispaces.com/Alpha-Beta
func alphaBeta(position *position, alpha int, beta int, depth int) int {
	if depth == 0 {
		return evaluate(*position)
	}

	moves := generateLegalMoves(*position)
	for index, move := range moves {

		artifacts := makeMove(position, move)

		score := -alphaBeta(position, -beta, -alpha, depth-1)

		if score >= beta {
			cutoffs++
			cutoffPositions[index]++

			unmakeMove(position, move, artifacts)
			return beta
		}

		if score > alpha {
			alpha = score
		}

		unmakeMove(position, move, artifacts)
	}

	return alpha
}
