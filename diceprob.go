// Package diceprob - Calculating outcome distributions and probabilities for complicated dice expressions.
package diceprob

// DiceProb - Base data structure.
type DiceProb struct {
	expression    string             // Expression provided when creating the instance.
	parsed        *Expression        // Parsed expression data structure.
	outcomes      *[]int64           // List of outcome values.
	permutations  int64              // Total number of outcomes.
	distribution  *map[int64]int64   // Distribution of summed outcomes and their frequency.
	probabilities *map[int64]float64 // Probability of each outcome.
	bounds        *[]int64           // Min/Max Bounds of the outcomes.
}
