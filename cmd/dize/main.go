package main

import (
	"fmt"
	"os"

	dice_probability "github.com/jason-dour/go-dice-probability"
)

func main() {
	dize, err := dice_probability.New("2d6")
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Expression = %s", dize.Expression())
}
