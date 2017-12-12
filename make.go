package main

type moveArtifacts struct {
	halfmove          int
	castling          castleMap
	enPassantPosition int
	captured          byte
}

func (m move) isQuiet() bool {
	return (m&MoveTypeMask == 0)
}

func (m move) isPromotion() bool {
	return (m&Promotion != 0)
}

func (m move) isPromotionCapture() bool {
	return m.isPromotion() && (m&Capture != 0)
}

func (m move) getPromotedPiece(piece byte) byte {
	var promotedPiece byte

	switch m & PromotionTypeMask {
	case BishopPromotion:
		promotedPiece |= Bishop
	case KnightPromotion:
		promotedPiece |= Knight
	case QueenPromotion:
		promotedPiece |= Queen
	case RookPromotion:
		promotedPiece |= Rook
	}

	promotedPiece |= getColor(piece)

	return promotedPiece
}

func (m move) isCapture() bool {
	return (m&Capture != 0)
}

func (m move) isCastle() bool {
	return m.isKingCastle() || m.isQueenCastle()
}

func (m move) isKingCastle() bool {
	return ((m&MoveTypeMask)>>16 == 2)
}

func (m move) isQueenCastle() bool {
	return ((m&MoveTypeMask)>>16 == 3)
}

func (m move) isDoublePawnPush() bool {
	//fmt.Printf("Checking double pawn push with move %b, typemask %b")
	return ((m&MoveTypeMask)>>16 == 1)
}

func (m move) isEnPassantCapture() bool {
	return m.isCapture() && (m&EnPassant != 0)
}

func makeQuietMove(position *position, from byte, to byte) {
	pieceMoved := position.board[from]

	position.board[from] = 0
	position.board[to] = pieceMoved
}

func makeMove(position *position, move move) moveArtifacts {
	var castleCopy = castleMap{
		White: map[int]bool{KingCastle: position.castling[White][KingCastle],
			QueenCastle: position.castling[White][QueenCastle]},
		Black: map[int]bool{KingCastle: position.castling[Black][KingCastle],
			QueenCastle: position.castling[Black][QueenCastle]},
	}

	var artifacts = moveArtifacts{
		halfmove:          position.halfmove,
		castling:          castleCopy,
		enPassantPosition: position.enPassantTarget,
		captured:          0,
	}

	position.enPassantTarget = -1

	if move&Capture != 0 {
		position.halfmove = 0
	} else {
		position.halfmove++
	}

	// If the king or rook are moving, remove castle rights
	if getPieceType(position.board[move.From()]) == King {
		position.castling[position.toMove][KingCastle] = false
		position.castling[position.toMove][QueenCastle] = false
	}

	if getPieceType(position.board[move.From()]) == Rook {
		if position.toMove == White && move.From() == 0 {
			position.castling[White][QueenCastle] = false
		} else if position.toMove == White && move.From() == 7 {
			position.castling[White][KingCastle] = false
		} else if position.toMove == Black && move.From() == 112 {
			position.castling[Black][QueenCastle] = false
		} else if position.toMove == Black && move.From() == 119 {
			position.castling[Black][KingCastle] = false
		}
	}

	if move.isQuiet() {
		makeQuietMove(position, move.From(), move.To())

		if isPawn(position.board[move.From()]) {
			position.halfmove = 0
		}

	} else if move.isCastle() {
		position.castling[position.toMove][KingCastle] = false
		position.castling[position.toMove][QueenCastle] = false

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
			artifacts.captured = position.board[move.To()]

			position.board[move.From()] = 0
			position.board[move.To()] = promotionPiece
		} else if move.isPromotion() {
			promotionPiece := move.getPromotedPiece(pieceMoved)

			position.board[move.From()] = 0
			position.board[move.To()] = promotionPiece

			position.halfmove = 0
		} else if move.isEnPassantCapture() {
			position.board[move.From()] = 0
			position.board[move.To()] = pieceMoved

			var captureIndex int
			if getColor(pieceMoved) == White {
				captureIndex = int(move.To()) - 16
			} else {
				captureIndex = int(move.To()) + 16
			}

			artifacts.captured = position.board[captureIndex]
			position.board[captureIndex] = 0
		} else if move.isDoublePawnPush() {
			position.board[move.From()] = 0
			position.board[move.To()] = pieceMoved

			position.enPassantTarget = int(move.From()+move.To()) / 2
			position.halfmove = 0
		} else if move.isCapture() {
			artifacts.captured = position.board[move.To()]

			position.board[move.From()] = 0
			position.board[move.To()] = pieceMoved
		}

	}

	// Check if the rook was captured for castling purposes
	if artifacts.captured != 0 && getPieceType(artifacts.captured) == Rook {
		color := getColor(artifacts.captured)

		if position.castling[color][QueenCastle] == true {
			if color == White && move.To() == 0 {
				position.castling[White][QueenCastle] = false
			} else if color == Black && move.To() == 112 {
				position.castling[Black][QueenCastle] = false
			}
		}

		if position.castling[color][KingCastle] == true {
			if color == White && move.To() == 7 {
				position.castling[White][KingCastle] = false
			} else if color == Black && move.To() == 119 {
				position.castling[Black][KingCastle] = false
			}
		}
	}

	if position.toMove == White {
		position.toMove = Black
	} else {
		position.toMove = White
		position.fullmove++
	}

	return artifacts
}

func unmakeMove(position *position, move move, artifacts moveArtifacts) {
	position.halfmove = artifacts.halfmove
	position.castling = artifacts.castling
	position.enPassantTarget = artifacts.enPassantPosition

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

		king := position.board[kingFinal]
		rook := position.board[rookFinal]
		position.board[kingFinal] = 0
		position.board[rookFinal] = 0

		position.board[kingOrigin] = king
		position.board[rookOrigin] = rook
	} else {
		pieceMoved := position.board[move.To()]

		if move.isPromotionCapture() {
			// recreate pawn
			var pawn byte
			pawn |= getColor(pieceMoved)
			pawn |= Pawn

			position.board[move.From()] = pawn
			position.board[move.To()] = artifacts.captured
		} else if move.isPromotion() {
			// recreate pawn
			var pawn byte
			pawn |= getColor(pieceMoved)
			pawn |= Pawn

			position.board[move.From()] = pawn
			position.board[move.To()] = 0
		} else if move.isEnPassantCapture() {
			position.board[move.To()] = 0
			position.board[move.From()] = pieceMoved

			var captureIndex int
			if getColor(pieceMoved) == White {
				captureIndex = int(move.To()) - 16
			} else {
				captureIndex = int(move.To()) + 16
			}

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
