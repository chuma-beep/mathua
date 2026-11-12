# Reproducing kernel Hilbert spaces

A reproducing kernel Hilbert space is a Hilbert space in which evaluation at a point is a bounded linear functional. The kernel is the representer of that functional.

## The reproducing property

Evaluation satisfies \(f(x)=\langle f,k(\cdot,x)\rangle\). The kernel reproduces the value of every function in the space.

## The representer theorem

A regularised empirical risk minimiser lies in the span of the kernel sections at the data. An \(n\)-point problem therefore has \(n\) coefficients.

## The norm

The squared norm of a combination equals the coefficient vector applied to the Gram matrix. That identity turns function penalties into finite quadratic forms.

To check your understanding, state the identity, count the coefficients and write the quadratic form.

The RKHS makes infinite-dimensional function classes computable through the kernel.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
