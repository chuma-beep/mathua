> Content sourced from [Algebrica](https://algebrica.org/absolute-value-equations/) — CC BY-NC 4.0

## Introduction

Absolute value equations are a particular class of [equations](<../equations>) in which the variable \\( x \\) appears within an [absolute value](<../absolute-value>) expression. The absolute value represents how far a [number](<../types-of-numbers>) is from zero on the number line, regardless of its sign. It transforms any real number into its non-negative counterpart according to the following rule:

\\[|x| = \begin{cases} +x & \text{if } x \geq 0 \quad \forall \, x \in \mathbb{R} \\\\[0.6em] -x & \text{if } x < 0 \quad \forall \, x \in \mathbb{R} \end{cases} \\]

When solving an equation that involves an absolute value, it is essential to consider both possible cases (positive and negative) since the sign of the expression inside the bars determines which definition must be applied. This process often leads to two separate linear equations that need to be examined individually.

## Properties

The [absolute value function](<../absolute-value-function>) \\( |x| \\) obeys several essential algebraic properties that describe how it interacts with basic arithmetic operations and comparisons. Understanding these rules is fundamental when simplifying expressions or solving equations that include absolute values. The following list summarizes the key properties that characterize the behavior of absolute values in real numbers:

\\[|x| = |-x| \quad \forall x \in \mathbb{R} \\]


\\[|x \cdot y| = |x| \cdot |y| \quad \forall x, y \in \mathbb{R} \\]


\\[|x| = |y| \iff x = \pm y \quad \forall x, y \in \mathbb{R} \\]


\\[|x| \leq |y| \iff x^2 \leq y^2 \quad \forall x, y \in \mathbb{R} \\]


\\[\left| \frac{x}{y} \right| = \frac{|x|}{|y|} \quad \forall x, y \in \mathbb{R},\ y \ne 0 \\]


\\[\sqrt{x^2} = |x| \quad \forall x \in \mathbb{R} \\]

##### These properties are useful for manipulating expressions involving absolute values and for simplifying calculations when solving related equations.

## Solving absolute value equations

To solve equations involving absolute values, it is necessary to analyze the expression inside the bars and consider the possible signs it can assume. The general approach depends on the type of equation: whether the absolute value is equal to a constant or to another expression. In each case, the solution process involves separating the equation into distinct cases, reflecting the definition of the absolute value function.

Let’s consider the basic case: \\[|A(x)| = a \\]

In general, if \\( a \geq 0 \\), the equation is equivalent to \\( A(x) = a \\) or \\( A(x) = -a \\). If \\( a < 0 \\), the equation has no solution. For example, the equation \\( |3 + 2x| = -2 \\) has no solution because the absolute value of an expression can never be negative.


Let’s now consider the case:

\\[|A(x)| = B(x) \\]

and solving it requires taking into account the sign of the absolute value when performing the calculations.

## Example 1

Let’s solve the equation: \\[|2x - 4| = x + 1 \\]


First, we analyze the sign of the expression inside the absolute value. We have \\( 2x - 4 \geq 0 \\), which leads to \\( x \geq 2 \\).

According to (1), the equation becomes:

\\[|2x - 4| = \begin{cases} 2x-4 & \text{if } x \geq 2\\\\[0.5em] -2x+4 & \text{if } x < 2 \end{cases} \\]


Let’s now solve the first system given by:

\\[\begin{cases} x \geq 2 \\\\[0.5em] 2x - 4 = x + 1 \end{cases} \\]

\\[\begin{cases} x \geq 2 \\\\[0.5em] 2x - x = + 1 + 4 \rightarrow x = 5 \end{cases} \\]

This solution is acceptable because it satisfies the condition \\( x \geq 2 \\).

##### What we have just seen is a [system of inequalities](<../systems-of-inequalities>) with one variable and two inequalities. Explore the related entry to learn more about the solving process and how to handle more complex cases.


Let’s now solve the second system given by:

\\[\begin{cases} x < 2 \\\\[0.5em] -2x + 4 = x + 1 \end{cases} \\]

\\[\begin{cases} x < 2 \\\\[0.5em] -2x -x =+ 1 - 4\rightarrow -3x = -3 \rightarrow x = 1 \end{cases} \\]

This solution is acceptable because it satisfies the condition \\( x < 2 \\).

The solution to the equation is:

\\[x= 1 \quad x=5\\]

## Example 2

Let’s solve the equation: \\[\frac{|3x|}{|x+1|} = |x|\\]

We are dealing with a rational equation for which it is necessary to determine the conditions of existence, which in this case correspond to the values of \\( x \\) that make the denominator zero. The denominator becomes zero when \\( x = -1 \\), so this value must be excluded from the set of solutions.


Using the properties of absolute value, we can rewrite the equation as:

\\[3 \left| \frac{x}{x+1} \right| = |x|\\]


We must therefore analyze the signs of the expressions. Let’s start by considering the following case:

\\[\frac{3x}{x+1} = x \rightarrow 3x = x^2 + x \rightarrow x^2 -2x \rightarrow x(x - 2)\\]

We have obtained a straightforward [quadratic equation](<quadratic-equations>) whose solutions are:

\\[x=0 \quad x= 2\\]

Both solutions are acceptable because they are different from \\( -1 \\).

##### We solved the quadratic equation without using the [quadratic formula](<quadratic-formula>), by [factoring](<../factoring-quadratic-equations/>) the corresponding [polynomial](<../polynomials>) and finding the values of \\( x \\) that make each linear factor equal to zero.


Let’s now consider the case:

\\[\frac{3x}{x+1} = -x \rightarrow 3x = -x^2 - x \rightarrow x^2 + 4x \rightarrow x(x +4)\\]

Proceeding as above we have obtained a straightforward quadratic equation whose solutions are:

\\[x=0 \quad x= -4\\]

Both solutions are acceptable because they are different from \\( -1 \\).

The solution to the equation is:

\\[x=-4 \quad x=0 \quad x= 2 \\]
