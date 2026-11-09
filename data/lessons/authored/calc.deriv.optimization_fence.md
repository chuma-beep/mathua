# Optimization: fencing

A fencing problem fixes the amount of fence and asks for the largest enclosed area. The constraint reduces the problem to one variable.

## Writing the constraint

A pen against a river needs fence on three sides. If the two perpendicular sides have length \(x\), the parallel side has length \(P-2x\).

The constraint is the total fence length. It eliminates one variable from the area formula.

The river side is free, which is why only three sides appear.

## The area function

The area is \(A=x(P-2x)\). It is a downward parabola in \(x\).

The maximum occurs at the vertex, where the derivative is zero.

The endpoints of the feasible interval give zero area, so the vertex is the maximum.

## The optimal dimensions

Differentiating gives \(P-4x=0\), so \(x=P/4\). The parallel side is \(P/2\).

The optimal pen is twice as long as it is wide in the perpendicular direction. The maximum area is \(P^2/8\).

Doubling the available fence quadruples the maximum area, which is a useful sanity check.
