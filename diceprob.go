// Package diceprob - Calculating outcome distributions and probabilities for complicated dice expressions.
package diceprob

// DiceProb - Base data structure.
type DiceProb struct {
	Expression    string             // Expression provided when creating the instance.
	Parsed        *Expression        // Parsed expression data structure.
	Outcome       *[]int64           // List of outcome values.
	Outcomes      int64              // Total number of outcomes.
	Distribution  *map[int64]int64   // Distribution of summed outcomes and their frequency.
	Probabilities *map[int64]float64 // Probability of each outcome.
	Bounds        *[]int64           // Min/Max Bounds of the outcomes.
}
