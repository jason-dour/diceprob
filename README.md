# diceprob

Go Library for calculating combinations and probabilities for complex dice expressions.

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

From there you can calculate the combinations, probabilities, and other metrics.

``` golang
// COMING SOON
```

Or you can just "roll" the dice expression and retrieve a value.

``` golang
dize.ParsedExpression().Roll()
```

## Notes

* golang memoize
  * import "golang.org/x/tools/internal/memoize"
* golang parser
  * import "github.com/alecthomas/participle"
* binom coeff
  * import "math/big"
  * Binomial(n,k int64)
