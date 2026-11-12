# Bayesian optimisation

Bayesian optimisation fits a surrogate to an expensive objective and chooses the next point by maximising an acquisition function. Exploration and exploitation are balanced by that choice rather than by a fixed schedule.

## The surrogate

A Gaussian process is the usual surrogate. It supplies a posterior mean and variance at every candidate point.

## Acquisition

Expected improvement integrates the improvement over the posterior. Upper confidence bounds combine mean and variance with a single weight.

## The cost of a round

Each round costs one objective evaluation plus the surrogate update. The update is cheap relative to a costly experiment.

To check your understanding, fit the surrogate, maximise expected improvement and count evaluations.

Bayesian optimisation spends evaluations where the surrogate says the payoff or the information is largest.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
