# Caratheodory's criterion

Caratheodory's criterion decides which sets are measurable. A set is measurable when it splits every test set additively.

## Splitting test sets

A set \(E\) is measurable when \(m^*(A) = m^*(A \cap E) + m^*(A \cap E^c)\) for every set \(A\). Take \(E = [0,1]\) and \(A = [0,2]\). The two pieces have measures 1 and 1, which add to 2. Every interval passes the test.

## The sigma-algebra

The sets passing the test form a sigma-algebra. It contains the empty set and is closed under complements. It is also closed under countable unions. A null set passes trivially, because \(m^*(A \cap E) = 0\) and the other piece already has measure \(m^*(A)\).

## Measurable sets

Countable additivity follows for measurable sets, and intervals, open sets and closed sets are all measurable. A Vitali set fails the test: some test set \(A\) is split into pieces whose outer measures add to more than \(m^*(A)\). The criterion is therefore a genuine restriction rather than a formality.

Caratheodory's criterion converts additivity into a local test, and the test is what makes the measurable sets a sigma-algebra.
