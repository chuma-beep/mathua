# Limit of a sequence

## The epsilon definition

A sequence \(a_n\) tends to \(L\) when every positive \(\epsilon\) admits an index \(N\) such that \(|a_n-L|<\epsilon\) for all \(n>N\). Take \(a_n = 1/n\) and \(L = 0\). The inequality \(1/n < \epsilon\) holds once \(n > 1/\epsilon\), so \(\epsilon = 0.01\) needs \(N = 100\) while \(\epsilon = 0.001\) needs \(N = 1000\). The index depends on the tolerance.

## Limit laws

Sums, products and quotients of convergent sequences converge to the matching combination, provided the denominator limit is not zero. If \(a_n \to 3\) and \(b_n \to 5\), then \(2a_n - b_n \to 6 - 5 = 1\), the product \(a_n b_n \to 15\) and the quotient \(a_n/b_n \to 3/5\). Each step quotes a limit already known, so no new epsilon argument is needed.

## A rational limit

Divide numerator and denominator by the highest power of \(n\) present. For \(a_n = (3n^2+2n)/(n^2-7)\) division by \(n^2\) gives \((3+2/n)/(1-7/n^2)\), which tends to 3. At \(n = 100\) the value is \(30200/9993 \approx 3.022\), just above 3 as the positive correction predicts. Degrees decide the outcome: a larger numerator diverges and a larger denominator gives zero.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
