## Introduction

> The geometric construction of the cosecant from the unit circle is developed in secant and cosecant. Here, the cosecant is treated as a real function of a real variable.

The cosecant function $f(x) = \csc(x)$ assigns to each angle $x,$ measured in radians, the reciprocal of its sine value, defined wherever $\sin(x) \neq 0.$ Its graph is a periodic curve with period $2\pi$ and has vertical asymptotes where the sine of $x$ vanishes, at $x = k\pi$ with $k \in \mathbb{Z}.$ The domain is the set of all real numbers except these points, and the range is $(-\infty, -1] \cup [1, +\infty).$

![IMG. 1](/diagrams/algebrica/secant-and-cosecant-4.svg)


The cosecant is the reciprocal of the sine, so it stays bounded where the sine is near $\pm 1$ and grows without bound as the sine approaches zero:

$$\csc(x) = \frac{1}{\sin(x)}$$

Because the sine never exceeds 1 in absolute value, its reciprocal never falls between $-1$ and 1, and each branch reaches a single extremum of 1 or $-1$ where the sine equals $\pm 1$ before diverging toward the neighbouring asymptotes.

## Properties

The following properties of the cosecant function follow from its definition as the reciprocal of the sine.

+ Domain: $\{\ x \in \mathbb{R} \mid x \neq k\pi \ \text{ for all } k \in \mathbb{Z} \ \}$
+ Range: $y \in (-\infty, -1] \cup [1, +\infty)$
+ Periodicity: periodic in $x$ with period $2\pi$
+ Parity: odd, with $\csc(-x) = -\csc(x)$
+ The graph has vertical asymptotes at $x = k\pi$ with $k \in \mathbb{Z}$

## Relation with the cotangent

The cosecant and the cotangent are tied by a Pythagorean identity. Dividing $\sin^2(x) + \cos^2(x) = 1$ by $\sin^2(x)$ and using the reciprocal definitions gives:

$$\csc^2(x) = 1 + \cot^2(x)$$

The two functions share the same vertical asymptotes, and the cosecant grows together with the cotangent.

## Limits, derivatives, and integrals of the cosecant function

Near $x = \pi/2$ the sine reaches its maximum, so the cosecant takes its least positive value:

$$\lim_{x \to \frac{\pi}{2}} \csc(x) = 1$$

The behaviour near the vertical asymptote at the origin is described by one-sided limits. As $x$ approaches 0 from the right the sine is positive and tends to zero, so the function grows without bound,

$$\lim_{x \to 0^+} \csc(x) = +\infty$$

while from the left the sine is negative and the values diverge to negative infinity,

$$\lim_{x \to 0^-} \csc(x) = -\infty$$

The function is continuous and differentiable on its domain. Its derivative is:

$$\frac{d}{dx}\csc(x) = -\csc(x)\cot(x)$$

The indefinite integral is:

$$\int \csc(x) \ dx = -\ln\left|\csc(x) + \cot(x)\right| + c$$

> A broader treatment of trigonometric integrals, with the transformation and substitution techniques for the more complex cases, is given in trigonometric function integrals.

The cosecant function can also be written using imaginary numbers. With $e^{ix}$ the exponential function of base $e$ and $i$ the imaginary unit, Euler's formula gives:

$$\csc(x) = \frac{2i}{e^{ix} - e^{-ix}}$$
