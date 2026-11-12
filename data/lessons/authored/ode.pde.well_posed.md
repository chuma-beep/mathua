# Well-posed problems

A boundary-value problem is well posed when it has exactly one solution and that solution depends continuously on its data. Hadamard separated the three requirements because a problem can fail any one of them while still looking harmless.

## The three requirements

Existence asks that a solution be available, uniqueness asks that there be only one, and continuous dependence asks that a small change of data produce a small change of solution. All three are needed before a numerical answer can be trusted.

## Failure by data

Hadamard's Cauchy problem for Laplace's equation is the standard counterexample. The boundary data can be changed by an arbitrarily small amount while the solution changes by an arbitrarily large amount, so the third requirement fails even though a solution exists.

## Why the definition matters

A discretisation inherits the conditioning of the continuous problem. An ill-posed problem needs regularisation before discretisation, because a finer mesh makes the discrete answer worse rather than better.

To check your understanding, state the three Hadamard requirements, state the data for the Cauchy problem and relate conditioning to discretisation.

Well-posedness is a property of the problem and its data, not of the method used to solve it.
