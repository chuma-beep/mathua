# The Dirichlet distribution

The Dirichlet distribution lives on the simplex of positive weights summing to one. It is the conjugate prior for categorical probabilities.

## Density

With parameters \(\alpha_i>0\), the density is proportional to the product of the coordinates raised to \(\alpha_i-1\).

## Moments

The mean of coordinate \(i\) is \(\alpha_i/\alpha_0\) and its variance involves the total concentration \(\alpha_0\).

## Aggregation

Merging coordinates adds their parameters, which keeps the family Dirichlet.

To check your understanding, normalise the weights, divide by the total and merge two coordinates.

The Dirichlet is the natural distribution on proportions.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
