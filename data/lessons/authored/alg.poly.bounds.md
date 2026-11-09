# Upper and lower bounds for roots

Before a root is found, its possible size can be bounded. Two classical tools do that work. The rational root theorem narrows the search to a finite list of candidates, and Cauchy's bound puts every real root inside a symmetric interval.

## Divisors of the constant term

For a monic polynomial with integer coefficients, every integer root divides the constant term. If \(P(x) = x^3 - 4x + 6\), then any integer root must be a divisor of 6. The positive divisors are 1, 2, 3 and 6, and each one has a negative partner.

The test is quick because a divisor list is short. It also tells you when to stop looking: if none of the divisors makes the polynomial zero, the polynomial has no integer roots at all.

The statement needs the polynomial to be monic. When the leading coefficient is not one, the rational root theorem replaces the integer test with a fraction test.

## Cauchy's bound

Cauchy's bound says every real root satisfies \(|x| < 1 + \max |a_i| / |a_n|\), where \(a_n\) is the leading coefficient and the maximum runs over the other coefficients. The right side is an upper bound on the absolute value of every root.

For \(P(x) = x^3 - 4x + 6\) the largest lower coefficient in absolute value is 6, and the leading coefficient is 1. The bound is \(1 + 6 = 7\), so every real root lies between -7 and 7.

The bound is often loose. A polynomial can satisfy it comfortably while its roots sit much closer to zero. Its value is that it turns an unbounded search into a finite interval.

## Candidates from the rational root theorem

For a polynomial with integer coefficients, every rational root written in lowest terms has a numerator that divides the constant term and a denominator that divides the leading coefficient. Each numerator and denominator pair produces one positive and one negative candidate.

If the constant term has 4 positive divisors and the leading coefficient has 2, the theorem produces \(2 \times 4 \times 2 = 16\) signed candidates. The list may contain repeats after reduction, but it is still finite.

Testing candidates is the last step. A candidate that survives is a root. A candidate that fails is discarded, and the bound and divisor counts explain why the list was complete in the first place.
