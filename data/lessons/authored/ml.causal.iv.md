# Instrumental variables

An instrument affects the treatment but reaches the outcome only through it. Under that exclusion restriction the instrument identifies a causal effect despite unmeasured confounding.

## The three conditions

Relevance, exclusion and independence must all hold. Relevance is testable, and the other two are assumptions.

## Two-stage least squares

The first stage predicts the treatment from the instrument. The second stage regresses the outcome on the predicted treatment.

## Weak instruments

A weak first stage inflates the variance of the estimate. The first-stage F statistic is the usual diagnostic.

To check your understanding, test relevance, predict the treatment and read the F statistic.

An instrument is a source of variation that is as good as random.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
