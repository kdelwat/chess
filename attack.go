package main

func newBuildAttackMap(position position, toMove byte, index int) uint64 {

	var empty uint64
	var queens uint64
	var rooks uint64
	var bishops uint64

	for i := 0; i < BoardSize; i++ {
		if isOnBoard(i) {
			piece := position.board[i]
			if piece.is(Empty) {
				empty |= 1 << map0x88ToStandard(i)
			} else if piece.is(Rook) {
				rooks |= 1 << map0x88ToStandard(i)
			} else if piece.is(Bishop) {
				bishops |= 1 << map0x88ToStandard(i)
			} else if piece.is(Queen) {
				queens |= 1 << map0x88ToStandard(i)
			}
		}
	}

	showBitboard(empty)
	showBitboard(queens)
	showBitboard(rooks)
	showBitboard(bishops)

	return empty
}

func map0x88ToStandard(index int) uint {
	rank := index / 16
	file := index % 16
	return uint(rank*8 + file)
}
