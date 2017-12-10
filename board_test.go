package main

import "testing"

func TestFEN(t *testing.T) {
	// random FEN strings from Wikipedia and custom
	cases := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",

		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b Kkq - 1 2",
		"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b - - 1 2",
	}

	for _, fen := range cases {
		convertedPosition := fromFen(fen)
		convertedFen := toFEN(convertedPosition)

		if convertedFen != fen {
			t.Errorf("FEN roundtrip failed!\nInput: %v\nOutput: %v\n", fen, convertedFen)
		}
	}
}
