package main

type castleMap map[byte]map[int]bool

type position struct {
	board           [128]piece
	castling        castleMap
	toMove          byte
	enPassantTarget int
	halfmove        int
	fullmove        int
}

func isOnBoard(index int) bool {
	if index&OffBoard != 0 {
		return false
	}

	return true

}

func piecePresent(position position, index int) bool {
	return position.board[index].exists()
}

func finalRank(index int, color byte) bool {
	if color == White && index >= 112 && index <= 119 {
		return true
	}

	if color == Black && index >= 0 && index <= 7 {
		return true
	}

	return false
}

func isEnPassantTarget(position position, index int, direction int) bool {
	leftTarget := 15 * direction
	rightTarget := 17 * direction

	return position.enPassantTarget != -1 && (position.enPassantTarget == index+leftTarget || position.enPassantTarget == index+rightTarget)
}

func isOnStartingRow(index int, color byte) bool {
	if color == White && index >= 16 && index <= 23 {
		return true
	}

	if color == Black && index >= 96 && index <= 103 {
		return true
	}

	return false
}
