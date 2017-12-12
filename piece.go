package main

type piece byte

func (p piece) isSliding() bool {
	return (p&Sliding != 0)
}

func (p piece) is(target piece) bool {
	return (p&Piece == target)
}

func (p piece) exists() bool {
	return (p&Piece != 0)
}

func (p piece) color() byte {
	return byte(p & Color)
}

func (p piece) identity() piece {
	return p & Piece
}
