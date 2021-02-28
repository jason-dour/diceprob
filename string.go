package diceprob

import (
	"fmt"
	"strings"
)

// String - Output the dice Expression as a string; top level of recursive output functions.
func (e *Expression) String() string {
	out := []string{e.Left.String()}
	for _, r := range e.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

// String - Output the Operator as a string; part of the recursive output functions.
func (o Operator) String() string {
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
func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

// String - Output the Term as a string; part of the recursive output functions.
func (t *Term) String() string {
	out := []string{t.Left.String()}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

// String - Output the Operator and Atom as a string; part of the recursive output functions.
func (o *OpAtom) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Atom)
}

// String - Output the Atom as a string; part of the recursive output functions.
func (a *Atom) String() string {
	if a.Modifier != nil {
		return fmt.Sprintf("%d", *a.Modifier)
	}
	if a.RollExpr != nil {
		return a.RollExpr.String()
	}
	return "(" + a.SubExpression.String() + ")"
}

// String - Output the DiceRoll as a string; the deepest of the recursive output functions.
func (s *DiceRoll) String() string {
	ret := string(*s)
	return ret
}
