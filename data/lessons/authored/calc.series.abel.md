# Abel's test

Abel's test multiplies a convergent series by a bounded monotone factor. The product still converges, which lets a difficult factor be separated from the series that carries the cancellation.

## The Abel test

Let \(\sum a_n\) converge and let \(b_n\) be bounded and monotone. Then
\[
\sum a_n b_n
\]
converges. The factor \(b_n\) may have a nonzero limit, so the test is not limited to factors that tend to zero.

The proof uses summation by parts. The partial sums of \(\sum a_n\) are bounded because the series converges, and the monotone factor has bounded variation. The boundary terms and the remaining sum both stay controlled.

## Bounded monotone factors

A factor such as \(1/n\) or \(e^{-n}\) decreases to zero. A factor such as \(1+1/n\) decreases to 1, while \(3-1/n\) increases to 3. All of them are bounded and monotone, so all of them satisfy the hypothesis.

What matters is not the limit of \(b_n\) but the two properties. An unbounded factor such as \(n\) fails the hypothesis even when it is monotone. An oscillating factor such as \((-1)^n\) is bounded but not monotone.

## Checking the hypothesis

To apply the test, check the two conditions separately. First decide whether the factor stays inside a finite interval. Then decide whether it moves in one direction. A factor that passes both tests can be attached to any convergent series.

A factor that fails either condition needs a different argument. Abel's test is a sufficient condition, so a failure does not prove divergence. It only means this test cannot decide the question.
