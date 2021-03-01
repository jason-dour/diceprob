package diceprob

import (
	"reflect"
	"sort"
	"testing"

	"github.com/alecthomas/repr"
)

func TestNew(t *testing.T) {
	expected := "3d6"
	t.Logf("expected=%v", expected)
	d, err := New(expected)
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	actual := d.InputExpression()
	t.Logf("actual=%v", actual)
	if actual != expected {
		t.Errorf("Actual value does not match expected value.")
	}
}

func TestParsed(t *testing.T) {
	expr := DiceRoll("3d6")
	expected := repr.String(&Expression{
		Left: &Term{
			Left: &Atom{
				RollExpr: &expr,
			},
		},
	})
	t.Logf("expected=%v", expected)
	d, err := New("3d6")
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	actual := repr.String(d.ParsedExpression())
	t.Logf("actual=%v", actual)
	if actual != expected {
		t.Errorf("Actual value does not match expected value.")
	}
}

func TestRoll(t *testing.T) {
	d, err := New("3d6")
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	d.Calculate()
	t.Logf("expected=(%v <= Roll() <= %v)", d.Min(), d.Max())
	actual := d.Roll()
	t.Logf("expected=(%v <= %v <= %v)", d.Min(), actual, d.Max())
	if !((actual >= d.Min()) && (actual <= d.Max())) {
		t.Errorf("Rolled value outside of bounds.")
	}
}

func TestDistributionNdS(t *testing.T) {
	d, err := New("3d6")
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	d.Calculate()
	expected := map[int64]int64{
		3:  1,
		4:  3,
		5:  6,
		6:  10,
		7:  15,
		8:  21,
		9:  25,
		10: 27,
		11: 27,
		12: 25,
		13: 21,
		14: 15,
		15: 10,
		16: 6,
		17: 3,
		18: 1,
	}
	t.Logf("expected=%v", expected)
	actual := *d.Distribution
	t.Logf("actual=%v", actual)
	eq := reflect.DeepEqual(expected, actual)
	t.Logf("DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Calculated distribution does not match the control distribution.")
	}
}

func TestDistributionNdF(t *testing.T) {
	d, err := New("3dF")
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	d.Calculate()
	expected := map[int64]int64{
		-3: 1,
		-2: 3,
		-1: 6,
		0:  7,
		1:  6,
		2:  3,
		3:  1,
	}
	t.Logf("expected=%v", expected)
	actual := *d.Distribution
	t.Logf("actual=%v", actual)
	eq := reflect.DeepEqual(expected, actual)
	t.Logf("DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Calculated distribution does not match the control distribution.")
	}
}

func TestDistributionMidS(t *testing.T) {
	d, err := New("mid20")
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	d.Calculate()
	expected := map[int64]int64{
		1:  58,
		2:  166,
		3:  262,
		4:  346,
		5:  418,
		6:  478,
		7:  526,
		8:  562,
		9:  586,
		10: 598,
		11: 598,
		12: 586,
		13: 562,
		14: 526,
		15: 478,
		16: 418,
		17: 346,
		18: 262,
		19: 166,
		20: 58,
	}
	t.Logf("expected=%v", expected)
	actual := *d.Distribution
	t.Logf("actual=%v", actual)
	eq := reflect.DeepEqual(expected, actual)
	t.Logf("DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Calculated distribution does not match the control distribution.")
	}
}

func TestDistributionMidF(t *testing.T) {
	d, err := New("midF")
	if err != nil {
		t.Errorf("Could not create new instance.")
	}
	t.Logf("d=%v", d)
	d.Calculate()
	expected := map[int64]int64{
		-1: 7,
		0:  13,
		1:  7,
	}
	t.Logf("expected=%v", expected)
	actual := *d.Distribution
	t.Logf("actual=%v", actual)
	eq := reflect.DeepEqual(expected, actual)
	t.Logf("DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Calculated distribution does not match the control distribution.")
	}
}

func TestCombinedDistributions(t *testing.T) {
	d1, err := New("2d6+2d6")
	if err != nil {
		t.Errorf("Could not create new d1 instance.")
	}
	t.Logf("d1=%v", d1)
	d1.Calculate()
	t.Logf("d1.Distribution()=%v", d1.Distribution)
	t.Logf("d1.Probabilities()=%v", d1.Probabilities)

	d2, err := New("4d6")
	if err != nil {
		t.Errorf("Could not create new d2 instance.")
	}
	t.Logf("d2=%v", d2)
	d2.Calculate()
	t.Logf("d2.Distribution()=%v", d2.Distribution)
	t.Logf("d2.Probabilities()=%v", d2.Probabilities)

	eq := reflect.DeepEqual(d1.Distribution, d2.Distribution)
	t.Logf("Distribution.DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Distribution of (%s) does not match distribution of (%s).", d1.Parsed.String(), d2.Parsed.String())
	}

	eq = reflect.DeepEqual(d1.Probabilities, d2.Probabilities)
	t.Logf("Probabilities.DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Probabilities of (%s) do not match Probabilities of (%s).", d1.Parsed.String(), d2.Parsed.String())
	}
}

func TestSameProbabilities(t *testing.T) {
	d1, err := New("2d6-2d6")
	if err != nil {
		t.Errorf("Could not create new d1 instance.")
	}
	t.Logf("d1=%v", d1)
	d1.Calculate()
	t.Logf("d1.Distribution()=%v", d1.Distribution)

	d2, err := New("4d6")
	if err != nil {
		t.Errorf("Could not create new d2 instance.")
	}
	t.Logf("d2=%v", d2)
	d2.Calculate()
	t.Logf("d2.Distribution()=%v", d2.Distribution)

	eq := reflect.DeepEqual(d1.Distribution, d2.Distribution)
	t.Logf("Distribution.DeepEqual?=%v", eq)
	if eq {
		t.Errorf("Distribution of (%s) matches distribution of (%s).", d1.Parsed.String(), d2.Parsed.String())
	}

	k1 := reflect.ValueOf(*d1.Probabilities).MapKeys()
	k11 := make([]int64, len(k1))
	for i := 0; i < len(k1); i++ {
		k11[i] = k1[i].Int()
	}
	sort.Slice(k11, func(i, j int) bool { return k11[i] < k11[j] })

	k2 := reflect.ValueOf(*d2.Probabilities).MapKeys()
	k22 := make([]int64, len(k1))
	for i := 0; i < len(k2); i++ {
		k22[i] = k2[i].Int()
	}
	sort.Slice(k22, func(i, j int) bool { return k22[i] < k22[j] })

	p1 := make([]float64, len(k1))
	p2 := make([]float64, len(k2))
	for i := 0; i < len(k1); i++ {
		p1[i] = (*d1.Probabilities)[k11[i]]
		p2[i] = (*d2.Probabilities)[k22[i]]
	}
	t.Logf("d1.Probabilities()->Values=%v", p1)
	t.Logf("d2.Probabilities()->Values=%v", p2)
	eq = reflect.DeepEqual(p1, p2)
	t.Logf("Probabilities.DeepEqual?=%v", eq)
	if !eq {
		t.Errorf("Probabilities of (%s) do not match Probabilities of (%s).", d1.Parsed.String(), d2.Parsed.String())
	}
}
