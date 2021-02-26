// dizeroll - "Roll" a given dice expression and display the outcome.
package main

import (
	"os"

	"github.com/alecthomas/repr"
	"github.com/jason-dour/diceprob"
)

func main() {
	dize, err := diceprob.New(os.Args[1])
	if err != nil {
		panic(err)
	}
	// repr.Println(dize.ParsedExpression())
	repr.Println(dize.ParsedExpression().Roll())
}
