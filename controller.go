package main

import (
	"bufio"
	"fmt"
	"os"
)

func startEngine() {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	for err == nil {
		fmt.Print(input)
		input, err = reader.ReadString('\n')
	}
}
