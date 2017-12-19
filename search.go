package main

import "fmt"

func search(position position, depth int) move {
	moves := generateLegalMoves(position)

	bestScore := -1000
	var bestMove move

	for _, move := range moves {
		artifacts := makeMove(&position, move)

		negamaxScore := negamax(&position, depth)

		if negamaxScore > bestScore {
			bestScore = negamaxScore
			bestMove = move
		}

		unmakeMove(&position, move, artifacts)
	}

	fmt.Printf("Best move is %v with score %v\n", toAlgebraic(position, bestMove), bestScore)
	return bestMove
}

// negamax algorithm from https://chessprogramming.wikispaces.com/Minimax
func negamax(position *position, depth int) int {
	if depth == 0 {
		return evaluate(*position)
	}

	max := -1000

	moves := generateMoves(*position)

	for _, move := range moves {

		artifacts := makeMove(position, move)

		score := negamax(position, depth-1)
		if score > max {
			max = score
		}

		unmakeMove(position, move, artifacts)
	}

	return max

}
