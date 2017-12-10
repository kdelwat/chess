package main

type position struct {
	board           [128]byte
	castling        map[byte]map[int]bool
	toMove          byte
	enPassantTarget int
	halfmove        int
	fullmove        int
}

var startBoard = [128]byte{
	5, 2, 4, 7, 3, 4, 2, 5, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	65, 65, 65, 65, 65, 65, 65, 65, 0, 0, 0, 0, 0, 0, 0, 0,
	69, 66, 68, 71, 67, 68, 66, 69, 0, 0, 0, 0, 0, 0, 0, 0,
}

func isOnBoard(index int) bool {
	if index&OffBoard != 0 {
		return false
	}

	return true

}

func piecePresent(position position, index int) bool {
	return isPiece(position.board[index])
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
