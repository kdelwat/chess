package main

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

func isOnBoard(index int) bool {
	if index&OffBoard != 0 {
		return false
	} else {
		return true
	}
}

func piecePresent(position position, index int) bool {
	return isPiece(position.board[index])
}
