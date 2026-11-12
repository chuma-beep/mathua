# The error function

The error function is the integral of the Gaussian, with values pinned at zero and infinity.

## Definition

\(\operatorname{erf}(z)=\frac{2}{\sqrt{\pi}}\int_{0}^{z}e^{-t^{2}}dt\).

## Special values

\(\operatorname{erf}(0)=0\), \(\operatorname{erf}(\infty)=1\) and \(\operatorname{erf}(-z)=-\operatorname{erf}(z)\).

## Gaussian integral

The total Gaussian integral is \(\sqrt{\pi}\), which is exactly the normalisation in the definition.

To check your understanding, write the integral, use the origin and integrate the Gaussian.

The error function is the Gaussian integral with its total mass normalised to one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
