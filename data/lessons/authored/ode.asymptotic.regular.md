# Regular perturbation

A regular perturbation expands the solution in powers of a small parameter when the parameter multiplies a term that does not vanish at leading order. Each order solves the same type of problem with a known right-hand side.

## The expansion

Write the solution as a leading term plus the parameter times a correction plus higher orders. Substituting into the equation and collecting powers of the parameter produces one problem per order.

## Order by order

The leading problem is the unperturbed equation. Each correction solves a linear problem driven by the previously computed terms, so the expansion is a sequence of linear solves.

## Validity

The expansion is regular when the corrections stay smaller than the leading term across the whole domain. Near a point where the leading term vanishes the ordering can fail, which is the first warning of a singular limit.

To check your understanding, write the series, solve the leading problem and compare correction to leading term.

Regular perturbation is ordinary Taylor expansion in a parameter, and its failure is more informative than its success.
