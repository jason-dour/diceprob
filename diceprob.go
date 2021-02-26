// Package diceprob - Calculating probabilities and outcomes for complicated dice expressions.
package diceprob

import (
	"fmt"
	"strings"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var (
	// Our dice expression lexer.
	diceLexer = stateful.MustSimple([]stateful.Rule{
		{Name: "DiceRoll", Pattern: `(\d+|[mM][iI])[dD](\d+|[fF])`, Action: nil},
		{Name: "Modifier", Pattern: `\d+`, Action: nil},
		{Name: "+", Pattern: `\+`, Action: nil},
		{Name: "-", Pattern: `-`, Action: nil},
		{Name: "*", Pattern: `\*`, Action: nil},
		{Name: "/", Pattern: `/`, Action: nil},
		{Name: "(", Pattern: `\(`, Action: nil},
		{Name: ")", Pattern: `\)`, Action: nil},
	})

	// Parser for our dice expressions.
	diceParser = participle.MustBuild(&Expression{}, participle.Lexer(diceLexer), participle.UseLookahead(2))
)

// DiceProb - Base data structure.
type DiceProb struct {
	expression string
	parser     *participle.Parser
	parsed     *Expression
}

//
// Operators - Capture and parse operators properly.

// Operator type
type Operator int

// Operator constants
const (
	OpMul Operator = iota
	OpDiv
	OpAdd
	OpSub
)

// operartorMap - Map parsed operators to constants.
var operatorMap = map[string]Operator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv}

// Capture - Capture the costants while parsing.
func (o *Operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

//
// Parser structures.

// Expression - Top level parsing unit.
type Expression struct {
	Left  *Term     `parser:"@@"`
	Right []*OpTerm `parser:"@@*"`
}

// OpTerm - Expression Operator and Term.
type OpTerm struct {
	Operator Operator `parser:"@('+' | '-')"`
	Term     *Term    `parser:"@@"`
}

// Term - Expression Term
type Term struct {
	Left  *Atom     `parser:"@@"`
	Right []*OpAtom `parser:"@@*"`
}

// OpAtom - Expression Operator and Atom.
type OpAtom struct {
	Operator Operator `parser:"@('*' | '/')"`
	Atom     *Atom    `parser:"@@"`
}

// Atom - Smallest unit of an expression.
type Atom struct {
	Modifier      *float64    `parser:"@Modifier"`
	DiceRoll      *string     `parser:"| @DiceRoll"`
	SubExpression *Expression `parser:"| '(' @@ ')'"`
}

//
// String functions; print the expression components as normalized text.
func (e *Expression) String() string {
	out := []string{e.Left.String()}
	for _, r := range e.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

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

func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

func (t *Term) String() string {
	out := []string{t.Left.String()}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (o *OpAtom) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Atom)
}

func (v *Atom) String() string {
	if v.Modifier != nil {
		return fmt.Sprintf("%g", *v.Modifier)
	}
	if v.DiceRoll != nil {
		return *v.DiceRoll
	}
	return "(" + v.SubExpression.String() + ")"
}

//
// Roll functions; calculate a "roll" of the expression.
func (e *Expression) Roll() float64 {
	left := e.Left.Roll()
	for _, right := range e.Right {
		left = right.Operator.Roll(left, right.Term.Roll())
	}
	return left
}

func (o Operator) Roll(l, r float64) float64 {
	switch o {
	case OpMul:
		return l * r
	case OpDiv:
		return l / r
	case OpAdd:
		return l + r
	case OpSub:
		return l - r
	}
	panic("unsupported operator")
}

func (t *Term) Roll() float64 {
	left := t.Left.Roll()
	for _, right := range t.Right {
		left = right.Operator.Roll(left, right.Atom.Roll())
	}
	return left
}

func (a *Atom) Roll() float64 {
	switch {
	// case a.Number != nil:
	// 	return *a.Number
	// case a.Variable != nil:
	// 	value, ok := ctx[*a.Variable]
	// 	if !ok {
	// 		panic("no such variable " + *a.Variable)
	// 	}
	// 	return value
	default:
		return a.SubExpression.Roll()
	}
}

// Calculate // TODO - Need to write.

// New - Create a new instance.
func New(s string) (*DiceProb, error) {
	obj := &DiceProb{expression: s, parser: diceParser, parsed: &Expression{}}

	err := obj.parser.ParseString("", obj.expression, obj.parsed)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// InputExpression - Return the original expression for the instance.
func (d *DiceProb) InputExpression() string {
	return d.expression
}

// ParsedExpression - Return the parsed expression
func (d *DiceProb) ParsedExpression() *Expression {
	return d.parsed
}

// Roll - Perform a "roll" of the expression and return the outcome
func (d *DiceProb) Roll() float64 {
	return d.parsed.Roll()
}
