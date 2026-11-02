# The Binomial Distribution

**The binomial distribution** counts successes in a fixed number of independent trials with the same success probability.

## Bernoulli trials

A Bernoulli trial has two outcomes, usually called success and failure. Repeat the trial $n$ times, independently, with success probability $p$ on every trial. The random variable $X$ records the total number of successes. Its possible values are the integers from 0 to $n$.

The order of the successes does not matter for the count, but different orders are different outcomes. For example, with three trials the sequences $SFF$, $FSF$ and $FFS$ all produce one success, so the probability of exactly one success is three times the probability of any one sequence.

## The binomial probability formula

The number of sequences with exactly $k$ successes is the binomial coefficient
$$\binom{n}{k} = \frac{n!}{k!(n-k)!}.$$
Each such sequence has probability $p^k(1-p)^{n-k}$, so
$$P(X=k) = \binom{n}{k} p^k(1-p)^{n-k}.$$
For five trials with success probability 0.4, the chance of exactly two successes is
$$\binom{5}{2}(0.4)^2(0.6)^3 = 10 \cdot 0.16 \cdot 0.216 = 0.3456.$$
The binomial coefficients account for all orderings and the powers account for the successes and failures.

## Mean and variance

The mean number of successes is $np$. The variance is $np(1-p)$, and the standard deviation is its square root. These formulas come from adding $n$ independent Bernoulli variables, each with mean $p$ and variance $p(1-p)$.

The variance is largest when $p=0.5$ and shrinks toward zero as $p$ approaches 0 or 1. A binomial count is therefore most predictable when the trials are strongly biased and least predictable when success and failure are equally likely.
