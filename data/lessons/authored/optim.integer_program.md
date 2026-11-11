# Integer programming

An integer program is a linear program whose variables may only take integer values. The restriction destroys the geometry that makes linear programming easy, because the feasible set becomes a set of lattice points and the linear relaxation is no longer exact. Modelling an integer program is mostly the art of choosing variables whose integer values mean something.

## Classifying a model

A model is an integer program when at least one variable is restricted to integers. If all variables are integer it is pure, and if only some are, it is mixed. Binary variables are the special case restricted to zero and one, and they encode yes-or-no decisions. The classification tells a solver which family of methods to use.

## Testing integrality

Take a candidate solution and check each variable that should be integral for a zero fractional part. A rounding error must be distinguished from a genuine fraction, so the test uses a tolerance. A solution that is integral is feasible for the integer problem only if it also satisfies every constraint, which rounding can break.

## Writing an integer objective

Choose variables that count or select, write the objective as a linear combination of them and add constraints that express the logical conditions. A yes-or-no choice becomes a binary variable with a linking constraint. The objective must not contain products of variables, since a product would leave the linear class.

Integer programming is where linear models meet combinatorial decisions. The modelling discipline, not the solver, is usually what decides whether an instance is tractable.
