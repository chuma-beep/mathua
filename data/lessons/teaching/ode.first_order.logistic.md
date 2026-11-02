# The Logistic Equation

The logistic equation models a population that grows quickly at first and then slows as it approaches a limit set by its environment.

## The logistic model

The equation is

$$\frac{dP}{dt} = rP\left(1 - \frac{P}{K}\right),$$

where $P$ is the population, $r > 0$ is the intrinsic growth rate and $K > 0$ is the carrying capacity. When $P$ is small, the factor $1 - P/K$ is close to 1 and the population grows almost exponentially. As $P$ approaches $K$, the factor approaches 0 and growth slows.

## Carrying capacity and equilibria

Set the right-hand side equal to zero. The constant solutions are $P = 0$ and $P = K$. The zero solution is unstable: a small positive population moves away from it. The carrying capacity $K$ is stable, because populations below $K$ increase and populations above $K$ decrease.

## Maximum growth and long-run behavior

The growth rate is a downward parabola in $P$. Its maximum occurs halfway to the carrying capacity, at $P = K/2$. Substituting gives a maximum growth rate of $rK/4$. Every positive solution with $0 < P(0) < K$ increases and approaches $K$ as $t \to \infty$; every solution above $K$ decreases toward $K$.
