# Minimum feedback arc sets

## Directed cycles

A directed cycle is a closed directed walk that returns to its start without repeating an arc. Removing a feedback arc set breaks every directed cycle. The minimum size is the feedback arc number.

A graph with no directed cycles is a directed acyclic graph. The feedback arc set is the smallest set of arcs whose removal leaves a DAG.

## Removing arcs

Check a proposed set by deleting its arcs and searching for a remaining directed cycle. A single remaining cycle means the set is not a feedback arc set.

The minimum size is not always one. Disjoint directed cycles each need at least one removed arc, so the minimum is at least the number of disjoint cycles.

## Acyclic remainder

After removing one arc from a directed cycle, the remaining arcs form a directed path. A path has no directed cycle. Several cycles require one removal each.

Count the remaining arcs to check the arithmetic. If the original graph has m arcs and k disjoint cycles, removing one arc per cycle leaves m - k arcs.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
