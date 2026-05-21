> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Non-homogeneous_differential_equation) — CC BY-SA 4.0

# Non-homogeneous second-order ODEs

A **non-homogeneous** ordinary differential equation is one that includes a non-zero forcing function. For a second-order linear ODE:

\[
y'' + p(x)y' + q(x)y = f(x)
\]

The general solution is the sum of the complementary solution (solving the homogeneous case \(f(x) = 0\)) and a particular solution:

\[
y(x) = y_c(x) + y_p(x)
\]

### Methods for finding particular solutions

**Method of undetermined coefficients:** When \(f(x)\) is a polynomial, exponential, sine, cosine, or a product thereof, guess a solution of similar form with unknown coefficients and solve.

**Variation of parameters:** For any continuous \(f(x)\), the particular solution is:

\[
y_p(x) = -y_1 \int \frac{y_2 f}{W} \, dx + y_2 \int \frac{y_1 f}{W} \, dx
\]

where \(y_1, y_2\) are linearly independent solutions of the homogeneous equation and \(W\) is their Wronskian.
