package main

import (
	"fmt"
)

/*
perft is a performance test for the move generation function. It generates a
move tree to a given depth, recording various information. This information can
be compared to known results to determine if the move generation is behaving
correctly.
*/
type perftResults struct {
	nodes           uint64
	quiet           uint64
	captures        uint64
	enpassant       uint64
	promotion       uint64
	promoCapture    uint64
	castleKingSide  uint64
	castleQueenSide uint64
	pawnJump        uint64
	checks          uint64
}

// Run a perft analysis of the position to the given depth. This function is
// based on the C code at https://chessprogramming.wikispaces.com/Perft
func perft(position position, depth int) perftResults {
	results := perftResults{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// If the end of the tree is reached, increment the number of nodes found.
	if depth == 0 {
		return perftResults{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	// Generate all moves for the position.
	moves := generateMoves(position)

	checked := 0

	// Make each move, recording information about the move.
	for _, move := range moves {

		artifacts := makeMove(&position, move)

		if !isKingInCheck(position, position.toMove) {
			if move.isQuiet() {
				results.quiet++
			} else if move.isQueenCastle() {
				results.castleQueenSide++
			} else if move.isKingCastle() {
				results.castleKingSide++
			} else if move.isPromotionCapture() {
				results.promoCapture++
			} else if move.isPromotion() {
				results.promotion++
			} else if move.isEnPassantCapture() {
				results.enpassant++
			} else if move.isDoublePawnPush() {
				results.pawnJump++
			} else if move.isCapture() {
				results.captures++
			}

			if isKingInCheck(position, position.toMove) {
				results.checks++
			}

			perftResults := perft(position, depth-1)
			results.nodes += perftResults.nodes
			results.quiet += perftResults.quiet
			results.captures += perftResults.captures
			results.enpassant += perftResults.enpassant
			results.promotion += perftResults.promotion
			results.promoCapture += perftResults.promoCapture
			results.castleKingSide += perftResults.castleKingSide
			results.castleQueenSide += perftResults.castleQueenSide
			results.pawnJump += perftResults.pawnJump
			results.checks += perftResults.checks
		} else {
			checked++
		}

		unmakeMove(&position, move, artifacts)
	}

	// If every move results in a check, return an empty statistics list, since
	// the game has ended and there are no nodes to record.
	if checked == len(moves) {
		return perftResults{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	return results
}

// Run a perft analysis, but divide the initial level of the move tree. This
// allows for debugging the problematic paths of move generation.
func dividePerft(position position, depth int) {
	moves := generateLegalMoves(position)
	var total uint64

	for _, move := range moves {
		artifacts := makeMove(&position, move)
		results := perft(position, depth-1)

		fmt.Printf("%v: %v\n", toAlgebraic(position, move), results.nodes)

		total += results.nodes

		unmakeMove(&position, move, artifacts)
	}

	fmt.Printf("TOTAL: %v\n", total)
}

// Convert an index to a number representing its file.
func indexToFile(index byte) string {
	return fmt.Sprintf("%v", rune((index%16)+'a'))
}

// Convert an index to a rank and file coordinate.
func indexToSquare(index byte) string {
	rank := index/16 + 1
	file := rune((index % 16) + 'a')

	return fmt.Sprintf("%c%v", file, rank)
}

// Convert a move to algebraic notation.
func toAlgebraic(position position, move move) string {

	if move.isKingCastle() {
		return "0-0"
	}

	if move.isQueenCastle() {
		return "0-0-0"
	}

	if move.isPromotion() || move.isPromotionCapture() {
		pieceMoved := position.board[move.From()]
		promotionPiece := move.getPromotedPiece(pieceMoved)

		return fmt.Sprintf("%v%v%v", indexToSquare(move.From()), indexToSquare(move.To()), pieceToString(promotionPiece))
	}

	return fmt.Sprintf("%v%v", indexToSquare(move.From()), indexToSquare(move.To()))
}
