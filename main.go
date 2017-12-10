package main

import "fmt"

func main() {
	fmt.Printf("Hello, world\n")

	startPosition := position{board: startBoard, canCastleKingside: false, canCastleQueenside: false}

	showPosition(startPosition)
	showSliding(startPosition)

	generateMoves(startPosition, White)

	//showMoves(moves)
}
