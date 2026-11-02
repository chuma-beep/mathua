# Sets of Measure Zero

## What measure zero means

A subset of the real line has measure zero when it can be covered by countably many open intervals whose total length is smaller than any prescribed positive number. A single point has measure zero, and so does every finite set. The definition is about covering efficiency rather than about the number of points.

Countable sets also have measure zero: enumerate the points and place a short interval around the n-th point with length proportional to a geometric sequence. The total length can be made arbitrarily small. This is the first sign that cardinality and measure are different scales.

## Covering by intervals

To prove that a set does not have measure zero, one must rule out every countable interval cover. Compactness helps: on a closed bounded interval, any open cover has a finite subcover, and a finite cover has a definite total length. If every finite cover must be long, the set cannot be null.

The measure of an interval is its length, and the measure of a countable disjoint union of intervals is the sum of their lengths when the union is measurable. These facts make many null-set computations routine. They also explain why a null set can be dense, as the rationals are.

## Countable unions

A countable union of null sets is null. Given a total error budget, split it into countably many smaller budgets and cover the n-th set within its share. The geometric series then absorbs the sum. This closure property is what makes measure zero a practical notion rather than a single example.

Null sets are invisible to the integral: changing a function on a null set does not change its Riemann integral when the function remains bounded, and it does not change a Lebesgue integral at all. That insensitivity is the reason null sets appear in almost-everywhere statements and in criteria for integrability.
