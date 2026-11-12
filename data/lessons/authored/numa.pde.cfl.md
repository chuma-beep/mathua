# The CFL condition

The Courant-Friedrichs-Lewy condition says a numerical domain of dependence must contain the physical one. For advection it caps the time step at one cell crossing.

## The number

With speed \(c\), spacing \(\Delta x\) and step \(\Delta t\), the Courant number is \(C=c\Delta t/\Delta x\). Explicit advection requires \(C\le 1\).

## What it means

A signal may travel at most one cell per step. A larger step lets information outrun the stencil, and the scheme becomes unstable.

## Refinement

Halving the spacing halves the admissible time step, so the work per unit of simulated time doubles in one dimension. In two dimensions the constraint is stricter still.

To check your understanding, compute the Courant number, compare domains of dependence and halve the spacing.

The Courant number compares how far information travels with how far the stencil reaches.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
