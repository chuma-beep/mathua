# Homogeneous linear recurrences

A homogeneous linear recurrence expresses each term as a fixed linear combination of earlier terms with no forcing term. Its solutions form a vector space spanned by characteristic-root sequences.

## The form

\(a_n=c_1a_{n-1}+\cdots+c_da_{n-d}\) with constant coefficients. The order is \(d\), and \(d\) initial values determine the sequence.

## The characteristic equation

Substituting \(a_n=r^n\) gives \(r^d=c_1r^{d-1}+\cdots+c_d\). Its roots generate the solution space.

## Superposition

Linear combinations of solutions are solutions. That is why a basis of \(d\) independent sequences describes every solution once the initial values are fixed.

To check your understanding, identify the order, substitute r^n and combine two solutions.

Homogeneous recurrences are linear algebra in sequence form. The characteristic roots replace the matrix in the finite-dimensional case.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
