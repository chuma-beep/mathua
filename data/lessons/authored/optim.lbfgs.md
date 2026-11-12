# Limited-memory quasi-Newton

A quasi-Newton method stores an approximate inverse Hessian, which is a full matrix. Limited memory keeps only the last few update pairs and applies them as a two-loop recursion.

## The storage

Each iteration stores the step and the gradient change, two vectors of the problem dimension. With a memory of ten pairs the cost is twenty vectors regardless of the dimension.

## The two-loop recursion

The product of the approximate inverse with a vector is computed by walking the pairs backward and then forward. The result equals the product of the implied dense matrix without forming it.

## The scaling

A scalar scaling of the initial matrix is chosen from the latest pair so that the implied step has a sensible length. That choice is what makes the method competitive with a well-tuned gradient method on poorly scaled problems.

To check your understanding, keep recent step pairs, walk pairs backward and use the latest pair.

Limited-memory quasi-Newton keeps the curvature information that fits in a fixed budget and discards the rest.
