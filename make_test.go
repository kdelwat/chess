package main

import "testing"

type testMove struct {
	name   string
	fen    string
	newFen string
	move   move
}

func TestMakeUnmake(t *testing.T) {
	cases := []testMove{
		{"Quiet move", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", "rnbqkbnr/pppp1ppp/8/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR b KQkq - 1 2", createQuietMove(5, 50)},
		{"Capture", "rnbqkb1r/pppp1ppp/5n2/1B2p3/4P3/8/PPPP1PPP/RNBQK1NR b KQkq - 3 3", "rnbqkb1r/pppp1ppp/8/1B2p3/4n3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 4", createCaptureMove(85, 52)},
		{"Double pawn push", "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", createDoublePawnPush(20, 52)},
	}

	for _, test := range cases {
		position := fromFen(test.fen)

		artifacts := makeMove(&position, test.move)

		newFen := toFEN(position)

		if newFen != test.newFen {
			t.Errorf("Make move test failed (%v)!\nExpected: %v\nActual: %v\n", test.name, test.newFen, newFen)
		}

		unmakeMove(&position, test.move, artifacts)

		newFen = toFEN(position)

		if newFen != test.fen {
			t.Errorf("Unmake move test failed (%v)!\nExpected: %v\nActual: %v\n", test.name, test.fen, newFen)
		}
	}
}
