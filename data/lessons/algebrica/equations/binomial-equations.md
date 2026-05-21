> Content sourced from [Algebrica](https://algebrica.org/binomial-equations/) — CC BY-NC 4.0

## What are binomial equations

Binomial equations are a specific type of algebraic [equations](<../equations>) that involve only two terms, typically expressed in the form:

\\[ax^n + b = 0\\]

  * \\(a\\) and \\(b\\) are constants.
  * \\(x\\) is a variable raised to a power of \\(n\\), where \\(n\\) is a positive [integer](<../integers/>).


##### Different strategies exist to solve these equations based on the value of \\( n \\). When \\( n \\) is an odd integer, the equation typically has a unique real solution. If \\( n \\) is even, the presence of real solutions depends on the sign of the right-hand side: for example, the equation has no real solutions if it requires taking an even root of a negative number. In more advanced contexts, solutions may also include complex numbers, especially when working with negative or non-real results.

## Equations with a degree less than two

If the value of \\(n\\) equals \\(1\\), the given equation reduces to a simple [linear equation](<../linear-equations>) in the form of \\(ax + b = 0\\) that can be solved by isolating the variable \\(x\\) and finding its corresponding value.

\\[ax+b= 0 \quad \rightarrow \quad x = \frac{-b}{a} \\]

If the value of \\(n\\) equals \\(2\\), the given equation reduces to a [quadratic equation](<../quadratic-equations>) in the form of \\(ax^2 + b = 0\\) that can be solved using the [quadratic formula](<../quadratic-formula>) or in a more straightforward calculating the square root of the term \\(\large{\frac{-b}{a}}\\):

\\[x = \pm \sqrt{\frac{-b}{a}} \\]

##### This formula yields real solutions only when \\( -\dfrac{b}{a} \geq 0 \\). Otherwise, the equation has [complex solutions](<../quadratic-equations-with-complex-solutions/>).

## Equations with a degree greater the than two

If \\(n\\) is greater than 2, we are dealing with a relatively simple case of an equation with a degree higher than two. Generally, such [equations](<../equations>) can be solved by calculating the nth root of the value \\(\large{\frac{-b}{a}}\\) while considering two two distinct cases, depending on whether \\(n\\) os even or \\(n\\) odd:

When \\(n\\) is even and \\(\large{\frac{-b}{a}}\\) is positive, we have two distinct solutions, or a single solution if \\(\large{\frac{-b}{a}}\\) equals \\(0\\) in the form:

\\[x = \pm \sqrt[\Large{n}]{\frac{-b}{a}} \\]

When \\(n\\) is even and \\(\large{\frac{-b}{a}}\\) is negative, the equation has no solutions in the [real number](<../types-of-numbers>) field.


When \\(n\\) is odd, there is always a single solution in the form:

\\[x = \sqrt[\Large{n}]{\frac{-b}{a}} \\]

In fact, the nth [root](<../roots>) of a negative number is feasible only when the root’s index is odd. This arithmetic operation is mathematically valid and can be executed using the same methodology as real positive numbers. The result of extracting the nth root of a negative number is also negative, except when the index is an even number. In such cases, there is no real solution.

## Example 1

Solve the binomial equation \\(x^3-27 = 0\\).

We have:

\begin{align*} x^3 &= 27 \\\\[0.6em] x &= \sqrt[\Large{3}]{27} x & = 3 \end{align*}

##### In this case, \\(n\\) is odd, and the number under the root is positive, so the equation admits a single real solution.

The solution to the equation is: \\[x = 3 \\]

## Example 2

Solve the binomial equation \\(x^3 + 8 = 0\\).

We have:

\begin{align*} x^3 &= -8 \\\\[0.6em] x &= \sqrt[\large{3}]{-8} \end{align*}

In this case, \\(n\\) is odd, and the number under the root is negative, so the equation admits a single real solution.

The solution to the equation is: \\[x= -2\\]

## Example 3

Solve the binomial equation \\(x^4 +5 = 0\\).

We have:

\begin{align*} x^4 &= -5 \\\\[0.6em] x &= \sqrt[\large{4}]{-5} \end{align*}

In this case, \\(n\\) is even, and the number under the root is negative, so the equation does not admit real solutions.

The solution to the equation is: \\[\nexists x \in \mathbb{R}\\]

##### Indeed, the examples presented are merely generic illustrations of easily solvable equations. In practice, it is possible to encounter instances where complex [polynomial equations](<../polynomial-equations/>) may be reduced to binomial equations, which, as demonstrated, can be resolved using straightforward methods
