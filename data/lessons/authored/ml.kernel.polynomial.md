# Polynomial kernels

A polynomial kernel raises an affine combination of the inner product to a power. It corresponds to a finite feature space of monomials up to a fixed degree.

## The formula

The kernel is \((x\cdot y+c)^d\) with degree \(d\) and offset \(c\). The offset adds monomials of lower degree.

## The feature count

Homogeneous monomials of degree \(d\) in \(n\) variables number \(\binom{n+d-1}{d}\). The full polynomial space adds the lower degrees.

## The degree

A higher degree widens the hypothesis class. The cost grows with the number of monomials, so degree is a capacity control.

To check your understanding, evaluate the kernel, count homogeneous monomials and raise the degree.

A polynomial kernel makes the inner product compute a monomial expansion implicitly.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
