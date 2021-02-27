// Package diceprob - Calculating probabilities and outcomes for complicated dice expressions.
package diceprob

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"regexp"
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
	Expression    string             // Expression provided when creating the instance.
	Parser        *participle.Parser // Participle parser for the expression.
	Parsed        *Expression        // Parsed expression data structure.
	Outcome       *[]int64           // List of outcome values.
	Outcomes      int64              // Total number of outcomes.
	Distribution  *map[int64]int64   // Distribution of summed outcomes and their frequency.
	Probabilities *map[int64]float64 // Probability of each outcome.
	Bounds        *[]int64           // Min/Max Bounds of the outcomes.
}

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

// Distribution - Determine the outcomes' distribution for the Expression; top-level of the recursive distribution functions.
func (e *Expression) Distribution() *map[int64]int64 {
	left := e.Left.Distribution()
	for _, right := range e.Right {
		left = right.Operator.Distribution(left, right.Term.Distribution())
	}
	return left
}

// Distribution - Determine the outcomes' distribution around an Operator; part of the recursive distribution functions.
func (o Operator) Distribution(left, right *map[int64]int64) *map[int64]int64 {
	combined := map[int64]int64{}

	for outcome1, freq1 := range *left {
		for outcome2, freq2 := range *right {
			outcomeNew := int64(0)
			switch o {
			case OpMul:
				outcomeNew = outcome1 * outcome2
			case OpDiv:
				outcomeNew = outcome1 / outcome2
			case OpAdd:
				outcomeNew = outcome1 + outcome2
			case OpSub:
				outcomeNew = outcome1 - outcome2
			}
			combined[outcomeNew] = combined[outcomeNew] + (freq1 * freq2)
		}
	}

	return &combined
}

// Distribution - Determine the outcomes' distribution for the Term; part of the recursive distribution functions.
func (t *Term) Distribution() *map[int64]int64 {
	left := t.Left.Distribution()
	for _, right := range t.Right {
		left = right.Operator.Distribution(left, right.Atom.Distribution())
	}
	return left
}

// Distribution - Determine the outcomes' distribution for the Atom; part of the recursive distribution functions.
func (a *Atom) Distribution() *map[int64]int64 {
	switch {
	case a.Modifier != nil:
		return &map[int64]int64{*a.Modifier: 1}
	case a.RollExpr != nil:
		return a.RollExpr.Distribution()
	default:
		return a.SubExpression.Distribution()
	}
}

// Distribution - Determine the outcomes' distribution for the DiceRoll; deepest of the recursive distribution functions.
func (s *DiceRoll) Distribution() *map[int64]int64 {
	// Convert s to a string.
	sActual := strings.ToLower(string(*s))

	// Prepare for the distribution.
	retDist := map[int64]int64{}

	// Prepare a regex to parse the dice roll.
	re := regexp.MustCompile(`(?P<left>\d+|mi)d(?P<right>\d+|f)`)

	// Parse the roll syntax.
	m := re.FindStringSubmatch(sActual)
	left := m[re.SubexpIndex("left")]
	right := m[re.SubexpIndex("right")]

	switch right {
	case "f":
		// TODO - All Fudge/FATE dice code needs an overhaul. Ugh.
		switch left {
		case "mi":
			// If the dice roll is a "middle" roll of 3 dice.
			// TODO - midF?  Haven't done that before; need to craft that logic.
			// if sActual[0:3] == "mid" {
			// 	// For each outcome in 1..s, calculate the number of combinations giving outcome as the middle value.
			// 	for outcome := int64(1); outcome <= right; outcome++ {
			// 		retDist[outcome] = 1 + (3 * (right - 1)) + (6 * (outcome - 1) * (right - outcome))
			// 	}
			// 	return &retDist
			// }
			break
		default:
			// Standard Fudge/Fate dice roll.
			break
		}
		break
	default:
		rightInt, err := strconv.ParseInt(right, 10, 64)
		if err != nil {
			panic(err)
		}
		switch left {
		case "mi":
			// If the dice roll is a "middle" roll of 3 dice.
			if sActual[0:3] == "mid" {
				// For each outcome in 1..s, calculate the number of combinations giving outcome as the middle value.
				for outcome := int64(1); outcome <= rightInt; outcome++ {
					retDist[outcome] = 1 + (3 * (rightInt - 1)) + (6 * (outcome - 1) * (rightInt - outcome))
				}
			}
			break
		default:
			// Standard dice roll. Convert left to Int64.
			leftInt, err := strconv.ParseInt(left, 10, 64)
			if err != nil {
				panic(err)
			}

			// Save effort if only one die...
			if leftInt == 1 {
				// Each outcome only occurs once.
				for outcome := int64(1); outcome <= rightInt; outcome++ {
					retDist[outcome] = 1
				}
				break
			}

			// More than 1 die!
			// Calculate min, max
			min := leftInt
			max := leftInt * rightInt
			// And peak around which we will mirror the distribution.
			peak := (min + max) / 2

			// For every outcome from min to peak...
			for outcome := min; outcome <= peak; outcome++ {
				// Caculated the mirrored outcome.
				reflected := min + max - outcome
				// Determine the ceiling of the sum function.
				ceiling := (outcome - leftInt) / rightInt
				// Initialize the frequency.
				frequency := int64(0)
				// For 0 to ceiling, sum the frequencies.
				for i := int64(0); i <= ceiling; i++ {
					part1 := big.NewInt(0)
					part2 := big.NewInt(0)
					frequency = frequency + (int64(math.Pow(-1, float64(i))) *
						part1.Binomial(leftInt, i).Int64() *
						part2.Binomial((outcome-(rightInt*i)-1), (leftInt-1)).Int64())
				}
				// Assign the outcome...
				retDist[outcome] = frequency
				// ...and its mirror.
				retDist[reflected] = frequency
			}

			break
		}
		break
	}
	// Return the distribution.
	return &retDist
}

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
		for i := int64(1); i <= n; i++ {
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
	// Create our object.
	obj := &DiceProb{
		Expression:    s,
		Parser:        diceParser,
		Parsed:        &Expression{},
		Distribution:  &map[int64]int64{},
		Probabilities: &map[int64]float64{},
		Bounds:        &[]int64{},
		Outcome:       &[]int64{},
		Outcomes:      int64(0),
	}

	// Parse the expression and put it into the object.
	err := obj.Parser.ParseString("", obj.Expression, obj.Parsed)
	if err != nil {
		return nil, err
	}

	// Return the object.
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
