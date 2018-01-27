package main

/*
position contains the complete game state after a turn.

board is the board state as an array of pieces. The array is 128 elements long,
rather than 64, because it is in 0x88 form. This essentially places a junk board
to the right of the main board, like so:

0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x
0 0 0 0 0 0 0 0 x x x x x x x x

The bottom left hand corner is index 0, while the top right hand corner is 127.

0x88 form has the advantage of allowing very fast checks to see if a position is
on the board, which is used in move generation.

castling is a byte that represents castling rights for both players. Only the
lower 4 bits are used, with 1 indicating castling is allowed.

x x x x _ _ _ _
        ^ ^ ^ ^
        | | | |
        | | | |
        + | | |
White king| | |
          + | |
White queen | |
            | |
Black king+-+ |
              |
Black queen+--+

toMove is the colour of the player who is next to move.

enPassantTarget is the index of a square where there is an en passant
opportunity. If a pawn was double pushed in the previous turn, its jumped
position will appear as the en passant target.

halfmove and fullmove represent the time elapsed in the game.
*/
type position struct {
	board           [128]piece
	castling        byte
	toMove          byte
	enPassantTarget byte
	halfmove        byte
	fullmove        int
}

// Set castling rights in the castle byte.
func setCastle(castling byte, side int, color byte, canCastle bool) byte {
	var offset uint8

	if side == QueenCastle {
		offset++
	}

	if color == Black {
		offset += 2
	}

	if canCastle {
		castling |= 1 << offset
	} else {
		castling &= ^(1 << offset)
	}

	return castling
}

// Get castling rights from the castle byte.
func getCastle(castling byte, side int, color byte) bool {
	var offset uint8

	if side == QueenCastle {
		offset++
	}

	if color == Black {
		offset += 2
	}

	return (castling&(1<<offset) != 0)
}

// Returns true if the index is on the physical board, false otherwise, using
// the 0x88 form for a fast check.
func isOnBoard(index int) bool {
	return index&OffBoard == 0
}

// Determines if there is a piece present at the index.
func piecePresent(position position, index int) bool {
	return position.board[index].exists()
}

// Determines if a piece is on the rank, from 0 to 7, relative to the color of
// the player passed. So 0 will be the row closest to the player, regardless on
// the color selected.
func isOnRelativeRank(index int, color byte, rank int) bool {
	var start int
	if color == White {
		start = 16 * rank
	} else {
		start = 112 - 16*rank
	}

	end := start + 7
	return (index >= start && index <= end)
}

// Determines if the piece is on the final rank (the opposite side of their
// starting position), for purposes of promotion.
func isOnFinalRank(index int, color byte) bool {
	return isOnRelativeRank(index, color, 7)
}

// Determines if a pawn is on its starting row.
func isOnStartingRow(index int, color byte) bool {
	return isOnRelativeRank(index, color, 1)
}

// Given the index of an attacking pawn, and a direction of movement (1 for
// white, -1 for black), returns whether the pawn can attack the current en
// passant target, if it exists.
func isEnPassantTarget(position position, index int, direction int) bool {
	leftTarget := 15 * direction
	rightTarget := 17 * direction

	return position.enPassantTarget != NoEnPassant && (position.enPassantTarget == byte(index+leftTarget) || position.enPassantTarget == byte(index+rightTarget))
}
