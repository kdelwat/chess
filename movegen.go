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
	return byte((m & (0xFF << 8)) >> 8)
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

func createDoublePawnPush(from int, to int) move {
	var m move = 0

	m = m | move(to)
	m = m | (move(from) << 8)
	m = m | DoublePawnPush

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
				// try double push
				if isStartingPawn(i, color) {
					// do double push
					var newIndex int

					if color == White {
						newIndex = i + 32
					} else {
						newIndex = i - 32
					}

					if !piecePresent(position, newIndex) {
						newMove := createDoublePawnPush(i, newIndex)
						showMove(newMove)
						moves = append(moves, newMove)
					}
				}

				// try normal move forwards
				var newIndex int

				if color == White {
					newIndex = i + 16
				} else {
					newIndex = i - 16
				}

				if !piecePresent(position, newIndex) {
					newMove := createQuietMove(i, newIndex)
					showMove(newMove)
					moves = append(moves, newMove)
				}

				// try attacks

				var leftAttack int
				var rightAttack int

				if color == White {
					leftAttack = i + 15
					rightAttack = i + 17
				} else {
					leftAttack = i - 15
					rightAttack = i - 17
				}

				if isOnBoard(leftAttack) && piecePresent(position, leftAttack) && getColor(position.board[leftAttack]) != color {
					newMove := createCaptureMove(i, leftAttack)
					showMove(newMove)
					moves = append(moves, newMove)
				}

				if isOnBoard(rightAttack) && piecePresent(position, rightAttack) && getColor(position.board[rightAttack]) != color {
					newMove := createCaptureMove(i, rightAttack)
					showMove(newMove)
					moves = append(moves, newMove)
				}

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

			fmt.Print("\n")
		}
	}

	return moves
}
