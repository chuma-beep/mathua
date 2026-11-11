# Cutting planes

A cutting plane is an inequality that every integer feasible point satisfies and that the current fractional relaxation solution violates. Adding cuts tightens the relaxation without removing integer points. Repeating the process moves the relaxation optimum toward the integer hull.

## Testing a cut

Evaluate the cut at the current relaxation solution. If the left side violates the inequality, the cut separates that point. Then check that no integer feasible point violates it, which is what makes the cut valid. A cut that separates but is invalid can remove the true optimum.

## Adding a valid inequality

A valid inequality can come from a logical implication of the model, from rounding a constraint with integer coefficients or from a known family such as cover inequalities. It is added to the relaxation and the linear program is re-solved. The objective value can only improve for a minimization.

## Deciding whether a cut is violated

Compute the left side minus the right side at the candidate point. A positive difference for a less-than-or-equal cut means violation. The size of the difference measures how deep the cut cuts, and a tiny violation may be numerical noise rather than a useful cut.

Cuts are how a solver learns global facts about the integer set. They reduce the tree that branch and bound would otherwise have to explore.
