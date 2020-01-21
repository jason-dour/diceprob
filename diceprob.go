package diceprob

// DiceProb - Structure of our entire module.
type DiceProb struct {
	expression string
}

type DiceExpr struct {
	Nodes []*Nodes
}

// New - Create a new instance.
func New(s string) (*DiceProb, error) {
	return &DiceProb{expression: s}, nil
}

// Expression - Return the original expression for the instance.
func (d *DiceProb) Expression() string {
	return d.expression
}
