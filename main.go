package main

import "fmt"

func main() {
	fmt.Printf("Hello, world\n")

	var castling map[byte]map[int]bool

	castling[White][KingCastle] = true
	castling[Black][KingCastle] = true
	castling[White][QueenCastle] = true
	castling[Black][QueenCastle] = true

	startPosition := position{board: startBoard, toMove: White, castling: castling, enPassantTarget: -1, halfmove: 0, fullmove: 1}

	showPosition(startPosition)
	showSliding(startPosition)

	generateMoves(startPosition)

	fmt.Printf("FEN: %v\n", toFEN(startPosition))
}
