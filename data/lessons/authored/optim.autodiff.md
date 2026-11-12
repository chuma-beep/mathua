# Automatic differentiation

Automatic differentiation computes exact derivatives of a program by applying the chain rule to each operation. It is neither symbolic nor numerical: it propagates derivative values through the same code that computes the function.

## Forward and reverse

Forward mode carries a derivative with each value and costs a constant multiple per input direction. Reverse mode records the computation and propagates adjoints backward, costing one pass for a scalar output regardless of input count.

## The graph

The recorded graph is a sequence of elementary operations with known local derivatives. The chain rule composes them, and the result is exact up to floating-point rounding.

## Memory against time

Reverse mode needs the intermediate values, so memory grows with the depth of the computation. Checkpointing trades recomputation for memory and is the standard remedy for deep graphs.

To check your understanding, carry derivatives forward, compose local derivatives and store intermediates.

Automatic differentiation gives exact derivatives at a cost that depends on the direction count and on the memory allowed for the recorded graph.
