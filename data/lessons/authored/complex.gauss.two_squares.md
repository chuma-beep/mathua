# The two-squares theorem

An integer is a sum of two squares exactly when every prime congruent to 3 mod 4 appears to an even power. The count of representations is four times the excess of divisors congruent to 1 mod 4 over those congruent to 3.

## The criterion

Factor \(n\). If some \(p\equiv3\pmod4\) has an odd exponent, \(n\) is not a sum of two squares.

## Counting representations

If \(d_{1}(n)\) counts divisors congruent to 1 and \(d_{3}(n)\) counts divisors congruent to 3, then \(r_{2}(n)=4(d_{1}(n)-d_{3}(n))\) counts signed representations in order.

## Small examples

\(25=5^{2}+0^{2}=3^{2}+4^{2}\) has twelve signed representations, while 3 has none because its exponent of 3 is odd.

To check your understanding, factor the integer, count divisors in each class and count the representations of a square.

The Gaussian norm turns the sum-of-two-squares question into a question about divisors.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
