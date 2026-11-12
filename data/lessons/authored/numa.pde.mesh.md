# Mesh refinement

Mesh refinement reduces the cell size where accuracy is needed. For a smooth solution the error falls at a rate set by the polynomial degree and the norm used.

## h and p

h-refinement subdivides cells and keeps the degree fixed. p-refinement raises the degree on fixed cells; both reduce error, with different sparsity consequences.

## Error rate

For piecewise linear elements in the energy norm the error is \(O(h)\). In the \(L^{2}\) norm it is \(O(h^{2})\).

## Adaptivity

An error estimator marks cells with the largest local error. Refining only those cells reaches a target accuracy with far fewer unknowns than uniform refinement.

To check your understanding, subdivide the cells, halve the cell size and estimate the local error.

Refinement is a choice between smaller cells, higher degree and better placement.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
