# Conditioning of errors

The condition number amplifies data perturbations into solution perturbations. It is the factor that separates a stable algorithm from an accurate answer.

## Amplification

A relative data perturbation of size delta can produce a relative solution perturbation of size kappa times delta. The condition number is the worst-case ratio.

## Digits

If the condition number has d decimal digits, roughly d digits of the answer can be lost. This is the practical reading of the bound.

## Geometry

The condition number measures how close the matrix is to singular. A large value means some directions are nearly collapsed.

To check your understanding, perturb the data, count digits of kappa and measure near-singularity.

Conditioning is a property of the problem, and no algorithm can remove it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
