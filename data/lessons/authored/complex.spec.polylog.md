# Polylogarithms

The polylogarithm extends the logarithm and the dilogarithm, and at 1 it returns the zeta function.

## Definition

\(\operatorname{Li}_{s}(z)=\sum_{n\ge1}z^{n}/n^{s}\) for \(|z|<1\), continued analytically beyond the disk.

## At one

\(\operatorname{Li}_{s}(1)=\zeta(s)\) whenever the series converges, so the polylogarithm packages the zeta values.

## Dilogarithm

\(\operatorname{Li}_{1}(z)=-\log(1-z)\), and \(\operatorname{Li}_{2}(1)=\pi^{2}/6\).

To check your understanding, write the series, set z to one and use the logarithm case.

The polylogarithm is a one-parameter family that starts at the logarithm.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
