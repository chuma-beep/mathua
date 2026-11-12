# Differentiable layers

A differentiable layer is a building block whose output is defined by an optimisation or a fixed point rather than by an explicit formula. The layer is trained end to end because its derivative is available.

## Construction

The layer states an inner problem whose data depend on the input. The forward pass solves it, and the backward pass differentiates the solution through the optimality conditions.

## The examples

Sparse coding, projection onto a convex set and equilibrium models are all differentiable layers. A projection layer is a quadratic program whose derivative is a linear system on the active set.

## The design constraint

The inner problem must be strictly convex or otherwise have a unique solution, because a nonunique solution has no well-defined derivative. Regularisation is what makes the layer well posed.

To check your understanding, define an inner problem, name sparse coding and require uniqueness.

A differentiable layer is an optimisation embedded in a network, and uniqueness of its inner solution is the condition that makes it trainable.
