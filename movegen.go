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

var castlingBlocks = map[byte]map[int][]int{
	White: map[int][]int{KingCastle: {5, 6}, QueenCastle: {1, 2, 3}},
	Black: map[int][]int{KingCastle: {117, 118}, QueenCastle: {113, 114, 115}},
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

func clearToCastle(position position, side int) bool {
	for _, index := range castlingBlocks[position.toMove][side] {
		if piecePresent(position, index) {
			return false
		}
	}

	return true
}

func generateCastlingMoves(position position) []move {
	var moves []move

	if position.castling[position.toMove][KingCastle] && clearToCastle(position, KingCastle) {
		moves = append(moves, move(KingCastle))
	}

	if position.castling[position.toMove][QueenCastle] && clearToCastle(position, QueenCastle) {
		moves = append(moves, move(QueenCastle))
	}

	return moves
}

func generatePawnMoves(position position, index int) []move {
	var moves []move

	// to change offset based on playing color
	var direction int
	if position.toMove == White {
		direction = 1
	} else {
		direction = -1
	}

	// try double push
	if isStartingPawn(index, position.toMove) {
		// do double push
		newIndex := index + 32*direction

		if !piecePresent(position, newIndex) {
			moves = append(moves, createDoublePawnPush(index, newIndex))
		}
	}

	// try normal move forwards
	newIndex := index + 16*direction

	if !piecePresent(position, newIndex) {
		// check promotions
		if finalRank(newIndex, position.toMove) {
			moves = append(moves, createPromotionMove(index, newIndex, Knight))
			moves = append(moves, createPromotionMove(index, newIndex, Rook))
			moves = append(moves, createPromotionMove(index, newIndex, Queen))
			moves = append(moves, createPromotionMove(index, newIndex, Bishop))
		} else {
			moves = append(moves, createQuietMove(index, newIndex))
		}

	}

	leftAttack := index + 15*direction
	rightAttack := index + 17*direction

	attackIndices := [2]int{leftAttack, rightAttack}

	for _, attackIndex := range attackIndices {

		if isOnBoard(attackIndex) && piecePresent(position, attackIndex) && getColor(position.board[attackIndex]) != position.toMove {
			// check promotions
			if finalRank(attackIndex, position.toMove) {
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Knight))
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Bishop))
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Rook))
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Queen))
			} else {
				moves = append(moves, createCaptureMove(index, attackIndex))
			}
		}

	}

	// opt target - must be on 4th or 5th rank
	if isEnPassantTarget(position, index) {
		moves = append(moves, createEnPassantCaptureMove(index, position.enPassantTarget))

	}

	return moves
}

func generateRegularMoves(position position, index int, piece byte) []move {
	var moves []move

	for _, offset := range moveOffsets[getPieceType(piece)] {
		newIndex := index

		for {
			newIndex = newIndex + offset
			// skip if new position is off the board
			if !isOnBoard(newIndex) {
				break
			}

			var newMove move
			if piecePresent(position, newIndex) {
				if getColor(position.board[newIndex]) == position.toMove {
					break
				}

				newMove = createCaptureMove(index, newIndex)
				break
			} else {
				newMove = createQuietMove(index, newIndex)
			}

			moves = append(moves, newMove)

			if !isSliding(piece) {
				break
			}
		}
	}

	return moves
}

func generateMoves(position position) []move {
	var moves []move

	var piece byte
	for i := 0; i < BoardSize; i++ {

		piece = position.board[i]

		if isOnBoard(i) && isPiece(position.board[i]) && getColor(piece) == position.toMove {
			fmt.Printf("Generating move for piece: %v\n", pieceToString(position.board[i]))

			var pieceMoves []move

			// handle castling
			if isKing(piece) {
				pieceMoves = append(pieceMoves, generateCastlingMoves(position)...)
			}

			if isPawn(piece) {
				pieceMoves = append(pieceMoves, generatePawnMoves(position, i)...)

			} else {
				pieceMoves = append(pieceMoves, generateRegularMoves(position, i, piece)...)
			}

			showMoves(pieceMoves)
			fmt.Print("\n")

			moves = append(moves, pieceMoves...)
		}

	}

	return moves
}
