package diceprob

import (
	"testing"

	"github.com/alecthomas/repr"
)

func TestDiceProbNew(t *testing.T) {
	expected := "3d6"
	d, _ := New(expected)
	actual := d.InputExpression()
	if actual != expected {
		t.Fatalf("TestDiceProbNew: expected [%s], got [%s]\n", expected, actual)
	}
}

func TestDiceProbNewParsed(t *testing.T) {
	expr := DiceRoll("3d6")
	expected := repr.String(&Expression{
		Left: &Term{
			Left: &Atom{
				RollExpr: &expr,
			},
		},
	})
	d, _ := New("3d6")
	actual := repr.String(d.ParsedExpression())
	if actual != expected {
		t.Fatalf("TestDiceProbNewParsed:\n  expected: %s\n  actual: %s", expected, actual)
	}
}
