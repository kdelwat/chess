package main

import "fmt"

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

func showMove(move move) {
	fmt.Printf("Move from %v to %v\n", move.From(), move.To())
}

func showMoves(moves []move) {
	for i := 0; i < len(moves); i++ {
		showMove(moves[i])
	}
}
