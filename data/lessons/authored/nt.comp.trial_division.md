# Trial division

Trial division tests divisibility by every prime up to the square root. It is the baseline against which every faster factorisation method is measured.

## The bound

A composite number has a factor at most its square root. Testing beyond that bound cannot find a new factor, so the search stops there.

## The count

The number of trial divisions is the number of primes up to the bound, not the bound itself. Sieving the primes first is what makes the count small.

## When it wins

Trial division is fast for numbers with a small factor and hopeless for semiprimes with two large factors. That split is why hybrid algorithms start with it.

To check your understanding, take the square root, sieve the primes and test a small factor.

Trial division is complete but slow, and its square-root bound is the first fact every other factorisation method improves on.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
