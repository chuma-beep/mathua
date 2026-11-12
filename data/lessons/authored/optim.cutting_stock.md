# The cutting stock problem

The cutting stock problem cuts rolls of fixed length into ordered pieces while using as few rolls as possible. Each cutting pattern is a column, and the pattern set is far too large to list.

## The pattern model

A variable counts how many rolls are cut with a given pattern. The demand constraints require enough pieces of each length, and the objective counts rolls.

## Column generation

The linear relaxation is solved over a small subset of patterns. The reduced cost of a missing pattern is a knapsack problem, and a negative reduced cost means the pattern should be added.

## The integer step

The linear solution is usually fractional, and a branch-and-price step restores integrality. Rounding the pattern counts first is a fast heuristic, but it can miss the optimum by a wide margin.

To check your understanding, count rolls per pattern, solve a restricted master and notice fractional counts.

Cutting stock is the textbook column generation problem: the master is easy, the pricing problem is a knapsack, and integrality is the remaining work.
