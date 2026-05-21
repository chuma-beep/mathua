> Content sourced from [Algebrica](https://algebrica.org/linear-inequalities/) — CC BY-NC 4.0

## Introduction to inequalities

An inequality is a mathematical statement involving algebraic expressions for which we seek the values of the variables that make the inequality true. In general, an inequality between two algebraic expressions \\( A(x) \\) and \\( B(x) \\) is defined as the relation:

\\[A(x) > B(x) \\]

Solving an inequality involves determining the solution set, that is, all the values of \\( x \\) that satisfy the previous inequality. The solutions of inequalities are subsets of \\( \mathbb{R} \\) defined as intervals.


Given two real numbers \\( a \\) and \\( b \\) with \\( a < b \\), a bounded interval is defined as the set of real numbers between \\( a \\) and \\( b \\), where \\( a \\) and \\( b \\) are the lower and upper bounds, respectively. An unbounded interval is the set of numbers that either precede \\( a \\) or follow \\( a \\). If the endpoints of a bounded interval are included, the interval is called closed, otherwise, it is called open.

![](/diagrams/algebrica/inequalities-1.png)

The inequality sign determines whether the solution interval is open or closed at its boundary. A strict inequality, expressed with \\( > \\) or \\( < \\), excludes the boundary value, yielding an open interval. A non-strict inequality, expressed with \\( \geq \\) or \\( \leq \\), includes it, yielding a closed or half-open interval.

Inequality| Interval  
---|---  
\\( x > a \\)| \\( (a, +\infty) \\)  
\\( x \geq a \\)| \\( [a, +\infty) \\)  
\\( x < a \\)| \\( (-\infty, a) \\)  
\\( x \leq a \\)| \\( (-\infty, a] \\)  
\\( a < x < b \\)| \\( (a, b) \\)  
\\( a \leq x \leq b \\)| \\( [a, b] \\)  
  
###### The previous table summarizes the correspondence between the inequality sign and the resulting interval notation for a single-variable linear inequality.


The degree of an inequality corresponds to the degree of the [polynomial](<../polynomials>) \\( P(x) \\) obtained by rewriting the inequality in the form \\( P(x) > 0 \\). A linear or first-degree inequality is defined in the general form:

\\[ax > b\\]

Every first-degree inequality with \\( a \neq 0 \\) always has an interval of values as its solution. The condition \\( a \neq 0 \\) is essential. When \\( a = 0 \\), the variable disappears entirely from the inequality, and the problem reduces to comparing two constants. Consider the general form with \\( a = 0 \\).

\\[0 \cdot x > b \\]

This simplifies to \\( 0 > b \\). If \\( b < 0 \\), the inequality holds regardless of the value of \\( x \\), and the solution set is all of \\( \mathbb{R} \\). If \\( b \geq 0 \\), the inequality is never satisfied, and the solution set is empty. In either case, there is no interval determined by a boundary point: the outcome is a global truth or a contradiction, not a proper first-degree inequality.


Given an inequality, an equivalent inequality can be obtained by adding the same number or expression to both sides, provided that the expression is well-defined within the same [domain](<../determining-the-domain-of-a-function/>). The inequality \\( 5x - 2 > x + 3 \\) is equivalent to the inequality \\( 4x - 2 > 3 \\), by subtracting \\( x \\) from both sides.

An equivalent inequality can also be obtained by dividing or multiplying both sides by the same non-zero value. However, if the value is negative, the inequality sign must be reversed. It is always useful to rewrite a first-degree inequality like \\( -x + 5 < -3 \\) in the form \\( x - 5 > 3 \\). Changing the signs requires reversing the direction of the inequality.

## Geometric interpretation

The solution set of a linear inequality always corresponds to a half-line on the real number line, that is, an unbounded interval extending indefinitely in one direction from a boundary point.

  * When the inequality is strict, the boundary point is excluded and the half-line is open at that end.

  * When the inequality is non-strict, the boundary point belongs to the solution set and the half-line is closed.


Consider, for instance, the simple inequality \\( x \geq 1 \\). Its solution set is the closed half-line \\( [1, +\infty) \\), represented below.

| \\[1 \\]|   
---|---|---  
| |   
| |   
  
This geometric reading clarifies why the solution of a first-degree inequality is never an isolated point or a bounded interval: the linear structure of the expression \\( ax - b \\) ensures that its sign changes exactly once, at \\( x = b/a \\), dividing the real line into precisely two regions, one of which constitutes the solution.

## Example 1

Consider the following inequality:

\\[-\frac{x}{2} + 5 > 2x - 5\\]

The goal is to reduce it to the standard form \\( ax > b \\). Since the left-hand side contains a fraction with denominator \\( 2 \\), multiplying both sides by \\( 2 \\) clears the denominator without altering the direction of the inequality, as the multiplier is positive.

\\[-x + 10 > 4x - 10\\]

Collecting all terms involving \\( x \\) on the left and moving the constant terms to the right yields the following.

\\[-5x > -20\\]

Dividing both sides by \\( -5 \\) isolates the variable. Since the divisor is negative, the direction of the inequality must be reversed.

\\[x < 4\\]

The solution set is the open interval \\( (-\infty, 4) \\).

## How to solve literal first-degree inequalities

In the case of literal first-degree inequalities, the coefficient of the variable is not a fixed constant but depends on one or more parameters. Consider the following inequality.

\\[z(x - 1) < 2x + 1\\]

The first step is to expand and collect terms so as to rewrite the expression in one of the standard forms \\( Ax > B \\), \\( Ax < B \\), \\( Ax \geq B \\), or \\( Ax \leq B \\). Expanding the left-hand side and moving all terms involving \\( x \\) to the left yields the following.

\\[\begin{align} &zx - z < 2x + 1 \\\\[6pt] &zx - 2x < z + 1 \\\\[6pt] &(z - 2)x < z + 1 \end{align} \\]

The inequality is now in the form \\( Ax < B \\) with \\( A = z - 2 \\) and \\( B = z + 1 \\). Since dividing both sides by \\( A \\) requires knowing its sign, and the sign of \\( z - 2 \\) depends on the value of the parameter \\( z \\), it is necessary to distinguish three separate cases.


When \\( z > 2 \\), the coefficient \\( z - 2 \\) is positive. Dividing both sides by a positive quantity preserves the direction of the inequality, giving the following.

\\[x < \frac{z + 1}{z - 2}\\]


When \\( z = 2 \\), the coefficient \\( z - 2 \\) vanishes. Substituting this value reduces the inequality to \\( 0 \cdot x < 3 \\), which holds for every \\( x \in \mathbb{R} \\). The solution set is all of \\( \mathbb{R} \\).


When \\( z < 2 \\), the coefficient \\( z - 2 \\) is negative. Dividing both sides by a negative quantity reverses the direction of the inequality, giving the following.

\\[x > \frac{z + 1}{z - 2}\\]

## Absolute value inequalities

Absolute value inequalities are inequalities that follow the properties of [absolute value](<../absolute-value>) and are of the form:

\\[|x| < z \\]

In general, inequalities involving absolute values do not belong to the class of linear inequalities, since the presence of the absolute value function alters the behavior of the expression by introducing a discontinuity that breaks its linear structure. However, through an appropriate decomposition process, they can be transformed into one or more equivalent systems of linear inequalities, allowing them to be analyzed using the same methods applied to first-degree inequalities.


To solve them, we need to consider the sign of the absolute value by dividing it into two cases:

  * \\( |x| < z \\) if and only if \\(-z < x < z\\).
  * \\( |x| > z \\) if and only if \\( x < -z \quad \text{or} \quad x > z \\).


We must therefore solve a system of inequalities resulting from the possible cases considered in the analysis. The process is slightly more complicated than simple linear equations, and careful attention must be paid to the signs and the direction of the inequality during calculations.

## Example 2

Let’s try to solve the following absolute value inequality:

\\[|x - 1| < 2x + 4 \\]


First, we study the sign of \\(|x - 1|\\). We have:

\\[|x - 1| = \begin{cases} x - 1 & \quad x \ge 1 \\\\[0.8em] -x + 1 & \quad x < 1 \end{cases} \\]


Substituting the first branch into the inequality and solving yields the following system.

\\[\begin{cases} x \geq 1 \\\\[6pt] x - 1 < 2x + 4 \end{cases} \\]

| \\[-5\\]| \\[1\\]|   
---|---|---|---  
| | |   
| | |   
| | |   
| | |   
  
The second inequality reduces to \\( x > -5 \\). The solution of the system is the intersection of \\( x \geq 1 \\) and \\( x > -5 \\), which is \\( x \geq 1 \\).


Substituting the second branch into the inequality and solving yields the following system:

\\[\begin{cases} x < 1 \\\\[6pt] -x + 1 < 2x + 4 \end{cases} \\]

| \\( -1 \\)| \\( 1 \\)|   
---|---|---|---  
| | |   
| | |   
| | |   
  
The second inequality reduces to \\( x > -1 \\). The solution of the system is the intersection of \\( x < 1 \\) and \\( x > -1 \\), which is the open interval \\( (-1, 1) \\).

The complete solution is obtained by taking the union of the two partial solution sets. Since \\( x \geq 1 \\) and \\( (-1, 1) \\) are adjacent intervals that together cover all values greater than \\( -1 \\), the solution to the inequality is the following.

\\[x > -1\\]

## Selected references

  * **Stony Brook University**. [Linear Inequalities](https://www.math.stonybrook.edu/Videos/MAP103Online/Handouts/Lecture-17-Handout.pdf)

  * **Stony Brook University**. [Absolute Value Inequalities](https://www.math.stonybrook.edu/Videos/MAP103Online/Handouts/Lecture-18-Handout.pdf)

  * **Harvard University**. [Appendix H: Interpreting and Working with Inequalities](https://abel.math.harvard.edu/archive/xb_spring_03/icearchive1/webAPPH.pdf)

  * **University of Houston**. [Interval Notation and Linear Inequalities](https://online.math.uh.edu/Math1300-unpaid/ch1/s17/1300_Ch1_Section7.pdf)


Inequalities

Inequalities express order relations between quantities, defining ranges and constraints.

2.5k

[Quadratic Inequalities](https://algebrica.org/quadratic-inequalities/)

2.7k

[Sign Analysis in Inequalities](https://algebrica.org/sign-analysis-in-inequalities/)

1.2k

[Rational Inequalities](https://algebrica.org/rational-inequalities/)

1.1k

[Irrational Inequalities](https://algebrica.org/irrational-inequalities/)

1.2k

[Inequalities with Absolute Value](https://algebrica.org/inequalities-with-absolute-value/)

3.6k

[Logarithmic Inequalities](https://algebrica.org/logarithmic-inequalities/)

869

[Trigonometric Inequalities](https://algebrica.org/trigonometric-inequalities/)

1.3k

[Systems of Inequalities](https://algebrica.org/systems-of-inequalities/)
