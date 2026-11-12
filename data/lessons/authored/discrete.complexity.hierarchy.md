# The time hierarchy theorem

The time hierarchy theorem says more time buys strictly more problems: if \(f\) is time constructible and grows fast enough, then \(\mathrm{DTIME}(f(n))\) is a proper subset of \(\mathrm{DTIME}(f(n)\log f(n))\).

## Diagonalisation

Enumerate machines and define a language that disagrees with machine \(i\) on input \(i\). The disagreement is built in, so no machine in the class decides the language.

## Time constructibility

The bound must be computable within the time it names, otherwise the simulating machine cannot budget its steps.

## Consequences

\(\mathrm{P}\ne\mathrm{EXP}\) follows immediately. The theorem does not separate P from NP because their gap is not known to be large enough.

To check your understanding, enumerate the machines, simulate within the bound and separate P from EXP.

The hierarchy theorem is the diagonal argument applied to running time. It proves separation by construction rather than by conjecture.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
