# Weak derivatives

A weak derivative moves the derivative onto a test function.

## The integration by parts definition

A function \(v\) is the weak derivative of \(u\) when the integral of \(u\) against a test derivative equals minus the integral of \(v\) against the test function.

The test functions are smooth and vanish near the boundary, so no boundary term appears. Every classical derivative is also a weak derivative. The identity must hold for every test function.

For \(u(x)=x^{2}\) and the test \(\phi=x(1-x)\) both sides equal \(-1/6\).

## An example with a corner

Take \(u(x)=|x|\) on the interval from -1 to 1, which has no derivative at 0.

The weak derivative exists and equals -1 on the left half and 1 on the right half, because the corner cancels the boundary terms. Its square integrates to 2.

The norm in \(H^{1}\) is the square root of \(2/3+2\), about 1.633.

## Sobolev regularity

Each weak derivative in \(L^{2}\) raises the Sobolev order by one.

The function \(|x|\) is in \(H^{1}\) but not in \(H^{2}\), since its derivative jumps. A function with a continuous derivative, such as \(x^{2}\), is in \(H^{2}\) with second weak derivative 2.

The jump is what stops the ascent at a single point.
