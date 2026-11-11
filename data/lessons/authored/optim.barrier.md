# Logarithmic barrier method

A barrier method keeps the iterates strictly inside the feasible region by adding a logarithmic term that grows without bound at the boundary. Decreasing the barrier parameter traces a path of interior minimizers toward the constrained optimum. The method is naturally suited to inequalities and produces strictly feasible iterates throughout.

## Evaluating a barrier

Add minus the barrier parameter times the logarithm of each slack to the objective. The term is finite inside the region and tends to infinity at the boundary, so a minimizer cannot approach the boundary too closely. The barrier parameter controls how strong that repulsion is.

## Taking a barrier step

For a fixed parameter, minimize the barrier objective from a strictly interior point. The result is the next point on the central path. Reducing the parameter and repeating moves the path toward the boundary and toward the constrained optimum.

## Testing the central path

The central path is the set of minimizers as the parameter varies. A point is on it when the gradient of the barrier objective vanishes. The duality gap at a barrier point is bounded by the parameter times the number of inequalities, which gives a computable stopping rule.

Barriers replace a hard boundary with a smooth repulsion. They are the interior-point idea, and the parameter schedule is what connects the interior to the solution.
