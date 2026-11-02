# Convolution of Random Variables

**Convolution** gives the distribution of a sum of independent random variables from the distributions of the summands.

## Sums of independent variables

If $X$ and $Y$ are independent, then means add and variances add:
$$E[X+Y] = E[X] + E[Y], \qquad \operatorname{Var}(X+Y) = \operatorname{Var}(X) + \operatorname{Var}(Y).$$
Named families often stay inside the same family. Independent Poisson variables add their rates. Independent normal variables add their means and their variances. A sum of independent Bernoulli variables with the same success probability is binomial.

## The convolution formula

For discrete variables the mass of a sum is a sum over all splits of the target value:
$$P(X+Y=k) = \sum_x P(X=x)P(Y=k-x).$$
For continuous variables the same idea becomes an integral:
$$f_{X+Y}(z) = \int_{-\infty}^{\infty} f_X(x) f_Y(z-x)\,dx.$$
The formula uses independence: the joint probability of the pair factors into the product of the two marginal probabilities.

## Discrete and continuous examples

Let $X$ and $Y$ be independent Bernoulli variables with success probability $1/2$. Their sum is 0 with probability $1/4$, 1 with probability $1/2$, and 2 with probability $1/4$. The middle value has two possible splits, so its probability is doubled.

For continuous variables, let $X$ and $Y$ be independent uniform variables on the interval from 0 to 1. Their sum has a triangular density: it rises from 0 to 1 and then falls back to 0. At the midpoint the density is 1. At a point one quarter of the way from an endpoint it is $1/4$. Two independent exponential variables with rate $\lambda$ add to a gamma variable with shape 2 and rate $\lambda$.
