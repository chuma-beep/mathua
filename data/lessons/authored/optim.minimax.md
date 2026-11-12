# Minimax problems

A minimax problem asks for a choice that is good against the worst case. The objective is written as the minimum over one block of variables of the maximum over another, and the two blocks pull in opposite directions.

## The outer and inner problem

For each choice of the first block the inner problem maximises over the second. The outer objective is that maximum, and the solver minimises it. The inner maximum is usually not smooth, so a minimax problem inherits kinks even when the underlying functions are smooth.

## When the orders commute

Swapping the two operations is not free in general. A minimax inequality always holds in one direction, and equality needs a convex-concave structure, a compact set and a continuity condition. Without equality a minimax point and a saddle point are different objects.

## Why the formulation appears

Worst-case design, robust regression and game theory all produce an inner adversary. Writing the adversary explicitly keeps the guarantee visible: the reported value is attained against a concrete worst case rather than an average.

To check your understanding, fix the first block, state the one-sided inequality and model adversarial uncertainty.

A minimax problem is the right model whenever the uncertainty is adversarial, and the inner maximum is the price of the guarantee.
