# PSPACE

PSPACE is the class of problems a Turing machine can solve using polynomially many tape cells. Quantified Boolean formulas are its standard complete problem.

## Polynomial space

The machine may run for exponential time but may not touch more than a polynomial number of cells. Reuse makes this generous.

## QBF completeness

Evaluating a quantified Boolean formula with alternating quantifiers is in PSPACE and is PSPACE-complete under polynomial reductions.

## Inclusions

\(\mathrm{NP}\subseteq\mathrm{PSPACE}\subseteq\mathrm{EXP}\). Both inclusions are believed strict, and neither strictness is known.

To check your understanding, bound the tape, read the quantifiers and place NP inside.

PSPACE captures games with polynomial board size, because a winning strategy is a quantified formula. Space, not time, is the resource that makes this class natural.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
