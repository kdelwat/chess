package main

import "testing"

type evaluateTest struct {
	fen      string
	expected int
}

func TestEvaluate(t *testing.T) {
	cases := []evaluateTest{
		{"r1b1k1nr/1p3ppp/pn1b4/8/1p2P3/2PP4/P1P2PPP/2BQK2R w Kkq - 0 11", -3},
		{"r1b1k1nr/1p3ppp/pn1b4/8/1p2P3/2PP4/P1P2PPP/2BQK2R b Kkq - 0 11", 3},
	}

	for _, test := range cases {
		position := fromFEN(test.fen)

		result := evaluate(position)

		if result != test.expected {
			t.Errorf("Evaluate test failed!\nFEN: %v\nExpected: %v\nActual: %v\n", test.fen, test.expected, result)
		}

	}
}
