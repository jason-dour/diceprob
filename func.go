package diceprob

import (
	"math/rand"
	"sort"
	"time"
)

// New - Create a new DiceProb instance.
func New(s string) (*DiceProb, error) {
	// Create our object.
	obj := &DiceProb{
		expression:    s,
		parsed:        &Expression{},
		distribution:  &map[int64]int64{},
		probabilities: &map[int64]float64{},
		bounds:        &[]int64{},
		outcomes:      &[]int64{},
		permutations:  int64(0),
	}

	// Parse the expression and put it into the object.
	err := diceParser.ParseString("", obj.expression, obj.parsed)
	if err != nil {
		return nil, err
	}

	// Return the object.
	return obj, nil
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
