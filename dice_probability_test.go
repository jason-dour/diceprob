package dice_probability

import "testing"

func TestNewDiceProbability(t *testing.T) {
	expected := "3d6"
	d, _ := New(expected)
	actual := d.Expression()
	if actual != expected {
		t.Fatalf("TestNewDiceProbability: expected [%s], got [%s]\n", expected, actual)
	}
}
