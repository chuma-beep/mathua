> Content sourced from [Algebrica](https://algebrica.org/rational-inequalities/) — CC BY-NC 4.0

## Introduction

A rational inequality is an inequality that involves at least one rational expression, that is, a ratio in which both the numerator and the denominator are [polynomials](<../polynomials>). The natural domain of such an expression is the set of all real numbers for which the denominator does not vanish; any value of \\(x\\) that makes the denominator equal to zero is excluded from the domain and cannot belong to the solution set. Every rational inequality can be reduced to one of the following canonical forms.

\\[\frac{P(x)}{Q(x)} > 0 \qquad \frac{P(x)}{Q(x)} \geq 0\\] \\[\frac{P(x)}{Q(x)} < 0 \qquad \frac{P(x)}{Q(x)} \leq 0\\]

In each case, \\(P(x)\\) and \\(Q(x)\\) are polynomials with real coefficients, and the expression on the left-hand side is defined only where \\(Q(x) \neq 0\\). These inequalities take their name from the notion of a rational function, in the same way that [rational equations](<../rational-equations>) do. All four forms are handled by the same method, as illustrated in the examples below.


Solving a rational inequality means finding all values of \\(x\\) in the natural domain of the expression for which the inequality holds. The standard approach proceeds through the following steps.

  * The first step is to rewrite the inequality so that all terms appear on the left-hand side and the right-hand side is zero, thereby reducing the problem to one of the canonical forms above. This allows one to study the sign of a single rational expression rather than comparing two separate expressions.
  * The second step is to determine the natural [domain](<../determining-the-domain-of-a-function/>) of the rational expression by identifying all values of \\(x\\) for which the denominator vanishes. These values must be excluded from the solution set regardless of the inequality symbol.
  * The third step is to find the zeros of the numerator and the zeros of the denominator separately. These critical points, taken together, partition the real line into open intervals on which the rational expression maintains a constant sign, since the sign of a ratio can change only where the numerator or the denominator equals zero.
  * The fourth step is to construct a [sign chart](<../sign-analysis-in-inequalities/>), recording the sign of the numerator and the sign of the denominator in each interval, and then determining the sign of the ratio by the standard rule: the ratio is positive when numerator and denominator share the same sign, and negative when they have opposite signs.
  * The fifth step is to collect the intervals where the sign of the expression is consistent with the original inequality, taking care to include the zeros of the numerator if and only if the inequality is non-strict, and to exclude in every case the zeros of the denominator.


##### The key observation underlying this method is that a rational expression \\(P(x)/Q(x)\\) can change sign only at a zero of \\(P(x)\\) or at a zero of \\(Q(x)\\). Between any two consecutive critical points, the expression is either strictly positive or strictly negative throughout, which makes it sufficient to test a single representative value in each interval.

## Multiplicity and sign changes

The multiplicity of each zero is crucial in determining how the sign of a rational expression changes across its critical points. A zero of odd multiplicity, whether located in the numerator or denominator, causes the expression to change sign as \\(x\\) crosses that point. The expression is positive on one side and negative on the other.

In contrast, a zero of even multiplicity does not produce a sign change. The expression maintains the same sign on both sides of that point because the corresponding factor contributes an even power, which does not change sign in a neighbourhood of the zero.

This behaviour must be explicitly considered when constructing the sign chart. Treating all critical points as sign-change points, regardless of multiplicity, is a common error that yields incorrect solution sets when repeated factors appear in the numerator or denominator.


Consider the following inequality, which demonstrates the effect of even multiplicity on the sign chart.

\\[\frac{(x-1)^2}{x-3} \geq 0\\]

The denominator vanishes at \\(x = 3\\), which must therefore be excluded from the solution set regardless of the inequality symbol, since the expression is undefined at that point. The natural domain is thus the set of all real numbers with \\(x \neq 3\\). The numerator vanishes at \\(x = 1\\), which is a zero of multiplicity two. The critical points are \\(x = 1\\) and \\(x = 3\\), and they partition the real line into three intervals.

| | \\[1\\]| \\[3\\]  
---|---|---|---  
\\[(x-1)^2 \geq 0\\]| \\(\boldsymbol{+}\\)| \\(\boldsymbol{+}\\)| \\(\boldsymbol{+}\\)  
\\[x - 3 \geq 0\\]| \\(\boldsymbol{-}\\)| \\(\boldsymbol{-}\\)| \\(\boldsymbol{+}\\)  
\\[\frac{(x-1)^2}{x-3} \geq 0\\]| \\(\boldsymbol{-}\\)| \\(\boldsymbol{-}\\)| \\(\boldsymbol{+}\\)  
  
Since \\((x-1)^2\\) is non-negative for all real \\(x\\) and equals zero only at \\(x = 1\\), it does not change sign as \\(x\\) crosses \\(x = 1\\). The sign of the entire expression is therefore determined solely by the sign of the denominator \\(x - 3\\).

The expression is non-negative when \\(x > 3\\), and it equals zero at \\(x = 1\\). The point \\(x = 3\\) is excluded from the solution set because it does not belong to the natural domain of the expression. The solution set is as follows.

\\[\\{1\\} \cup (3, +\infty)\\]

The point \\(x = 1\\) appears as an isolated point in the solution set, a consequence of the even multiplicity of the corresponding zero.

## Structure of the solution set

The solution set of a rational inequality is always a union of intervals of the real line, possibly supplemented by isolated points. Each interval in the solution set is bounded by critical points, and its endpoints are either zeros of the numerator, which may be included if the inequality is non-strict, or zeros of the denominator, which are always excluded. When a zero of even multiplicity appears in the numerator, it contributes an isolated point to the solution set rather than an interval, as discussed in the preceding section on multiplicity and sign changes.

To make this concrete, consider a rational expression with critical points at \\(x = 1\\), \\(x = 2\\), and \\(x = 4\\), where \\(x = 2\\) is a zero of the denominator. If the sign chart shows that the expression is non-negative on \\([1, 2)\\) and on \\((2, 4]\\), the solution set is the following.

\\[[1, 2) \;\cup\; (2, 4]\\]

The point \\(x = 2\\) is absent from both intervals because it does not belong to the natural domain of the expression, even though the adjacent intervals extend up to it from both sides.

## Example 1

Determine the values of \\(x\\) that satisfy the following inequality.

\\[\frac{x - 1}{2 - x} < 0 \\]

The inequality is already in canonical form. The denominator \\(2 - x\\) vanishes at \\(x = 2\\), so this value must be excluded from the solution set. The natural domain of the expression is the following.

\\[D = \lbrace x \in \mathbb{R} : x \neq 2 \rbrace \\]

To determine the intervals where the inequality is satisfied, we use the fact that a rational expression is negative exactly when the numerator and the denominator have opposite signs. The numerator \\(x - 1\\) is negative when \\(x < 1\\), while the denominator \\(2 - x\\) is negative when \\(x > 2\\). By constructing the [sign chart](<../sign-analysis-in-inequalities/>), we obtain:

| | \\[1\\]| \\[2\\]  
---|---|---|---  
\\[x - 1 \lt 0 \\]| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\[2 - x \lt 0 \\]| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)  
\\[\frac{x-1}{2-x} \lt 0 \\]| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)  
  
