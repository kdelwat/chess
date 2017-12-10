package main

import "testing"

func TestFEN(t *testing.T) {
	cases := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	}

	for _, fen := range cases {
		convertedPosition := fromFen(fen)
		convertedFen := toFEN(convertedPosition)

		if convertedFen != fen {
			t.Errorf("FEN roundtrip failed!\nInput: %v\nOutput: %v\n", fen, convertedFen)
		}
	}
}
