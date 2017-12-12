package main

import "testing"

type testPerft struct {
	fen      string
	depth    int
	expected int
}

func TestPerft(t *testing.T) {
	cases := []testPerft{
		// Martin Sedlak's test positions
		// (http://www.talkchess.com/forum/viewtopic.php?t=47318)
		// code copied from Evert Glebbeek at http://www.talkchess.com/forum/viewtopic.php?topic_view=threads&p=657840&t=59046
		// avoid illegal ep
		{"3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", 6, 1134888},
		{"8/8/8/8/k1p4R/8/3P4/3K4 w - - 0 1", 6, 1134888},
		{"8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", 6, 1015133},
		{"8/b2p2k1/8/2P5/8/4K3/8/8 b - - 0 1", 6, 1015133},
		// en passant capture checks opponent:
		{"8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", 6, 1440467},
		{"8/5k2/8/2Pp4/2B5/1K6/8/8 w - d6 0 1", 6, 1440467},
		// short castling gives check:
		{"5k2/8/8/8/8/8/8/4K2R w K - 0 1", 6, 661072},
		{"4k2r/8/8/8/8/8/8/5K2 b k - 0 1", 6, 661072},
		// long castling gives check:
		//{"3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", 6, 803711},
		//{"r3k3/8/8/8/8/8/8/3K4 b q - 0 1", 6, 803711},
		// castling (including losing cr due to rook capture):
		{"r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", 4, 1274206},
		{"r3k2r/7b/8/8/8/8/1B4BQ/R3K2R b KQkq - 0 1", 4, 1274206},
		// castling prevented:
		{"r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", 4, 1720476},
		{"r3k2r/8/5Q2/8/8/3q4/8/R3K2R w KQkq - 0 1", 4, 1720476},
		// promote out of check:
		{"2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", 6, 3821001},
		{"3K4/8/8/8/8/8/4p3/2k2R2 b - - 0 1", 6, 3821001},
		// discovered check:
		{"8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", 5, 1004658},
		{"5K2/8/1Q6/2N5/8/1p2k3/8/8 w - - 0 1", 5, 1004658},
		// promote to give check:
		{"4k3/1P6/8/8/8/8/K7/8 w - - 0 1", 6, 217342},
		{"8/k7/8/8/8/8/1p6/4K3 b - - 0 1", 6, 217342},
		// underpromote to check:
		{"8/P1k5/K7/8/8/8/8/8 w - - 0 1", 6, 92683},
		{"8/8/8/8/8/k7/p1K5/8 b - - 0 1", 6, 92683},
		// self stalemate:
		{"K1k5/8/P7/8/8/8/8/8 w - - 0 1", 6, 2217},
		{"8/8/8/8/8/p7/8/k1K5 b - - 0 1", 6, 2217},
		// stalemate/checkmate:
		{"8/k1P5/8/1K6/8/8/8/8 w - - 0 1", 7, 567584},
		{"8/8/8/8/1k6/8/K1p5/8 b - - 0 1", 7, 567584},
		// double check:
		{"8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", 4, 23527},
		{"8/5k2/8/5N2/5Q2/2K5/8/8 w - - 0 1", 4, 23527},
		// short castling impossible although the rook never moved away from its corner
		{"1k6/1b6/8/8/7R/8/8/4K2R b K - 0 1", 5, 1063513},
		{"4k2r/8/8/7r/8/8/1B6/1K6 w k - 0 1", 5, 1063513},
		// long castling impossible although the rook never moved away from its corner
		{"1k6/8/8/8/R7/1n6/8/R3K3 b Q - 0 1", 5, 346695},
		{"r3k3/8/1N6/r7/8/8/8/1K6 w q - 0 1", 5, 346695},
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		results := perft(position, test.depth)

		if int(results.nodes) != test.expected {
			t.Errorf("Perft test failed!\nFEN: %v\nExpected: %v\nActual: %v\n", test.fen, test.expected, results.nodes)
		}

	}
}
