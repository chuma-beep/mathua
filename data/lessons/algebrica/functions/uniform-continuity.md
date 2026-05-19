> Content sourced from [Algebrica](https://algebrica.org/uniform-continuity/) — CC BY-NC 4.0

## Introduction

[Ordinary continuity](<../continuous-functions/>) describes the local behaviour of a [function](<../functions/>), where small changes in the input near each point result in small changes in the output. This property can depend on the specific point, and in many contexts, a more robust form of continuity is required to ensure consistent regulation of oscillations across the entire [domain](<../determining-the-domain-of-a-function/>). Uniform continuity addresses this need by requiring a single tolerance for input variations to apply everywhere in the set.

This stability property leads to several important consequences:

  * Continuous functions can be extended from dense subsets to their closures.
  * [Integrability](<../indefinite-integrals/>) and [limit](<../limits/>) operations behave in a controlled, predictable manner.
  * Uncontrolled oscillations near infinity are prevented whenever the domain is compact.


###### Oscillations refer to the variation in the output values of a function over a given region, specifically how much the function rises and falls within a small portion of the domain. For example, the function \\(\sin(1/x)\\) near the origin exhibits infinite oscillations within any small interval, which makes its behaviour increasingly difficult to control.

## Definition

Consider a function \\( f : A \subset \mathbb{R} \to \mathbb{R} \\). The function is uniformly continuous on \\( A \\) if, for every \\( \varepsilon > 0 \\) there exists a \\( \delta > 0 \\) such that for all \\( x, y \in A \\) the following implication holds:

\\[|x - y| < \delta \to |f(x) - f(y)| < \varepsilon \\]

The essential feature of this definition is that the value of \\( \delta \\) depends solely on \\( \varepsilon \\), and not on the particular points \\( x \\) and \\( y \\). The same \\( \delta \\) must be valid for every pair of points in the domain.

![Uniform continuity.](https://algebrica.org/wp-content/uploads/resources/images/uniform-continuity-1.png)

The graph shows the global nature of the condition, highlighting that the same \\( \delta \\) controls the function’s variation throughout the domain. Specifically, whenever \\( |x - y| < \delta \\), it follows that \\( |f(x) - f(y)| < \varepsilon \\) for all \\( x, y \in A \\).


The distinction between [continuity](https://algebrica.org/continuous-functions/) and uniform continuity is clarified by comparing the standard definition of continuity at a point \\( x_0 \\). A function is continuous at \\( x_0 \\) if, for every \\( \varepsilon > 0 \\), there exists a \\( \delta > 0 \\), which may depend on \\( x_0 \\), such that:

\\[|x - x_0| < \delta \to |f(x) - f(x_0)| < \varepsilon \\]

The distinction lies in the dependence of \\( \delta \\):

  * For ordinary continuity, \\( \delta \\) may vary with each point in the domain.
  * For uniform continuity, a single value of \\( \delta \\) is valid for all points in the set.


## Heine–Cantor Theorem

The Heine–Cantor theorem establishes the conditions under which continuity implies uniform continuity. Specifically, on compact sets, the local property of continuity suffices to guarantee the existence of a single global control parameter \\( \delta \\) applicable to the entire domain.

Let \\( K \subset \mathbb{R} \\) be a compact set. If a function \\( f : K \to \mathbb{R} \\) is continuous on \\( K \\), then it is uniformly continuous on \\( K \\). In particular, since compactness in \\( \mathbb{R} \\) is equivalent to the set being closed and bounded, the theorem admits the following equivalent formulation:

If \\( f : [a,b] \to \mathbb{R} \\) is continuous on a closed and bounded interval \\( [a,b] \\), then \\( f \\) is uniformly continuous on \\( [a,b] \\).

###### The theorem identifies compactness as the precise structural property of the domain that prevents local continuity from degenerating into purely pointwise behaviour, ensuring instead a global form of regularity.

## Sequential characterization

Uniform continuity can be equivalently formulated in terms of [sequences](<../sequences/>). Specifically, a function \\( f : A \to \mathbb{R} \\) is uniformly continuous on \\( A \\) if and only if, for every pair of sequences \\( (x_n) \\) and \\( (y_n) \\) in \\( A \\) that satisfy \\(|x_n - y_n| \to 0\\) we also have:

\\[|f(x_n) - f(y_n)| \to 0 \\]

This formulation highlights that uniformly continuous functions preserve infinitesimal closeness between sequences, which is also why they map [Cauchy sequences](<../cauchy-sequence/>) to Cauchy sequences.

###### This criterion is particularly useful for establishing that a function is not uniformly continuous, as shown in Example 1.

## Example 1

Consider the function \\( f(x) = x^2 \\) on the interval \\( \mathbb{R} \\). This function is continuous everywhere, since it is a [polynomial](<../polynomials/>), but it is not uniformly continuous on \\( \mathbb{R} \\). The reason is that the function grows faster and faster as \\( |x| \\) increases. For points far from the origin, even a small variation in ( x ) can produce a large variation in \\( f(x) \\) and no single \\( \delta \\) can control this behaviour simultaneously across the real line.

In contrast, when the function is restricted to a closed and bounded interval \\( [a, b] \\), it becomes uniformly continuous. On such domains, the growth of the function is effectively bounded. This can be made precise using the sequential characterization and considering the sequences defined by \\( x_n = n + \frac{1}{n} \\) and \\( y_n = n. \\) The difference between these sequences is given by:

\\[|x_n - y_n| = \frac{1}{n} \to 0\\]

The difference between the images of these sequences under the function satisfies:

\\[|f(x_n) - f(y_n)| = \left(n + \frac{1}{n}\right)^2 - n^2 = 2 + \frac{1}{n^2} \to 2\\]

Since \\( |f(x_n) - f(y_n)| \\) does not tend to zero, the function fails the sequential criterion for uniform continuity, therefore, the function is not uniformly continuous on \\( \mathbb{R} \\).


This example clarifies the relationship between continuity and uniform continuity. In general terms we have:

  * Continuity does not imply uniform continuity. A function may be continuous at every point of its domain without satisfying a single global \\( \delta \\) that works uniformly for all pairs of points.

  * Uniform continuity, however, is a stronger condition. If a function is uniformly continuous on \\( A \\), then it is necessarily continuous at every point of \\( A \\).


## Lipschitz Continuity

A particular form of uniform continuity is Lipschitz continuity. A function \\( f : A \subset \mathbb{R} \to \mathbb{R} \\) is said to be Lipschitz continuous on \\( A \\) if there exists a constant \\( L > 0 \\), called a Lipschitz constant, such that, for all ( x,y \in A ) we have:

\\[|f(x) - f(y)| \le L |x - y| \\]

This condition guarantees that nearby points are mapped to nearby values and establishes a global linear bound on the function’s oscillation: the output’s variation is proportional to the input’s variation. The same constant \\( L \\) must work everywhere on the domain. The implication for uniform continuity follows directly. Given any \\( \varepsilon > 0 \\), selecting \\( \delta = \varepsilon / L \\) guarantees that:

\\[|x - y| < \delta \to |f(x) - f(y)| \le L |x - y| < L \delta = \varepsilon \\]

Therefore, every Lipschitz continuous function is uniformly continuous.

However, the converse does not hold in general. A function may be uniformly continuous without satisfying any global linear bound of this form. Thus, Lipschitz continuity constitutes a strictly stronger condition.

A useful criterion links Lipschitz continuity to [differentiability](<../derivatives/>). If a function \\( f \\) is differentiable on an interval and its derivative is bounded, that is, there exists \\( M > 0 \\) such that for all \\( x \\) in the interval:

\\[|f’(x)| \le M \\]

Then, by the [Mean Value Theorem](<../lagrange-theorem/>), \\( f \\) is Lipschitz continuous with Lipschitz constant \\( L = M \\). For any two points \\( x, y \\) in the interval, there exists \\( c \\) between them such that:

\\[f(x) - f(y) = f’(c)(x - y) \\]

Taking [absolute values](<../absolute-value/>) then yields \\( |f(x) - f(y)| \leq M|x - y| \\), which is the Lipschitz condition with constant \\( L = M \\).

###### Geometrically, Lipschitz continuity restricts the steepness of the graph. The slopes of all secant lines are uniformly bounded in magnitude by \\( L \\). Thus, the function cannot oscillate more rapidly than a fixed linear rate.

## Example 2

###### The Heine–Cantor theorem guarantees uniform continuity whenever the [domain](<../determining-the-domain-of-a-function/>) is compact. However, compactness is a sufficient but not necessary condition, and uniform continuity can also hold on domains that are not compact. This example demonstrates this situation.

Consider the function \\( f(x) = \sin(x) \\) defined on the entire real line \\( \mathbb{R} \\). The domain is unbounded and not compact, so uniform continuity is not guaranteed a priori and must be verified directly.

Given \\( \varepsilon > 0 \\) and arbitrary \\( x, y \in \mathbb{R} \\), the objective is to find a single \\( \delta > 0 \\), depending only on \\( \varepsilon \\), such that whenever \\( |x - y| < \delta \\), it follows that \\( |\sin(x) - \sin(y)| < \varepsilon \\). Applying the [sum-to-product identity](<../trigonometric-identities/>) yields:

\\[|\sin(x) - \sin(y)| = 2\left|\cos\\!\left(\frac{x+y}{2}\right)\sin\\!\left(\frac{x-y}{2}\right)\right|\\]

Since \\( |\cos(\cdot)| \leq 1 \\) everywhere, this simplifies to:

\\[|\sin(x) - \sin(y)| \leq 2\left|\sin\\!\left(\frac{x-y}{2}\right)\right|\\]

Applying the standard inequality \\( |\sin(t)| \leq |t| \\), valid for all \\( t \in \mathbb{R} \\), gives:

\\[|\sin(x) - \sin(y)| \leq 2 \cdot \frac{|x - y|}{2} = |x - y|\\]

This satisfies the Lipschitz condition with constant \\( L = 1 \\), which holds for every pair of points in \\( \mathbb{R} \\) without restriction. By selecting \\( \delta = \varepsilon \\), whenever \\( |x - y| < \delta \\), it follows that:

\\[|\sin(x) - \sin(y)| \leq |x - y| < \delta = \varepsilon\\]

Here, \\( \delta \\) depends solely on \\( \varepsilon \\), therefore, \\( f \\) is uniformly continuous on \\( \mathbb{R} \\).

###### The underlying reason is geometric: the [sine function](<../sine-function/>) oscillates indefinitely, but its oscillations remain permanently bounded in amplitude, and its slope never exceeds one in absolute value, since \\( |f’(x)| = |\cos(x)| \leq 1 \\) everywhere. This global bound on the rate of change prevents the uncontrolled growth observed in Example 1 and ensures that a uniform \\( \delta \\) is attainable across the entire real line.

## Summary  
  
---  
Pointwise continuity does not imply uniform continuity.  
Uniform continuity implies pointwise continuity.  
Lipschitz continuity implies uniform continuity.  
  
## Selected references

  * **MIT, S. Rodriguez**. [Uniform Continuity and Derivative](https://ocw.mit.edu/courses/18-100a-real-analysis-fall-2020/mit18_100af20_lec17.pdf)

  * **University of Wisconsin–Madison, J.W. Robbin**. [Continuity and Uniform Continuity](https://people.math.wisc.edu/~jwrobbin/521dir/cont.pdf)

  * **University of Maryland, D. Cellarosi**. [Uniform Continuity](https://www.math.umd.edu/~immortal/MATH410/lecturenotes/ch3-4.pdf)

  * **University of California, N. Donaldson**. [Continuity and Uniform Continuity](https://www.math.uci.edu/~ndonalds/math140b/1contrev.pdf)

  * **Michigan State University, J. Shapiro**. [Notes on Uniform Continuity](https://users.math.msu.edu/users/shapiro/Teaching/classes/320/Handouts/Unif_contin_notes.pdf)


Functions

A function maps each input to a unique output.

17.6k

[Functions](https://algebrica.org/functions/)

3 comments

10.2k

[Determining the Domain of a Function](https://algebrica.org/determining-the-domain-of-a-function/)

1.6k

[Even and Odd Functions](https://algebrica.org/even-and-odd-functions/)

2.7k

[Increasing, Decreasing and Monotonic Functions](https://algebrica.org/increasing-and-decreasing-functions/)

1.9k

[Convexity and Concavity of Functions](https://algebrica.org/convexity-and-concavity-of-functions/)

1.3k

[Composite Functions](https://algebrica.org/composite-functions/)

1.2k

[Inverse Function](https://algebrica.org/inverse-function/)

1.9k

[Continuous Functions](https://algebrica.org/continuous-functions/)

2 comments

1.3k

[Discontinuities of Real Functions](https://algebrica.org/discontinuities-of-real-functions/)

1.9k

[Analyzing the Graphs of Functions](https://algebrica.org/analyzing-the-graphs-of-functions/)

1.1k

[Polynomial Function](https://algebrica.org/polynomial-function/)

2.6k

[Rational Functions](https://algebrica.org/rational-functions/)

2k

[Logarithmic Function](https://algebrica.org/logarithmic-function/)

3.1k

[Exponential Function](https://algebrica.org/exponential-function/)

1.5k

[Absolute Value Function](https://algebrica.org/absolute-value-function/)

847

[Sign Function](https://algebrica.org/sign-function/)

2.2k

[Sine Function](https://algebrica.org/sine-function/)

2k

[Cosine Function](https://algebrica.org/cosine-function/)

1.7k

[Tangent Function](https://algebrica.org/tangent-function/)

2k

[Cotangent Function](https://algebrica.org/cotangent-function/)

1.8k

[Secant Function](https://algebrica.org/secant-function/)

1.4k

[Cosecant Function](https://algebrica.org/cosecant-function/)

1.2k

[Dirichlet Function](https://algebrica.org/dirichlet-function/)

1.4k

[Sigmoid Function](https://algebrica.org/sigmoid-function/)
