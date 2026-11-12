# Dedekind cuts

A Dedekind cut is a way to build a real number out of rationals alone. It is a partition of the rationals into a lower set and an upper set, with the lower set closed downward and carrying no greatest element.

The construction matters because it fills the gaps that the rationals leave. The cut for \(\sqrt{2}\) is the set of rationals whose square is less than two. That set is a number even though no rational sits at its edge.

## Cuts of the rationals

Write \(A\) for the lower set. The cut for \(\sqrt{2}\) contains 1.4, since \(1.4^2 = 1.96 < 2\). It excludes 1.5, since \(2.25 > 2\). Every rational below a member is a member, and no member is greatest, because a rational just below \(\sqrt{2}\) always has a larger rational still below it.

## Completeness from cuts

A set of cuts that is bounded above has a least upper bound, namely the union of its lower sets. The cuts \(A_n\) for \(1 - 1/n\) have union the cut for 1, which is their supremum. This is the completeness property that the rationals lack.

## Comparing constructions

Cuts and Cauchy sequences both complete the rationals, and the resulting ordered fields are isomorphic. Addition of cuts is defined by adding lower sets, so the cut for \(\sqrt{2}\) added to itself gives the cut for \(2\sqrt{2}\). The cut construction keeps the order visible at every step.

A cut turns a gap into an object, and the union of cuts supplies the least upper bound that the rationals cannot produce.
