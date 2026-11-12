# Zero-knowledge proofs

A zero-knowledge proof convinces a verifier that a statement is true without revealing why. Repeated random challenges make the soundness error shrink geometrically.

## Commit and challenge

The prover commits to a random value, the verifier issues a challenge, and the prover answers in a way that only works if the secret is known.

## Soundness error

A cheating prover passes one challenge with probability one half. Repeating the protocol independently raises the failure cost exponentially.

## Repetition count

To reach a soundness error below one part in a chosen power of two, the protocol needs that many independent repetitions.

To check your understanding, commit to a random value, compute one-round error and choose a target exponent.

Zero knowledge trades certainty for probability, and the number of repetitions is the dial that buys the certainty.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
