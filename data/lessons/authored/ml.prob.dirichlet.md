# Dirichlet priors

A Dirichlet distribution is a prior over probability vectors. Each component has a Gamma marginal, and the total concentration controls how far the vector strays from its mean.

## The parameter vector

A Dirichlet on \(k\) categories has \(k\) positive parameters. The sum of the parameters is the concentration.

## The mean

Component \(i\) has mean \(\alpha_i/\alpha_0\), where \(\alpha_0\) is the total concentration. The means sum to one by construction.

## Conjugacy

Multinomial observations add to the corresponding parameters. The posterior stays Dirichlet, which is why the prior is convenient for topic models.

To check your understanding, count the parameters, compute the mean and add the counts.

The Dirichlet is the natural prior for weights that must sum to one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
