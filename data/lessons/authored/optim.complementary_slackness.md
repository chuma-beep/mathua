# Complementary slackness

Complementary slackness pairs each inequality constraint with its multiplier and requires at least one of them to be zero. A constraint with positive slack has a zero price, and a positive price means the constraint is tight. The condition is necessary and sufficient for optimality together with primal and dual feasibility.

## Testing a pair

For each inequality, multiply the slack by the multiplier and check that the product is zero. A positive product means the pair is not complementary. The test is arithmetic once the slack and the multiplier are known, and it is the quickest way to confirm a claimed solution.

## Finding a complementary solution

Given an optimal basis, the multipliers come from the dual equations and the slacks from the primal values. The complementary pattern says which constraints can be tight. Solving the resulting smaller system gives the candidate, which is then checked for feasibility and optimality.

## Using a positive multiplier

A positive multiplier marks a binding constraint and measures its marginal value. A zero multiplier on a binding constraint is allowed, but a positive multiplier on a slack constraint is not. The distinction matters when a price is quoted for relaxing a constraint.

Complementary slackness is the bridge between a solution and its prices. It explains which constraints matter at the optimum and which can be ignored locally.
