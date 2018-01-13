package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getBestMove(position position) string {
	rand.Seed(time.Now().Unix())

	moves := generateMoves(position)

	//fmt.Printf("Moves: %v\n", moves)
	move := moves[rand.Intn(len(moves))]
	alg := toAlgebraic(position, move)

	//fmt.Printf("Got best move %v, gives %v\n", move, alg)
	return alg
}

func pieceToAlgebraic(p piece) string {
	var code string

	switch p & Piece {
	case King:
		code = "k"
	case Queen:
		code = "q"
	case Rook:
		code = "r"
	case Bishop:
		code = "b"
	case Knight:
		code = "n"
	default:
		code = ""
	}

	return code
}

func indexToFile(index byte) string {
	return fmt.Sprintf("%v", rune((index%16)+'a'))
}

func indexToSquare(index byte) string {
	rank := index/16 + 1
	file := rune((index % 16) + 'a')

	return fmt.Sprintf("%c%v", file, rank)
}

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

		return fmt.Sprintf("%v%v%v", indexToSquare(move.From()), indexToSquare(move.To()), pieceToAlgebraic(promotionPiece))
	}

	return fmt.Sprintf("%v%v", indexToSquare(move.From()), indexToSquare(move.To()))
}

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

// Based on C code from https://chessprogramming.wikispaces.com/Perft
func perft(position position, depth int) perftResults {
	results := perftResults{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if depth == 0 {
		return perftResults{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	moves := generateMoves(position)

	checked := 0

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

	if checked == len(moves) {
		return perftResults{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	return results
}

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
