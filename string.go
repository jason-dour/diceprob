package diceprob

import (
	"fmt"
	"strings"
)

// String - Output the dice Expression as a string; top level of recursive output functions.
func (e *Expression) String() string {
	out := []string{e.Left.string()}
	for _, r := range e.Right {
		out = append(out, r.string())
	}
	return strings.Join(out, " ")
}

// String - Output the Operator as a string; part of the recursive output functions.
func (o Operator) string() string {
	switch o {
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpSub:
		return "-"
	case OpAdd:
		return "+"
	}
	panic("unsupported operator") // TODO - We can do better here.
}

// String - Output the Operator and Term as a string; part of the recursive output functions.
func (o *OpTerm) string() string {
	return fmt.Sprintf("%s %s", o.Operator.string(), o.Term.string())
}

// String - Output the Term as a string; part of the recursive output functions.
func (t *Term) string() string {
	out := []string{t.Left.string()}
	for _, r := range t.Right {
		out = append(out, r.string())
	}
	return strings.Join(out, " ")
}

// String - Output the Operator and Atom as a string; part of the recursive output functions.
func (o *OpAtom) string() string {
	return fmt.Sprintf("%s %s", o.Operator.string(), o.Atom.string())
}

// String - Output the Atom as a string; part of the recursive output functions.
func (a *Atom) string() string {
	if a.Modifier != nil {
		return fmt.Sprintf("%d", *a.Modifier)
	}
	if a.RollExpr != nil {
		return a.RollExpr.string()
	}
	return "(" + a.SubExpression.String() + ")"
}

// String - Output the DiceRoll as a string; the deepest of the recursive output functions.
func (s *DiceRoll) string() string {
	ret := string(*s)
	return ret
}
