# Solving Systems by Substitution

Substitution solves a system by turning two equations in two unknowns into one
equation in one unknown: isolate a variable in one equation, then substitute
that expression wherever the variable appears in the other equation. Use it
when one equation already isolates a variable — or isolates one cheaply.

## Isolate one variable

Pick the equation where a variable stands almost alone: a lone \(x\), a lone
\(y\), or anything with coefficient \(1\) or \(-1\). Solve that equation for
the variable, keeping the expression exact. For \(x + 2y = 7\), isolating
\(x\) gives \(x = 7 - 2y\). Do not round or approximate; the expression must
stay exact to feed the next step.

## Substitute and solve

Replace every occurrence of the isolated variable in the *other* equation with
the expression, then solve the resulting one-variable equation. With
\(x = 7 - 2y\) and \(3x - y = 7\): \(3(7 - 2y) - y = 7\), so \(21 - 7y = 7\)
and \(y = 2\). Back-substitute to recover the other variable:
\(x = 7 - 2(2) = 3\). A \(2 \times 2\) system collapses to one linear
equation, which always solves the same way.

## Check in both equations

A pair solves the system only if it satisfies *every* equation — checking one
is not enough. For \((3, 2)\): \(3 + 2(2) = 7\) and \(3(3) - 2 = 7\), so the
pair checks out in both. If a check fails, the substitution or the arithmetic
slipped; rework the same steps rather than starting over.
