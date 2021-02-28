package diceprob

import (
	"sort"
	"strconv"
)

// InputExpression - Return the original expression for the instance.
func (d *DiceProb) InputExpression() string {
	return d.Expression
}

// ParsedExpression - Return the parsed expression for the instance.
func (d *DiceProb) ParsedExpression() *Expression {
	return d.Parsed
}

// Roll - Perform a "roll" of the expression and return the outcome.
func (d *DiceProb) Roll() int64 {
	return d.Parsed.Roll()
}

// Min - Minimum outcome value for the expression's distribution.
func (d *DiceProb) Min() int64 {
	return (*d.Bounds)[0]
}

// Max - Maximum outcome value for the expression's distribution.
func (d *DiceProb) Max() int64 {
	return (*d.Bounds)[1]
}

// OutcomeList - Return list of outcomes for the expression.
func (d *DiceProb) OutcomeList() *[]int64 {
	return d.Outcome
}

// OutcomeListString - Return list of outcomes for the expression as strings.
func (d *DiceProb) OutcomeListString() *[]string {
	ret := []string{}
	for i := 0; i < len(*d.Outcome); i++ {
		ret = append(ret, strconv.FormatInt((*d.Outcome)[i], 10))
	}
	return &ret
}

// TotalOutcomes - Total outcomes for the expression.
func (d *DiceProb) TotalOutcomes() int64 {
	return d.Outcomes
}

// Calculate - Calculate the Distribution and Probabilities for the ParsedExpression.
func (d *DiceProb) Calculate() {
	d.Distribution = d.Parsed.Distribution()

	keys := make([]int64, 0, len(*d.Distribution))
	for k := range *d.Distribution {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	d.Outcome = &keys
	d.Bounds = &[]int64{keys[0], keys[len(keys)-1]}

	for _, frequency := range *d.Distribution {
		d.Outcomes = d.Outcomes + frequency
	}

	for outcome, frequency := range *d.Distribution {
		(*d.Probabilities)[outcome] = float64(frequency) / float64(d.Outcomes)
	}

	return
}
