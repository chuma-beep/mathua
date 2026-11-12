# The saddle-point gap

The saddle gap measures how far a pair is from equilibrium. It is the worst improvement either player could make by moving alone, so a small gap certifies that the pair is nearly optimal.

## The definition

The gap is the maximum over the second block of the value minus the minimum over the first block, evaluated at the candidate pair. Both terms are single-block optimisations, so the gap can often be bounded without solving the full problem.

## Certification

A gap below a tolerance proves that the candidate pair is within that tolerance of the equilibrium value, for any function, convex or not. That makes the gap a stopping rule that does not assume the method converges.

## Computing the gap

For a bilinear problem the two single-block maxima are linear programs, so the gap is computable exactly. For a general convex-concave function the two terms need a small inner solve, and that cost is the price of a certificate.

To check your understanding, maximise over the second block, bound the distance to equilibrium and solve two single-block problems.

The saddle gap turns a convergence claim into a number, and a number can be checked.
