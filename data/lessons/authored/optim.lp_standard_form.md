# Linear programming standard form

Standard form writes a linear program with a minimizing objective, equality constraints and non-negative variables. Any linear program can be rewritten this way by adding slack variables, splitting free variables and flipping the sign of a maximizing objective. The rewrite is mechanical, and it is what lets one solver accept every model.

## Converting a model

Turn a maximizing objective into a minimizing one by negating its coefficients. Turn an inequality into an equality by adding a slack variable with the right sign. Replace a free variable with the difference of two non-negative variables. Each move changes the matrix and keeps the set of attainable objective values the same.

## Reading the matrix

The constraint matrix has one row per equality and one column per variable, including slacks. The right-hand side is non-negative in standard form. The objective row lists the costs of every column. Reading the matrix back into a model is a check that the conversion preserved the meaning.

## Counting basic variables

A basic solution sets enough variables to zero that the remaining square system can be solved. The number of basic variables equals the number of equality constraints, and the number of non-basic variables is the total number of columns minus that number. This count is what makes the simplex method finite: there are finitely many bases.

Standard form is the interface between a model and a solver. Converting correctly is a prerequisite for every duality and simplex statement that follows.
