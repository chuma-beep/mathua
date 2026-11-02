# Real Integrals by Residues

## Contour and residues

Residues evaluate many definite integrals of real functions by closing the integration path with a complex contour. The real integral becomes one piece of a closed contour, and the residue theorem replaces the contour integral by a sum over enclosed singularities. The remaining pieces must tend to zero or be computable.

The standard rational case integrates a function of x over the whole real line. Close the contour with a large semicircle in the upper half-plane. If the degree of the denominator exceeds the degree of the numerator by at least two, the semicircle contribution vanishes as the radius grows, so the real integral equals two pi i times the sum of residues in the upper half-plane.

Trigonometric integrals over a full period use the unit circle instead. Substitute z on the unit circle, write cosine and sine in terms of z, and convert the real integral into a contour integral around the circle. The residues inside the unit disc then determine the answer, and the orientation is counterclockwise.

## Rational integrands

For an integrand with real coefficients, the poles occur in conjugate pairs. Only the upper half-plane poles contribute when the contour is closed above, and the lower half-plane poles would contribute if the contour were closed below with the opposite orientation. Choosing the half-plane consistently avoids sign errors.

A pole on the real axis is not enclosed by either semicircle in the usual limit. Such an integral is interpreted as a principal value, and the contour is indented around the pole. The indentation contributes half a residue, a detail that is easy to miss. Problems that avoid real poles are safer for a first pass.

Before computing residues, check that the denominator has no real zeros and that the decay condition holds. If the numerator and denominator have the same degree, the large semicircle does not vanish, and the integral may diverge or require a different method. Partial fractions can sometimes reduce the problem to simpler pieces.

## Choosing the contour

The contour should enclose the singularities you want and avoid the singularities you do not want. A semicircle is natural for the real line, a circle for a periodic integrand, and a rectangle for integrals involving exponentials or hyperbolic functions. The choice is guided by the decay of the integrand in different directions.

Exponential factors require Jordan's lemma. If the integrand contains e to the i a z with a positive, the upper half-plane is the correct side because the exponential decays there. For a negative frequency, close below. The lemma controls the semicircle even when the rational factor decays only like one over the modulus.

After closing the contour, write the real integral as the contour integral minus the other pieces. Estimate each remaining piece and show that it tends to zero. The residue sum then gives the answer. Checking the estimate is not optional, because a contour that does not vanish leads to a wrong value even if the residues are computed correctly.

