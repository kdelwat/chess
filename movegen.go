package main

var moveOffsets = map[piece][]int{
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

var castlingChecks = map[byte]map[int][]int{
	White: map[int][]int{KingCastle: {5, 6}, QueenCastle: {2, 3}},
	Black: map[int][]int{KingCastle: {117, 118}, QueenCastle: {114, 115}},
}

// Attack map and associated method created by Jonatan Pettersson
// https://mediocrechess.blogspot.com.au/2006/12/guide-attacked-squares.html
var attackNone = 0
var attackKQR = 1
var attackQR = 2
var attackKQBwP = 3
var attackKQBbP = 4
var attackQB = 5
var attackN = 6

var attackArray = []int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,
	0, 0, 0, 5, 0, 0, 5, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 5, 0,
	0, 0, 0, 5, 0, 0, 0, 0, 2, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0,
	5, 0, 0, 0, 2, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0,
	2, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 6, 2, 6, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 4, 1, 4, 6, 0, 0, 0, 0, 0,
	0, 2, 2, 2, 2, 2, 2, 1, 0, 1, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0,
	0, 0, 6, 3, 1, 3, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 6,
	2, 6, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 2, 0, 0, 5,
	0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 2, 0, 0, 0, 5, 0, 0, 0,
	0, 0, 0, 5, 0, 0, 0, 0, 2, 0, 0, 0, 0, 5, 0, 0, 0, 0, 5, 0,
	0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 5, 0, 0, 5, 0, 0, 0, 0, 0, 0,
	2, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0}

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

		// this is such an easy opt target - just generate the attack map once
		if piecePresent(position, index) {
			return false
		}
	}

	for _, index := range castlingChecks[position.toMove][side] {
		var attackingColor byte
		if position.toMove == White {
			attackingColor = Black
		} else {
			attackingColor = White
		}

		if isAttacked(position, attackingColor, index) {
			return false
		}
	}

	return true
}

func generateCastlingMoves(position position) []move {
	var moves []move

	var attackingColor byte
	if position.toMove == White {
		attackingColor = Black
	} else {
		attackingColor = White
	}

	if position.castling[position.toMove][KingCastle] && clearToCastle(position, KingCastle) && !isKingInCheck(position, attackingColor) {
		moves = append(moves, move(KingCastle))
	}

	if position.castling[position.toMove][QueenCastle] && clearToCastle(position, QueenCastle) && !isKingInCheck(position, attackingColor) {
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
	if isOnStartingRow(index, position.toMove) {
		// do double push
		newIndex := index + 32*direction
		jumpIndex := index + 16*direction

		if !piecePresent(position, jumpIndex) && !piecePresent(position, newIndex) { // reorder for opt?
			moves = append(moves, createDoublePawnPush(index, newIndex))
		}
	}

	// try normal move forwards
	newIndex := index + 16*direction

	if !piecePresent(position, newIndex) {
		// check promotions
		if isOnFinalRank(newIndex, position.toMove) {
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

		if isOnBoard(attackIndex) && piecePresent(position, attackIndex) && position.board[attackIndex].color() != position.toMove {
			// check promotions
			if isOnFinalRank(attackIndex, position.toMove) {
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Knight))
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Bishop))
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Rook))
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Queen))
			} else {
				moves = append(moves, createCaptureMove(index, attackIndex))
			}
		}

	}

	// opt target - must be on 4th or 5th rank
	if isEnPassantTarget(position, index, direction) {
		moves = append(moves, createEnPassantCaptureMove(index, position.enPassantTarget))
	}

	return moves
}

func generateRegularMoves(position position, index int, piece piece) []move {
	var moves []move

	for _, offset := range moveOffsets[piece.identity()] {
		newIndex := index

		for {
			newIndex = newIndex + offset
			// skip if new position is off the board
			if !isOnBoard(newIndex) {
				break
			}

			var newMove move
			if piecePresent(position, newIndex) {
				if position.board[newIndex].color() == position.toMove {
					break
				}

				newMove = createCaptureMove(index, newIndex)
				moves = append(moves, newMove)

				break
			} else {
				newMove = createQuietMove(index, newIndex)
			}

			moves = append(moves, newMove)

			if !piece.isSliding() {
				break
			}
		}
	}

	return moves
}

func generateMoves(position position) []move {
	var moves []move

	var piece piece
	for i := 0; i < BoardSize; i++ {

		piece = position.board[i]

		if isOnBoard(i) && position.board[i].exists() && piece.color() == position.toMove {
			var pieceMoves []move

			if piece.is(King) {
				pieceMoves = append(pieceMoves, generateCastlingMoves(position)...)
			}

			if piece.is(Pawn) {
				pieceMoves = append(pieceMoves, generatePawnMoves(position, i)...)

			} else {
				pieceMoves = append(pieceMoves, generateRegularMoves(position, i, piece)...)
			}

			moves = append(moves, pieceMoves...)
		}

	}

	return moves
}

func generateLegalMoves(position position) []move {
	var legal []move
	moves := generateMoves(position)

	for _, move := range moves {
		artifacts := makeMove(&position, move)
		if !isKingInCheck(position, position.toMove) {
			legal = append(legal, move)
		}
		unmakeMove(&position, move, artifacts)
	}

	return legal
}

// doesn't handle en passant because currently only used for checking checks
func buildAttackMap(position position, toMove byte, index int) [128]byte {
	var attackMap [128]byte
	for i := 0; i < BoardSize; i++ {
		piece := position.board[i]

		if isOnBoard(i) && position.board[i].exists() && piece.color() == toMove {
			canAttack := attackArray[index-i+128]

			switch piece.identity() {
			case Queen:
				if canAttack == attackNone || canAttack == attackN {
					continue
				}
			case Bishop:
				if !(canAttack == attackKQBbP || canAttack == attackKQBwP || canAttack == attackQB) {
					continue
				}
			case Rook:
				if !(canAttack == attackKQR || canAttack == attackQR) {
					continue
				}
			case Knight:
				if canAttack != attackN {
					continue
				}
			case Pawn:
				if !((toMove == White && canAttack == attackKQBwP) || (toMove == Black && canAttack == attackKQBbP)) {
					continue
				}
			}

			if piece.is(Pawn) {
				// to change offset based on playing color
				var direction int
				if toMove == White {
					direction = 1
				} else {
					direction = -1
				}

				leftAttack := i + 15*direction
				rightAttack := i + 17*direction

				if isOnBoard(leftAttack) {
					attackMap[leftAttack] = 1
				}
				if isOnBoard(rightAttack) {
					attackMap[rightAttack] = 1
				}
			} else {
				for _, offset := range moveOffsets[piece.identity()] {
					newIndex := i

					for {
						newIndex = newIndex + offset
						// skip if new position is off the board
						if !isOnBoard(newIndex) {
							break
						}

						attackMap[newIndex] = 1

						if piecePresent(position, newIndex) {
							break
						}

						if !piece.isSliding() {
							break
						}
					}
				}

			}
		}

	}

	return attackMap
}

func isAttacked(position position, attackingColor byte, index int) bool {
	attackMap := buildAttackMap(position, attackingColor, index)
	return attackMap[index] == 1
}

func isKingInCheck(position position, attackingColor byte) bool {
	// find the king
	var kingIndex int

	for i := 0; i < BoardSize; i++ {
		piece := position.board[i]

		if isOnBoard(i) && piece.exists() && piece.is(King) && piece.color() != attackingColor {
			kingIndex = i
			break
		}
	}

	return isAttacked(position, attackingColor, kingIndex)
}
