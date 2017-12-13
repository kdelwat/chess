package main

var NotA uint64 = 0xfefefefefefefefe
var NotH uint64 = 0x7f7f7f7f7f7f7f7f

func newBuildAttackMap(position position, toMove byte, index int) uint64 {

	var empty uint64
	var queens uint64
	var rooks uint64
	var bishops uint64

	for i := 0; i < BoardSize; i++ {
		if isOnBoard(i) {
			piece := position.board[i]
			if piece.color() != toMove {
				continue
			}

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
	showBitboard(southAttacks(queens|rooks, empty))
	showBitboard(northAttacks(queens|rooks, empty))
	showBitboard(eastAttacks(rooks, empty))
	showBitboard(westAttacks(rooks, empty))
	showBitboard(northEastAttacks(queens|bishops, empty))
	showBitboard(northWestAttacks(queens|bishops, empty))
	showBitboard(southEastAttacks(queens|bishops, empty))
	showBitboard(southWestAttacks(queens|bishops, empty))

	return empty
}

func map0x88ToStandard(index int) uint {
	rank := index / 16
	file := index % 16
	return uint(rank*8 + file)
}

func southAttacks(start uint64, empty uint64) uint64 {
	var flood uint64

	for start != 0 {
		flood |= start
		start = (start >> 8) & empty
	}

	return flood >> 8
}

func northAttacks(start uint64, empty uint64) uint64 {
	var flood uint64

	for start != 0 {
		flood |= start
		start = (start << 8) & empty
	}

	return flood << 8
}

func eastAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotA
	for start != 0 {
		flood |= start
		start = (start << 1) & empty
	}

	return (flood << 1) & NotA
}

func westAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotH
	for start != 0 {
		flood |= start
		start = (start >> 1) & empty
	}

	return (flood >> 1) & NotH
}

func northEastAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotA
	for start != 0 {
		flood |= start
		start = (start << 9) & empty
	}

	return (flood << 9) & NotA

}

func northWestAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotH
	for start != 0 {
		flood |= start
		start = (start << 7) & empty
	}

	return (flood << 7) & NotH
}

func southEastAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotA
	for start != 0 {
		flood |= start
		start = (start >> 9) & empty
	}

	return (flood >> 9) & NotA
}

func southWestAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotH
	for start != 0 {
		flood |= start
		start = (start >> 7) & empty
	}

	return (flood >> 7) & NotA
}
