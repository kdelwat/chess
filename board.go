package main

type castleMap map[byte]map[int]bool

type position struct {
	board           [128]piece
	castling        byte
	toMove          byte
	enPassantTarget int
	halfmove        int
	fullmove        int
}

func setCastle(castling byte, side int, color byte, canCastle bool) byte {
	var offset uint8

	if side == QueenCastle {
		offset++
	}

	if color == Black {
		offset += 2
	}

	if canCastle {
		castling |= 1 << offset
	} else {
		castling &= ^(1 << offset)
	}

	return castling
}

func getCastle(castling byte, side int, color byte) bool {
	var offset uint8

	if side == QueenCastle {
		offset++
	}

	if color == Black {
		offset += 2
	}

	return (castling&(1<<offset) != 0)
}

func isOnBoard(index int) bool {
	return index&OffBoard == 0
}

func piecePresent(position position, index int) bool {
	return position.board[index].exists()
}

func isOnRelativeRank(index int, color byte, rank int) bool {
	var start int
	if color == White {
		start = 16 * rank
	} else {
		start = 112 - 16*rank
	}

	end := start + 7
	return (index >= start && index <= end)
}

func isOnFinalRank(index int, color byte) bool {
	return isOnRelativeRank(index, color, 7)
}

func isOnStartingRow(index int, color byte) bool {
	return isOnRelativeRank(index, color, 1)
}

func isEnPassantTarget(position position, index int, direction int) bool {
	leftTarget := 15 * direction
	rightTarget := 17 * direction

	return position.enPassantTarget != -1 && (position.enPassantTarget == index+leftTarget || position.enPassantTarget == index+rightTarget)
}
