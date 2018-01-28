package main

/*
A piece is represented by a single byte.

          off the board
          +       +
          |       | +----+sliding
          v       v v
          _ _ _ _ _ _ _ _
            ^ ^     ^ ^ ^
  colour+---+ |     + + +
              |     identity
double pushed++

If the piece is white, the colour bit is 0. Otherwise, it is 1.

+--------+----------+
| Piece  | Identity |
+--------+----------+
| empty  |      000 |
| pawn   |      001 |
| knight |      010 |
| bishop |      100 |
| rook   |      101 |
| king   |      011 |
| queen  |      111 |
+--------+----------+
*/
type piece byte

const pieceIdentityMask = 0x0F
const color = 0x40

// Is the piece a sliding piece (bishop, rook, and queen)?
const sliding = 0x04

func (p piece) isSliding() bool {
	return (p&sliding != 0)
}

// Is the piece the same as the target type?
func (p piece) is(target piece) bool {
	return (p&pieceIdentityMask == target)
}

// Is there a piece on the square?
func (p piece) exists() bool {
	return (p&pieceIdentityMask != 0)
}

// Find the colour of the piece.
func (p piece) color() byte {
	return byte(p & color)
}

// Find the type of the piece.
func (p piece) identity() piece {
	return p & pieceIdentityMask
}
