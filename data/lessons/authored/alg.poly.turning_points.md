# Turning points of a polynomial

A turning point is a place where a polynomial changes from increasing to decreasing, or the other way around. A smooth graph can bend only so many times, and its degree says how many bends are possible.

## Degree sets the ceiling

The derivative of a degree \(n\) polynomial has degree \(n-1\). That single fact gives the ceiling: a polynomial of degree \(n\) has at most \(n-1\) turning points.

A line has degree one and no turns. A parabola has degree two and one turn. A cubic can turn twice. The pattern continues, and no choice of coefficients can beat it.

The ceiling is a maximum, not a promise. A cubic may have two turns, one turn, or none, depending on whether its derivative has two real zeros, one repeated zero, or no real zeros. What the degree fixes is the largest number that is possible, not the number the graph must have.

## The derivative counts turns

At a turning point the graph is momentarily flat, so the derivative is zero there. Turning points therefore live among the solutions of \(f'(x)=0\), which is why the degree of the derivative controls them.

A root of the derivative does not always produce a turn. If the derivative touches zero and stays on the same side, as with \(f(x)=x^3\) at the origin, the graph flattens but keeps moving in the same direction. Such a repeated root costs degree without adding a turn, so the count can be lower than the ceiling.

For the greatest possible number, every zero of the derivative is a simple zero and each one changes the sign of the slope. Then the turns occur one after another, and the count reaches the degree of the derivative.

## From turns back to degree

The relationship also runs backward. If a polynomial has \(t\) turning points, its derivative must have at least \(t\) real zeros, so the derivative has degree at least \(t\). The polynomial itself then has degree at least \(t+1\).

This is the least possible degree, not the exact degree. A quartic can have three turning points only if its derivative is a cubic with three real zeros, but the same quartic might be written with extra terms that do not add turns. Extras such as a y-intercept shift the graph without changing how many times it bends.

To use the rule, count the turns and add one. That is the smallest degree that can produce the picture. To go the other way, read the degree and subtract one. That is the largest number of turns the polynomial can have.
