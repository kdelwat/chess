package main

import "testing"
import "strings"

type evaluateTest struct {
	name     string
	expected int
	fen      string
}

func TestEvaluate(t *testing.T) {
	cases := []evaluateTest{
		{"Pawn testing", 330, "8/8/8/8/4P3/3P4/2P5/8 w KQkq - 0 11"},
		{"Knight, rook, bishop", -685, "8/5n2/r2r4/8/8/6B1/3B4/8 w KQkq - 0 1"},
		{"Asymetrical kings", 60, "8/8/8/8/8/8/8/K4k2 w KQkq - 0 1"},
		{"​Starting position", 0, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"},
	}

	for _, test := range cases {

		position := fromFEN(test.fen)
		result := evaluate(position)

		if result != test.expected {
			t.Errorf("Evaluate test failed! (%v)\nFEN: %v\nExpected: %v\nActual: %v\n", test.name, test.fen, test.expected, result)
		}

		mirroredPosition := fromFEN(strings.Replace(test.fen, "w", "b", 1))

		mirroredResult := evaluate(mirroredPosition)

		if mirroredResult != -test.expected {
			t.Errorf("Mirrored evaluate test failed! (%v)\nFEN: %v\nExpected: %v\nActual: %v\n", test.name, test.fen, -test.expected, result)
		}
	}
}
