package main

type move uint32

func (m move) isQuiet() bool {
	return (m&MoveTypeMask == 0)
}

func (m move) isPromotion() bool {
	return (m&Promotion != 0)
}

func (m move) isPromotionCapture() bool {
	return m.isPromotion() && (m&Capture != 0)
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
	return ((m&MoveTypeMask)>>16 == 1)
}

func (m move) isEnPassantCapture() bool {
	return m.isCapture() && (m&EnPassant != 0)
}

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
