# Zorn's Lemma

## Partial orders

A partial order is a relation that is reflexive, antisymmetric and transitive. A chain is a subset in which every two elements are comparable, and an upper bound of a subset is an element of the ambient set that lies above every member of the subset. These three ideas are enough to state the lemma precisely.

Finite partial orders already display the phenomenon: a chain of divisors of a number has an upper bound, namely the number itself. Zorn's lemma extends that observation to orders that need not be finite, provided every chain has an upper bound.

## Maximal elements

An element is maximal when no strictly larger element exists. Maximal does not mean largest: a partial order can have several maximal elements that are incomparable. In the divisibility order on the positive integers, the maximal elements below a fixed number are its maximal proper divisors.

Zorn's lemma says that if every chain has an upper bound, then a maximal element exists. The hypothesis is checked chain by chain, which is often easier than constructing a maximal element directly. The lemma is equivalent to the axiom of choice in the usual foundations.

## An application

A typical application builds a maximal object by taking the union of a chain of smaller objects. If the union remains inside the class under study, the chain has an upper bound and the lemma supplies a maximal member. This pattern proves the existence of maximal ideals, bases and maximal independent sets.

To use the lemma safely, verify that the union of an arbitrary chain still satisfies every required property. A union of nested independent sets is independent, while a union of nested finite sets can be infinite. The verification step is where most proofs succeed or fail.
