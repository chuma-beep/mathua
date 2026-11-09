# Domain of logarithmic functions

A logarithm is defined only for positive arguments. The domain of a logarithmic function is therefore the set of inputs that make its argument positive, and finding that set is an inequality problem before it is a logarithm problem.

## A shifted argument

For \(f(x) = \log_b(x - k)\), the argument is \(x - k\). The condition is \(x - k > 0\), which rearranges to \(x > k\). The domain is the open interval \((k, \infty)\).

The endpoint is excluded because the argument would be zero at \(x = k\), and the logarithm of zero is not defined. The smallest integer in the domain is the first integer strictly greater than \(k\), which is \(k + 1\).

A graph makes the shift visible. The vertical asymptote of \(y = \log_b x\) at \(x = 0\) moves to \(x = k\), and the curve lives to the right of that line.

## A reversed argument

For \(f(x) = \log_b(k - x)\), the argument is \(k - x\). The condition is \(k - x > 0\), which rearranges to \(x < k\). The domain is \((-\infty, k)\), and the largest integer in it is \(k - 1\).

The direction of the inequality flips because the coefficient of \(x\) is negative. Dividing or multiplying an inequality by a negative number reverses its sign, and forgetting that step produces the complement of the correct domain.

Here the graph approaches the asymptote from the left. The curve exists only on the side where the argument is positive.

## A quadratic argument

When the argument is quadratic, the logarithm is defined wherever the quadratic is positive. For \(f(x) = \log_b(x^2 - k)\), the condition is \(x^2 > k\).

If \(k\) is positive, the solution is \(x < -\sqrt{k}\) or \(x > \sqrt{k}\). The interval between the two square roots is excluded because the quadratic is negative there.

The smallest positive integer in the domain is the first integer greater than \(\sqrt{k}\). Since the inequality is strict, that value is \(\lfloor \sqrt{k} \rfloor + 1\). The same reasoning applies to the negative branch, where the largest negative integer is one less than \(-\sqrt{k}\).
