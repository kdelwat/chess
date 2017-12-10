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

func createQuietMove(from int, to int) move {
	var m move = 0

	m = m | move(to)
	m = m | (move(from) << 8)

	return m
}

func createCaptureMove(from int, to int) move {
	var m move = 0

	m = m | move(to)
	m = m | (move(from) << 8)
	m = m | Capture

	fmt.Printf("Capture byte is: %b", Capture)

	return m
}

func generateMoves(position position, color byte) []move {
	var moves []move

	var piece byte
	for i := 0; i < BoardSize; i++ {

		piece = position.board[i]

		if isOnBoard(i) && isPiece(position.board[i]) && getColor(piece) == color {
			fmt.Printf("Generating move for piece:")
			showPiece(position.board[i])
			fmt.Print("\n")

			if isPawn(piece) {
				// PAWN LOGIC
			} else if isSliding(piece) {
				// SLIDING
			} else {
				for offset := range moveOffsets[getPieceType(piece)] {
					newIndex := i + offset

					// skip if new position is off the board
					if !isOnBoard(newIndex) {
						continue
					}

					var newMove move
					if piecePresent(position, newIndex) {
						if getColor(position.board[newIndex]) == color {
							continue
						}

						newMove = createCaptureMove(i, newIndex)
					} else {
						newMove = createQuietMove(i, newIndex)
					}

					moves = append(moves, newMove)
					showMove(newMove)
				}
			}
		}
	}

	return moves
}
