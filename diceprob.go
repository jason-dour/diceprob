// Package diceprob - Calculating probabilities and outcomes for complicated dice expressions.
package diceprob

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

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
	Expression string             // Expression provided when creating the instance.
	Parser     *participle.Parser // Participle parser for the expression.
	Parsed     *Expression        // Parsed expression data structure.
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
// Parser definitions.

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

//
// String functions; print the expression components as normalized text.

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

//
// Roll functions; calculate a "roll" of the expression.

// Roll - Roll a random value for the Expression; top-level of the recursive roll functions.
func (e *Expression) Roll() int64 {
	left := e.Left.Roll()
	for _, right := range e.Right {
		left = right.Operator.Roll(left, right.Term.Roll())
	}
	return left
}

// Roll - Roll a random values around the Operator; part of the recursive roll functions.
func (o Operator) Roll(l, r int64) int64 {
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
	sActual := string(*s)

	// Find the D in the roll.
	dToken := strings.Index(sActual, "d")
	if dToken == 0 {
		dToken = strings.Index(sActual, "D")
	}
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

//
// Calculate functions; calculate combinations, probabilities, bounds, etc.
// TODO - Need to write.

// rollIt - Using the selected method, roll n dice of s faces, and return the sum.
func rollIt(method string, n int64, s int64) int64 {
	// Seed the randomizer.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Depending on the method...
	switch method {
	// Mid rolling method.
	case "m":
		// Initialize the array of rolls.
		ret := []int64{}
		// Loop three times...
		for i := int64(1); i <= 3; i++ {
			// Appending a roll to the array.
			ret = append(ret, (r.Int63n(s) + 1))
		}
		// Sort the array numerically.
		sort.Slice(ret, func(i, j int) bool { return ret[i] < ret[j] })
		// Return the middle value.
		return ret[1]
	// "Standard" rolling method.
	case "d":
		// Initialize the return value.
		ret := int64(0)
		// Loop from 1 to n...
		for i := n; i <= n; i++ {
			// Add the value of the roll to the return value.
			ret = ret + (r.Int63n(s) + 1)
		}
		// Return the summed roll.
		return ret
	// Should not reach.
	default:
		panic("invalid rollIt method")
	}
}

// New - Create a new DiceProb instance.
func New(s string) (*DiceProb, error) {
	obj := &DiceProb{Expression: s, Parser: diceParser, Parsed: &Expression{}}

	err := obj.Parser.ParseString("", obj.Expression, obj.Parsed)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

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
