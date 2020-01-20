package dice_probability

import "testing"

func TestNewDiceProbability(t *testing.T) {
	expected := "3d6"
	d, _ := New(expected)
	actual, err := d.Expression()
	if err != nil {
		t.Fatalf("TestNewDiceProbability: error returning expression: %s", err.Error())
	}
	if actual != expected {
		t.Fatalf("TestNewDiceProbability: expected [%s], got [%s]\n", expected, actual)
	}
}
