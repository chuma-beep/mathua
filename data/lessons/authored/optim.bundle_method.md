# Bundle methods

A bundle method keeps a collection of subgradients and builds a cutting-plane model. A trust region keeps the model honest.

## The cutting model

Each evaluated point contributes a linear piece below the function. The maximum of those pieces is a lower bound that improves as the bundle grows.

## The trust region

The model is trusted only near the current point. A radius limits the step, and the ratio of actual to predicted decrease updates it.

## The master

The master minimises the model over the trust region. Its solution proposes the next evaluation point.

To check your understanding, evaluate each linear piece, count the stored subgradients and minimise the model with the radius.

The bundle method is a cutting-plane method with a stability term. The term prevents the wild steps that a pure cutting plane can take.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
