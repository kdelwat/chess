package main

func isSliding(piece byte) bool {
	if piece&Sliding != 0 {
		return true
	} else {
		return false
	}
}

func isPawn(piece byte) bool {
	return (getPieceType(piece) == Pawn)
}

func isKing(piece byte) bool {
	return (getPieceType(piece) == King)
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

func isStartingPawn(index int, color byte) bool {
	if color == White && index >= 16 && index <= 23 {
		return true
	}

	if color == Black && index >= 96 && index <= 103 {
		return true
	}

	return false
}
