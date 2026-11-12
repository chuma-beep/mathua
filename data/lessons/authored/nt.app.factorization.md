# Factorisation algorithms

Integer factorisation has no known efficient classical algorithm. Practical methods exploit smooth numbers, cycles and algebraic relations.

## Trial division

Dividing by every prime up to the square root always works. Its cost grows exponentially in the number of digits, so it is only a first step.

## Fermat's method

Writing n as a difference of squares reveals a factor. The method is fast when the two factors are close, and slow when they are far apart.

## Modern methods

The quadratic sieve and the number field sieve collect congruences of squares. Their complexity is subexponential, which is why key sizes keep growing.

To check your understanding, divide by small primes, take the ceiling root and collect congruences.

Factorisation algorithms are searches for structure, and each one is fast on a different shape of input.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