From the sign chart, the expression is negative on the intervals \\(x < 1\\) and \\(x > 2\\). Since the inequality is strict, the endpoints are excluded, and the solution set is the following:

\\[(-\infty,\, 1) \;\cup\; (2,\, +\infty) \\]

## Example 2

Determine the values of \\(x\\) that satisfy the following inequality.

\\[\frac{x - 1}{x^2 - 5x + 6} \geq 0 \\]

The denominator is a quadratic polynomial whose zeros are found by solving the associated [quadratic equation](<../quadratic-equations/>). The polynomial [factors](<../factoring-quadratic-equations/>) as follows.

\\[(x-3)(x-2) = 0\\]

This gives \\(x = 2\\) and \\(x = 3\\), which are the values that make the denominator zero and must therefore be excluded from the solution set. The natural domain is the following.

\\[D = \lbrace \, x \in \mathbb{R} : x \neq 2, \, x \neq 3 \,\rbrace \\]

The numerator vanishes at \\(x = 1\\), which together with \\(x = 2\\) and \\(x = 3\\) gives three critical points partitioning the real line into four open intervals. We now represent their signs on a [sign chart](<../sign-analysis-in-inequalities/>):

| | \\[1 \\]| \\[2 \\]| \\[3 \\]  
---|---|---|---|---  
\\[x-1 \geq 0 \\]| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\[x^2 - 5x + 6 \geq 0 \\]| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\[\frac{x - 1}{x^2 - 5x + 6} \geq 0\\]| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
  
