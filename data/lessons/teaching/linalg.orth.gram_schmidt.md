# The Gram-Schmidt Process

Gram-Schmidt turns a basis into an orthogonal basis without changing the subspace it spans. Each new vector has its projection onto the earlier directions removed, so what remains is perpendicular to all of them.

## Orthogonalising Two Vectors

Start with independent vectors $v_1$ and $v_2$. Keep $u_1 = v_1$. The projection of $v_2$ onto $u_1$ is

$$\operatorname{proj}_{u_1}(v_2) = \frac{v_2 \cdot u_1}{u_1 \cdot u_1} u_1.$$

Subtract it to get $u_2 = v_2 - \operatorname{proj}_{u_1}(v_2)$. Then $u_2 \cdot u_1 = 0$, and $u_1, u_2$ span the same plane as $v_1, v_2$.

## The Gram-Schmidt Formula

For a list $v_1, \ldots, v_k$, set $u_1 = v_1$ and for each later index

$$u_j = v_j - \sum_{i=1}^{j-1} \frac{v_j \cdot u_i}{u_i \cdot u_i} u_i.$$

Each denominator is nonzero because the earlier vectors are independent. The resulting list is orthogonal and has the same span at every stage.

## From Orthogonal to Orthonormal Bases

Divide each $u_j$ by its norm to obtain $q_j = u_j / \|u_j\|$. The vectors $q_j$ are unit length and mutually perpendicular, so they form an orthonormal basis. Collecting them as columns gives a matrix $Q$ with $Q^T Q = I$, which is the first factor in a QR decomposition.
