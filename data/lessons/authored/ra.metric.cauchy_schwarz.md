# The Cauchy-Schwarz inequality

Cauchy-Schwarz bounds an inner product by the product of two norms. It is the inequality behind the triangle inequality and the notion of angle in a general inner product space.

## The inner product

An inner product is linear in each argument, symmetric and positive definite. In \(\mathbb{R}^3\) take \(x = (1,2,2)\) and \(y = (2,0,1)\). Their inner product is 4, and the norms are 3 and \(\sqrt{5} \approx 2.236\).

## The inequality

The inequality states \(|\langle x,y\rangle| \le \|x\|\,\|y\|\). Here \(4 \le 3\sqrt{5} \approx 6.708\), so the bound holds with room to spare. The proof expands \(\|x - ty\|^2\) as a quadratic in \(t\) and uses its nonnegative discriminant.

## Equality and angle

Equality holds exactly when the vectors are dependent. The cosine of the angle is \(\langle x,y\rangle/(\|x\|\|y\|) \approx 0.596\), so the angle is about 53.4 degrees. For \(y = 2x\) the cosine is 1 and the inequality is an equality.

Cauchy-Schwarz converts a product of norms into a bound on the inner product, and that bound is what makes angle and distance agree in every dimension.
