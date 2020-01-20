package diceprob

// DiceProbability - Structure of our entire module.
type DiceProbability struct {
	expression string
}

// New - Create a new instance.
func New(s string) (*DiceProbability, error) {
	retval := &DiceProbability{
		expression: s,
	}
	return retval, nil
}

// Expression - Return the original expression for the instance.
func (d *DiceProbability) Expression() string {
	retval := d.expression

	return retval
}
