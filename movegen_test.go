package main

import (
	"testing"
)

type testCase struct {
	name          string
	fen           string
	expectedMoves int
}

type testCheck struct {
	fen           string
	expectedCheck bool
}

func TestCheck(t *testing.T) {
	cases := []testCheck{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", false},
		{"rnbqkb1r/ppp1pp1p/6p1/1B1n4/3P4/2N5/PP2PPPP/R1BQK1NR b KQkq - 0 1", true},
		{"rnbq1b1r/pppkpppp/3pPn2/8/2PP4/8/PP3PPP/RNBQKBNR w KQkq - 0 1", false},
		{"rnbq1b1r/pppkpppp/3pPn2/8/2PP4/8/PP3PPP/RNBQKBNR b KQkq - 0 1", true},
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		var attackingColor byte
		if position.toMove == White {
			attackingColor = Black
		} else {
			attackingColor = White
		}

		checked := isKingInCheck(position, attackingColor)

		if checked != test.expectedCheck {
			t.Errorf("Check test failed!\nFEN: %v\nExpected: %v\nActual: %v\n", test.fen, test.expectedCheck, checked)
		}
	}
}

func TestLegalMoveGen(t *testing.T) {
	cases := []testCase{
		// test cases generated from JetChess
		{"JetChess 1", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 48},
		{"JetChess 2", "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 14},
		{"JetChess 3", "n1rb4/1p3p1p/1p6/1R5K/8/p3p1PN/1PP1R3/N6k w - - 0 1", 28},
		{"JetChess 4", "5RKb/4P1n1/2p4p/3p2p1/3B2Q1/5B2/r6k/4r3 w - - 0 1", 47},
		{"JetChess 5", "7r/3B4/k7/8/6Qb/8/Kn6/6R1 w - - 0 1", 41},
		{"JetChess 6", "b1N1rb2/3p4/r6p/2Pp1p1K/3Pk3/2PN1p2/2B2P2/8 w - - 0 1", 17},
		{"Jetchess 7", "1kN2bb1/4r1r1/Q1P1p3/8/6n1/8/8/2B1K2B w - - 0 1", 34},
		{"Jetchess 8", "8/8/7K/6p1/NN5k/8/6PP/8 w - - 0 1", 16},
		{"Jetchess 9", "8/2p5/2Pb4/2pp3R/1ppk1pR1/2n2P1p/1B2PPpP/K7 w - - 0 1", 22},
		{"Jetchess 10", "8/6pp/8/p5p1/Pp6/1P3p2/pPK4P/krQ4R w - - 0 1", 18},
		{"Jetchess 11", "rnbqkbnr/ppp1pppp/8/3p4/2P5/8/PP1PPPPP/RNBQKBNR w KQkq - 0 2", 23},
		{"Jetchess 12", "rnbqkbnr/ppp1pppp/8/3p4/Q1P5/8/PP1PPPPP/RNB1KBNR b KQkq - 1 2", 6},
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		moves := generateLegalMoves(position)

		if len(moves) != test.expectedMoves {
			var formattedMoves []string

			for _, move := range moves {
				formattedMoves = append(formattedMoves, "("+toAlgebraic(position, move)+" "+toMoveString(move)+")")
			}
			t.Errorf("Legal move generation test (%v) failed!\nExpected: %v\nActual: %v\nMoves: %v\n", test.name, test.expectedMoves, len(moves), formattedMoves)

		}
	}

}
func TestMoveGen(t *testing.T) {
	cases := []testCase{
		{"Starting position", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 20},
		{"King movement", "8/5k2/8/8/3K4/8/8/8 w - - 0 1", 8},
		{"Rook sliding normal", "8/5k2/8/8/3R4/8/8/8 w KQkq - 0 1", 14},
		{"Rook attacking", "3p4/5k2/8/8/1p1R2p1/8/8/8 w KQkq - 0 1", 12},
		{"Bishop sliding normal", "8/7k/8/8/3B4/8/8/8 w KQkq - 0 1", 13},
		{"Queen sliding normal", "8/7k/8/8/3Q4/8/8/8 w KQkq - 0 1", 27},
		{"Knight in centre", "8/7k/8/8/3N4/8/8/8 w KQkq - 0 1", 8},
		{"Knight on edge", "8/7k/8/8/N7/8/8/8 w KQkq - 0 1", 4},

		{"Pawn at start", "8/7k/8/8/8/8/3P4/8 w KQkq - 0 1", 2},
		{"Pawn after moving", "8/7k/8/8/8/3P4/8/8 w KQkq - 0 1", 1},
		{"Pawn captures", "7k/8/8/8/8/2p1p3/3P4/8 w KQkq - 0 1", 4},
		{"En passant capture", "8/7k/8/3Pp3/8/8/8/8 w KQkq e6 0 1", 2},
		{"Two en passant options", "8/7k/8/3PpP2/8/8/8/8 w KQkq e6 0 1", 4},

		{"Black pawn at start", "8/4p3/8/8/8/8/8/4K3 b - - 0 1", 2},
		{"Black pawn captures", "8/4p3/3P1P2/8/8/8/8/4K3 b - - 0 1", 4},
		{"Black two en passant options", "8/8/8/8/4pPp1/8/8/4K3 b - f3 0 1", 4},

		{"Promotion", "8/2P4k/8/8/8/8/8/8 w KQkq - 0 1", 4},
		{"Capture promotion", "2q4k/3P4/8/8/8/8/8/8 w KQkq - 0 1", 8},

		{"Black promotion", "8/8/8/8/8/8/2p5/4K3 b - - 0 1", 4},

		{"Kingside castle", "8/5k2/8/8/8/8/8/4K2R w K - 0 1", 15},
		{"Both side castle", "8/5k2/8/8/8/8/8/R3K2R w KQ - 0 1", 26},
		{"No castle", "8/5k2/8/8/8/8/8/1R2K1R1 w - - 0 1", 24},

		{"Black kingside castle", "r2k4/8/8/8/8/8/8/4K3 b k - 0 1", 15},

		{"Black king", "8/5k2/8/8/8/8/8/1R2K1R1 b - - 0 1", 8},

		{"Can't castle through check", "6q1/8/8/8/8/8/8/4K2R w Kkq - 0 1", 14},
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		moves := generateMoves(position)

		if len(moves) != test.expectedMoves {
			t.Errorf("Move generation test (%v) failed!\nExpected: %v\nActual: %v\n", test.name, test.expectedMoves, len(moves))
		}
	}
}
