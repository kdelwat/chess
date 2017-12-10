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

func generateCastlingMoves(position position) []move {
	var moves []move

	if position.toMove == White && position.whiteCanCastleKingside && !piecePresent(position, 5) && !piecePresent(position, 6) {
		moves = append(moves, move(KingCastle))
	} else if position.toMove == Black && position.blackCanCastleKingside && !piecePresent(position, 117) && !piecePresent(position, 118) {
		moves = append(moves, move(KingCastle))
	}

	if position.toMove == White && position.whiteCanCastleQueenside && !piecePresent(position, 1) && !piecePresent(position, 2) && !piecePresent(position, 3) {
		moves = append(moves, move(QueenCastle))
	} else if position.toMove == Black && position.blackCanCastleQueenside && !piecePresent(position, 113) && !piecePresent(position, 114) && !piecePresent(position, 115) {
		moves = append(moves, move(QueenCastle))
	}

	return moves
}

func generatePawnMoves(position position, index int) []move {
	var moves []move

	// try double push
	if isStartingPawn(index, position.toMove) {
		// do double push
		var newIndex int

		if position.toMove == White {
			newIndex = index + 32
		} else {
			newIndex = index - 32
		}

		if !piecePresent(position, newIndex) {
			moves = append(moves, createDoublePawnPush(index, newIndex))
		}
	}

	// try normal move forwards
	var newIndex int

	if position.toMove == White {
		newIndex = index + 16
	} else {
		newIndex = index - 16
	}

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

	var leftAttack int
	var rightAttack int

	var leftEnPassant int
	var rightEnPassant int

	// try attacks

	if position.toMove == White {
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

		if isOnBoard(attackIndex) && piecePresent(position, attackIndex) && getColor(position.board[attackIndex]) != position.toMove {

			// check promo
			if finalRank(attackIndex, position.toMove) {
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Knight))
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Bishop))
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Rook))
				moves = append(moves, createPromotionCaptureMove(index, newIndex, Queen))
			} else {
				moves = append(moves, createCaptureMove(index, attackIndex))
			}
		} else if isOnBoard(attackIndex) && piecePresent(position, enPassantIndex) && getColor(position.board[enPassantIndex]) != position.toMove && pawnHasDoubledAdvanced(position.board[enPassantIndex]) {
			moves = append(moves, createEnPassantCaptureMove(index, enPassantIndex))
		}

	}

	return moves
}

func generateRegularMoves(position position, index int, piece byte) []move {
	var moves []move

	for _, offset := range moveOffsets[getPieceType(piece)] {
		newIndex := index

		for {
			newIndex = index + offset

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
