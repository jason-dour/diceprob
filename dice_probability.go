package dice_probability

type DiceProbability struct {
	expression string
}

func New(s string) (*DiceProbability, error) {
	retval := &DiceProbability{
		expression: s,
	}

	return retval, nil
}

func (d *DiceProbability) Expression() string {
	retval := d.expression

	return retval
}
