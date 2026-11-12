# The well-ordering theorem

Every set can be well ordered. The statement is equivalent to the axiom of choice, so it is a theorem about the reach of that axiom rather than a construction anyone performs by hand.

## Well-orders

A well-order is a linear order in which every nonempty subset has a least element. The natural numbers are well ordered. The integers are not, because the set of negative integers has no least element. The interval \((0,1)\) is not well ordered either, since it has no least element at all.

## From choice to well-order

Choice supplies a selection function on the nonempty subsets of a set. Transfinite recursion uses that function to pick the next element at each stage, and the process stops only when the set is exhausted. Zorn's lemma reaches the same conclusion through a maximal chain.

## Transfinite recursion

Recursion over ordinals builds a sequence indexed by ordinals rather than integers. The pairs of natural numbers can be ordered by maximum first and sum second, so \((2,1)\) precedes \((3,1)\). Each ordinal stage appends one more element. The order type of the resulting well-order is a cardinal, and every cardinal is an initial ordinal.

The theorem converts choice into an order, and the order is what makes transfinite induction available on an arbitrary set.
