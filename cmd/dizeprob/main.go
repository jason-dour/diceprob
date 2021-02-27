// dizeprob - Calculate and display probabilities for a given dice expression.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jason-dour/diceprob"
)

func main() {
	dize, err := diceprob.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	dize.Calculate()

	fmt.Printf("Expression: %s\n", dize.InputExpression())
	fmt.Printf("Bounds: %v..%v\n", dize.Min(), dize.Max())
	fmt.Printf("Outcomes: %v\n", dize.TotalOutcomes())
	fmt.Printf("Outcome Set: %s\n", strings.Join(*dize.OutcomeListString(), ","))
	fmt.Printf("Distribution:\n  Outcome | Frequency | Probability\n")

	for _, i := range *dize.OutcomeList() {
		fmt.Printf("  %-8d  %-8d    %.6g\n", i, (*dize.Distribution)[i], (*dize.Probabilities)[i])
	}
}
