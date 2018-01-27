package main

/*
When making a move, some information about the previous state cannot be
recovered from the next state. The moveArtifacts type contains this information.
*/
type moveArtifacts struct {
	halfmove          byte
	castling          byte
	enPassantPosition byte
	captured          piece
}

// Makes a quiet move (a regular move with no captures) given the position,
// origin, and destination.
func makeQuietMove(position *position, from byte, to byte) {
	pieceMoved := position.board[from]

	position.board[from] = 0
	position.board[to] = pieceMoved
}

/*
Performs a move on the given position. This function takes a pointer to the
current position, so it modifies it in-place. While this isn't ideal from a
debugging point of view, it would be impossible to copy the position each
time due to memory constraints.

makeMove returns the artifacts required to reverse the move later.
*/
func makeMove(position *position, move move) moveArtifacts {

	// Record the current state in artifacts.
	var artifacts = moveArtifacts{
		halfmove:          position.halfmove,
		castling:          position.castling,
		enPassantPosition: position.enPassantTarget,
		captured:          0,
	}

	// The new position by default has no en passant target.
	position.enPassantTarget = NoEnPassant

	// The halfmove counter is rest on a capture, and incremented otherwise.
	if move&Capture != 0 {
		position.halfmove = 0
	} else {
		position.halfmove++
	}

	// If the king or rook are moving, remove castling rights
	if position.board[move.From()].is(King) {
		position.castling = setCastle(position.castling, KingCastle, position.toMove, false)
		position.castling = setCastle(position.castling, QueenCastle, position.toMove, false)
	}

	if position.board[move.From()].is(Rook) {
		if position.toMove == White && move.From() == 0 {
			position.castling = setCastle(position.castling, QueenCastle, White, false)
		} else if position.toMove == White && move.From() == 7 {
			position.castling = setCastle(position.castling, KingCastle, White, false)
		} else if position.toMove == Black && move.From() == 112 {
			position.castling = setCastle(position.castling, QueenCastle, Black, false)
		} else if position.toMove == Black && move.From() == 119 {
			position.castling = setCastle(position.castling, KingCastle, Black, false)
		}
	}

	// Determine which type of move to make.
	if move.isQuiet() {
		makeQuietMove(position, move.From(), move.To())

		// A quiet move resets the halfmove counter if it is made by a pawn.
		if position.board[move.From()].is(Pawn) {
			position.halfmove = 0
		}

	} else if move.isCastle() {
		// If the player is castling, remove all castle rights in the future.
		position.castling = setCastle(position.castling, KingCastle, position.toMove, false)
		position.castling = setCastle(position.castling, QueenCastle, position.toMove, false)

		// Determine the starting and ending location of the pieces involved.
		var kingOrigin int
		var rookOrigin int
		var kingFinal int
		var rookFinal int

		if position.toMove == Black {
			kingOrigin = 116
		} else {
			kingOrigin = 4
		}

		if move.isQueenCastle() {
			rookOrigin = kingOrigin - 4
			kingFinal = kingOrigin - 2
			rookFinal = kingOrigin - 1
		} else {
			rookOrigin = kingOrigin + 3
			kingFinal = kingOrigin + 2
			rookFinal = kingOrigin + 1
		}

		// Swap the pieces in the board.
		king := position.board[kingOrigin]
		rook := position.board[rookOrigin]
		position.board[kingOrigin] = 0
		position.board[rookOrigin] = 0

		position.board[kingFinal] = king
		position.board[rookFinal] = rook
	} else {
		pieceMoved := position.board[move.From()]

		if move.isPromotionCapture() {
			promotionPiece := move.getPromotedPiece(pieceMoved)

			// Save the captured piece in the move artifacts.
			artifacts.captured = position.board[move.To()]

			position.board[move.From()] = 0
			position.board[move.To()] = promotionPiece
		} else if move.isPromotion() {
			promotionPiece := move.getPromotedPiece(pieceMoved)

			position.board[move.From()] = 0
			position.board[move.To()] = promotionPiece

			// The halfmove counter is reset on a promotion.
			position.halfmove = 0
		} else if move.isEnPassantCapture() {
			position.board[move.From()] = 0
			position.board[move.To()] = pieceMoved

			// Determine the en passant target, depending on the direction of
			// movement.
			var captureIndex int
			if pieceMoved.color() == White {
				captureIndex = int(move.To()) - 16
			} else {
				captureIndex = int(move.To()) + 16
			}

			artifacts.captured = position.board[captureIndex]
			position.board[captureIndex] = 0
		} else if move.isDoublePawnPush() {
			position.board[move.From()] = 0
			position.board[move.To()] = pieceMoved

			// A double pawn push creates an en passant target, which must be
			// saved in the new position.
			position.enPassantTarget = byte(move.From()+move.To()) / 2
			position.halfmove = 0
		} else if move.isCapture() {
			artifacts.captured = position.board[move.To()]

			position.board[move.From()] = 0
			position.board[move.To()] = pieceMoved
		}

	}

	// If the rook was captured, remove castling rights for that side.
	if artifacts.captured != 0 && artifacts.captured.is(Rook) {
		color := artifacts.captured.color()

		if getCastle(position.castling, QueenCastle, color) {
			if color == White && move.To() == 0 {
				position.castling = setCastle(position.castling, QueenCastle, White, false)
			} else if color == Black && move.To() == 112 {
				position.castling = setCastle(position.castling, QueenCastle, Black, false)
			}
		}

		if getCastle(position.castling, KingCastle, color) {
			if color == White && move.To() == 7 {
				position.castling = setCastle(position.castling, KingCastle, White, false)
			} else if color == Black && move.To() == 119 {
				position.castling = setCastle(position.castling, KingCastle, Black, false)
			}
		}
	}

	// Increment the fullmove counter when black finishes their turn.
	if position.toMove == White {
		position.toMove = Black
	} else {
		position.toMove = White
		position.fullmove++
	}

	return artifacts
}

