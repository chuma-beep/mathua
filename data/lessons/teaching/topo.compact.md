# Compactness

**Compact set:** A set $K$ is compact if every open cover of $K$ has a finite subcover. Compact sets behave like finite sets: continuous images stay compact, and in $\mathbb{R}^n$ compactness is just closed plus bounded.

## Definition of compact sets

Cover [0, 1] by the opens $U_n = (-1, 1 - 1/n)$ together with $V = (1/2, 2)$:

1. Every $x$ below 1 lies in some $U_n$ once $1 - 1/n$ passes $x$; points near 1 lie in $V$.
2. Take $N$ with $1 - 1/N$ above 1/2; then $U_N$ plus $V$ already cover all of [0, 1].
3. So this arbitrary cover has a 2-set finite subcover, as the definition demands.

So compactness means no cover can force infinitely many sets to do essential work: finitely many always suffice.

## Heine-Borel theorem

Test whether $(0, 1]$ is compact in $\mathbb{R}$:

1. Check boundedness: $(0, 1]$ sits inside $[-1, 2]$, so it is bounded.
2. Check closedness: 0 is a limit point (take $1/n$) but 0 is missing, so it is not closed.
3. Heine-Borel says compact needs both, so $(0, 1]$ is not compact.

So in $\mathbb{R}^n$ the abstract cover condition collapses to two concrete checks: closed and bounded, nothing more.

## Examples and non-examples

Sort four sets with the two-check test:

1. [0, 1]: closed and bounded, so compact.
2. $(0, 1)$: bounded but not closed, so not compact.
3. $\mathbb{R}$: closed but not bounded, so not compact.
4. [0, 1] union [2, 3]: a finite union of compact sets, closed and bounded, so compact.

So openness fails closedness, unboundedness fails boundedness, and finite unions of compact sets stay compact.
