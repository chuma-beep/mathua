# Metropolis-Hastings

Metropolis-Hastings builds a Markov chain whose stationary distribution is the target. A proposal is accepted with a probability that corrects for the proposal's asymmetry.

## Acceptance

Accept with probability \(\min(1,\pi(y)q(x\mid y)/\pi(x)q(y\mid x))\). A symmetric proposal cancels the ratio of proposal densities.

## Detailed balance

The acceptance rule is chosen so that the chain satisfies detailed balance. That condition makes the target stationary without knowing its normalising constant.

## Tuning

A proposal that is too small accepts almost everything but moves slowly. A proposal that is too large is rejected constantly; an acceptance rate near one quarter is a common target in high dimension.

To check your understanding, form the target ratio, state detailed balance and diagnose slow mixing.

The acceptance ratio is what turns an arbitrary proposal into the target distribution.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
