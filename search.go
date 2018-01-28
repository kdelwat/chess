package main

import (
	"context"
)

/* Runs a search for the best move, given a context, which determines when the
search will end, a position, the depth to search until, and a move channel for
passing the current best move.

runSearch uses iterative deepening. It will search to a progressively greater
depth, returning each best move until it is signalled to stop.
*/
func runSearch(ctx context.Context, position position, depth int, ch chan move) {
	// Declare the initial cutoffs for the alpha-beta pruning.
	alpha := -100000
	beta := 100000

	for i := 1; i <= depth; i++ {
		result := search(position, i, alpha, beta)

		select {
		case <-ctx.Done():
			return
		default:
			ch <- result
		}
	}
}

// Search for the best move for a position, to a given depth.
func search(position position, depth int, alpha int, beta int) move {
	// Generate all legal moves for the current position.
	moves := generateLegalMoves(position)

	bestScore := -100000
	var bestMove move

	// For each move available, run a search of its tree to the given depth, to
	// identify the best outcome.
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

/* Run a negamax search of the move tree from a given position, to a given
depth. The negamax search finds the "least-bad" move; the move that minimises
the opponents advantage no matter how they play.

An alpha-beta cutoff algorithm prunes the search tree to save time that would be
wasted exploring moves that have proven already to be worst than the best
candidate.

This funciton was implemented from the pseudocode at
https://chessprogramming.wikispaces.com/Alpha-Beta.
*/
func alphaBeta(position *position, alpha int, beta int, depth int) int {
	// At the bottom of the tree, return the score of the position for the attacking player.
	if depth == 0 {
		return evaluate(*position)
	}

	// Otherwise, generate all possible moves.
	moves := generateLegalMoves(*position)
	for _, move := range moves {

		// Make the move.
		artifacts := makeMove(position, move)

		// Recursively call the search function to determine the move's score.
		score := -alphaBeta(position, -beta, -alpha, depth-1)

		// If the score is higher than the beta cutoff, the rest of the search
		// tree is irrelevant and the cutoff is returned.
		if score >= beta {
			unmakeMove(position, move, artifacts)
			return beta
		}

		// Otherwise, replace the alpha if the new score is higher.
		if score > alpha {
			alpha = score
		}

		// Restore the pre-move state of the board.
		unmakeMove(position, move, artifacts)
	}

	return alpha
}
