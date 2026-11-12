# Martingale representation

Every square-integrable martingale adapted to a Brownian filtration is a stochastic integral against that Brownian motion.

## Statement

For such a martingale \(M_t\), there is an adapted \(H\) with \(M_t=M_0+\int_0^t H_s dB_s\).

## Consequence

The Brownian motion is the only source of randomness, so hedging strategies can replicate any square-integrable claim.

## Uniqueness

The integrand is unique up to equality almost everywhere in time and probability.

To check your understanding, start from the initial value, replicate a claim and compare two integrands.

Brownian filtrations have a single generator.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
