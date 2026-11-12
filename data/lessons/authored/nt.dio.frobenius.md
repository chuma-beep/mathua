# The Frobenius coin problem

The Frobenius coin problem asks for the largest amount that cannot be formed from given coin denominations. For two coprime coins there is a closed formula.

## Two coins

With coprime denominations a and b, the largest unrepresentable amount is a b minus a minus b. Every larger amount is representable.

## The count

The number of unrepresentable amounts is a minus one times b minus one divided by two. The formula counts the holes in the numerical semigroup.

## More coins

With three or more denominations no closed formula is known in general. The problem becomes a computation in a numerical semigroup.

To check your understanding, multiply the coins, subtract one from each and compare with two coins.

Two coprime coins produce a finite gap set with an exact size, and the largest gap is the most useful single number in it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
