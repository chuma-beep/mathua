# PAC learning

Probably approximately correct learning asks for a hypothesis that is close to the best in the class with high probability. The sample complexity states how much data that guarantee needs.

## The guarantee

The learner returns a hypothesis with error at most \(\epsilon\) with probability at least \(1-\delta\). Both parameters are chosen by the user.

## The sample bound

For a finite class the bound grows like \(\log|H|/\epsilon\) plus a confidence term. The dependence on \(\epsilon\) is inverse, not inverse square.

## Realizability

When the class contains a perfect hypothesis the bound improves. The agnostic case pays for the possibility that no hypothesis fits.

To check your understanding, name epsilon, count the class and assume a perfect fit.

PAC learning converts a desired error and confidence into a sample size.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
