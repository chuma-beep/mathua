# Optional Stopping

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Stopping times

A stopping time is a random time whose value can be determined from the information available up to that time. The first moment a random walk reaches a level is a stopping time, because at each step you know whether the level has been reached. A time chosen by looking into the future is not.

Optional stopping asks when the identity $E[X_T]=E[X_0]$ holds for a martingale $X$ and a stopping time $T$. The answer is not automatic: the stopped process must be well behaved, and the usual sufficient conditions bound the stopped variables or the stopping time.

## The gambler's ruin calculation

A symmetric random walk on the integers starts at $k$ and stops when it reaches zero or $N$. The probability of reaching $N$ before zero is

$$P=\frac{k}{N}.$$

The expected number of steps until absorption is $k(N-k)$. Both results follow from martingales: the position is a martingale, and the squared position minus the step count is also a martingale.

These formulas are a useful check on a simulation. They also show why a fair game cannot turn a profit in expectation when the stopping rule is honest.

## When the identity fails

Optional stopping can fail when the stopping time is unbounded or the stopped variable is not integrable. A martingale that doubles its stake after every loss can be stopped at the first win with a guaranteed profit, but the expected stake at the stopping time is infinite, so the hypotheses do not hold.

The practical lesson is to verify the conditions rather than apply the identity mechanically. Bounded stopping times and bounded increments are the simplest safe cases.

