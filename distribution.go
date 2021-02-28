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
