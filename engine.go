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

func pieceToAlgebraic(piece byte) string {
	var code string

	switch piece & Piece {
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

	if move.isPromotion() || move.isPromotionCapture() {
		pieceMoved := position.board[move.From()]
		promotionPiece := move.getPromotedPiece(pieceMoved)

		return fmt.Sprintf("%v%v%v", indexToSquare(move.From()), indexToSquare(move.To()), pieceToAlgebraic(promotionPiece))
	} else {
		return fmt.Sprintf("%v%v", indexToSquare(move.From()), indexToSquare(move.To()))
	}
}
