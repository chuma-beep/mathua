# Resolvent kernels

The resolvent kernel sums the Neumann series of an integral operator. It turns a Fredholm equation into an explicit formula whenever the series converges.

## The geometric case

For the kernel \(\lambda K(x,t)\) with a constant \(\lambda=1/n\), the resolvent factor is \(n/(n-1)\). The series converges because \(|\lambda|\) is below one.

## The series

Truncating the Neumann series after \(N\) terms gives the \(N\)-th approximation. Each term adds one more application of the integral operator.

## The determinant

The Fredholm determinant is \(1-\lambda\) in the constant case, so its numerator is \(n-1\) when \(\lambda=1/n\). A zero determinant marks a value where the resolvent fails to exist.

To check your understanding, set the constant kernel, choose a truncation and evaluate one minus lambda.

The resolvent is the closed form of the Neumann series.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
