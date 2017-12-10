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
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		moves := generateMoves(position)

		if len(moves) != test.expectedMoves {
			t.Errorf("Move generation test (%v) failed!\nExpected: %v\nActual: %v\n", test.name, test.expectedMoves, len(moves))
		}
	}
}
