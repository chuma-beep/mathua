# Gibbs sampling

Gibbs sampling updates one coordinate at a time from its full conditional distribution. Every proposal is accepted, so no acceptance ratio is needed.

## The sweep

One sweep replaces each coordinate in turn with a draw from its conditional given the others. The stationary distribution of the chain is the joint target.

## Relation to Metropolis

A full conditional draw is a Metropolis step with acceptance probability one. Gibbs is therefore a special case with a particular proposal.

## Limitations

Strong dependence between coordinates makes the chain move slowly. Reparametrisation or blocking restores mixing.

To check your understanding, draw each coordinate, set the acceptance to one and spot dependent coordinates.

Gibbs sampling is Metropolis where the proposal is the exact conditional.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
