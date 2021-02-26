# diceprob

Go Library for calculating combinations and probabilities for complex dice expressions.

## Purpose

`diceprob` is a personal combinatorics and game design hobby turned into Golang package.

For decades I've been fascinated with game design and specifically dice as the randomizer
element of those designs.  

I spent years researching combinatorics and probabilities of dice as a personal hobby.
Maybe that's weird.  Regardless, that eventually led to me writing code, and that initially
took form as a Perl module `Games::Dice::Probability`.

The `G::D::P` module was blazing fast, worked well for what I used it for, and seemed
to have enough users to have made the coding worthwhile.

Fast forward many many years, and I'm rewriting it as a Golang package to allow me to
continue with my hobby, and maybe enable others in the process.

## Notes

* golang memoize
  * import "golang.org/x/tools/internal/memoize"
* golang parser
  * import "github.com/alecthomas/participle"
* binom coeff
  * import "math/big"
  * Binomial(n,k int64)
