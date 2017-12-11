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
		{"Promotion", "rnbq1bnr/pppBP1p1/6kp/5p2/3Q4/8/PPP2PPP/RNB1K1NR w KQ - 0 9", "rnbqQbnr/pppB2p1/6kp/5p2/3Q4/8/PPP2PPP/RNB1K1NR b KQ - 0 9", createPromotionMove(100, 116, Queen)},
		{"En passant capture", "rnbqkbnr/pp1p2pp/5p2/2pPp3/4P3/8/PPP2PPP/RNBQKBNR w KQkq c6 0 4", "rnbqkbnr/pp1p2pp/2P2p2/4p3/4P3/8/PPP2PPP/RNBQKBNR b KQkq - 0 4", createEnPassantCaptureMove(67, 82)},
		{"Promotion capture", "rnbqkbnr/pP4pp/5p2/3pp3/4P3/8/PPP2PPP/RNBQKBNR w KQkq - 0 6", "rnNqkbnr/p5pp/5p2/3pp3/4P3/8/PPP2PPP/RNBQKBNR b KQkq - 0 6", createPromotionCaptureMove(97, 114, Knight)},
		{"Losing castle rights from king", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPPKPPP/RNBQ1BNR b kq - 1 2", createQuietMove(4, 20)},
		{"Losing castle rights from rook", "rnbqkbnr/ppppppp1/7p/8/8/6PP/PPPPPP2/RNBQKBNR b KQkq - 0 2", "rnbqkbn1/pppppppr/7p/8/8/6PP/PPPPPP2/RNBQKBNR w KQq - 1 3", createQuietMove(119, 103)},
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
