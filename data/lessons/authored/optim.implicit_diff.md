# Implicit differentiation of layers

Many layers compute their output by solving an equation rather than by a fixed formula. Implicit differentiation differentiates through that solution without unrolling the solver.

## The implicit function theorem

If the output satisfies an equation with a nonsingular derivative in the output, the derivative of the output with respect to the parameters exists and solves a linear system. That system is built from the same Jacobian the solver already factorises.

## The adjoint solve

The gradient of a loss with respect to the parameters needs one linear solve with the transposed Jacobian. The cost equals one solver step regardless of how many iterations the forward solve took.

## The fixed-point case

When the layer outputs a fixed point, the derivative is obtained from the fixed-point equation directly. This is what makes implicit layers cheap and what makes their memory independent of depth.

To check your understanding, state the defining equation, use the transposed Jacobian and differentiate the fixed-point equation.

Implicit differentiation removes the memory cost of unrolling and replaces it with one linear solve, which is the same system the forward pass already handles.
