# Chebyshev polynomials

Chebyshev polynomials satisfy T_n(cos θ)=cos(nθ), which pins their values at 1, -1 and 0.

## Definition

The recurrence \(T_{0}=1\), \(T_{1}=z\), \(T_{n+1}=2zT_{n}-T_{n-1}\) produces a polynomial of degree \(n\).

## Values at 1 and -1

\(T_{n}(1)=1\) and \(T_{n}(-1)=(-1)^{n}\), because the cosine takes those endpoint values.

## Values at zero

\(T_{n}(0)=0\) for odd \(n\), while for even \(n\) it is \((-1)^{n/2}\).

To check your understanding, apply the recurrence, evaluate at one and test odd degree.

The cosine identity converts trigonometric parity into polynomial values.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
