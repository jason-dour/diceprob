package diceprob

import (
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// Distribution - Determine the outcomes' distribution for the Expression; top-level of the recursive distribution functions.
func (e *Expression) Distribution() *map[int64]int64 {
	left := e.Left.distribution()
	for _, right := range e.Right {
		left = right.Operator.distribution(left, right.Term.distribution())
	}
	return left
}

// Distribution - Determine the outcomes' distribution around an Operator; part of the recursive distribution functions.
func (o Operator) distribution(left, right *map[int64]int64) *map[int64]int64 {
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
func (t *Term) distribution() *map[int64]int64 {
	left := t.Left.distribution()
	for _, right := range t.Right {
		left = right.Operator.distribution(left, right.Atom.distribution())
	}
	return left
}

// Distribution - Determine the outcomes' distribution for the Atom; part of the recursive distribution functions.
func (a *Atom) distribution() *map[int64]int64 {
	switch {
	case a.Modifier != nil:
		return &map[int64]int64{*a.Modifier: 1}
	case a.RollExpr != nil:
		return a.RollExpr.distribution()
	default:
		return a.SubExpression.Distribution()
	}
}

// Distribution - Determine the outcomes' distribution for the DiceRoll; deepest of the recursive distribution functions.
func (s *DiceRoll) distribution() *map[int64]int64 {
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

	// Convert the right side to an integer.
	rightInt := int64(0)
	if right == "f" {
		// Fudge/FATE dice equate to d3.
		rightInt = 3
	} else {
		// Otherwise convert the string to integer.
		err := error(nil)
		rightInt, err = strconv.ParseInt(right, 10, 64)
		if err != nil {
			panic(err)
		}
	}

	// Determine which kind of roll it is...
	switch left {
	case "mi":
		// "Middle" roll.

		// For each outcome in the set...
		for outcome := int64(1); outcome <= rightInt; outcome++ {
			// Calculate the number of combinations giving outcome as the middle value.
			retDist[outcome] = 1 + (3 * (rightInt - 1)) + (6 * (outcome - 1) * (rightInt - outcome))
		}

		if right == "f" {
			for i := int64(1); i <= 3; i++ {
				retDist[i-2] = retDist[i]
			}
			delete(retDist, 2)
			delete(retDist, 3)
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
			for outcome := int64(1); outcome <= rightInt; outcome++ {
				if right == "f" {
					retDist[outcome-2] = 1
				} else {
					retDist[outcome] = 1
				}
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
			// Calculated the mirrored outcome.
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

		// If Fudge/FATE dice, adjust the outcomes.
		if right == "f" {
			for i, j := int64(leftInt*-1), min; i <= leftInt; i, j = i+1, j+1 {
				retDist[i] = retDist[j]
			}
			for i := leftInt + 1; i <= max; i++ {
				delete(retDist, i)
			}
		}

		break
	}
	// Return the distribution.
	return &retDist
}
