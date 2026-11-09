# Dynamic programming

Dynamic programming solves a problem by storing the answers to overlapping subproblems. The table is filled in an order that makes every dependency available.

## Optimal substructure

A problem has optimal substructure when an optimal solution contains optimal solutions to its subproblems. This is the condition that lets a recurrence compute the answer from smaller cases. Without it, a table of subproblem values would not compose.

## The table recurrence

The state records the part of the problem that remains. The recurrence expresses the value at one state through smaller states. Filling the table from small to large gives every value once.

## Counting and reconstruction

A table can count solutions as well as find an optimum. A separate choice table records which transition produced each value. Following the choices reconstructs an optimal solution, while the final table entry gives its value.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
