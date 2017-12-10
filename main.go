package main

import "fmt"

func main() {
	fmt.Printf("Hello, world\n")

	var castling = map[byte]map[int]bool{
		White: map[int]bool{KingCastle: true, QueenCastle: true},
		Black: map[int]bool{KingCastle: true, QueenCastle: true},
	}

	startPosition := position{board: startBoard, toMove: White, castling: castling, enPassantTarget: -1, halfmove: 0, fullmove: 1}

	showPosition(startPosition)
	//showSliding(startPosition)

	//generateMoves(startPosition)

	startFen := toFEN(startPosition)
	fmt.Printf("FEN: %v\n", toFEN(startPosition))
	fmt.Print("FEN gives...\n")

	anotherPosition := fromFen(startFen)
	showPosition(anotherPosition)

}
