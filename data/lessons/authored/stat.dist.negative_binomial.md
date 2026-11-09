# Negative binomial distribution

The negative binomial distribution counts the trials needed to reach a fixed number of successes. It is the natural extension of the geometric distribution, which waits for one success, and it appears whenever a process is repeated until a quota is met. The two parameters are the success probability \(p\) and the required number of successes \(r\).

## Trials until r successes

Let \(X\) be the total number of trials, including the final success, needed to obtain \(r\) successes. The last trial must be a success, and the remaining \(r-1\) successes are placed among the first \(k-1\) trials. That counting gives

\[
P(X = k) = \binom{k-1}{r-1} p^r (1-p)^{k-r}, \qquad k = r, r+1, r+2, \dots
\]

The binomial coefficient counts the arrangements of the earlier successes. The factor \(p^r\) accounts for all \(r\) successes, and \((1-p)^{k-r}\) accounts for the failures. Every sequence with the same counts has the same probability, which is why counting sequences is enough.

The support starts at \(r\). Fewer than \(r\) trials cannot contain \(r\) successes, so the probability is zero there. The distribution is skewed to the right, and the skew is strongest when \(p\) is small.

## The probability of a given trial count

To evaluate the formula, first identify the last trial as a success. Then choose the positions of the other \(r-1\) successes among the first \(k-1\) positions. For example, with \(r = 2\), \(k = 4\) and \(p = 1/2\), there are \(\binom{3}{1} = 3\) arrangements of the first success among the first three trials. Each arrangement has probability \((1/2)^4 = 1/16\), so the total is \(3/16\).

Keeping the answer as an exact fraction avoids the rounding that makes two nearby probabilities look equal. The denominator is a power of the success-probability denominator, and the numerator is a product of a binomial coefficient and a power of the failure probability. Reducing the fraction is the final step, not a change of value.

The formula also works when \(r = 1\). Then the binomial coefficient is one, the last trial is the only success, and the result is the geometric distribution on the number of trials.

## Mean and failures

The expected number of trials is

\[
E[X] = \frac{r}{p}.
\]

If each trial succeeds with probability \(1/5\) and three successes are required, the expected count is \(3 \times 5 = 15\) trials. The expected number of failures before the \(r\)-th success is the expected trial count minus the \(r\) successes themselves:

\[
E[X - r] = \frac{r(1-p)}{p} = r\left(\frac{1}{p} - 1\right).
\]

The failure form is often more natural in applications. A quality inspector who needs three acceptable parts expects \(r(1-p)/p\) rejects before the quota is filled, and that number can be much larger than \(r\) when \(p\) is small. The variance is \(r(1-p)/p^2\), which grows quickly as \(p\) falls, so the waiting time becomes less predictable as success becomes rarer.
