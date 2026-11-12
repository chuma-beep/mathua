# Gauss-Seidel splitting

Gauss-Seidel updates coordinates one at a time and uses each new value immediately. The triangular structure replaces the diagonal splitting of Jacobi.

## Splitting

Write the matrix as a lower triangular part plus an upper remainder. The iteration solves a triangular system at every sweep.

## Ordering

The result depends on the coordinate order, because new values enter later updates. A forward sweep and a backward sweep give different iterates.

## Convergence

Gauss-Seidel converges for symmetric positive definite matrices. Its asymptotic rate is often about twice that of Jacobi for the same problem.

To check your understanding, take the lower triangle, choose a sweep direction and use positive definiteness.

Using new values immediately is free in storage and usually halves the iteration count.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
