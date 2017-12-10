package main

import "testing"

type testCase struct {
	name          string
	fen           string
	expectedMoves int
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
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		moves := generateMoves(position)

		if len(moves) != test.expectedMoves {
			t.Errorf("Move generation test (%v) failed!\nExpected: %v\nActual: %v\n", test.name, test.expectedMoves, len(moves))
		}
	}
}
