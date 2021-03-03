# diceprob

Go package for calculating combinations and probabilities for complex dice expressions.

## Purpose

`diceprob` is a personal combinatorics and game design hobby turned into Golang package.

For decades I've been fascinated with game design and specifically dice as the randomizer
element of those designs.  

I spent years researching combinatorics and probabilities of dice as a personal hobby.
Maybe that's weird.  Regardless, that eventually led to me writing code, and that initially
took form as a [Perl module `Games::Dice::Probability`](https://metacpan.org/pod/Games::Dice::Probability).

The `G::D::P` module was blazing fast, worked well for what I used it for, and seemed
to have enough users to have made the coding worthwhile.

Fast forward many many years, and I'm rewriting it as a Golang package to allow me to
continue with my hobby, and maybe enable others in the process.

## Install

Install with:

``` shell
go get github.com/jason-dour/diceprob
```

## Dice Expression Syntax

* `N`
  * The number of dice in a single dice roll.
  * Can be any integer number.
* `S`
  * The number of sides on the dice in a single dice roll.
  * Can be any integer number, as well as `F` or `f` for Fudge/FATE dice.
* `NdS`
  * Roll `N` dice, each with same number of sides `S`.
  * Examples:
    * 1d6
    * 3d6
    * 1d4
    * 1d20
    * 3df
* `midS`
  * Roll 3 dice, each with same number of sides `S`, and return the middle value of the three values.
  * Examples:
    * mid20
    * mid10
    * midf
* `[+ | - | * | /]`
  * Math operators; will add/subtract/multiply/divide the left and right terms.
  * Example:
    * 2d6+1d4
* `[0-9+]`
  * Modifier; a fixed number.
  * Example:
    * 2d6+1
    * 3d6-4
* `( expression )`
  * Grouping; you may use parentheses to enclose sub-expressions, to ensure proper calculation.
  * Example:
    * (1d6+2)*3

## Usage

Everything is driven through the `DiceProb` type, and its methods.

Create a new instance by providing your dice expression:

``` golang
d, err := diceprob.New("3d6")
if err != nil {
  panic(err)
}
```

Creating the instance will automatically parse the expression into an object tree.

``` golang
repr.Println(d.ParsedExpression())
```

``` text
&diceprob.Expression{
  Left: &diceprob.Term{
    Left: &diceprob.Atom{
      RollExpr: &diceprob.DiceRoll("3d6"),
    },
  },
}
```

From there you can calculate the outcomes, distribution, probabilities, et al.

``` golang
d.Calculate()
```

And call them for display or computation.

``` golang
fmt.Printf("Expression: %s\n", d.InputExpression())
fmt.Printf("Bounds: %v..%v\n", d.Min(), d.Max())
fmt.Printf("Permutations: %v\n", d.Permutations())
fmt.Printf("Outcome Set: %s\n", strings.Join(*d.OutcomesString(), ","))
fmt.Printf("Distribution:\n  Outcome | Frequency | Probability\n")

for _, i := range *d.Outcomes() {
  fmt.Printf("  %-8d  %-8d    %.6g\n", i, (*d.Distribution())[i], (*d.Probabilities())[i])
}
```

Or you can just "roll" the dice expression and retrieve a value.

``` golang
d.Roll()
```

## Notes

* Memoize for speed?
  * "golang.org/x/tools/internal/memoize"
