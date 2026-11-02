# Liouville's Theorem

## Bounded entire functions

Liouville's theorem states that a bounded entire function is constant. An entire function is analytic on the whole complex plane, and bounded means that its modulus is at most some fixed number everywhere. The theorem is surprising because real differentiable functions can be bounded and nonconstant, such as sine and cosine.

The proof uses Cauchy's estimate for the first derivative. On a circle of radius r around any point, the derivative is bounded by the maximum modulus on that circle divided by r. Letting r grow makes the bound tend to zero, so the derivative vanishes everywhere. A function with zero derivative on a connected domain is constant.

An equivalent statement is that a nonconstant entire function is unbounded. This form is often easier to apply: if a problem gives a nonconstant entire function and claims a global bound, the claim is false. Polynomials of positive degree, the exponential and sine are all unbounded on the complex plane.

## The proof

The estimate is the heart of the proof. Analyticity gives a power series at each point, and Cauchy's integral formula gives the derivative as an integral over a circle. The modulus of the integral is at most the length of the circle times the maximum of the integrand, which produces the factor of r in the denominator.

Letting the radius tend to infinity is legitimate because the function is entire, so the circle can be chosen as large as we like. The bound on the derivative tends to zero, and since the point was arbitrary, the derivative is identically zero. Connectedness of the plane then forces the function to be constant.

The theorem has a converse in spirit: constants are exactly the bounded entire functions. That classification is the first step in many arguments. For instance, if an entire function misses a value, one can compose with a Möbius map and apply Liouville to prove a stronger result, which leads toward Picard's theorem.

## Polynomial growth

A useful generalization says that if an entire function grows at most like a constant times the n-th power of the modulus, then it is a polynomial of degree at most n. The proof applies Cauchy's estimates to derivatives of order n plus one and lets the radius tend to infinity, forcing those derivatives to vanish.

This polynomial-growth version explains why entire functions with slow growth are so rigid. A bounded entire function is the case n equal to zero. Linear growth forces an affine function, quadratic growth forces a quadratic, and so on. The bound only needs to hold for large modulus, since a finite initial region is compact and can be absorbed into the constant.

When solving problems, identify the growth exponent first. If the hypothesis gives a bound of the form C times modulus to the n, the conclusion is a polynomial of degree at most n. If the function is also known to have degree exactly n, then it must be a polynomial of that degree, and its leading coefficient can often be recovered from a limit.

