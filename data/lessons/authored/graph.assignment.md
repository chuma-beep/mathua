# The assignment problem

The assignment problem chooses one column for each row of a cost matrix so that the total cost is minimised. For three rows there are only six permutations to check.

## Permutations

An assignment is a permutation of the columns, one column per row. The cost is the sum of the chosen entries. Every permutation is a valid perfect matching.

## The optimum

Compare the six permutation costs and keep the smallest. The minimum is the optimal total cost. If the optimum is unique, the column chosen for each row is also unique.

## Checking a candidate

A proposed assignment is optimal exactly when its cost equals the minimum. Compute its cost and compare. This is a direct check that does not require re-solving the problem.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
