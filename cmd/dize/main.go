package main

import (
	"fmt"
	"os"

	"github.com/jason-dour/diceprob"
)

func main() {
	dize, err := diceprob.New("2d6")
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("Expression = %s\n", dize.Expression())
}
