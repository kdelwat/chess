package main

/*
moveOffsets maps each piece to the directions it can move. For example, the king
can move left, with an index offset of -1, or directly upwards, with a move
offset of 16.
*/
var moveOffsets = map[piece][]int{
	King:   {15, 16, 17, -1, 1, -15, -16, -17},
	Queen:  {15, 16, 17, -1, 1, -15, -16, -17},
	Bishop: {15, 17, -15, -17},
	Rook:   {16, -16, 1, -1},
	Knight: {14, 31, 33, 18, -14, -31, -33, -18},
}

/*
castlingBlocks declares the indices which, if they contain a piece, can block
castling to that side for each colour.

castlingChecks declares the indices which, if in check, can block castling.
*/
var castlingBlocks = map[byte]map[int][]int{
	White: map[int][]int{KingCastle: {5, 6}, QueenCastle: {1, 2, 3}},
	Black: map[int][]int{KingCastle: {117, 118}, QueenCastle: {113, 114, 115}},
}

var castlingChecks = map[byte]map[int][]int{
	White: map[int][]int{KingCastle: {5, 6}, QueenCastle: {2, 3}},
	Black: map[int][]int{KingCastle: {117, 118}, QueenCastle: {114, 115}},
}

// Create a quiet move between two indices.
func createQuietMove(from int, to int) move {
	var m move

	m = m | move(to)
	m = m | (move(from) << 8)

	return m
}

// Create a capture move between two indices.
func createCaptureMove(from int, to int) move {
	m := createQuietMove(from, to)

	m |= Capture

	return m
}

// Create a promotion move between two indices, promoting to the given piece
// type.
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

// Create a promotion capture move between two indices, promoting to the given
// piece type.
func createPromotionCaptureMove(from int, to int, pieceType byte) move {
	m := createPromotionMove(from, to, pieceType)

	m = m | Capture

	return m
}

// Create an en passant capture between two indices.
func createEnPassantCaptureMove(from int, to int) move {
	m := createCaptureMove(from, to)

	m = m | EnPassant

	return m
}

// Create a double pawn push between two indices.
func createDoublePawnPush(from int, to int) move {
	m := createQuietMove(from, to)

	m = m | DoublePawnPush

	return m
}

// Given a position and the side to castle (either KingSide or QueenSide),
// determine if the side is able to legally castle.
func clearToCastle(position position, side int) bool {
	// For each index in the potential blockers, check that there is no piece
	// present.
	for _, index := range castlingBlocks[position.toMove][side] {
		if piecePresent(position, index) {
			return false
		}
	}

	// For each index in the potentially-checked indices, ensure that the index
	// is not attacked. This could be optimised by only generating the attack
	// map once.
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

// Generate a slice of legal moves involving castling.
func generateCastlingMoves(position position) []move {
	var moves []move

	var attackingColor byte
	if position.toMove == White {
		attackingColor = Black
	} else {
		attackingColor = White
	}

	if getCastle(position.castling, KingCastle, position.toMove) && clearToCastle(position, KingCastle) && !isKingInCheck(position, attackingColor) {
		moves = append(moves, move(KingCastle))
	}

	if getCastle(position.castling, QueenCastle, position.toMove) && clearToCastle(position, QueenCastle) && !isKingInCheck(position, attackingColor) {
		moves = append(moves, move(QueenCastle))
	}

	return moves
}

// Generate a slice of legal pawn moves.
func generatePawnMoves(position position, index int) []move {
	var moves []move

	// The offset of pawn moves depends on the colour of the pawn, since they
	// can only move forwards.
	var direction int
	if position.toMove == White {
		direction = 1
	} else {
		direction = -1
	}

	// If the pawn is on the starting row, it can perform a double push and move
	// forward two spaces.
	if isOnStartingRow(index, position.toMove) {
		newIndex := index + 32*direction
		jumpIndex := index + 16*direction

		if !piecePresent(position, jumpIndex) && !piecePresent(position, newIndex) {
			moves = append(moves, createDoublePawnPush(index, newIndex))
		}
	}

	// Generate a regular move forwards, and check that the target square is not
	// occupied.
	newIndex := index + 16*direction

	if !piecePresent(position, newIndex) {
		// If the pawn is moving to the final rank, generate promotions.
		if isOnFinalRank(newIndex, position.toMove) {
			moves = append(moves, createPromotionMove(index, newIndex, Knight))
			moves = append(moves, createPromotionMove(index, newIndex, Rook))
			moves = append(moves, createPromotionMove(index, newIndex, Queen))
			moves = append(moves, createPromotionMove(index, newIndex, Bishop))
		} else {
			// Otherwise, generate a quiet move.
			moves = append(moves, createQuietMove(index, newIndex))
		}

	}

	// Generate attacks.
	leftAttack := index + 15*direction
	rightAttack := index + 17*direction

	attackIndices := [2]int{leftAttack, rightAttack}

	// For each attack, check if a capture is possible.
	for _, attackIndex := range attackIndices {
		if isOnBoard(attackIndex) && piecePresent(position, attackIndex) && position.board[attackIndex].color() != position.toMove {
			// If the pawn is capturing a piece on the final rank, generate
			// promotion captures.
			if isOnFinalRank(attackIndex, position.toMove) {
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Knight))
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Bishop))
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Rook))
				moves = append(moves, createPromotionCaptureMove(index, attackIndex, Queen))
			} else {
				// Otherwise, generate a regular capture.
				moves = append(moves, createCaptureMove(index, attackIndex))
			}
		}

	}

	// If the en passant target saved in the current position is capturable by
	// the pawn, generate an en passant move.
	if isEnPassantTarget(position, index, direction) {
		moves = append(moves, createEnPassantCaptureMove(index, int(position.enPassantTarget)))
	}

	return moves
}