From the sign chart, the expression is non-negative on the interval \\(1 \leq x < 2\\) and on the interval \\(x > 3\\). The endpoints \\(x = 2\\) and \\(x = 3\\) are excluded because they do not belong to the natural domain, while \\(x = 1\\) is included because the inequality is non-strict and the numerator vanishes there.

The solution set is the following:

\\[[1, 2) \;\cup\; (3, +\infty) \\]

## Reduction to canonical form

The method described above assumes that the inequality is already in canonical form, with a single rational expression on the left-hand side and zero on the right. In practice, inequalities do not always appear in this form. A common case is one in which a rational expression appears on both sides, such as the following.

\\[\frac{2x + 1}{x - 3} \geq \frac{1}{x + 1}\\]

The correct approach is to move all terms to the left-hand side and combine them into a single rational expression. Subtracting the right-hand side from both sides gives the following.

\\[\frac{2x + 1}{x - 3} - \frac{1}{x + 1} \geq 0\\]

The two fractions are then combined over a common denominator.

\\[\frac{(2x + 1)(x + 1) - (x - 3)}{(x - 3)(x + 1)} \geq 0\\]

Expanding the numerator yields the following.

\\[\frac{2x^2 + 3x + 1 - x + 3}{(x - 3)(x + 1)} \geq 0\\]

which simplifies to the following.

\\[\frac{2x^2 + 2x + 4}{(x - 3)(x + 1)} \geq 0\\]

The numerator \\(2x^2 + 2x + 4\\) can be written as \\(2(x^2 + x + 2)\\). The discriminant of \\(x^2 + x + 2\\) is \\(\Delta = 1 - 8 = -7 < 0\\), so the quadratic has [no real roots](<../quadratic-equations-with-complex-solutions/>) and is strictly positive for all real \\(x\\). The sign of the entire expression is therefore determined solely by the sign of the denominator \\((x-3)(x+1)\\), which is positive when \\(x < -1\\) or \\(x > 3\\). Since the numerator is never zero, neither endpoint can be included, and the solution set is the following.

\\[(-\infty, -1) \;\cup\; (3, +\infty) \\]

A common error in this type of problem is to multiply both sides of the original inequality by \\(x - 3\\) or \\(x + 1\\) without accounting for the sign of those factors. Since the sign of a linear expression depends on \\(x\\), such a multiplication would require a case analysis and is more error-prone than the reduction to canonical form illustrated above.

## Selected references

  * **University of Rhode Island, L. Barnes**. [Solving Polynomial and Rational Inequalities](https://www.math.uri.edu/~lbarnes/classes/2019-2020/spring/Inequality%20Lecture%20Notes.pdf)
  * **Florida State University, P. Kirby**. [Solving Polynomial and Rational Inequalities](https://www.math.fsu.edu/~pkirby/mac1140/5_4Inequalities.pdf)
  * **Stitz-Zeager Open Mathematics, C. Stitz, J. Zeager**. [Rational Inequalities and Applications](https://stitz-zeager.com/szprecalc08262010.pdf)
  * **Lamar University, P. Dawkins**. [Rational Inequalities](https://tutorial.math.lamar.edu/classes/alg/SolveRationalInequalities.aspx)
