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

func createPromotionCaptureMove(from int, to int, pieceType byte) move {
	m := createPromotionMove(from, to, pieceType)

	m = m | Capture

	return m
}

func createPromotionMove(from int, to int, pieceType byte) move {
	var m move = 0

	m = m | move(to)
	m = m | (move(from) << 8)

	m = m | Promotion

	switch pieceType {
	case Knight:
		{
			m = m | KnightPromotion
		}
	case Bishop:
		{
			m = m | BishopPromotion
		}
	case Rook:
		{
			m = m | RookPromotion
		}
	case Queen:
		{
			m = m | QueenPromotion
		}
	}

	return m
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

			// handle castling
			if isKing(piece) {
				if position.canCastleKingside {
					if color == White && !piecePresent(position, 5) && !piecePresent(position, 6) {
						newMove := move(KingCastle)
						showMove(newMove)
						moves = append(moves, newMove)
					} else if color == Black && !piecePresent(position, 117) && !piecePresent(position, 118) {
						newMove := move(KingCastle)
						showMove(newMove)
						moves = append(moves, newMove)
					}
				}

				if position.canCastleQueenside {
					if color == White && !piecePresent(position, 1) && !piecePresent(position, 2) && !piecePresent(position, 3) {
						newMove := move(QueenCastle)
						showMove(newMove)
						moves = append(moves, newMove)
					} else if color == Black && !piecePresent(position, 113) && !piecePresent(position, 114) && !piecePresent(position, 115) {
						newMove := move(QueenCastle)
						showMove(newMove)
						moves = append(moves, newMove)
					}
				}
			}

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

					// check promotions
					if finalRank(newIndex, color) {
						// do promo
						var newMove move

						newMove = createPromotionMove(i, newIndex, Knight)
						showMove(newMove)
						moves = append(moves, newMove)

						newMove = createPromotionMove(i, newIndex, Bishop)
						showMove(newMove)
						moves = append(moves, newMove)

						newMove = createPromotionMove(i, newIndex, Rook)
						showMove(newMove)
						moves = append(moves, newMove)

						newMove = createPromotionMove(i, newIndex, Queen)
						showMove(newMove)
						moves = append(moves, newMove)

					} else {
						newMove := createQuietMove(i, newIndex)
						showMove(newMove)
						moves = append(moves, newMove)
					}

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

				attackIndices := [2]int{leftAttack, rightAttack}

				for _, attackIndex := range attackIndices {
					if isOnBoard(attackIndex) && piecePresent(position, attackIndex) && getColor(position.board[attackIndex]) != color {

						// check promo
						if finalRank(attackIndex, color) {
							// do promo
							var newMove move

							newMove = createPromotionCaptureMove(i, newIndex, Knight)
							showMove(newMove)
							moves = append(moves, newMove)

							newMove = createPromotionCaptureMove(i, newIndex, Bishop)
							showMove(newMove)
							moves = append(moves, newMove)

							newMove = createPromotionCaptureMove(i, newIndex, Rook)
							showMove(newMove)
							moves = append(moves, newMove)

							newMove = createPromotionCaptureMove(i, newIndex, Queen)
							showMove(newMove)
							moves = append(moves, newMove)
						} else {
							newMove := createCaptureMove(i, attackIndex)
							showMove(newMove)
							moves = append(moves, newMove)
						}
					}

				}

			} else if isSliding(piece) {
				// SLIDING
				for _, offset := range moveOffsets[getPieceType(piece)] {
					newIndex := i
					for {
						newIndex = i + offset

						// skip if new position is off the board
						if !isOnBoard(newIndex) {
							break
						}

						var newMove move
						if piecePresent(position, newIndex) {
							if getColor(position.board[newIndex]) == color {
								break
							}

							newMove = createCaptureMove(i, newIndex)
							break
						} else {
							newMove = createQuietMove(i, newIndex)
						}

						moves = append(moves, newMove)
						showMove(newMove)
					}
				}

			} else {
				for _, offset := range moveOffsets[getPieceType(piece)] {
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
