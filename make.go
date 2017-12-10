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

func (m move) isCapture() bool {
	return (m&Capture != 0)
}

func (m move) isDoublePawnPush() bool {
	//fmt.Printf("Checking double pawn push with move %b, typemask %b")
	return ((m&MoveTypeMask)>>16 == 1)
}

func makeMove(position *position, move move) moveArtifacts {
	var artifacts = moveArtifacts{
		halfmove:          position.halfmove,
		castling:          position.castling,
		enPassantPosition: position.enPassantTarget,
		captured:          0,
	}

	position.enPassantTarget = -1

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
