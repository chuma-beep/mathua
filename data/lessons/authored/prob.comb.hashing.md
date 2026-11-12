# Hashing collisions

Hashing maps keys to buckets. Collisions are unavoidable, and their number is an occupancy problem in disguise.

## Expected collisions

With \(m\) keys and \(n\) buckets, the expected number of colliding pairs is \(\binom{m}{2}/n\).

## Load factor

The ratio \(m/n\) controls the average bucket length and therefore the cost of a lookup.

## Universal hashing

A family is universal when every distinct pair collides with probability at most \(1/n\), which bounds the expected collision count.

To check your understanding, count the key pairs, divide keys by buckets and bound a pair collision.

Hashing is balls into bins with a cost model.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
