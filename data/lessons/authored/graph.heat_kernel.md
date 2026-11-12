# The heat kernel on a graph

The heat kernel diffuses a function over the vertices, and its long-time limit is the degree distribution.

## The equation

The heat equation on a graph is \(\frac{d}{dt}f=-Lf\), with solution \(f(t)=e^{-Lt}f(0)\). The matrix exponential is defined through the Laplacian eigenvalues.

## Spectral form

Writing \(L=\sum_i \lambda_i u_i u_i^{T}\) gives \(e^{-Lt}=\sum_i e^{-\lambda_i t}u_i u_i^{T}\). The zero eigenvalue survives forever, and every other mode decays at its own rate.

## The limit

As \(t\to\infty\) the solution approaches the projection onto the constant vector, which for a connected graph is the degree distribution divided by \(2m\). The total heat is conserved throughout.

To check your understanding, write the heat equation, expand in eigenvectors and send time to infinity.

Diffusion smooths a function until only the stationary mode remains, and the spectrum sets the rate.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
