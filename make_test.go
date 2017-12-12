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
		{"Castle queenside", "rnbqkb1r/ppp1pppp/8/3p2B1/3Pn3/2N5/PPPQPPPP/R3KBNR w KQkq - 2 5", "rnbqkb1r/ppp1pppp/8/3p2B1/3Pn3/2N5/PPPQPPPP/2KR1BNR b kq - 3 5", move(QueenCastle)},
		{"Castle kingside", "rnbqk2r/ppp2ppp/3bpB2/3p4/3PN3/8/PPPQPPPP/2KR1BNR b kq - 2 7", "rnbq1rk1/ppp2ppp/3bpB2/3p4/3PN3/8/PPPQPPPP/2KR1BNR w - - 3 8", move(KingCastle)},

		{"Buggy case", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q2/PPPBBPpP/2R1K2R b Kkq - 1 2", "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q2/PPPBBP1P/2R1K1qR w Kkq - 0 3", createPromotionMove(22, 6, Queen)},
		{"Another bug", "8/8/8/8/k7/8/2Kp4/2R5 b - - 1 3", "8/8/8/8/8/1k6/2Kp4/2R5 w - - 2 4", createQuietMove(48, 33)},

		{"Promotion capture", "8/8/8/8/k7/8/2Kp4/2R5 b - - 1 3", "8/8/8/8/k7/8/2K5/2b5 w - - 0 4", createPromotionCaptureMove(19, 2, Bishop)},
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
