# Stars and Bars

**Stars and bars:** Number of nonnegative solutions to $x_1+\cdots+x_k=n$ is $\binom{n+k-1}{k-1}$ (identical objects into distinct boxes).

## Counting Distributions

### Nonnegative Solutions
$\binom{n+k-1}{k-1}$ counts $n$ stars and $k-1$ bars. For $x_1+x_2=3$, $\binom{4}{1}=4$ solutions: $(3,0),(2,1),(1,2),(0,3)$.

### Positive Solutions
If $x_i\ge1$, set $y_i=x_i-1$: $y_1+\cdots+y_k=n-k$, so $\binom{n-1}{k-1}$ solutions.

## Example

$4$ identical balls into $3$ boxes: $\binom{4+3-1}{3-1}=\binom{6}{2}=15$ distributions.
