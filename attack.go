package main

var NotA uint64 = 0xfefefefefefefefe
var NotH uint64 = 0x7f7f7f7f7f7f7f7f

// Attack map and associated method created by Jonatan Pettersson
// https://mediocrechess.blogspot.com.au/2006/12/guide-attacked-squares.html
var attackNone = 0
var attackKQR = 1
var attackQR = 2
var attackKQBwP = 3
var attackKQBbP = 4
var attackQB = 5
var attackN = 6

var attackArray = []int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0,
	0, 0, 0, 5, 0, 0, 5, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 5, 0,
	0, 0, 0, 5, 0, 0, 0, 0, 2, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0,
	5, 0, 0, 0, 2, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0,
	2, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 6, 2, 6, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 4, 1, 4, 6, 0, 0, 0, 0, 0,
	0, 2, 2, 2, 2, 2, 2, 1, 0, 1, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0,
	0, 0, 6, 3, 1, 3, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 6,
	2, 6, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 2, 0, 0, 5,
	0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 2, 0, 0, 0, 5, 0, 0, 0,
	0, 0, 0, 5, 0, 0, 0, 0, 2, 0, 0, 0, 0, 5, 0, 0, 0, 0, 5, 0,
	0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 5, 0, 0, 5, 0, 0, 0, 0, 0, 0,
	2, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// Uses the Dumb7Fill algorithm for sliding piece attacks
// implemented based on https://chessprogramming.wikispaces.com/Dumb7Fill
func isAttacked(position position, toMove byte, index int) bool {
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

			canAttack := attackArray[index-i+128]

			index := map0x88ToStandard(i)

			switch piece.identity() {
			case Queen:
				if canAttack == attackNone || canAttack == attackN {
					continue
				} else {
					queens |= 1 << index
				}
			case Bishop:
				if !(canAttack == attackKQBbP || canAttack == attackKQBwP || canAttack == attackQB) {
					continue
				} else {
					bishops |= 1 << index
				}
			case Rook:
				if !(canAttack == attackKQR || canAttack == attackQR) {
					continue
				} else {
					rooks |= 1 << index
				}
			case Knight:
				if canAttack == attackN {
					return true
				}
			case Pawn:
				if (toMove == White && canAttack == attackKQBwP) || (toMove == Black && canAttack == attackKQBbP) {
					return true
				}
			case King:
				if canAttack == attackKQR || canAttack == attackKQBbP || canAttack == attackKQBwP {
					return true
				}
			default:
				empty |= 1 << index
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

	return attackMap&(1<<map0x88ToStandard(index)) != 0
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
