# Separation of variables

Separation of variables turns a partial differential equation into ordinary equations by assuming a product solution.

## Splitting the variables

Write \(u(t,x)=T(t)X(x)\) for \(u_{t}=u_{xx}\), which gives \(T'X=TX''\) on substitution.

Dividing by \(TX\) puts each variable on its own side, so both sides equal a constant, written \(-\lambda\). The result is \(T'=-\lambda T\) and \(X''=-\lambda X\). The sign of the constant decides whether the factors oscillate or decay.

For \(\lambda=4\) the time factor is \(e^{-4t}\), which falls to 0.5 at \(t=(\ln 2)/4\), about 0.1733.

## Boundary conditions

The end values of \(u\) become conditions on \(X\), so on the interval from 0 to \(\pi\) we need \(X(0)=X(\pi)=0\).

Only special values of \(\lambda\) allow a nonzero solution. They are the squares \(n^{2}\). The matching functions are \(\sin nx\).

The first three are 1, 4 and 9, with functions \(\sin x\), \(\sin 2x\) and \(\sin 3x\).

## The separated solutions

Each pair gives a solution \(e^{-n^{2}t}\sin nx\). Any sum of them solves the equation.

Take the term with \(n=1\) and coefficient 3. At \(t=\ln 2\) the exponential is 0.5. At \(x=\pi/2\) the sine is 1, so the value is 1.5.

Separation gives the building blocks. A Fourier series chooses the mixture. The first term decays slowest, so it dominates after a while.
