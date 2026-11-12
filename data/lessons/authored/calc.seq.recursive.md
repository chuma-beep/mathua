# Limits of recursive sequences

A recursive sequence is analysed by finding its fixed points and then proving the iteration approaches one. Monotonicity and boundedness usually supply the proof.

## Halving

If \(x_{n+1}=x_{n}/2\), every step halves the distance to zero. Starting from \(k\), the number of steps needed to drop below one is the number of halvings that fit into \(k\).

## Averaging

If \(x_{n+1}=(x_{n}+a)/2\), the limit solves \(x=(x+a)/2\). The fixed point is \(a\), so the sequence converges to the constant being averaged in.

## A linear map

If \(x_{n+1}=rx_{n}+c\) with \(|r|<1\), the limit is \(c/(1-r)\). For \(r=1/2\) the limit is \(2c\), twice the constant term.

To check your understanding, halve repeatedly, solve the fixed-point equation and solve for the limit.

Find the fixed point first, then prove the iteration reaches it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
