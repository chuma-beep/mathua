# The signless Laplacian

The signless Laplacian adds the adjacency to the degree matrix, and its smallest eigenvalue detects bipartiteness.

## The definition

The signless Laplacian is \(Q=D+A\). It is positive semidefinite, and its row sums are the degrees, so the all-ones vector is an eigenvector with eigenvalue zero only when the graph has no edges.

## The spectrum

The eigenvalues are nonnegative and sum to \(2m\), the total degree. The smallest eigenvalue is zero exactly when the graph has a bipartite component.

## Comparison

The Laplacian is \(D-A\) and the signless Laplacian is \(D+A\). For a bipartite graph the two spectra coincide after a diagonal sign change, which is why bipartiteness is visible in the smallest signless eigenvalue.

To check your understanding, add the adjacency, sum the eigenvalues and write both Laplacians.

Changing the sign of the adjacency changes which structures the smallest eigenvalue detects.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
