# Certified robustness

Certified robustness proves that no perturbation inside a region changes a prediction. The certificate is a radius, and it holds for every input in that radius.

## The margin bound

A classifier with margin \(m\) and Lipschitz constant \(L\) is certified within \(m/L\). A larger margin or a smaller constant widens the certificate.

## Randomised smoothing

Gaussian smoothing certifies a radius from the probability of the top class. More samples tighten the estimate of that probability.

## The cost

Certification is usually done per input. A test set of \(n\) inputs needs \(n\) certificates, each with its own optimisation or sampling.

To check your understanding, read the margin, sample the noise and certify one input.

A certificate replaces an empirical attack with a proof over a region.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