// Generate a slice of moves for a non-pawn piece.
func generateRegularMoves(position position, index int, piece piece) []move {
	var moves []move

	// For each offset in the piece's offset map, attempt to make a move.
	for _, offset := range moveOffsets[piece.identity()] {
		newIndex := index

		// Slide along the offset for as long as possible, generating the attack
		// rays for sliding pieces (such as queens).
		for {
			newIndex = newIndex + offset

			// If the new position is off the board, stop sliding.
			if !isOnBoard(newIndex) {
				break
			}

			var newMove move
			if piecePresent(position, newIndex) {
				// If the sliding piece encounters another piece of its own
				// colour, stop sliding.
				if position.board[newIndex].color() == position.toMove {
					break
				}

				// If it encounters a piece of a different colour, capture that
				// piece.
				newMove = createCaptureMove(index, newIndex)
				moves = append(moves, newMove)

				break
			} else {
				// If there is no piece present, generate a quiet move to the
				// current index.
				newMove = createQuietMove(index, newIndex)
			}

			moves = append(moves, newMove)

			// If the piece isn't a sliding piece (i.e. the king and knight),
			// only slide once.
			if !piece.isSliding() {
				break
			}
		}
	}

	return moves
}

// Given a position, generate a slice of moves representing all the possible
// moves for the attacking player. This function generates pseudo-legal moves,
// meaning that the move may cause the player to move into check, which is
// illegal. These moves are filtered in a later step.
func generateMoves(position position) []move {
	var moves []move

	// For each piece of the player's colour, generate all possible moves.
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

// Given a position, generate a slice of moves representing all the possible
// legal moves for the attacking player.
func generateLegalMoves(position position) []move {
	var legal []move

	// Generate pseudo-legal moves.
	moves := generateMoves(position)

	// For each pseudo-legal move, make the move, then see if the king is in
	// check. If it isn't, the move is legal.
	for _, move := range moves {
		artifacts := makeMove(&position, move)
		if !isKingInCheck(position, position.toMove) {
			legal = append(legal, move)
		}
		unmakeMove(&position, move, artifacts)
	}

	return legal
}

// Determine whether the king is in check, given a position and an attacking
// colour.
func isKingInCheck(position position, attackingColor byte) bool {
	// Find the index of the king on the board.
	var kingIndex int
	for i := 0; i < BoardSize; i++ {
		piece := position.board[i]

		if isOnBoard(i) && piece.exists() && piece.is(King) && piece.color() != attackingColor {
			kingIndex = i
			break
		}
	}

	// Determine whether the index is attacked.
	return isAttacked(position, attackingColor, kingIndex)
}
