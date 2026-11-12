# Space complexity

Space complexity counts the working tape cells a Turing machine uses. PSPACE collects the problems solvable in polynomial space, and it contains NP because a polynomial computation cannot write down more than polynomially many configurations.

## The resource

A machine's space is the number of tape cells it visits, excluding the read-only input. Reusing cells is allowed, so the count is the maximum over the run.

## Inclusions

\(\mathrm{P}\subseteq\mathrm{NP}\subseteq\mathrm{PSPACE}\). The second inclusion holds because a polynomial-time computation touches only polynomially many cells.

## Savitch's theorem

Nondeterministic polynomial space equals deterministic polynomial space, so PSPACE is closed under nondeterminism in a way NP is not known to be.

To check your understanding, count working cells, place P inside NP and square the space.

Space is a more forgiving resource than time because configurations can be reused. That is why PSPACE contains NP and why Savitch can remove nondeterminism.
