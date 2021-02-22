package main

import (
	"os"

	"github.com/alecthomas/repr"
	"github.com/jason-dour/diceprob"
)

func main() {
	dize, err := diceprob.New(os.Args[1])
	if err != nil {
		os.Exit(1)
	}
	repr.Println(dize.ParsedExpression())
}
