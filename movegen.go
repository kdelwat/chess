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

func createQuietMove(from int, to int) move {
	var m move

	m = m | move(to)
	m = m | (move(from) << 8)

	return m
}

func createCaptureMove(from int, to int) move {
	m := createQuietMove(from, to)

	m |= Capture

	return m
}

func createPromotionMove(from int, to int, pieceType byte) move {
	m := createQuietMove(from, to)

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

func createPromotionCaptureMove(from int, to int, pieceType byte) move {
	m := createPromotionMove(from, to, pieceType)

	m = m | Capture

	return m
}

func createEnPassantCaptureMove(from int, to int) move {
	m := createCaptureMove(from, to)

	m = m | EnPassant

	return m
}

func createDoublePawnPush(from int, to int) move {
	m := createQuietMove(from, to)

	m = m | DoublePawnPush

	return m
}

func generateCastlingMoves(position position, color byte) []move {
	var moves []move

	if color == White && position.whiteCanCastleKingside && !piecePresent(position, 5) && !piecePresent(position, 6) {
		newMove := move(KingCastle)
		showMove(newMove)
		moves = append(moves, newMove)
	} else if color == Black && position.blackCanCastleKingside && !piecePresent(position, 117) && !piecePresent(position, 118) {
		newMove := move(KingCastle)
		showMove(newMove)
		moves = append(moves, newMove)
	}

	if color == White && position.whiteCanCastleQueenside && !piecePresent(position, 1) && !piecePresent(position, 2) && !piecePresent(position, 3) {
		newMove := move(QueenCastle)
		showMove(newMove)
		moves = append(moves, newMove)
	} else if color == Black && position.blackCanCastleQueenside && !piecePresent(position, 113) && !piecePresent(position, 114) && !piecePresent(position, 115) {
		newMove := move(QueenCastle)
		showMove(newMove)
		moves = append(moves, newMove)
	}

	return moves
}

func generatePawnMoves(position position, color byte, index int) []move {
	var moves []move

	// try double push
	if isStartingPawn(index, color) {
		// do double push
		var newIndex int

		if color == White {
			newIndex = index + 32
		} else {
			newIndex = index - 32
		}

		if !piecePresent(position, newIndex) {
			newMove := createDoublePawnPush(index, newIndex)
			showMove(newMove)
			moves = append(moves, newMove)
		}
	}

	// try normal move forwards
	var newIndex int

	if color == White {
		newIndex = index + 16
	} else {
		newIndex = index - 16
	}

	if !piecePresent(position, newIndex) {

		// check promotions
		if finalRank(newIndex, color) {
			// do promo
			var newMove move

			newMove = createPromotionMove(index, newIndex, Knight)
			showMove(newMove)
			moves = append(moves, newMove)

			newMove = createPromotionMove(index, newIndex, Bishop)
			showMove(newMove)
			moves = append(moves, newMove)

			newMove = createPromotionMove(index, newIndex, Rook)
			showMove(newMove)
			moves = append(moves, newMove)

			newMove = createPromotionMove(index, newIndex, Queen)
			showMove(newMove)
			moves = append(moves, newMove)

		} else {
			newMove := createQuietMove(index, newIndex)
			showMove(newMove)
			moves = append(moves, newMove)
		}

	}

	var leftAttack int
	var rightAttack int

	var leftEnPassant int
	var rightEnPassant int

	// try attacks

	if color == White {
		leftAttack = index + 15
		rightAttack = index + 17
	} else {
		leftAttack = index - 15
		rightAttack = index - 17
	}

	leftEnPassant = index - 1
	rightEnPassant = index + 1

	attackIndices := [2]int{leftAttack, rightAttack}

	for _, attackIndex := range attackIndices {

		var enPassantIndex int

		if attackIndex == leftAttack {
			enPassantIndex = leftEnPassant
		} else {
			enPassantIndex = rightEnPassant
		}

		if isOnBoard(attackIndex) && piecePresent(position, attackIndex) && getColor(position.board[attackIndex]) != color {

			// check promo
			if finalRank(attackIndex, color) {
				// do promo
				var newMove move

				newMove = createPromotionCaptureMove(index, newIndex, Knight)
				showMove(newMove)
				moves = append(moves, newMove)

				newMove = createPromotionCaptureMove(index, newIndex, Bishop)
				showMove(newMove)
				moves = append(moves, newMove)

				newMove = createPromotionCaptureMove(index, newIndex, Rook)
				showMove(newMove)
				moves = append(moves, newMove)

				newMove = createPromotionCaptureMove(index, newIndex, Queen)
				showMove(newMove)
				moves = append(moves, newMove)
			} else {
				newMove := createCaptureMove(index, attackIndex)
				showMove(newMove)
				moves = append(moves, newMove)
			}
		} else if isOnBoard(attackIndex) && piecePresent(position, enPassantIndex) && getColor(position.board[enPassantIndex]) != color && pawnHasDoubledAdvanced(position.board[enPassantIndex]) {
			newMove := createEnPassantCaptureMove(index, enPassantIndex)
			showMove(newMove)
			moves = append(moves, newMove)
		}

	}

	return moves
}

func generateMoves(position position, color byte) []move {
	var moves []move

	var piece byte
	for i := 0; i < BoardSize; i++ {

		piece = position.board[i]

		if isOnBoard(i) && isPiece(position.board[i]) && getColor(piece) == color {
			fmt.Printf("Generating move for piece: %v\n", pieceToString(position.board[i]))

			// handle castling
			if isKing(piece) {
				moves = append(moves, generateCastlingMoves(position, color)...)
			}

			if isPawn(piece) {
				moves = append(moves, generatePawnMoves(position, color, i)...)

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
