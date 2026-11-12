# Division in the Gaussian integers

Gaussian division rounds the real and imaginary parts of a quotient to the nearest integers. The remainder is smaller than the divisor in norm.

## Rounding the quotient

Divide \(z\) by \(w\) in the complex numbers, then round both coordinates to the nearest integer to obtain the Gaussian quotient.

## The remainder bound

The remainder \(r=z-qw\) satisfies \(N(r)\le N(w)/2\), which is the analogue of a remainder smaller than the divisor.

## Euclidean algorithm

Replacing the pair by the divisor and the remainder decreases the norm, so the process terminates at a greatest common divisor.

To check your understanding, round a real quotient, square the remainder and take one Euclidean step.

Rounding the quotient makes the Gaussian integers a Euclidean domain.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
