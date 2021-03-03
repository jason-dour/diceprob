package diceprob

import (
	"sort"
	"strconv"
)

// Expression - Return the original expression for the instance.
func (d *DiceProb) Expression() string {
	return d.expression
}

// ParsedExpression - Return the parsed expression for the instance.
func (d *DiceProb) ParsedExpression() *Expression {
	return d.parsed
}

// Roll - Perform a "roll" of the expression and return the outcome.
func (d *DiceProb) Roll() int64 {
	return d.parsed.Roll()
}

// Min - Minimum outcome value for the expression's distribution.
func (d *DiceProb) Min() int64 {
	return (*d.bounds)[0]
}

// Max - Maximum outcome value for the expression's distribution.
func (d *DiceProb) Max() int64 {
	return (*d.bounds)[1]
}

// Bounds - Range min to max of  outcome values for the expression's distribution.
func (d *DiceProb) Bounds() *[]int64 {
	return d.bounds
}

// Outcomes - Return list of outcomes for the expression.
func (d *DiceProb) Outcomes() *[]int64 {
	return d.outcomes
}

// OutcomeListString - Return list of outcomes for the expression as strings.
func (d *DiceProb) OutcomeListString() *[]string {
	ret := []string{}
	for i := 0; i < len(*d.outcomes); i++ {
		ret = append(ret, strconv.FormatInt((*d.outcomes)[i], 10))
	}
	return &ret
}

// Permutations - Total outcomes for the expression.
func (d *DiceProb) Permutations() int64 {
	return d.permutations
}

// Distribution - Distribution of summed outcomes and their frequency.
func (d *DiceProb) Distribution() *map[int64]int64 {
	return d.distribution
}

// Probabilities - Probability of each outcome.
func (d *DiceProb) Probabilities() *map[int64]float64 {
	return d.probabilities
}

// Calculate - Calculate the Distribution and Probabilities for the ParsedExpression.
func (d *DiceProb) Calculate() {
	d.distribution = d.parsed.Distribution()

	keys := make([]int64, 0, len(*d.distribution))
	for k := range *d.distribution {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	d.outcomes = &keys
	d.bounds = &[]int64{keys[0], keys[len(keys)-1]}

	for _, frequency := range *d.distribution {
		d.permutations = d.permutations + frequency
	}

	for outcome, frequency := range *d.distribution {
		(*d.probabilities)[outcome] = float64(frequency) / float64(d.permutations)
	}

	return
}
