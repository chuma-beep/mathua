# GMRES

GMRES minimises the residual over a Krylov subspace. It is the standard Krylov method for nonsymmetric systems.

## Minimisation

At step k GMRES chooses the Krylov vector whose residual has the smallest norm. The least-squares problem is solved in the Arnoldi basis.

## Residual

The residual norm is available from the least-squares problem without forming the iterate. This makes the stopping test cheap.

## Restarting

Full GMRES stores k basis vectors. A restarted version keeps only m vectors and loses the optimality of the longer recurrence.

To check your understanding, minimise the residual, read the residual norm and bound the basis size.

GMRES trades storage for optimality, and restarting trades optimality back for storage.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
