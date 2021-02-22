package diceprob

import "testing"

func TestNewDiceProb(t *testing.T) {
	expected := "3d6"
	d, _ := New(expected)
	actual := d.InputExpression()
	if actual != expected {
		t.Fatalf("TestNewDiceProb: expected [%s], got [%s]\n", expected, actual)
	}
}
