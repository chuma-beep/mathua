# Eigenfunction expansions

Sturm-Liouville problems supply a complete orthogonal family, so a boundary-value solution is a sum of eigenfunctions. The coefficients come from integrating the data against each eigenfunction.

## The zeros

On an interval of length \(L\), the \(n\)-th sine eigenfunction has \(n-1\) interior zeros. Counting zeros is a quick way to identify which mode a plot shows.

## The coefficients

A truncated expansion with \(N\) terms has \(N\) coefficients. Each coefficient is an inner product of the data with its eigenfunction, divided by the norm of that eigenfunction.

## The norm

For \(\sin(n\pi x/L)\) on \([0,L]\), the integral of the square is \(L/2\). That constant normalises every coefficient and never depends on \(n\).

To check your understanding, count interior zeros, count the terms and square the sine.

A boundary-value solution is a sum of orthogonal eigenfunctions.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
