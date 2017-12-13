package main

var NotA uint64 = 0xfefefefefefefefe
var NotH uint64 = 0x7f7f7f7f7f7f7f7f

func buildAttackMap(position position, toMove byte) uint64 {

	var attackMap uint64

	var empty uint64
	var queens uint64
	var rooks uint64
	var bishops uint64

	for i := 0; i < BoardSize; i++ {
		if isOnBoard(i) {
			piece := position.board[i]
			if piece.exists() && piece.color() != toMove {
				continue
			}

			index := map0x88ToStandard(i)

			if piece.is(Empty) {
				empty |= 1 << index
			} else if piece.is(Rook) {
				rooks |= 1 << index
			} else if piece.is(Bishop) {
				bishops |= 1 << index
			} else if piece.is(Queen) {
				queens |= 1 << index
			} else if piece.is(King) {
				attackMap |= kingAttacks(index)
			} else if piece.is(Knight) {
				attackMap |= knightAttacks(index)
			} else if piece.is(Pawn) {
				attackMap |= pawnAttacks(index, toMove)
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

	// showBitboard(westAttacks(queens|rooks, empty))

	return attackMap
}

func map0x88ToStandard(index int) uint {
	rank := index / 16
	file := index % 16
	return uint(rank*8 + file)
}

func pawnAttacks(index uint, color byte) uint64 {
	var pawn uint64

	if color == White && index%8 != 0 {
		pawn |= 1 << (index + 7)
	}

	if color == Black && index%8 != 0 {
		pawn |= 1 << (index - 9)
	}

	if color == White && index%8 != 7 {
		pawn |= 1 << (index + 9)
	}

	if color == Black && index%8 != 7 {
		pawn |= 1 << (index - 7)
	}

	return pawn
}

func knightAttacks(index uint) uint64 {
	var knight uint64

	knight |= 1 << (index + 10)
	knight |= 1 << (index + 6)
	knight |= 1 << (index - 10)
	knight |= 1 << (index - 6)
	knight |= 1 << (index + 17)
	knight |= 1 << (index + 15)
	knight |= 1 << (index - 17)
	knight |= 1 << (index - 15)

	if index%8 == 0 || index%8 == 1 {
		knight &= 0x3f3f3f3f3f3f3f3f
	} else if index%8 == 6 || index%8 == 7 {
		knight &= 0xfcfcfcfcfcfcfcfc

	}

	return knight
}

func kingAttacks(index uint) uint64 {
	var king uint64

	king |= 1 << (index + 8)
	king |= 1 << (index - 8)

	if index%8 != 0 {
		king |= 1 << (index - 1)
		king |= 1 << (index + 7)
		king |= 1 << (index - 9)
	}

	if index%8 != 7 {
		king |= 1 << (index + 1)
		king |= 1 << (index + 9)
		king |= 1 << (index - 7)
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
		start = (start >> 7) & empty
	}

	return (flood >> 7) & NotA
}

func southWestAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= NotH

	for start != 0 {
		flood |= start
		start = (start >> 9) & empty
	}

	return (flood >> 9) & NotH
}
