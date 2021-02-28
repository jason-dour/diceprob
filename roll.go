package diceprob

import (
	"strconv"
	"strings"
)

// Roll - Roll a random value for the Expression; top-level of the recursive roll functions.
func (e *Expression) Roll() int64 {
	left := e.Left.Roll()
	for _, right := range e.Right {
		left = right.Operator.Roll(left, right.Term.Roll())
	}
	return left
}

// Roll - Roll a random values around the Operator; part of the recursive roll functions.
func (o Operator) Roll(left, right int64) int64 {
	switch o {
	case OpMul:
		return left * right
	case OpDiv:
		return left / right
	case OpAdd:
		return left + right
	case OpSub:
		return left - right
	}
	panic("unsupported operator") // TODO - We can do better here.
}

// Roll - Roll a random value for the Term; part of the recursive roll functions.
func (t *Term) Roll() int64 {
	left := t.Left.Roll()
	for _, right := range t.Right {
		left = right.Operator.Roll(left, right.Atom.Roll())
	}
	return left
}

// Roll - Roll a random value for the Atom; part of the recursive roll functions.
func (a *Atom) Roll() int64 {
	switch {
	case a.Modifier != nil:
		return *a.Modifier
	case a.RollExpr != nil:
		return a.RollExpr.Roll()
	default:
		return a.SubExpression.Roll()
	}
}

// Roll - Roll a random value for the DiceRoll; deepest of the recursive roll functions.
func (s *DiceRoll) Roll() int64 {
	// Convert s to a string.
	sActual := strings.ToLower(string(*s))

	// Find the D in the roll.
	dToken := strings.Index(sActual, "d")
	if dToken == -1 {
		panic("invalid dice roll atomic expression")
	}

	// Grab the digits to the right of the D.
	right, err := strconv.ParseInt(sActual[dToken+1:], 10, 64)
	if err != nil {
		panic(err)
	}

	// If the dice roll is a "middle" roll of 3 dice.
	if sActual[0:3] == "mid" {
		// Return a middle rolled value.
		return rollIt("m", 3, right)
	}
	// Not a "middle" roll, therefore a standard roll.

	// Grab the number of dice from the left of the D.
	left, err := strconv.ParseInt(sActual[0:dToken], 10, 64)
	if err != nil {
		panic(err)
	}

	// Return a standard rolled value.
	return rollIt("d", left, right)
}
