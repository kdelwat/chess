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

// Is the piece a sliding piece (bishop, rook, and queen)?
func (p piece) isSliding() bool {
	return (p&Sliding != 0)
}

// Is the piece the same as the target type?
func (p piece) is(target piece) bool {
	return (p&Piece == target)
}

// Is there a piece on the square?
func (p piece) exists() bool {
	return (p&Piece != 0)
}

// Find the colour of the piece.
func (p piece) color() byte {
	return byte(p & Color)
}

// Find the type of the piece.
func (p piece) identity() piece {
	return p & Piece
}
