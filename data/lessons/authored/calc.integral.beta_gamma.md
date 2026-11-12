# The beta-gamma relation

## The beta integral

The beta integral is \(B(p,q) = \int_0^1 t^{p-1}(1-t)^{q-1} dt\), defined for positive \(p\) and \(q\). At \(p = 2\) and \(q = 3\) the integrand is \(t(1-t)^2 = t - 2t^2 + t^3\), whose integral is \(1/2 - 2/3 + 1/4 = 1/12\). The integral converges because both endpoints stay integrable.

## The gamma relation

The gamma relation is \(B(p,q) = \Gamma(p)\Gamma(q)/\Gamma(p+q)\). At \(p = 2\) and \(q = 3\) this reads \(1 \cdot 2/24 = 1/12\), matching the direct integral. The relation converts a two-parameter integral into factorials when the arguments are integers. It also turns a product of gamma values into a single beta value, which is why a table of one function serves the other.

## A worked value

A worked value: \(B(1/2,1/2) = \Gamma(1/2)^2/\Gamma(1)\). Since \(\Gamma(1/2) = \sqrt{\pi}\) and \(\Gamma(1) = 1\), the value is \(\pi \approx 3.1416\). The direct integral \(\int_0^1 t^{-1/2}(1-t)^{-1/2} dt\) is improper at both ends and still equals \(\pi\). The substitution \(t = \sin^2\theta\) reduces it to the integral of 2 over the first quadrant, which is \(\pi\).

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
