# Open and closed sets

Inside a metric space \\((X,d)\\), the most important regions are described by distances to a chosen point.

## Open Sets

A set \\(U \subseteq X\\) is **open** when every point of \\(U\\) is surrounded by a small ball that still fits inside \\(U\\). Formally: for every \\(x \in U\\) there is an \\(\varepsilon > 0\\) such that all points within distance \\(\varepsilon\\) of \\(x\\) also lie in \\(U\\). An interval like \\((0,1)\\) on the real line is open — nudge any point slightly and you stay inside — while \\([0,1]\\) is not, because no ball around the endpoint \\(0\\) avoids negative numbers.

## Closed Sets

A set \\(F\\subseteq X\\) is **closed** when its complement \\(X \setminus F\\) is open. Closed sets contain all their limit points: sequences drawn from \\([0,1]\\) can only converge to points inside \\([0,1]\\).

## Properties

Two useful facts follow immediately from the definitions. First, a set can be both open and closed (in a discrete space every set is), or neither — "open" and "closed" are not opposites. Second, arbitrary unions of open sets are open, and finite intersections of closed sets are closed; these two structural rules are what make topology run.
