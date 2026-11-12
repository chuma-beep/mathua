# The Cantor function

The Cantor function is a continuous nondecreasing map of \([0,1]\) onto itself that is constant on each interval removed from the Cantor set. It is the standard example of a continuous function that is not absolutely continuous.

## The staircase construction

At the first stage remove the middle third and set the function to \(1/2\) there. At the next stage set it to \(1/4\) and \(3/4\) on the two new middle thirds. This gives \(C(1/3) = C(2/3) = 1/2\), and \(C(1/9) = 1/4\). Continuing on the remaining intervals defines \(C\) everywhere as a limit.

## Flat intervals

The function is constant on every removed interval. Those intervals have total length 1, so the derivative is zero on a set of full measure. On the Cantor set, which has measure zero, the function still climbs from 0 to 1.

## A singular function

A singular function is continuous and nondecreasing. Its derivative is zero almost everywhere. The Cantor function is the model case, and it satisfies \(C(1) - C(0) = 1\) while \(\int_0^1 C'(t)\,dt = 0\). It also maps the Cantor set onto \([0,1]\), so it is not injective.

The Cantor function separates continuity from absolute continuity, and its flat steps show that a derivative can vanish almost everywhere while the function still climbs.
