package main

import "fmt"

func main() {
	fmt.Println("Welcome to ultimate engine")

	position := fromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	depth := 5
	results := perft(position, depth)
	fmt.Printf("Perft with depth=%v\n", depth)
	fmt.Printf("nodes          : %v\n", results.nodes)
	fmt.Printf("quiet          : %v\n", results.quiet)
	fmt.Printf("captures       : %v\n", results.captures)
	fmt.Printf("enpassant      : %v\n", results.enpassant)
	fmt.Printf("promotion      : %v\n", results.promotion)
	fmt.Printf("promoCapture   : %v\n", results.promoCapture)
	fmt.Printf("castleKingSide : %v\n", results.castleKingSide)
	fmt.Printf("castleQueenSide: %v\n", results.castleQueenSide)
	fmt.Printf("pawnJump       : %v\n", results.pawnJump)
	fmt.Printf("checks         : %v\n", results.checks)
}
