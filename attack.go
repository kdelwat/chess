package main

var NotA uint64 = 0xfefefefefefefefe
var NotH uint64 = 0x7f7f7f7f7f7f7f7f

func newBuildAttackMap(position position, toMove byte, index int) uint64 {

	var attackMap uint64

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
			} else if piece.is(King) {
				attackMap |= kingAttacks(i)
			} else if piece.is(Knight) {
				attackMap |= knightAttacks(i)
			}
		}
	}

	attackMap |= southAttacks(queens|rooks, empty)
	attackMap |= northAttacks(queens|rooks, empty)
	attackMap |= eastAttacks(queens|rooks, empty)
	attackMap |= westAttacks(queens|rooks, empty)
	attackMap |= northEastAttacks(queens|bishops, empty)
	attackMap |= northWestAttacks(queens|bishops, empty)
	attackMap |= southEastAttacks(queens|bishops, empty)
	attackMap |= southWestAttacks(queens|bishops, empty)

	showBitboard(attackMap)

	return empty
}

func map0x88ToStandard(index int) uint {
	rank := index / 16
	file := index % 16
	return uint(rank*8 + file)
}

func knightAttacks(index int) uint64 {
	var knight uint64

	standardIndex := map0x88ToStandard(index)

	knight |= 1 << (standardIndex + 10)
	knight |= 1 << (standardIndex + 6)
	knight |= 1 << (standardIndex - 10)
	knight |= 1 << (standardIndex - 6)
	knight |= 1 << (standardIndex + 17)
	knight |= 1 << (standardIndex + 15)
	knight |= 1 << (standardIndex - 17)
	knight |= 1 << (standardIndex - 15)

	if standardIndex%8 == 0 || standardIndex%8 == 1 {
		knight &= 0x3f3f3f3f3f3f3f3f
	} else if standardIndex%8 == 6 || standardIndex%8 == 7 {
		knight &= 0xfcfcfcfcfcfcfcfc

	}

	return knight
}

func kingAttacks(index int) uint64 {
	var king uint64

	standardIndex := map0x88ToStandard(index)

	king |= 1 << (standardIndex + 8)
	king |= 1 << (standardIndex - 8)

	if standardIndex%8 != 0 {
		king |= 1 << (standardIndex - 1)
		king |= 1 << (standardIndex + 7)
		king |= 1 << (standardIndex - 9)
	}

	if standardIndex%8 != 7 {
		king |= 1 << (standardIndex + 1)
		king |= 1 << (standardIndex + 9)
		king |= 1 << (standardIndex - 7)
	}

	return king
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
