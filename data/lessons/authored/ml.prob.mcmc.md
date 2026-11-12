# Markov chain Monte Carlo

Markov chain Monte Carlo draws correlated samples whose long-run distribution is the target posterior. A transition kernel leaves the target invariant, and the chain mixes towards it.

## The transition matrix

A chain on \(n\) states has an \(n \times n\) transition matrix. Each row sums to one because every state chooses a next state.

## Stationarity

A distribution is stationary when one transition leaves it unchanged. The target is stationary by construction, not by luck.

## Burn-in and autocorrelation

Early samples depend on the starting point, so they are discarded as burn-in. The autocorrelation time sets how many samples are needed for a given effective sample size.

To check your understanding, count the matrix entries, test invariance and choose a burn-in.

A Markov chain explores the posterior, and diagnostics decide whether it has mixed.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
