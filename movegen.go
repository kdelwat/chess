package main

import "fmt"

type move uint32

var moveOffsets = map[byte][]int{
	King:   {15, 16, 17, -1, 1, -15, -16, -17},
	Queen:  {15, 16, 17, -1, 1, -15, -16, -17},
	Bishop: {15, 17, -15, -17},
	Rook:   {16, -16, 1, -1},
	Knight: {14, 31, 33, 18, -14, -31, -33, -18},
}

func (m move) From() byte {
	return byte((m & 0xFF) >> 8)
}

func (m move) To() byte {
	return byte(m & 0xFF)
}

func (m *move) setFrom(from byte) {
	*m = *m | (move(from) << 8)
}

func (m *move) setTo(to byte) {
	*m = *m | move(to)
}

func generateMoves(position position) []move {
	var moves []move

	for i := 0; i < BoardSize; i++ {
		if isOnBoard(i) && isPiece(position.board[i]) {
			fmt.Printf("Generating move for piece:")
			showPiece(position.board[i])
			fmt.Print("\n")
		}
	}

	return moves
}
