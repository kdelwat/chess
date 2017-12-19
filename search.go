package main

import "fmt"

func search(position position, depth int) move {
	moves := generateLegalMoves(position)

	bestScore := -1000
	var bestMove move

	for _, move := range moves {
		artifacts := makeMove(&position, move)

		negamaxScore := alphaBeta(&position, -1000, 1000, depth)
		if negamaxScore > bestScore {
			bestScore = negamaxScore
			bestMove = move
		}

		unmakeMove(&position, move, artifacts)
	}

	fmt.Printf("Best move is %v with score %v\n", toAlgebraic(position, bestMove), bestScore)
	return bestMove
}

// alpha beta algorithm from pseudocode on
// https://chessprogramming.wikispaces.com/Alpha-Beta
func alphaBeta(position *position, alpha int, beta int, depth int) int {
	if depth == 0 {
		return evaluate(*position)
	}

	moves := generateMoves(*position)

	for _, move := range moves {

		artifacts := makeMove(position, move)

		score := -alphaBeta(position, -beta, -alpha, depth-1)
		if score >= beta {
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