# Maximum mean discrepancy

Maximum mean discrepancy measures the distance between two distributions by comparing their mean embeddings in a reproducing kernel Hilbert space. It is zero exactly when the embeddings agree.

## The embedding

A distribution maps to the mean of the kernel sections over its support. Two distributions with the same embedding are indistinguishable to the kernel.

## The estimator

The biased estimator averages kernel values within and across the two samples. With \(m\) and \(n\) points it uses \(m^2+n^2+2mn\) kernel evaluations.

## The test

A permutation test calibrates the statistic under the null. The kernel choice sets which differences the test can detect.

To check your understanding, define the mean map, average within each sample and permute the labels.

MMD turns a two-sample question into a distance between mean embeddings.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
