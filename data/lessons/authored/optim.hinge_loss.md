# The hinge loss

The hinge loss penalises a margin that is too small and is zero once the margin is wide enough. It is the loss behind the support vector machine.

## The margin

The margin is the signed distance of a point from the decision boundary. A correct classification with a comfortable margin has a large positive margin.

## The loss value

The loss is the positive part of one minus the margin. A margin of at least one costs nothing and a margin below one costs the shortfall.

## The subgradient

The slope is negative one while the margin is below one and zero above it. At the kink the subdifferential is the interval between them.

To check your understanding, take the signed distance, subtract the margin from one and compare the margin with one.

The zero region is what makes the solution sparse in the support vectors. Points already classified well stop pulling on the boundary.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
