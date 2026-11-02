# The Poisson Distribution

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Counting rare events

A Poisson random variable counts how many events occur in a fixed interval when the events arrive independently and at a constant average rate. The rate is written $\lambda$, and it is the expected count in one interval.

The distribution is appropriate for calls arriving at a desk, radioactive decays in a minute, or typing errors on a page. It is not appropriate when events cluster, because clustering makes one arrival more likely to bring another and breaks the independence the model assumes.

The interval matters. A rate of three events per hour becomes a rate of one event per twenty minutes, and the parameter must match the interval the question asks about.

## The Poisson formula

For a nonnegative integer $k$,

$$P(X=k)=e^{-\lambda}\frac{\lambda^k}{k!}.$$

The factor $e^{-\lambda}$ normalises the probabilities so that they sum to one. The factor $\lambda^k$ grows with the count, and $k!$ eventually grows faster, which is why the distribution has a peak and then a long right tail.

The mean and the variance are both $\lambda$. That equality is a quick check on a fitted model: data whose variance is much larger than its mean are overdispersed and need a different model.

## Working with the formula

To find $P(X=0)$, set $k=0$. Since $0!=1$ and $\lambda^0=1$, the result is $e^{-\lambda}$. To find $P(X\le 1)$, add $P(X=0)$ and $P(X=1)$.

For large $\lambda$ the distribution is approximately normal with mean $\lambda$ and variance $\lambda$. The approximation is better near the centre than in the tails, and a continuity correction improves it when a discrete count is compared with a continuous normal curve.

