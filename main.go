package main

import "fmt"

func main() {
	fmt.Printf("Hello, world\n")

	startPosition := position{board: startBoard, toMove: White, whiteCanCastleKingside: true, whiteCanCastleQueenside: true, blackCanCastleKingside: true, blackCanCastleQueenside: true, enPassantTarget: -1, halfmove: 0, fullmove: 1}

	showPosition(startPosition)
	showSliding(startPosition)

	generateMoves(startPosition, White)

	fmt.Printf("FEN: %v\n", toFEN(startPosition))
}
