package main

/*
These constants remove bits from a bitboard on the A and H file respectively.
For example, if we have a bitboard representing attacks of a queen that looks
like this:

    1 0 0 1 0 0 1 0
    0 1 0 1 0 1 0 0
    0 0 1 1 1 0 0 0
    1 1 1 1 1 1 1 1
    0 0 1 1 1 0 0 0
    0 1 0 1 0 1 0 0
    1 0 0 1 0 0 1 0
    0 0 0 1 0 0 0 1

A bitwise AND with notA will remove any attacks on the A file:

    0 0 0 1 0 0 1 0
    0 1 0 1 0 1 0 0
    0 0 1 1 1 0 0 0
    0 1 1 1 1 1 1 1
    0 0 1 1 1 0 0 0
    0 1 0 1 0 1 0 0
    0 0 0 1 0 0 1 0
    0 0 0 1 0 0 0 1
*/
var notA uint64 = 0xfefefefefefefefe
var notH uint64 = 0x7f7f7f7f7f7f7f7f

/*
The following declaration is an attack map, which represents the ability of
pieces to attack each other around the board. This is used for quick lookups -
we can skip generating attacks for a Queen, for exampe, if we know already that
its position couldn't possibly attack the King.

This attack map, and the associated method, was created by Jonatan Pettersson.
See his guide here:
https://mediocrechess.blogspot.com.au/2006/12/guide-attacked-squares.html
*/
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

/*
isAttacked determines if a piece is under attack. It takes the current game
position, the index of the piece in question (in 0x88 form), and the color of
the attacking side.
*/
func isAttacked(position position, toMove byte, targetIndex int) bool {
	/*
		Declare bitboards for representing the pieces present. While the normal board position is in 0x88 form, these bitboards don't require the extra squares and simple represent an 8x8 grid, meaning that they can fit in a 64-bit integer.

		attackMap represents the squares currently under attack. empty represents squares with no pieces. The other bitboards hold the positions of sliding pieces.
	*/
	var attackMap uint64
	var empty uint64
	var queens uint64
	var rooks uint64
	var bishops uint64

	// Loop through every index on the board, skipping over indices that fall
	// outside the visible playing area.
	for i := 0; i < BoardSize; i++ {
		if isOnBoard(i) {
			// Extract the piece in the current index.
			piece := position.board[i]

			// If there is a piece present, but it isn't on the attacking side,
			// we can skip the iteration.
			if piece.exists() && piece.color() != toMove {
				continue
			}

			// Look up the pieces that can attack the target from this index,
			// using the attack array declared above. The lookup returns a
			// constant representing the set of possible pieces.
			canAttack := attackArray[targetIndex-i+128]

			// Convert the index in 0x88 form to the standard 8x8 form.
			index := map0x88ToStandard(i)

			/*
				Moves are generated differently depending on the type of piece.

				Non-sliding pieces can simply be checked against the canAttack
				constant generated previously. If they are found to be
				attacking, we can return early and save computations.

				Sliding pieces are first checked against this constant, which
				saves costly move generation if it's impossible for them to ever
				attack the target square. If they could attack it, they are
				added to the relevant bitboard for later generation.
			*/
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

	/*
		Now that the bitboards have been filled by sliding pieces, we can generate
		the attack map. To do this, we use the Dumb7Fill algorithm, implemented
		based on the site https://chessprogramming.wikispaces.com/Dumb7Fill.

		For each direction a sliding piece could move, these moves are generated,
		following a three step process:

		    1. The pieces that can make that move have their bitboards combined with bitwise OR. For example, only queens and rooks can attack directly east, so we only combine their bitboards.

		    2. The bitboard is shifted according to the algorithm, which moves the pieces in the relevant direction until they hit a non-empty square (which blocks the attack.)

		    3. These moves are combined with the overall attack map using bitwise OR.
	*/
	attackMap |= southAttacks(queens|rooks, empty)
	attackMap |= northAttacks(queens|rooks, empty)
	attackMap |= eastAttacks(queens|rooks, empty)
	attackMap |= westAttacks(queens|rooks, empty)
	attackMap |= northEastAttacks(queens|bishops, empty)
	attackMap |= northWestAttacks(queens|bishops, empty)
	attackMap |= southEastAttacks(queens|bishops, empty)
	attackMap |= southWestAttacks(queens|bishops, empty)

	// Finally, we check whether the target index is attacked in the attack map
	// and return true if possible.
	return attackMap&(1<<map0x88ToStandard(targetIndex)) != 0
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
	empty &= notA
	for start != 0 {
		flood |= start
		start = (start << 1) & empty
	}

	return (flood << 1) & notA
}

func westAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= notH
	for start != 0 {
		flood |= start
		start = (start >> 1) & empty
	}

	return (flood >> 1) & notH
}

func northEastAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= notA

	for start != 0 {
		flood |= start
		start = (start << 9) & empty
	}

	return (flood << 9) & notA

}

func northWestAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= notH

	for start != 0 {
		flood |= start
		start = (start << 7) & empty
	}

	return (flood << 7) & notH
}

func southEastAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= notA

	for start != 0 {
		flood |= start
		start = (start >> 7) & empty
	}

	return (flood >> 7) & notA
}

func southWestAttacks(start uint64, empty uint64) uint64 {
	var flood uint64
	empty &= notH

	for start != 0 {
		flood |= start
		start = (start >> 9) & empty
	}

	return (flood >> 9) & notH
}
