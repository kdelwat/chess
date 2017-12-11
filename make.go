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

func (m move) isDoublePawnPush() bool {
	//fmt.Printf("Checking double pawn push with move %b, typemask %b")
	return ((m&MoveTypeMask)>>16 == 1)
}

func (m move) isEnPassantCapture() bool {
	return m.isCapture() && (m&EnPassant != 0)
}

func makeMove(position *position, move move) moveArtifacts {
	var castleCopy = castleMap{
		White: map[int]bool{KingCastle: position.castling[White][KingCastle], QueenCastle: position.castling[White][QueenCastle]},
		Black: map[int]bool{KingCastle: position.castling[Black][KingCastle], QueenCastle: position.castling[Black][QueenCastle]},
	}
	var artifacts = moveArtifacts{
		halfmove:          position.halfmove,
		castling:          castleCopy,
		enPassantPosition: position.enPassantTarget,
		captured:          0,
	}

	position.enPassantTarget = -1

	// checking castles
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

	if position.toMove == White {
		position.toMove = Black
	} else {
		position.toMove = White
		position.fullmove++
	}

	position.halfmove++

	if move.isQuiet() {
		// make quiet move
		pieceMoved := position.board[move.From()]

		if isPawn(pieceMoved) {
			position.halfmove = 0
		}

		position.board[move.From()] = 0
		position.board[move.To()] = pieceMoved
	} else if move.isPromotionCapture() {
		pieceMoved := position.board[move.From()]
		promotionPiece := move.getPromotedPiece(pieceMoved)

		artifacts.captured = position.board[move.To()]

		position.board[move.From()] = 0
		position.board[move.To()] = promotionPiece

		position.halfmove = 0
	} else if move.isPromotion() {
		pieceMoved := position.board[move.From()]
		promotionPiece := move.getPromotedPiece(pieceMoved)

		position.board[move.From()] = 0
		position.board[move.To()] = promotionPiece

		position.halfmove = 0
	} else if move.isEnPassantCapture() {
		pieceMoved := position.board[move.From()]

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

		position.halfmove = 0
	} else if move.isDoublePawnPush() {
		pieceMoved := position.board[move.From()]

		position.board[move.From()] = 0
		position.board[move.To()] = pieceMoved

		position.enPassantTarget = int(move.From()+move.To()) / 2
		position.halfmove = 0
	} else if move.isCapture() {
		// make capture
		pieceMoved := position.board[move.From()]
		artifacts.captured = position.board[move.To()]

		position.board[move.From()] = 0
		position.board[move.To()] = pieceMoved

		position.halfmove = 0
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

	position.halfmove = artifacts.halfmove

	if move.isQuiet() {
		// unmake quiet move
		pieceMoved := position.board[move.To()]

		position.board[move.To()] = 0
		position.board[move.From()] = pieceMoved
	} else if move.isEnPassantCapture() {
		pieceMoved := position.board[move.To()]

		position.board[move.To()] = 0
		position.board[move.From()] = pieceMoved

		var captureIndex int
		if getColor(pieceMoved) == White {
			captureIndex = int(move.To()) - 16
		} else {
			captureIndex = int(move.To()) + 16
		}

		position.board[captureIndex] = artifacts.captured
	} else if move.isPromotionCapture() {
		pieceMoved := position.board[move.To()]

		// reacreate pawn
		var pawn byte
		pawn |= getColor(pieceMoved)
		pawn |= Pawn

		position.board[move.From()] = pawn
		position.board[move.To()] = artifacts.captured

	} else if move.isPromotion() {
		pieceMoved := position.board[move.To()]

		// reacreate pawn
		var pawn byte
		pawn |= getColor(pieceMoved)
		pawn |= Pawn

		position.board[move.From()] = pawn
		position.board[move.To()] = 0
	} else if move.isDoublePawnPush() {
		pieceMoved := position.board[move.To()]

		position.board[move.To()] = 0
		position.board[move.From()] = pieceMoved

		position.enPassantTarget = artifacts.enPassantPosition
	} else if move.isCapture() {
		// make capture
		pieceMoved := position.board[move.To()]

		position.board[move.From()] = pieceMoved
		position.board[move.To()] = artifacts.captured
	}
}
