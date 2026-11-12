# The stiffness matrix

The stiffness matrix holds the integrals of derivatives of the basis functions. It is symmetric positive semidefinite, and with a mass matrix it forms the discrete operator.

## Entries

Entry \((i,j)\) is \(\int\nabla\phi_i\cdot\nabla\phi_j\). Neighbouring basis functions overlap, so the entry is nonzero only for adjacent nodes.

## Assembly cost

Each element contributes a small dense block. The total assembly cost is proportional to the number of elements times the square of the local basis size.

## Conditioning

The condition number grows like \(h^{-2}\) as the mesh is refined. Preconditioning or a multilevel solver is needed for large meshes.

To check your understanding, integrate the gradients, count element blocks and refine the mesh.

The stiffness matrix is the derivative form of the same mesh that carries the mass.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
