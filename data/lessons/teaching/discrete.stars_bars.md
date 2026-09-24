# Stars and Bars

**Stars and bars:** The number of nonnegative integer solutions to $x_1+\cdots+x_k=n$ is $\binom{n+k-1}{k-1}$: arrange $n$ stars and $k-1$ bars in a row, and the bars split the stars into $k$ boxes. This also counts distributions of $n$ identical objects into $k$ distinct boxes.

## Worked: x1+x2 = 3

List all nonnegative solutions directly:
1. $x_1$ can be 0, 1, 2, or 3, forcing $x_2=3,2,1,0$: pairs $(3,0),(2,1),(1,2),(0,3)$.
2. That is 4 solutions, and the formula gives $\binom{3+2-1}{2-1}=\binom{4}{1}=4$.
3. Each solution is one placement of a single bar among 4 slots: $3\star|\star$, $\star\star|\star\star$ up to ordering.

So one bar cutting a row of 3 stars enumerates every solution exactly once.

## Worked: 4 identical balls into 3 boxes

Encode each distribution as stars and bars:
1. Write 4 stars $\star\star\star\star$ and 2 bars splitting them into 3 groups, e.g. $\star\star||\star\star$ means 2, 0, 2.
2. Every arrangement is a choice of 2 bar positions among $4+2=6$ slots.
3. Count: $\binom{6}{2}=15$ distributions.

So identical balls into distinct boxes is $\binom{n+k-1}{k-1}$ with $n=4$, $k=3$.

## Worked: positive solutions to x1+x2+x3 = 5

Require $x_i\ge1$ and shift down to the nonnegative case:
1. Set $y_i=x_i-1$, so $y_i\ge0$ and $y_1+y_2+y_3=5-3=2$.
2. Nonnegative solutions in 3 variables: $\binom{2+3-1}{3-1}=\binom{4}{2}=6$.
3. Each $y$-solution lifts to exactly one $x$-solution by adding 1 back.

So positive solutions to $k$ variables summing to $n$ number $\binom{n-1}{k-1}$.
