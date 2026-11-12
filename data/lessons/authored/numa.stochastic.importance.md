# Importance sampling for integrals

Importance sampling draws from a proposal distribution and reweights each sample. It concentrates work where the integrand matters, which helps when rare events dominate.

## The weight

The estimator averages \(f(X)p(X)/q(X)\) with \(X\sim q\). The ratio \(p/q\) is the likelihood ratio or weight.

## Variance

The variance is small when q resembles the integrand times p. A proposal with thinner tails than the integrand gives infinite variance.

## Diagnostics

The effective sample size is \((\sum w_i)^{2}/\sum w_i^{2}\). It falls below n when a few weights dominate.

To check your understanding, form the likelihood ratio, match proposal to integrand and compute the effective size.

Importance sampling works when the proposal is shaped like the thing being measured.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
