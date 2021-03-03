package diceprob

import (
	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer/stateful"
)

// Dice expression lexer.
var diceLexer = stateful.MustSimple([]stateful.Rule{
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
var diceParser = participle.MustBuild(&Expression{}, participle.Lexer(diceLexer), participle.UseLookahead(2))

// Operator type
type Operator int

// Operator constants
const (
	OpMul Operator = iota
	OpDiv
	OpAdd
	OpSub
)

// operatorMap - Map parsed operators to constants.
var operatorMap = map[string]Operator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv}

// Capture - Capture the costants while parsing.
func (o *Operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

// DiceRoll - String representing a dice roll atomic expression.
type DiceRoll string

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
	Modifier      *int64      `parser:"@Modifier"`
	RollExpr      *DiceRoll   `parser:"| @DiceRoll"`
	SubExpression *Expression `parser:"| '(' @@ ')'"`
}
