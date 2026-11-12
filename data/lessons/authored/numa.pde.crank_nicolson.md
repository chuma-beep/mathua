# The Crank-Nicolson scheme

Crank-Nicolson averages the explicit and implicit operators at the two time levels. It is second order in time and unconditionally stable, which makes it a default for diffusion.

## The average

The scheme uses \(\theta=1/2\), so each side carries half the spatial operator. The update solves \((I-\tfrac{r}{2}L)u^{n+1}=(I+\tfrac{r}{2}L)u^{n}\).

## Order

The time truncation error cancels to second order. The price is a linear system at every step.

## Behaviour

Very large steps produce oscillatory transients near sharp initial data even though the scheme is stable. A few fully implicit steps at the start damp them.

To check your understanding, split the operator in half, cancel the first-order term and recognise the oscillation.

Crank-Nicolson is the midpoint rule applied to a stiff spatial operator.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
