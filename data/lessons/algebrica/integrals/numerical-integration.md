> Content sourced from [Algebrica](https://algebrica.org/numerical-integration/) — CC BY-NC 4.0

## Integrals without elementary antiderivativess359ebbb1

1.The [Fundamental Theorem of Calculus](<../fundamental-theorem-of-calculus/>) gives a way for evaluating a [definite integral](<../definite-integrals/>). If an antiderivative exists, the value of the integral becomes the difference between the interval endpoints. However, many integrands do not have a primitive expressible by elementary [functions](<../functions>) and method such as [substitution](<../integration-by-substitution/>), [integration by parts](<../integration-by-parts/>), and [Weierstrass substitution](<../the-weierstrass-substitution/>) often cannot yield a closed form.

2.A typical example is provided by the following integral, which appears in the theory of the [normal distribution](<../normal-distribution>) and whose integrand has no antiderivative within the elementary functions:

$$\int_0^1 e^{-x^2}\\\\\,dx$$

3.The same difficulty occurs with integrals containing $\sin(x)/x$, with elliptic integrals, and with a wide class of expressions involving combinations of algebraic and transcendental terms.

4.The numerical evaluation of such integrals follows a different conceptual route. Instead of looking for an exact symbolic expression, one builds an approximation whose accuracy can be controlled and improved as needed. The branch of analysis that studies these procedures is called numerical integration, the latter term being a remembrance of the geometric origin of the integral as the area of a planar region.

> 5.
> 
> Numerical methods become indispensable when the analytical route is unavailable or prohibitively expensive but when a primitive exists and is reasonably accessible, the exact evaluation through the Fundamental Theorem of Calculus remains the method of choice.

## General principle of a quadrature formulas345e11a4

6.Numerical integration rests on the same construction that defines the [Riemann integral](<../riemann-integrability-criteria/>): the interval of integration is partitioned into a finite number of subintervals, the integrand is replaced on each subinterval by a simpler function whose integral is known exactly, and the total area is approximated by summing the contributions of the individual pieces. The quality of the approximation is determined by two factors: the width of the subintervals and the order of accuracy of the local rule.

x y y=f(x) a h b The interval [a, b] has been divided  into 6 subintervals of width  h , whose  sum provides a lower approximation  of the area of the curvilinear  trapezoid.

7.Consider a uniform partition of the interval $[a,b]$:

$$ a = x_0 < x_1 < \cdots < x_n = b \quad x_k = a + k h $$

8.with constant step size:

$$ h = \frac{b-a}{n} $$

9.A quadrature formula approximates the integral by a finite linear combination of values of the integrand at the nodes of the partition:

$$ F_1 = \int_a^b f(x)\\\\\,dx \approx \sum_{k=0}^{n} w_k\\\\\,f(x_k) $$

10.The coefficients $w_k$ are called weights, and their choice determines the specific method. The simplest schemes arise by interpolating the integrand on each subinterval, or on each group of consecutive subintervals, by a low-degree polynomial, and then integrating that polynomial exactly. The degree of the interpolating polynomial determines both the form of the resulting formula and its order of accuracy.

11.**Definition 1.** A quadrature rule is said to have degree of exactness $m$ if it integrates exactly every polynomial of degree at most $m$.

> 12.
> 
> This notion provides a useful theoretical measure of the precision of a method, and it is the starting point for the construction of more refined formulas such as the Gauss quadrature rules.

## The rectangle and midpoint ruless09425c34

13.The most elementary quadrature rule approximates the integrand on each subinterval by a constant. Depending on whether one chooses the value of the integrand at the left endpoint, the right endpoint, or the midpoint, three variants of the rectangle rule are obtained. The left and right versions reproduce the Riemann sums already encountered in the construction of the definite integral. The midpoint version deserves particular attention because of its superior accuracy and the symmetry of its construction. Introducing the midpoints of the subintervals:

$$ \bar{x}_k = \frac{x_{k-1} + x_k}{2} $$

14.the midpoint rule on a single subinterval reads:

$$ \int_{x_{k-1}}^{x_k} f(x)\\\\\,dx \approx h\\\\\,f(\bar{x}_k) $$

15.Summing the local contributions over all subintervals, one obtains the composite midpoint formula:

$$ F_2 = \int_a^b f(x)\\\\\,dx \approx h \sum_{k=1}^{n} f(\bar{x}_k) $$

16.**Theorem 1.** When the integrand is of class $C^2$ on $[a,b]$, the error committed by the formula $F_2$ satisfies the bound:

$$ \left| \int_a^b f(x)\\\\\,dx - h\sum_{k=1}^{n} f(\bar{x}_k) \right| \le \frac{(b-a)\\\\\,h^2}{24} \max_{x \in [a,b]} |f''(x)| $$

17.The error therefore decreases as $h^2$ and halving $h$ reduces the bound by a factor of 4. This is already a clear improvement over the left and right rectangle rules, whose error scales linearly with $h$ and which therefore converge much more slowly to the exact value of the integral.

> 18.
> 
> The midpoint rule integrates exactly every polynomial of degree at most one. The geometric reason is that the area of a rectangle whose base coincides with the subinterval and whose height equals the value of an affine function at the midpoint is exactly the integral of the affine function over the subinterval. This degree of exactness is the source of the second-order convergence.

## The trapezoidal rules17635cbc

19.A more accurate approximation is obtained by replacing the integrand on each subinterval not by a constant but by an affine function, namely the segment joining the two graph points $(x_{k-1}, f(x_{k-1}))$ and $(x_k, f(x_k))$. The region underneath this segment is a trapezoid, and its area equals the average of the two ordinates multiplied by the base.

x y y=f(x) a b h Each subinterval generates a trapezoid: h is the  base, while f(x ₖ₋ ₁) and f(x ₖ ) are the two heights  used to approximate the area under the curve. x ₖ f(x ₖ ) x ₖ₋ ₁ f(x ₖ₋ ₁)

20.The local formula reads:

$$ \int_{x_{k-1}}^{x_k} f(x)\\\\\,dx \approx \frac{h}{2}\bigl[f(x_{k-1}) + f(x_k)\bigr] $$

21.Summing over all subintervals and observing that each internal node belongs to two adjacent subintervals, and therefore appears with multiplicity two in the global sum, one obtains the composite trapezoidal formula:

$$ F_3 = \int_a^b f(x)\\\\\,dx \approx \frac{h}{2}\Bigl[f(a) + f(b) + 2 \sum_{k=1}^{n-1} f(x_k)\Bigr] $$

22.**Theorem 2.** If the integrand is of class $C^2$ on $[a,b]$, the error of the composite trapezoidal formula admits the bound:

$$ \left| \int_a^b f(x)\\\\\,dx - \frac{h}{2}\Bigl[f(a) + f(b) + 2 \sum_{k=1}^{n-1} f(x_k)\Bigr] \right| \le \frac{(b-a)\\\,h^2}{12} \max_{x \in [a,b]} |f''(x)| $$

23.The decay is again of order $h^2$, and the asymptotic behaviour coincides with that of the midpoint rule. The constant appearing in the trapezoidal bound is, however, twice as large as that for the midpoint rule, so the two methods are comparable in convergence order, but the midpoint rule is slightly more accurate in the leading constant.

> 24.
> 
> The trapezoidal and midpoint rules share the same degree of exactness, namely one. The reason why the trapezoidal formula carries a larger error constant lies in the fact that it samples the integrand only at the endpoints of each subinterval, where the error of polynomial interpolation tends to be larger, while the midpoint rule samples at the centre, where this error is naturally smaller. This observation already suggests that not all sampling strategies are equally efficient, an idea that lies at the heart of the Gauss quadrature rules.

## Simpson's rules22a6fce0

25.A substantial gain in accuracy is achieved by approximating the integrand on a pair of consecutive subintervals by a polynomial of degree two, rather than by a piecewise affine function. The construction proceeds as follows. Consider three consecutive equally spaced nodes $x_{k-1}, x_k, x_{k+1}$, and let $P(x)$ denote the unique polynomial of degree at most two passing through the three points $(x_{k-1}, f(x_{k-1}))$, $(x_k, f(x_k))$, $(x_{k+1}, f(x_{k+1}))$. The integral of this polynomial over the pair of subintervals can be computed in closed form, for instance by integrating the Lagrange representation of $P(x)$, and the result reads:

$$ \int_{x_{k-1}}^{x_{k+1}} P(x)\\\\\,dx = \frac{h}{3}\bigl[f(x_{k-1}) + 4 f(x_k) + f(x_{k+1})\bigr] $$

26.Substituting this expression in place of the exact integral one obtains Simpson's formula on a single pair of subintervals. For the composite version, the number of subdivisions $n$ must be even, so that the whole interval can be covered by $n/2$ non-overlapping consecutive pairs. Grouping the nodes accordingly, one arrives at the composite Simpson formula:

$$ F_4 = \int_a^b f(x)\\\\\,dx \approx \frac{h}{3}\Bigl[f(a) + f(b) + 4 \sum_{k \text{ odd}} f(x_k) + 2 \sum_{k \text{ even}} f(x_k)\Bigr] $$

27.where the odd-indexed sum extends over $k = 1, 3, \dots, n-1$ and the even-indexed sum extends over the interior even nodes $k = 2, 4, \dots, n-2$. When the integrand is of class $C^4$ on $[a,b]$, the error of the composite Simpson formula satisfies:

$$ \left| \int_a^b f(x)\\\\\,dx - S_n \right| \le \frac{(b-a)\\\,h^4}{180} \max_{x \in [a,b]} |f^{(4)}(x)| $$

28.where $S_n$ stands for the right-hand side of formula $F_4$. The error now decreases as the fourth power of the step size, and halving $h$ reduces the bound by a factor of sixteen. This represents a dramatic improvement over the second-order behaviour of the trapezoidal and midpoint rules, and it explains why Simpson's formula has historically occupied a central role in numerical quadrature.

> 29.
> 
> Although Simpson's rule is built by interpolating a [polynomial](<../polynomials>) of degree two, an additional cancellation of symmetric terms in the error expansion makes the rule exact on polynomials of degree three as well. The degree of exactness is therefore three, not two, and the gain of two orders of accuracy with respect to the trapezoidal rule is the manifestation of this fortunate phenomenon.

## Comparison and order of convergences3d833505

30.The three formulas just discussed belong to the family of Newton-Cotes rules of low order, characterised by the use of equally spaced nodes and by the integration of a polynomial interpolant of fixed degree. The order of convergence summarises in a single number the rate at which the error decreases as the number of subdivisions grows. The following table collects the relevant information.

Rule| Degree of exactness| Global error  
---|---|---  
Midpoint| 1| $O(h^2)$  
Trapezoidal| 1| $O(h^2)$  
Simpson| 3| $O(h^4)$  
  
31.The qualitative consequence is striking. To divide the error by a factor of one hundred it is sufficient to roughly triple the number of subdivisions when using Simpson's method, while the trapezoidal rule would require multiplying that number by ten. Since the cost of a quadrature method is dominated by the number of evaluations of the integrand, and since each evaluation can be expensive in real applications, the practical advantage of Simpson's rule is considerable.

> 32.
> 
> The qualitative difference between second-order and fourth-order convergence is decisive in any setting where high precision is required from a limited budget of function evaluations. This is why Simpson's rule, and the higher-order Newton-Cotes formulas obtained by analogous constructions, occupy a privileged position among elementary quadrature methods.

## Example 1se93b3fa4

33.We illustrate the use of the trapezoidal and Simpson rules on the integral:

$$ \int_0^1 e^{-x^2}\\\\\,dx $$

34.The integrand is the kernel of the Gauss error function and admits no antiderivative expressible in elementary form, so the integral cannot be evaluated by the [Fundamental Theorem of Calculus](<../fundamental-theorem-of-calculus/>) in any direct way. The reference value, obtained with arbitrary precision, is approximately $0.7468241328$, and it provides a yardstick against which our approximations can be tested.

35.We choose $n = 4$ subintervals of equal width $h = 1/4$. The nodes of the partition and the corresponding values of the integrand are listed in the following table.

k| $x_k$| $f(x_k) = e^{-x_k^2}$  
---|---|---  
0| 0.00| 1.000000  
1| 0.25| 0.939413  
2| 0.50| 0.778801  
3| 0.75| 0.569783  
4| 1.00| 0.367879  
  
36.Applying the composite trapezoidal formula $(3)$, the endpoint values contribute with weight one and the three interior values contribute with weight two. A direct substitution gives:

$$ \begin{align} T_4 &= \frac{0.25}{2} \, \Bigl[1.000000 + 0.367879 + 2(0.939413 + 0.778801 + 0.569783)\Bigr] \\\\\\\\[0pt] &= 0.125 \, (1.367879 + 4.575994) \\\\\\\\[3pt] &= 0.742984 \end{align} $$

37.The discrepancy with the reference value is approximately $3.84 \times 10^{-3}$, which is consistent with the second-order error bound that governs the trapezoidal rule.

38.Applying the composite Simpson formula $(4)$, the odd-indexed nodes $x_1$ and $x_3$ contribute with weight four and the even-indexed interior node $x_2$ contributes with weight two. A direct substitution gives:

$$\begin{align} S_4 &= \frac{0.25}{3} \, \Bigl[1.000000 + 0.367879 + 4(0.939413 + 0.569783) + 2(0.778801)\Bigr]\\\\\\\\[0pt] &= \frac{1}{12} \, (1.367879 + 6.036784 + 1.557602) \\\\\\\\[3pt] &= 0.746855 \end{align}$$

39.The discrepancy with the reference value is now approximately $3.1 \times 10^{-5}$, more than two orders of magnitude smaller than the trapezoidal error obtained with the same number of nodes. The comparison confirms, in a concrete setting, the theoretical prediction that Simpson's rule produces results of much higher quality than the trapezoidal rule at the same computational cost.

## Beyond Newton-Cotess58a556c9

40.The methods presented above are the simplest representatives of a much larger landscape of quadrature techniques. The Newton-Cotes family can be extended to higher-degree interpolants, although the resulting formulas develop oscillatory weights and lose stability beyond degree seven or so, and are therefore rarely used in practice. A more powerful idea is to abandon the assumption of equally spaced nodes and to choose both the nodes and the weights so as to maximise the degree of exactness for a fixed number of evaluations.

41.The result is the family of Gauss quadrature rules, in which $n$ nodes are sufficient to integrate exactly every polynomial of degree at most $2n - 1$. Closely related are the adaptive methods, which subdivide the interval of integration more finely in the regions where the integrand varies rapidly and more coarsely where it is well behaved, and the Romberg method, which combines several trapezoidal estimates with different step sizes to extract a much more accurate value through Richardson extrapolation.

42.Numerical quadrature also requires care when applied to [improper integrals](<../improper-integrals/>), to integrands with rapidly oscillating behaviour, or to functions that approach a singularity within or at the boundary of the interval. In all these cases, the elementary rules introduced in this lecture lose accuracy and must be modified, either by a preliminary analytical transformation that removes the singularity or by specialised methods tailored to the specific situation.

43.The framework outlined here, nonetheless, contains the essential conceptual ingredients on which the more advanced techniques are built, and it is the natural starting point for any deeper study of numerical analysis.

Technique

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Advanced

1

Requires

0

Enables

The following concepts, [Definite Integrals](https://algebrica.org/definite-integrals/), are required as prerequisites for this entry.

Integrals

An integral represents a function’s accumulation over an interval, often as area under a curve.

9.3k

[Indefinite Integrals](https://algebrica.org/indefinite-integrals/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/indefinite-integrals.md?plain=1)

8.6k

[Definite Integrals](https://algebrica.org/definite-integrals/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/definite-integrals.md?plain=1)

1.5k

[Fundamental Theorem of Calculus](https://algebrica.org/fundamental-theorem-of-calculus/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/fundamental-theorem-of-calculus.md?plain=1)

6.2k

[Integration by Substitution](https://algebrica.org/integration-by-substitution/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/integration-by-substitution.md?plain=1)

4k

[Integration by Parts](https://algebrica.org/integration-by-parts/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/integration-by-parts.md?plain=1)

3.9k

[Finding Areas by Integration](https://algebrica.org/finding-areas-by-integration/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/finding-areas-by-integration.md?plain=1)

4.5k

[Integral of the Exponential Function](https://algebrica.org/integral-of-the-exponential-function/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/integral-of-the-exponential-function.md?plain=1)

1.7k

[Integral of Trigonometric Functions](https://algebrica.org/integral-of-trigonometric-functions/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/integral-of-trigonometric-functions.md?plain=1)

3.4k

[Integral of Rational Functions](https://algebrica.org/integral-of-rational-functions/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/integral-of-rational-functions.md?plain=1)

1.9k

[Trigonometric Substitution for Integrals](https://algebrica.org/trigonometric-substitution-for-integrals/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/trigonometric-substitution-for-integrals.md?plain=1)

1.3k

[The Weierstrass Substitution](https://algebrica.org/the-weierstrass-substitution/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/the-weierstrass-substitution.md?plain=1)

5.2k

[Riemann Integrability Criteria](https://algebrica.org/riemann-integrability-criteria/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/riemann-integrability-criteria.md?plain=1)

3.3k

[Improper Integrals](https://algebrica.org/improper-integrals/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/improper-integrals.md?plain=1)

914

[Arc Length of a Curve](https://algebrica.org/arc-length-of-a-curve/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/integrals/arc-length-of-a-curve.md?plain=1)
