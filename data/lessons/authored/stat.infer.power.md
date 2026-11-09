# Power of a hypothesis test

Power is the probability that a test rejects a null hypothesis when that hypothesis is false. A test with high power detects a real effect most of the time, while a low-power test misses it and reports no evidence. Power is therefore the complement of the Type II error probability, and it is the quantity to design into a study before the data are collected.

## Power and the two errors

A hypothesis test can fail in two ways. A Type I error rejects a true null hypothesis, and its probability is the significance level \(\alpha\). A Type II error fails to reject a false null hypothesis, and its probability is \(\beta\). Power is

\[
\text{power} = 1 - \beta.
\]

The two error probabilities trade against each other. Lowering \(\alpha\) makes a test more conservative, which usually raises \(\beta\) and lowers power. Raising the sample size is the way to improve both at once, because a larger sample narrows the sampling distribution and separates the null and alternative hypotheses more clearly.

Power is not the probability that a significant result is true. That is a different quantity and it depends on the prior probability that the null is false. Power is a property of the test and the alternative, measured before any data arrive.

## Recovering beta

The relation is symmetric in the two quantities. Given power, the Type II error probability is

\[
\beta = 1 - \text{power}.
\]

A test with power 0.80 has a 20 percent chance of missing a true effect of the size it was designed to detect. That convention is common in sample-size planning, but it is a convention rather than a law. A study with power 0.80 still fails to detect a real effect one time in five.

The complement is only meaningful when the alternative is specified. Power depends on the size of the effect and on the variability of the data, so a single number without its alternative is incomplete. Reporting power always means reporting the effect size and the sample size that produced it.

## Power from a critical region

For a test with a known critical region, power is computed under the alternative distribution. If the test rejects when a sample mean exceeds a critical value \(c\), and the true mean is \(\mu_1\) with standard error \(SE\), then

\[
\text{power} = P\left(Z > \frac{c - \mu_1}{SE}\right).
\]

The standardized distance measures how far the critical value sits from the true mean in standard error units. A negative distance means the critical value lies below the true mean, so rejection is likely and power is high. A positive distance means the critical value lies above the true mean, so rejection is unlikely.

A small table of standard normal tail probabilities is enough for the cases used here: \(P(Z > -2) = 0.9772\), \(P(Z > -1) = 0.8413\) and \(P(Z > 0) = 0.5\). Those values show the pattern without requiring a full normal table, and they make the dependence on the standardized distance visible.
