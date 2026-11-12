# Duality over cones

Cone programs have a dual built from the dual cone. The pairing reproduces weak duality and complementary slackness with a cone in place of the nonnegative orthant.

## The dual cone

The dual cone holds every vector whose inner product with every member of the original cone is nonnegative. For the nonnegative orthant it is the orthant itself.

## Weak duality for cones

Any dual feasible point bounds the primal value from below. The proof is one inner-product inequality and does not depend on the cone being polyhedral.

## Complementarity

At an optimum the product of each cone variable and its dual slack is zero. The statement is the same as in linear programming, with the cone product replacing the scalar one.

To check your understanding, take the inner product with each generator, evaluate the dual objective and multiply the variable by its slack.

The dual cone is the only new object. Once it is known, the algebra of linear programming carries over unchanged.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
