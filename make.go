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

func makeMove(position *position, move move) moveArtifacts {
	var artifacts = moveArtifacts{
		halfmove:          position.halfmove, // decrement this on capture or pawn move
		castling:          position.castling,
		enPassantPosition: position.enPassantTarget, // change this
		captured:          0,
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

		position.board[move.From()] = 0
		position.board[move.To()] = pieceMoved
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
	} else if move.isCapture() {
		// make capture
		pieceMoved := position.board[move.To()]

		position.board[move.From()] = pieceMoved
		position.board[move.To()] = artifacts.captured
	}
}
