# Online gradient descent

Online gradient descent makes one step per arriving example and is judged by regret against the best fixed point.

## Regret

After many rounds the regret is the total loss minus the loss of the best fixed decision. A sublinear regret means the average gap vanishes.

## The bound

With a diminishing step proportional to one over the square root of the round count the regret grows as that square root for convex losses. The bound needs no distributional assumption.

## Static versus dynamic

A static comparator is one fixed point. A dynamic comparator changes with time and admits a larger regret.

To check your understanding, compute a regret, evaluate a root step and pick a comparator.

Online learning changes the question from convergence to regret, and the step schedule answers it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
