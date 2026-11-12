# Universal hashing

Universal hashing chooses a hash function at random from a family with a proven collision bound. Randomness replaces the adversary's knowledge of the function.

## The family

For a prime p and a table size m, the functions a x plus b modulo p, then modulo m, form a universal family when a is nonzero.

## Collision probability

For any two distinct keys, the probability that a random member of the family collides is at most one over m. That bound is what makes the expected chain length constant.

## Choosing parameters

The prime must exceed every key, and the pair of parameters ranges over p times p minus one choices. Sampling uniformly from that set is the whole implementation.

To check your understanding, choose a prime, fix two keys and bound the keys.

Universal hashing makes worst-case behaviour improbable rather than impossible, and the collision bound is its only promise.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