/*
Reverses a move on the given position. This function takes a position, the move
which was applied, and the artifacts generated by makeMove, and restores the
position in-place to the state before the move was applied.
*/
func unmakeMove(position *position, move move, artifacts moveArtifacts) {
	// Restore state information from artifacts.
	position.halfmove = artifacts.halfmove
	position.castling = artifacts.castling
	position.enPassantTarget = artifacts.enPassantPosition

	// Decrement the fullmove counter if black made the last move.
	if position.toMove == White {
		position.fullmove--
		position.toMove = Black
	} else {
		position.toMove = White
	}

	if move.isQuiet() {
		pieceMoved := position.board[move.To()]

		position.board[move.To()] = 0
		position.board[move.From()] = pieceMoved
	} else if move.isCastle() {
		// Determine the starting and ending location of the pieces involved.
		var kingOrigin int
		var rookOrigin int
		var kingFinal int
		var rookFinal int

		if position.toMove == Black {
			kingOrigin = 116
		} else {
			kingOrigin = 4
		}

		if move.isQueenCastle() {
			rookOrigin = kingOrigin - 4
			kingFinal = kingOrigin - 2
			rookFinal = kingOrigin - 1
		} else {
			rookOrigin = kingOrigin + 3
			kingFinal = kingOrigin + 2
			rookFinal = kingOrigin + 1
		}

		// Swap the pieces in the board.
		king := position.board[kingFinal]
		rook := position.board[rookFinal]
		position.board[kingFinal] = 0
		position.board[rookFinal] = 0

		position.board[kingOrigin] = king
		position.board[rookOrigin] = rook
	} else {
		pieceMoved := position.board[move.To()]

		if move.isPromotionCapture() {
			// If the move was a promotion capture, recreate the pawn piece that
			// was replaced by the promoted piece.
			var pawn piece
			pawn |= piece(pieceMoved.color())
			pawn |= Pawn

			position.board[move.From()] = pawn
			position.board[move.To()] = artifacts.captured
		} else if move.isPromotion() {
			// If the move was a promotion , recreate the pawn piece that was
			// replaced by the promoted piece.
			var pawn piece
			pawn |= piece(pieceMoved.color())
			pawn |= Pawn

			position.board[move.From()] = pawn
			position.board[move.To()] = 0
		} else if move.isEnPassantCapture() {
			position.board[move.To()] = 0
			position.board[move.From()] = pieceMoved

			// Determine the index captured by the en passant, depending on the
			// direction of movement.
			var captureIndex int
			if pieceMoved.color() == White {
				captureIndex = int(move.To()) - 16
			} else {
				captureIndex = int(move.To()) + 16
			}

			// Restore the captured piece.
			position.board[captureIndex] = artifacts.captured
		} else if move.isDoublePawnPush() {
			position.board[move.To()] = 0
			position.board[move.From()] = pieceMoved

			position.enPassantTarget = artifacts.enPassantPosition
		} else if move.isCapture() {
			position.board[move.From()] = pieceMoved
			position.board[move.To()] = artifacts.captured
		}
	}
}
