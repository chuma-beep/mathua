# Boolean Functions

## From Truth Tables to Functions

A Boolean function takes n input bits and returns one output bit. Its behaviour is completely described by a truth table with 2 to the power n rows, one for each assignment of the inputs.

Two functions are equal when they agree on every row. A formula is only a description of a function; different formulas can describe the same function.

The number of possible functions grows very quickly with n, which is why Boolean functions are usually described by formulas or by a list of the rows where the output is one.

## Minterms and Normal Forms

Each row where the output is one contributes a minterm, a conjunction that is true on exactly that assignment. Joining the minterms with OR gives the disjunctive normal form.

The number of minterms is the number of ones in the truth table. A function with no minterms is the constant zero, and a function with every row a minterm is the constant one.

Simplifying a normal form means finding a shorter formula for the same set of rows.

## Counting and Evaluating

Counting the ones in a table is the most direct way to identify a Boolean function. Evaluating a formula at an assignment is the same process in the other direction.

Use the standard precedence NOT, then AND, then OR, and parenthesise anything unusual. A single wrong row is enough to change the function.

These skills support circuits, satisfiability and coding, where the exact set of accepting assignments is the object of interest.
