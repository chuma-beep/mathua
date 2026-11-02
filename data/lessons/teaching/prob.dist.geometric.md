# The geometric distribution

The geometric distribution is part of the probability theory curriculum. This lesson develops the idea in three sections, from the defining picture through the main computation to the interpretation.

## The Geometric Model

A geometric random variable counts the number of independent Bernoulli trials needed to obtain the first success. If each trial succeeds with probability p, then the first success occurs on trial k exactly when the first k-1 trials fail and trial k succeeds. This gives the probability mass function P(X = k) = (1-p)^(k-1) p for k = 1, 2, 3, and so on. The parameter p must lie strictly between zero and one, because p = 0 never produces a success and p = 1 succeeds on the first trial.

The distribution is a model for waiting times in a repeated experiment. It describes the number of coin flips until the first head, the number of items inspected until the first defective one, or the number of attempts before a connection succeeds. Each trial must be independent of the others and must have the same success probability. When those assumptions fail, the geometric model is not appropriate, even if the outcomes look similar.

A short calculation sums the geometric series and confirms that the probabilities add to one. The tail probability is especially simple: P(X > k) = (1-p)^k, because the first k trials must all fail. That formula is useful whenever a question asks for at least or more than a stated number of trials, and it avoids summing a long finite series.

## Mean and Variance

The expected value of a geometric random variable is 1/p. This matches intuition: if a trial succeeds one time in four, the average wait is four trials. The variance is (1-p)/p^2, so the spread grows quickly as p becomes small. A rare success produces a long and highly variable waiting time, while a common success produces a short and predictable one.

These moments can be derived by differentiating the geometric series, but the memoryless property gives a quicker route. After any number of failures the remaining waiting time has the same distribution as the original waiting time. The past does not change the chances of future success, so the expected additional wait is still 1/p. This is the discrete analogue of the exponential distribution.

When a problem gives the mean rather than the success probability, invert the formula to get p = 1 / mean. That step is common in applied questions. Always check that the resulting probability lies between zero and one before using it, because a mean below one is impossible for a geometric count.

## The Memoryless Property

The defining feature of the geometric distribution is memorylessness: P(X > m + n | X > m) = P(X > n). Given that the first m trials failed, the chance of waiting at least n more trials is the same as the chance of waiting at least n trials from the start. The condition only removes a known run of failures, and independent trials make that information irrelevant to the future.

This property distinguishes the geometric model from most other counting distributions. A sum of geometric variables is not geometric, and a distribution with a changing success probability is not geometric either. The property also provides a quick test: if a proposed waiting-time model does not satisfy the conditional identity, it cannot be geometric.

The continuous counterpart is the exponential distribution, which is the only memoryless continuous distribution. Both models appear in reliability and queueing problems. Recognising the shared structure helps when translating a discrete count into a continuous time, and it explains why the same exponential formulas keep appearing in both settings.
