# Continued fraction expansions

A continued fraction expands a number as a nested sequence of integer parts and reciprocals. Truncating the expansion gives convergents, which are the best rational approximations with small denominators.

## The expansion

The algorithm repeatedly takes the integer part and inverts the fractional part. Rational inputs terminate and irrational inputs continue forever.

## Convergents

The truncations satisfy a two-term recurrence. Even convergents lie below the value and odd convergents above, alternating toward it.

## Best approximation

No fraction with a smaller denominator is closer to the target. That optimality is what makes continued fractions useful for rational reconstruction.

To check your understanding, take the integer part, truncate the expansion and compare denominators.

Continued fractions encode a number by its integer-part sequence. The convergent recurrence is the same second-order linear recurrence that appears throughout the domain.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
