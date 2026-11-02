# Resolution

## Clauses and Resolvents

Resolution works with clauses, which are disjunctions of literals. A literal is a variable or its negation. A formula in conjunctive normal form is a set of clauses that must all be true.

Two clauses can be resolved when one contains a variable and the other contains its negation. Delete the complementary pair and join the remaining literals to form the resolvent.

The resolvent is a logical consequence of the two parent clauses, so adding it preserves satisfiability.

## Deriving the Empty Clause

Repeated resolution is a complete refutation method for propositional logic. If a set of clauses is unsatisfiable, some sequence of resolutions derives the empty clause.

The empty clause represents a contradiction. Deriving it proves that the original formula has no satisfying assignment.

Unit resolution is the special case where one parent is a single literal, and it is often the fastest way to expose a contradiction.

## Reading a Refutation

Count the literals in a resolvent by adding the sizes of the parents and subtracting two for the complementary pair. A unit resolvent has one literal, and the empty resolvent has none.

When several pairs are complementary, resolve one pair at a time and keep track of the clauses produced. The order affects the length of the derivation but not whether a refutation exists.

Resolution is the basis of practical satisfiability solvers, which add heuristics for choosing the next pair.
