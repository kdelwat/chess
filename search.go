package main

import (
	"context"
	"fmt"
)

var cutoffs int
var cutoffPositions = make(map[int]int)

// aspiration window code from pseudocode at https://mediocrechess.blogspot.com.au/2007/01/guide-aspiration-windows-killer-moves.html
var aspirationWindow = 10

func runSearch(ctx context.Context, position position, depth int, ch chan move) {
	cutoffs = 0

	alpha := -1000
	beta := 1000

	for i := 1; i <= depth; {
		result, bestAlpha := search(position, i, alpha, beta)

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Searching to depth %v\n", i)
			ch <- result
		}

		if bestAlpha <= alpha || bestAlpha >= beta {
			alpha = -1000
			beta = 1000
			continue
		} else {
			alpha = bestAlpha - aspirationWindow
			beta = bestAlpha + aspirationWindow
			i++
		}
	}
}

func search(position position, depth int, alpha int, beta int) (move, int) {
	moves := generateLegalMoves(position)

	bestScore := -1000
	var bestMove move

	for _, move := range moves {
		artifacts := makeMove(&position, move)
		negamaxScore := alphaBeta(&position, alpha, beta, depth)

		if negamaxScore >= bestScore {
			bestScore = negamaxScore
			bestMove = move
		}

		unmakeMove(&position, move, artifacts)
	}

	fmt.Printf("Best move is %v with score %v\nCutoffs: %v (@ %v)\n", toAlgebraic(position, bestMove), bestScore, cutoffs, cutoffPositions)
	return bestMove, bestScore
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
