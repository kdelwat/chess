package main

/*
These constants specify the value, in centipawns, of each piece when evaluating a position.

A piece's value is a combination of its base weight and a modifier based on its position on the board. This reflects more subtle information about a position. For example, a knight in the centre of the board is more effective than one on the side, so it recieves a bonus.

The values and piece tables used are from Tomasz Michniewski and can be found at https://chessprogramming.wikispaces.com/Simplified+evaluation+function.
*/
const kingWeight = 10000
const queenWeight = 900
const rookWeight = 500
const bishopWeight = 300
const knightWeight = 300
const pawnWeight = 100

var pawnPositions = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var knightPositions = []int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 0, 0, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 15, 20, 20, 15, 0, -30,
	-30, 5, 10, 15, 15, 10, 5, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var bishopPositions = []int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
}

var rookPositions = []int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var queenPositions = []int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
}

var kingPositions = []int{
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-30, -40, -40, -50, -50, -40, -40, -30,
	-20, -30, -30, -40, -40, -30, -30, -20,
	-10, -20, -20, -20, -20, -20, -20, -10,
	20, 20, 0, 0, 0, 0, 20, 20,
	20, 30, 10, 0, 0, 10, 30, 20,
}

/*
evaluate returns an objective score representing the game's current result. A
game starts at 0, with no player having the advantage. As it progresses, if
white were to start taking pieces, the score would increase. If black were to
instead perform strongly, the score would decrease.

Since the move search function uses the Negamax algorithm, this evaluation is
symmetrical. A position for black is the same as the identical one for white,
but negated.
*/
func evaluate(position position) int {
	var score int

	var direction int
	if position.toMove == White {
		direction = 1
	} else {
		direction = -1
	}

	// Loop through the board, finding the score for each piece present. If the
	// piece is white, add it to the total; if black, subtract it.
	for i := 0; i < BoardSize; i++ {

		piece := position.board[i]
		if isOnBoard(i) && position.board[i].exists() {
			var increment int
			if piece.color() == White {
				increment = 1
			} else {
				increment = -1
			}

			piecemapIndex := map0x88ToPiecemap(i, increment)

			switch piece.identity() {
			case King:
				score += (kingWeight + kingPositions[piecemapIndex]) * increment
			case Queen:
				score += (queenWeight + queenPositions[piecemapIndex]) * increment
			case Bishop:
				score += (bishopWeight + bishopPositions[piecemapIndex]) * increment
			case Rook:
				score += (rookWeight + rookPositions[piecemapIndex]) * increment
			case Knight:
				score += (knightWeight + knightPositions[piecemapIndex]) * increment
			case Pawn:
				score += (pawnWeight + pawnPositions[piecemapIndex]) * increment
			}
		}
	}

	return score * direction

}

// Map a 0x88 index to the required position in the score piecemaps. The
// direction (1 for white, -1 for black) is required because the indices are
// asymmetrical in 0x88 but symmetrical in the piecemap.
func map0x88ToPiecemap(index int, direction int) uint {
	rank := index / 16
	file := index % 16

	if direction == 1 {
		rank = 7 - rank
	}

	return uint(rank*8 + file)
}
