# A gradient descent step

Gradient descent moves a point against the gradient by an amount called the step size. The rule is one line, but everything that matters is in the size and in what the next iterate does to the objective. One step is easy to compute and easy to check, which makes it a good unit for learning the mechanics.

## Taking one step

Compute the gradient at the current point, multiply it by the step size and subtract. The result is the next point. Keeping the intermediate values visible matters, because an arithmetic slip in the gradient propagates into a step that still looks plausible. The new objective value should be lower, and comparing it with the old one is the cheapest sanity check.

## Predicting the next iterate

Applying the rule twice gives the second iterate, and the pattern can be predicted for a quadratic without recomputing the whole sequence. For a one-dimensional quadratic the error is multiplied by a fixed factor each step, so the iterate is a geometric sequence. Predicting it tests both the gradient and the step size at once.

## Comparing two directions

Two candidate directions can be compared by their inner product with the gradient, or by the objective after a trial step of the same length. The direction with the more negative inner product makes more progress per unit of movement, to first order. This comparison is the intuition behind choosing a preconditioner or a better descent direction.

A single step exposes the whole loop in miniature: gradient, scaling, update, check. Methods that look complicated are usually this step with a different choice of direction or size.
