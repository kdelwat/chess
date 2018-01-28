package main

/*
A move is encoded as a 32-bit integer.

        _ _ _ _  _ _ _ _   unused
castle  +-----+ +------+   special
        +--------------+   from index
        +--------------+   to index

The special byte encodes information about captures, promotions, and other
non-standard moves. Its schema is taken from
https://chessprogramming.wikispaces.com/Encoding+Moves.

+----------------------+-----------+---------+-----------+-----------+
|      Move type       | Promotion | Capture | Special 1 | Special 2 |
+----------------------+-----------+---------+-----------+-----------+
| quiet                |         0 |       0 |         0 |         0 |
| double pawn push     |         0 |       0 |         0 |         1 |
| kingside castle      |         0 |       0 |         1 |         0 |
| queenside castle     |         0 |       0 |         1 |         1 |
| capture              |         0 |       1 |         0 |         0 |
| en passant           |         0 |       1 |         0 |         1 |
| knight promotion     |         1 |       0 |         0 |         0 |
| bishop promotion     |         1 |       0 |         0 |         1 |
| rook promotion       |         1 |       0 |         1 |         0 |
| queen promotion      |         1 |       0 |         1 |         1 |
| knight promo-capture |         1 |       1 |         0 |         0 |
| bishop promo-capture |         1 |       1 |         0 |         1 |
| rook promo-capture   |         1 |       1 |         1 |         0 |
| queen promo-capture  |         1 |       1 |         1 |         1 |
+----------------------+-----------+---------+-----------+-----------+
*/
type move uint32

// Extract the index the piece is moving from.
func (m move) From() byte {
	return byte((m & (0xFF << 8)) >> 8)
}

// Extract the index the piece is moving to.
func (m move) To() byte {
	return byte(m & 0xFF)
}

// Is the move a quiet move?
func (m move) isQuiet() bool {
	return (m&MoveTypeMask == 0)
}

// Is the move a promotion?
func (m move) isPromotion() bool {
	return (m&Promotion != 0)
}

// Is the move a promotion capture?
func (m move) isPromotionCapture() bool {
	return m.isPromotion() && (m&Capture != 0)
}

// Is the move a capture?
func (m move) isCapture() bool {
	return (m&Capture != 0)
}

// Is the move a castle?
func (m move) isCastle() bool {
	return m.isKingCastle() || m.isQueenCastle()
}

// Is the move a kingside castle?
func (m move) isKingCastle() bool {
	return ((m&MoveTypeMask)>>16 == 2)
}

// Is the move a queenside castle?
func (m move) isQueenCastle() bool {
	return ((m&MoveTypeMask)>>16 == 3)
}

// Is the move a double pawn push?
func (m move) isDoublePawnPush() bool {
	return ((m&MoveTypeMask)>>16 == 1)
}

// Is the move an en passant?
func (m move) isEnPassantCapture() bool {
	return m.isCapture() && (m&EnPassant != 0)
}

// Extract the piece that the move promotes to, in the case of promotions or
// promotion capture. This information is encoded in the special section of the
// move, as in the above table.
func (m move) getPromotedPiece(p piece) piece {
	var promotedPiece piece

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

	promotedPiece |= piece(p.color())

	return promotedPiece
}
