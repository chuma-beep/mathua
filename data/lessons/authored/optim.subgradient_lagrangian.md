# Subgradient optimisation of the Lagrangian

The Lagrangian dual is concave and usually nonsmooth, so its multipliers are updated by subgradient steps. Each step uses the constraint residual at the relaxed solution as the subgradient.

## The dual function

The dual value is the minimum of the Lagrangian over the integer variables. It is concave in the multipliers and its subgradient is the vector of constraint residuals at the minimiser.

## The step rule

A subgradient step moves the multipliers in the direction of the residual. The step size must decrease or be normalised against the estimated gap, otherwise the iterates oscillate.

## The stopping rule

The dual value rises toward the true bound, so a small increase over several iterations is the signal to stop. The best multipliers found are used, not the last ones.

To check your understanding, minimise over integers, move along the residual and track the best dual value.

Subgradient optimisation of the Lagrangian is a multiplier search, and its step rule is the only place where quality is decided.
