> Content sourced from [Algebrica](https://algebrica.org/fourier-series/) — CC BY-NC 4.0

## Definition

A Fourier series represents a periodic function as an infinite sum of [sine](<../sine-function/>) and [cosine functions](<../cosine-function>). More precisely, it shows that periodic behavior can be decomposed into elementary harmonic oscillations. This result expresses a structural property of periodic [functions](<../functions>): oscillatory components form a natural coordinate system for describing repetition.


Let \\( f : \mathbb{R} \to \mathbb{R} \\) be a function that is periodic with period \\( 2\pi \\), meaning:

\\[f(x + 2\pi) = f(x) \, \forall x \in \mathbb{R} \\]

Assume that \\( f \\) is [integrable](<../definite-integrals/>) on the [interval](<../intervals/>) \\( [-\pi,\pi] \\). The Fourier series of \\( f \\) is the formal trigonometric expansion:

\\[f(x) \sim \frac{a_0}{2} \+ \sum_{n=1}^{\infty} a_n \cos(nx) + b_n \sin(nx) \\]

  * The symbol \\( \sim \\) emphasizes that we are not yet asserting equality (we are defining a trigonometric series associated with \\( f \\)).
  * The question of whether the series converges to \\( f \\) will be addressed later.
  * Each term \\( \cos(nx) \\) and \\( \sin(nx) \\) represents an oscillation of frequency \\( n \\).
  * The expansion therefore decomposes \\( f \\) into its harmonic components.


## Fourier Coefficients

The coefficients \\( a_n \\) and \\( b_n \\) are defined by the following integrals:

\\[\begin{align} a_0 &= \frac{1}{\pi} \int_{-\pi}^{\pi} f(x)\,dx \\\ a_n &= \frac{1}{\pi} \int_{-\pi}^{\pi} f(x)\cos(nx)\,dx \\\ b_n &= \frac{1}{\pi} \int_{-\pi}^{\pi} f(x)\sin(nx)\,dx \quad n \ge 1 \end{align} \\]

These expressions are not introduced by convention, and they are not chosen just because they work. They follow from a structural fact about [sines and cosines](<../sine-and-cosine/>): over a full period they are orthogonal to one another. On the interval \\( [-\pi,\pi] \\), trigonometric waves with different frequencies remain independent under integration, which is exactly what lets us isolate one harmonic at a time and read off the corresponding coefficient.

\\[\begin{align} \int_{-\pi}^{\pi} \cos(nx)\cos(mx)\,dx &= \begin{cases} \pi & n=m\neq 0 \\\ 0 & n\ne m \end{cases} \\\\[6pt] \int_{-\pi}^{\pi} \sin(nx)\sin(mx)\,dx &= \begin{cases} \pi & n=m \\\ 0 & n\ne m \end{cases} \\\\[6pt] \int_{-\pi}^{\pi} \sin(nx)\cos(mx)\,dx &= 0 \end{align} \\]

These relations imply that the trigonometric system behaves like an orthogonal basis under the inner product:

\\[\langle f, g \rangle = \int_{-\pi}^{\pi} f(x)g(x)\,dx \\]

Each coefficient measures how much of a specific harmonic direction is present in the function. In this sense, the Fourier expansion is a projection process in an infinite-dimensional space.

## Example 1

This example illustrates how even a simple linear function acquires a rich harmonic structure when periodically extended. Consider the function \\(f(x) = x\\), defined on \\( (-\pi,\pi) \\) and extended periodically with period \\( 2\pi \\). This function is odd. Therefore:

\\[a_0 = 0 \quad a_n = 0 \\]

We compute the [sine](<../sine-and-cosine/>) coefficients:

\\[b_n = \frac{1}{\pi} \int_{-\pi}^{\pi} x\sin(nx)\,dx \\]

Using [integration by parts](<../integration-by-parts/>), we obtain:

\\[b_n = \frac{2(-1)^{n+1}}{n} \\]

Hence the Fourier series is:

\\[x \sim 2 \sum_{n=1}^{\infty} \frac{(-1)^{n+1}}{n} \sin(nx) \\]

###### The coefficients decay like \\( \frac{1}{n} \\). The slower decay reflects the fact that \\( f \\) is continuous but not differentiable at the endpoints of the period. The periodic extension introduces jump discontinuities at multiples of \\( \pi \\), which influences convergence behavior.

## Convergence of Fourier series

The definition of the Fourier series does not automatically guarantee convergence to the original function. A classical result states that if \\( f \\) satisfies the following [Dirichlet conditions](<../dirichlet-function/>):

  * \\( f \\) is piecewise continuous
  * \\( f \\) has finitely many local extrema in \\( [-\pi,\pi] \\)
  * \\( f \\) has finitely many jump discontinuities


then the Fourier series converges at every point \\( x. \\) More precisely consider the \\(N\\)-th partial sum:

\\[S_N(x) = \frac{a_0}{2} \+ \sum_{n=1}^{N} a_n\cos(nx)+b_n\sin(nx) \\]

\\[\lim_{N\to\infty} S_N(x) = \frac{f(x^+)+f(x^-)}{2} \\]

At points where \\( f \\) is continuous, the series converges to \\( f(x) \\). At jump discontinuities, it converges to the midpoint of the left and right [limits](<../limits>). This behavior reveals a fundamental property of Fourier approximation: it respects average local behavior rather than pointwise values at discontinuities.

  * If \\( f \\) is continuously differentiable, coefficients decay faster.
  * If \\( f \\) has discontinuities, decay is slower.
  * The smoother the function, the more rapidly the harmonic amplitudes decrease.


## Selected references

  * **E. M. Stein, R. Shakarchi**. [Fourier Analysis: An Introduction](https://kryakin.site/am2/Stein-Shakarchi-1-Fourier_Analysis.pdf)
  * **L. Grafakos**. [Classical Fourier Analysis](https://www.math.stonybrook.edu/~bishop/classes/math638.F20/Grafakos_Classical_Fourier_Analysis.pdf)
  * **G. P. Tolstov**. [Fourier Series](https://archive.org/embed/tolstov-fourier-series-1962)
  * **G. B. Folland**. [Fourier Analysis and Its Applications](https://www-elec.inaoep.mx/~rogerio/Tres/FourierAnalysisUno.pdf)
  * **A. Zygmund**. [Trigonometric Series](https://archive.org/details/trigonometricseries)
  * **Y. Katznelson**. [An Introduction to Harmonic Analysis](https://archive.org/details/introductiontoha0000katz)
  * **Stanford University**. [The Fourier Transform and Its Applications](https://see.stanford.edu/materials/lsoftaee261/book-fall-07.pdf)
  * **Oxford University Press**. [Fourier Series and Fourier Transforms](https://academic.oup.com/book/54863/chapter/422699978)
  * **R. Herman**. [An Introduction to Fourier and Complex Analysis](https://people.uncw.edu/hermanr/mat367/fcabook/FCA)


Series

Series are infinite sums of sequence terms, typically indexed by natural numbers.

1.3k

[Series](https://algebrica.org/series/)

2.4k

[Cauchy’s Convergence Criterion for Series](https://algebrica.org/cauchy-convergence-criterion-series/)

1.4k

[Series with Positive Terms](https://algebrica.org/series-with-positive-terms/)

1.7k

[Harmonic Series](https://algebrica.org/harmonic-series/)

1.6k

[Geometric Series](https://algebrica.org/geometric-series/)

1.4k

[Integral Test for Series Convergence](https://algebrica.org/integral-test-for-series-convergence/)

1.1k

[Root Test for Series Convergence](https://algebrica.org/root-test-for-series-convergence/)

1.9k

[Leibniz’s Criterion](https://algebrica.org/leibniz-criterion/)

641

[Function Series](https://algebrica.org/function-series/)

636

[Power Series](https://algebrica.org/power-series/)

1.7k

[Taylor Series](https://algebrica.org/taylor-series/)

4 comments
