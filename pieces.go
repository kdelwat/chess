package main

func isSliding(piece byte) bool {
	if piece&Sliding != 0 {
		return true
	} else {
		return false
	}
}

func isPiece(piece byte) bool {
	if piece&Piece != 0 {
		return true
	} else {
		return false
	}
}
