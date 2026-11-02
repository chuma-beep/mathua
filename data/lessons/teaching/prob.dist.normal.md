# The Normal Distribution

This lesson develops the definition, the main calculation and the interpretation of the idea.

## The bell curve

A normal random variable with mean $\mu$ and variance $\sigma^2$ has density

$$f(x)=\frac{1}{\sigma\sqrt{2\pi}}e^{-(x-\mu)^2/(2\sigma^2)}.$$

The curve is symmetric about $\mu$. The standard deviation $\sigma$ sets the width: a larger $\sigma$ makes the curve lower and wider, while the total area stays one.

The normal family is closed under adding independent members and under affine changes. Those two facts explain why averages of many small independent effects are often approximately normal.

## Standardising

Every normal question can be reduced to the standard normal $Z$ with mean zero and variance one. The standardisation is

$$z=\frac{x-\mu}{\sigma}.$$

A value below the mean gives a negative $z$, and a value above the mean gives a positive $z$. The standard normal table or a numerical routine then supplies probabilities such as $P(Z\le z)$.

The transformation is reversible: $x=\mu+\sigma z$. Reversing it is how a probability is turned back into a measurement on the original scale.

## The empirical rule

For any normal distribution, about 68 percent of the probability lies within one standard deviation of the mean, about 95 percent within two, and about 99.7 percent within three. These figures are approximations, not exact identities, but they are accurate enough for a quick check.

A value more than three standard deviations from the mean is unusual under the model. That does not by itself prove an error or an outlier, because unusual values do occur, but it is a reason to inspect the observation and the model.

