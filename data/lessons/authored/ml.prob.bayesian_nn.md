# Bayesian neural networks

A Bayesian neural network puts a distribution over weights instead of a point estimate. Predictions average over that distribution, so the model reports uncertainty as well as a value.

## The posterior over weights

The posterior is proportional to the likelihood times the prior. It is intractable for any network with a hidden layer.

## Predictive averaging

The predictive distribution averages the network output over posterior weight samples. A Monte Carlo estimate uses \(m\) samples.

## Uncertainty

Aleatoric uncertainty comes from the noise model. Epistemic uncertainty comes from the spread of the posterior and shrinks with more data.

To check your understanding, write Bayes' rule, average over samples and separate the two kinds.

Bayesian neural networks turn weight fitting into posterior inference.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
