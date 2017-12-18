package main

const KingWeight = 100
const QueenWeight = 9
const RookWeight = 5
const BishopWeight = 3
const KnightWeight = 3
const PawnWeight = 1

// return the negamax score for the given position. Eventually calculate this as
// the same time as move making or generation
// TODO: add mobility
func evaluate(position position) int {
	var kings int
	var queens int
	var rooks int
	var bishops int
	var knights int
	var pawns int

	for i := 0; i < BoardSize; i++ {

		piece := position.board[i]

		if isOnBoard(i) && position.board[i].exists() {
			var increment int
			if piece.color() == White {
				increment = 1
			} else {
				increment = -1
			}

			switch piece.identity() {
			case King:
				kings += increment
			case Queen:
				queens += increment
			case Bishop:
				bishops += increment
			case Rook:
				rooks += increment
			case Knight:
				knights += increment
			case Pawn:
				pawns += increment

			}
		}
	}

	var direction int
	if position.toMove == White {
		direction = 1
	} else {
		direction = -1
	}

	// formula from https://chessprogramming.wikispaces.com/Evaluation
	return (KingWeight*kings + QueenWeight*queens + RookWeight*rooks + BishopWeight*bishops + KnightWeight*knights + PawnWeight*pawns) * direction

}
