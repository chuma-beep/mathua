# The Probability Integral Transform

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Uniform from any continuous CDF

If $X$ has a continuous distribution function $F$, then

$$U=F(X)$$

has a uniform distribution on the interval from zero to one. This is the probability integral transform. The reason is that $F$ maps the support of $X$ onto the unit interval and preserves cumulative probability.

The result fails at atoms. A discrete variable has jumps in $F$, and the transformed variable is not continuous, so the uniform conclusion needs the continuous case.

## Simulating from a distribution

The transform runs in reverse for simulation. Draw $U$ uniform on the unit interval and set

$$X=F^{-1}(U).$$

Then $X$ has distribution $F$. For an exponential variable with rate $\lambda$, the quantile function is $x=-\ln(1-u)/\lambda$, so a uniform draw becomes an exponential draw after a logarithm.

The inverse need not have a closed form. Numerical inversion, rejection sampling or specialised algorithms are used when it does not.

## Checking a model

The transform is also a diagnostic. If a fitted distribution is correct, the transformed observations should look uniform. A histogram or a uniform quantile plot that departs from flatness points to a mismatch in the tails or the centre.

The check is not a proof, especially with a small sample, but it is a fast way to see whether a fitted model is plausible before a formal test is chosen.

