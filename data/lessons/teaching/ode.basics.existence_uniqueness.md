# Existence and Uniqueness of Solutions

An initial value problem asks for a solution that passes through a given point. The existence and uniqueness theorem gives conditions under which that solution is guaranteed to exist and to be the only one near the point.

## The existence and uniqueness theorem

Consider $y' = f(x,y)$ with $y(x_0) = y_0$. Suppose $f$ and $\partial f/\partial y$ are continuous on a rectangle containing $(x_0,y_0)$. Then there is an interval around $x_0$ on which a solution exists, and the solution is unique on that interval.

## Checking continuity

Check both functions at the initial point and on a small rectangle around it. For $y' = \sqrt{y-a}$ with $y_0 > a$, the right-hand side is continuous near the point. Its $y$-partial is $1/(2\sqrt{y-a})$, which is also continuous when $y > a$. If $y_0 = a$, the partial derivative is undefined, so the theorem gives no conclusion.

## When the theorem does not apply

Failure of the hypotheses does not prove that a solution fails to exist or that it is nonunique. It only means the theorem cannot decide the question. For $y' = \sqrt{y}$ with $y(0)=0$, the right-hand side is continuous but its $y$-partial is not continuous at the initial point. Solutions still exist, and uniqueness needs a separate argument.
