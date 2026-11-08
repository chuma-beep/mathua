# Brent's method

Brent's method finds a root inside a bracket by combining three ideas: bisection for safety, the secant method for speed, and inverse quadratic interpolation for a better guess. It keeps the guarantee of bisection while often converging much faster.

## A bracket is a promise

A bracket is an interval whose endpoints give function values of opposite sign. By the intermediate value theorem, a continuous function has at least one root inside. Every step of Brent's method keeps a valid bracket, so the root cannot be lost.

Bisection halves the bracket each time. After \(n\) steps the width is the original width divided by \(2^n\). That is slow but certain, and it is the fallback whenever the faster step looks unreliable.

## Interpolation when it is safe

The secant method draws a line through two points and uses its intercept as the next guess. Inverse quadratic interpolation fits a quadratic through three points and solves it for \(x\) as a function of \(y\). When the points are well placed, these steps jump close to the root.

The danger is an extrapolated point outside the bracket or too close to an endpoint. Brent's method checks the proposed point and rejects it if it would break the bracket or fail to shrink the interval enough. Safety is checked before speed is used.

## The bisection fallback

When interpolation is rejected, the method takes a bisection step. That step is slower, but it guarantees progress. The alternation between fast interpolation and safe bisection is why Brent's method is a default choice in numerical libraries.

The worst case is linear convergence, order one, because the fallback only halves the interval. The typical case is superlinear, often close to the secant method. The guarantee is what matters: the bracket shrinks every step and the root is never abandoned.
