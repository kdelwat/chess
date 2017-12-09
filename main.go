package main

import "fmt"

func main() {
	fmt.Printf("Hello, world\n")

	startPosition := position{board: startBoard}

	showPosition(startPosition)
	showSliding(startPosition)

	generateMoves(startPosition)
}
