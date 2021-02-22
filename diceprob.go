package diceprob

import (
	"fmt"
	"strings"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var (
	diceLexer = stateful.MustSimple([]stateful.Rule{
		{Name: "DiceRoll", Pattern: `(\d+|[mM][iI])[dD](\d+|[fF])`, Action: nil},
		{Name: "Modifier", Pattern: `\d+`, Action: nil},
		{Name: "+", Pattern: "\\+", Action: nil},
		{Name: "-", Pattern: "\\-", Action: nil},
		{Name: "*", Pattern: "\\*", Action: nil},
		{Name: "/", Pattern: "\\/", Action: nil},
		{Name: "(", Pattern: "\\(", Action: nil},
		{Name: ")", Pattern: "\\)", Action: nil},
	})
	diceParser = participle.MustBuild(&Expression{}, participle.Lexer(diceLexer), participle.UseLookahead(2))
)

// DiceProb - Structure of our entire module.
type DiceProb struct {
	expression string
	parser     *participle.Parser
}

// Let's model the following syntax:
//
// expression:	add_sub end  { $item[1] }
// add_sub:	mult_div '+' add_sub { { left => $item[1], op => '+', right => $item[3] } }
// add_sub:	mult_div '-' add_sub { { left => $item[1], op => '-', right => $item[3] } }
// add_sub:	mult_div
// mult_div:	bracket '/' mult_div { { left => $item[1], op => '/', right => $item[3] } }
// mult_div:	bracket '*' mult_div { { left => $item[1], op => '*', right => $item[3] } }
// mult_div:	bracket
// bracket:	'(' add_sub ')' { $item[2] }
// bracket:	dicenode
// dicenode:	/(\d+|mi)d(\d+|f)/i
// dicenode:	/\d+/
// end:		/\s*$/

// Operator - TODO
type Operator int

// Operators
const (
	OpMul Operator = iota
	OpDiv
	OpAdd
	OpSub
)

var operatorMap = map[string]Operator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv}

// Capture - TODO
func (o *Operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

// Expression - Top level parsing unit.
type Expression struct {
	Left  *Term     `parser:"@@"`
	Right []*OpTerm `parser:"@@*"`
}

// OpTerm - TODO
type OpTerm struct {
	Operator Operator `parser:"@('+' | '-')"`
	Term     *Term    `parser:"@@"`
}

// Term - TODO
type Term struct {
	Left  *Atom     `parser:"@@"`
	Right []*OpAtom `parser:"@@*"`
}

// OpAtom - TODO
type OpAtom struct {
	Operator Operator `parser:"@('*' | '/')"`
	Atom     *Atom    `parser:"@@*"`
}

// Atom - Smallest unit of an expression.
type Atom struct {
	Modifier      *float64    `parser:"@Modifier"`
	DiceRoll      *string     `parser:"| @DiceRoll"`
	SubExpression *Expression `parser:"| '(' @@ ')'"`
}

// TODO

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
	panic("unsupported operator")
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

func (o *OpAtom) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Atom)
}

func (t *Term) String() string {
	out := []string{t.Left.String()}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

func (e *Expression) String() string {
	out := []string{e.Left.String()}
	for _, r := range e.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

// Evaluation

// func (o Operator) Eval(l, r float64) float64 {
// 	switch o {
// 	case OpMul:
// 		return l * r
// 	case OpDiv:
// 		return l / r
// 	case OpAdd:
// 		return l + r
// 	case OpSub:
// 		return l - r
// 	}
// 	panic("unsupported operator")
// }

// func (v *Atom) Eval(ctx Context) float64 {
// 	switch {
// 	case v.Number != nil:
// 		return *v.Number
// 	case v.Variable != nil:
// 		value, ok := ctx[*v.Variable]
// 		if !ok {
// 			panic("no such variable " + *v.Variable)
// 		}
// 		return value
// 	default:
// 		return v.Subexpression.Eval(ctx)
// 	}
// }

// func (t *Term) Eval(ctx Context) float64 {
// 	n := t.Left.Eval(ctx)
// 	for _, r := range t.Right {
// 		n = r.Operator.Eval(n, r.Factor.Eval(ctx))
// 	}
// 	return n
// }

// func (e *Expression) Eval(ctx Context) float64 {
// 	l := e.Left.Eval(ctx)
// 	for _, r := range e.Right {
// 		l = r.Operator.Eval(l, r.Term.Eval(ctx))
// 	}
// 	return l
// }

// TODO

// New - Create a new instance.
func New(s string) (*DiceProb, error) {
	return &DiceProb{expression: s, parser: diceParser}, nil
}

// InputExpression - Return the original expression for the instance.
func (d *DiceProb) InputExpression() string {
	return d.expression
}

// ParsedExpression - Return the parsed expression
func (d *DiceProb) ParsedExpression() *Expression {
	expr := &Expression{}
	err := d.parser.ParseString("", d.expression, expr)
	if err != nil {
		panic(err)
	}

	return expr
}
