package main

import "fmt"

type position struct {
	board [128]byte
}

var startBoard = [128]byte{
	69, 66, 68, 71, 67, 68, 66, 69, 0, 0, 0, 0, 0, 0, 0, 0,
	65, 65, 65, 65, 65, 65, 65, 65, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
	5, 2, 4, 7, 3, 4, 2, 5, 0, 0, 0, 0, 0, 0, 0, 0,
}

func isSliding(piece byte) bool {
	if piece&Sliding != 0 {
		return true
	} else {
		return false
	}
}

func showPiece(piece byte) {
	switch piece & Piece {
	case King:
		fmt.Print("K")
	case Queen:
		fmt.Print("Q")
	case Rook:
		fmt.Print("R")
	case Bishop:
		fmt.Print("B")
	case Knight:
		fmt.Print("H")
	case Pawn:
		fmt.Print("P")
	default:
		fmt.Print("_")
	}
}
func showPosition(position position) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			showPiece(position.board[i])
		}
		if (i+1)%16 == 0 {
			fmt.Print("\n")
		}
	}
}

func showSliding(position position) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			if isSliding(position.board[i]) {
				fmt.Print("T")
			} else {
				fmt.Print("F")
			}
		}

		if (i+1)%16 == 0 {
			fmt.Print("\n")
		}
	}

}
func main() {
	fmt.Printf("Hello, world\n")

	startPosition := position{board: startBoard}

	showPosition(startPosition)
	showSliding(startPosition)
}
