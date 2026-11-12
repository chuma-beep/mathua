# Circulations

A circulation is a flow with no source and no sink. Every vertex conserves flow, and the feasible circulations form a cone.

## Conservation

At every vertex the inflow equals the outflow. A circulation can carry flow around a cycle without any external supply, and the zero flow is always feasible.

## Lower bounds

With lower bounds \(l_e\) and upper bounds \(u_e\), a feasible circulation exists exactly when every cut has enough capacity to absorb the lower bounds crossing it. Subtracting the lower bounds turns the question into an ordinary flow feasibility problem.

## Cycle decomposition

Every circulation decomposes into flows around simple cycles. The decomposition is not unique, but the total flow on any cut is fixed by the circulation.

To check your understanding, balance inflow and outflow, subtract the lower bounds and find a positive cycle.

Remove the source and sink, and what remains is a flow that only moves in loops.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
