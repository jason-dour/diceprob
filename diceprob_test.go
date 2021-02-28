package diceprob

import (
	"reflect"
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

func TestDistribution(t *testing.T) {
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
	t.Logf("equal?=%v", eq)
	if !eq {
		t.Errorf("Distributions do not match.")
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

	d2, err := New("4d6")
	if err != nil {
		t.Errorf("Could not create new d2 instance.")
	}
	t.Logf("d2=%v", d2)
	d2.Calculate()
	t.Logf("d2.Distribution()=%v", d2.Distribution)

	eq := reflect.DeepEqual(d1.Distribution, d2.Distribution)
	t.Logf("equal?=%v", eq)
	if !eq {
		t.Errorf("Distribution of (2d6+2d6) does not match distribution of (4d6).")
	}
}
