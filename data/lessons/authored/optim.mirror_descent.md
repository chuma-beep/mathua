# Mirror descent

Projected gradient measures distance with the Euclidean norm. Mirror descent measures it with a divergence chosen for the feasible set, so a probability simplex or a positive cone is handled in its own geometry.

## The mirror map

The method maps the iterate to a dual space, takes a gradient step there, and maps back. The map is chosen so that the divergence between successive iterates is cheap and so that the feasible set becomes simple in the dual space.

## The regret bound

For a simplex the entropy divergence gives a regret that grows with the logarithm of the dimension rather than the dimension itself. That gap is the reason mirror descent appears in online learning and in experts problems.

## The trade

The right divergence can make a problem tractable and a wrong one makes every step an implicit subproblem. The divergence is part of the model, not an implementation detail.

To check your understanding, map to the dual space, use entropy on a simplex and match divergence to the set.

Mirror descent is gradient descent in the geometry that fits the constraint set, and the choice of divergence is where the method earns its advantage.
