package main

func isSliding(piece byte) bool {
	if piece&Sliding != 0 {
		return true
	} else {
		return false
	}
}

func isPawn(piece byte) bool {
	return (piece&Pawn != 0)
}

func isPiece(piece byte) bool {
	if piece&Piece != 0 {
		return true
	} else {
		return false
	}
}

func getPieceType(piece byte) byte {
	return piece & Piece
}

func getColor(piece byte) byte {
	return piece & Color
}
