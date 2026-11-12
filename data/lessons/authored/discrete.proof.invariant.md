# Loop invariants

A loop invariant is a statement true before the loop, preserved by every iteration, and strong enough at exit to prove the postcondition. It is the standard tool for proving iterative algorithms correct.

## Initialisation

The invariant must hold before the first iteration. For a search over an array, the interval bounds start at the full array.

## Preservation

Assuming the invariant before an iteration, the body must restore it afterwards. This is the inductive step of the correctness proof.

## Termination

At exit the invariant and the negation of the loop guard together imply the postcondition. A separate argument shows the loop terminates.

To check your understanding, state the invariant, assume before the body and combine with the exit guard.

An invariant is an induction hypothesis carried through the loop. Initialisation is the base case, preservation is the step, and exit combines the invariant with the guard.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
