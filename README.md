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

## Usage

Create a new instance:

``` golang
dize, err := diceprob.New("3d6")
if err != nil {
  panic(err)
}
```

Creating the instance will automatically parse the expression into an object tree.

``` golang
repr.Println(dize.ParsedExpression())
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
dize.Calculate()
```

And call them for display or computation.

``` golang
fmt.Printf("Expression: %s\n", dize.InputExpression())
fmt.Printf("Bounds: %v..%v\n", dize.Min(), dize.Max())
fmt.Printf("Outcomes: %v\n", dize.TotalOutcomes())
fmt.Printf("Outcome Set: %s\n", strings.Join(*dize.OutcomeListString(), ","))
fmt.Printf("Distribution:\n  Outcome | Frequency | Probability\n")

for _, i := range *dize.OutcomeList() {
  fmt.Printf("  %-8d  %-8d    %.6g\n", i, (*dize.Distribution)[i], (*dize.Probabilities)[i])
}
```

Or you can just "roll" the dice expression and retrieve a value.

``` golang
dize.Roll()
```

## Notes

* Memoize for speed?
  * "golang.org/x/tools/internal/memoize"
