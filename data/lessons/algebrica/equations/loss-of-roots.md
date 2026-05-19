> Content sourced from [Algebrica](https://algebrica.org/loss-of-roots/) — CC BY-NC 4.0

## Loss of Roots

Loss of roots occurs when an algebraic manipulation eliminates one or more [roots](<../roots-of-a-polynomial/>) of an [equation](<../equations/>), yielding a result that is only a partial solution set. The most frequent cause is dividing both sides of the equation by an expression containing the unknown, an operation that is valid only when that expression is guaranteed to be nonzero. Consider, for example, the equation:

\\[x(2x-5) = x\\]

A common mistake is to cancel the factor \\(x\\) appearing on both sides:

\\[\cancel{x}(2x-5) = \cancel{x}\\]

This reduces the equation to a linear equation in \\(x\\):

\\[2x - 5 = 1\\] \\[x = \frac{6}{2} = 3\\]

The cancellation is equivalent to dividing both sides by \\(x\\), an operation that presupposes \\(x \neq 0\\). The case \\(x = 0\\) is silently discarded, and one solution is lost.

More generally, any factor containing the unknown may be cancelled from both sides only after verifying that it cannot equal zero. When that condition cannot be established, the equation must be factored rather than simplified by division.

## Correct method

The correct approach consists in bringing the equation to standard form before any further manipulation. Rather than attempting to simplify factors, all terms are moved to one side so that the equation takes the form \\(ax^2 + bx + c = 0\\). Starting from the equation:

\\[x(2x-5) = x\\]

Expanding the left-hand side and subtracting \\(x\\) from both members yields the following.

\begin{align} 2x^2 - 5x - x &= 0 \\\\[6pt] 2x^2 - 6x &= 0 \\\\[6pt] 2x(x - 3) &= 0 \end{align}

The result is an [incomplete quadratic equation](<../incomplete-quadratic-equations/>) in which the constant term \\(c\\) is zero. Applying the zero-product property, the product \\(2x(x-3)\\) vanishes if and only if at least one of the two factors is zero.

Setting each factor equal to zero separately gives the two solutions of the equation.

\\[x_1 = 0 \qquad x_2 = 3\\]

##### The solution \\(x_1 = 0\\), discarded by the incorrect cancellation, is recovered as a direct consequence of the factored form. The equation admits two distinct real solutions, neither of which coincides with the value obtained by the flawed procedure.

## General principles

Several principles help prevent the unintentional loss of solutions when solving equations.

The most reliable strategy is to move all terms to one side and reduce the equation to zero, then factor the resulting expression. Once the equation is written as a product of factors equal to zero, the zero-product property guarantees that every solution corresponds to the vanishing of at least one factor, and none can be overlooked.

Division of both sides by an expression involving the unknown should be avoided unless that expression is known to be nonzero throughout the domain. When such a division appears necessary, the case in which the divisor equals zero must be analysed separately, as it may itself yield a valid solution.

When an equation involves even-index radicals, both the positive and the negative square root must be retained as candidates. Restricting attention to the principal root without verifying the other sign is a common source of incomplete solution sets.

Finally, substituting each candidate solution back into the original equation remains an essential step, in particular after squaring both sides or performing any transformation that is not strictly reversible. This verification detects both extraneous solutions introduced by the manipulation and any valid solutions that may have been dropped.

## Loss of roots by squaring

A second common source of root loss arises when solving [irrational equations](<../irrational-equations/>). Squaring both sides of an equation is a standard technique for eliminating radicals, but it is not an equivalence transformation: the squared equation may admit solutions that the original does not, and conversely, an incorrect application of the method can suppress valid ones. Consider the equation:

\\[\sqrt{2x + 1} = x - 1\\]

Squaring both sides yields the following:

\begin{align} 2x + 1 &= (x-1)^2 \\\\[6pt] 2x + 1 &= x^2 - 2x + 1 \\\\[6pt] x^2 - 4x &= 0 \\\\[6pt] x(x - 4) &= 0 \end{align}

The two candidates are \\(x_1 = 0\\) and \\(x_2 = 4\\). Substituting back into the original equation, \\(x_1 = 0\\) gives \\(\sqrt{1} = -1\\), which is false, so \\(x_1 = 0\\) is an extraneous solution introduced by squaring. The only valid solution is \\(x_2 = 4\\).

###### In this case no root is lost, but the example illustrates why verification is indispensable: squaring can both introduce spurious solutions and, if applied selectively or incorrectly, discard valid ones.

## Loss of roots in rational equations

A common mistake when solving [rational equations](<../rational-equations/>) is multiplying both sides by the denominator to eliminate fractions. If the denominator contains the variable, this step only works when the denominator is not zero. Just leaving out values that make the denominator zero does not mean you lose any solutions. But, if you don’t check whether a solution makes the denominator zero, you might include answers that don’t actually work.

On the other hand, if you simplify the equation by cancelling part of the denominator before solving, you might lose a solution. This happens because you remove any value that makes the cancelled part equal to zero from your list of possible answers. For example, consider the equation:

\\[\frac{x^2 - 4}{x - 2} = 0\\]

Cancelling the factor \\((x-2)\\) from both the numerator and denominator yields \\(x + 2 = 0\\), so \\(x = -2\\). The value \\(x = 2\\) is excluded from the domain since it causes the denominator to be zero. Therefore, the equation has exactly one solution, \\(x = -2\\). In this case, no root is lost, and the cancellation is valid because \\(x = 2\\) is not in the domain. However, the situation differs when the cancelled factor does not correspond to a domain restriction. For instance, consider:

\\[x(x - 3) = 2(x - 3)\\]

Cancelling \\((x - 3)\\) from both sides results in \\(x = 2\\). However, \\(x = 3\\) also satisfies the original equation, as both sides equal zero when \\(x = 3\\), but this solution is lost through cancellation. Thus, the equation has two solutions, \\(x_1 = 2\\) and \\(x_2 = 3\\), while the incorrect simplification retains only one.

Equations

Equations state equalities between quantities and are fundamental for solving mathematical problems.

9.8k

[Equations](https://algebrica.org/equations/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/equations.md?plain=1)

8.2k

[Polynomial Equations](https://algebrica.org/polynomial-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/equations.md?plain=1)

9.2k

[Linear Equations](https://algebrica.org/linear-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations.md?plain=1)

17.7k

[Quadratic Equations](https://algebrica.org/quadratic-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations.md?plain=1)

13.9k

[Quadratic Formula](https://algebrica.org/quadratic-formula/)

2 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/quadratic-formula.md?plain=1)

2.1k

[Factoring Quadratic Equations](https://algebrica.org/factoring-quadratic-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/factoring-quadratic-equations.md?plain=1)

3.9k

[Incomplete Quadratic Equations](https://algebrica.org/incomplete-quadratic-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/incomplete-quadratic-equations.md?plain=1)

1.7k

[The Geometric Interpretation of Quadratic Equations](https://algebrica.org/geometrical-meaning-quadratic-equations/)

1.4k

[Quadratic Equations with Complex Solutions](https://algebrica.org/quadratic-equations-with-complex-solutions/)

3.1k

[Binomial Equations](https://algebrica.org/binomial-equations/)

3.8k

[Trinomial Equations](https://algebrica.org/trinomial-equations/)

3.8k

[Rational Equations](https://algebrica.org/rational-equations/)

12.7k

[Irrational Equations](https://algebrica.org/irrational-equations/)

1.7k

[Absolute Value Equations](https://algebrica.org/absolute-value-equations/)

1.6k

[Exponential Equations](https://algebrica.org/exponential-equations/)

2.2k

[Logarithmic Equations](https://algebrica.org/logarithmic-equations/)

1.2k

[Homogeneous Trigonometric Equations](https://algebrica.org/homogeneous-trigonometric-equations/)

1.7k

[Trigonometric Equations: Part 1](https://algebrica.org/trigonometric-equations/)

651

[Equations with Parameters](https://algebrica.org/equations-with-parameters/)

1.3k

[Linear Equations with Parameters](https://algebrica.org/linear-equations-with-parameters/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations-with-parameters.md?plain=1)

1.6k

[Quadratic Equations with Parameters](https://algebrica.org/quadratic-equations-with-parameters/)
