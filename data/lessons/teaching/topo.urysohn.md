# Urysohn's Lemma

**Urysohn's lemma:** In a normal space $X$, for disjoint closed sets $A,B$ there is a continuous $f:X\to[0, 1]$ with $f=0$ on $A$ and $f=1$ on $B$. Normal spaces have enough continuous functions to separate closed sets numerically.

## Worked: separating two points on the line

Separate $A=\{0\}$ and $B=\{1\}$ in $\mathbb{R}$:

1. On [0, 1] set $f(x)=x$: then $f(0)=0$ and $f(1)=1$ as required.
2. Extend constantly: $f=0$ left of 0 and $f=1$ right of 1, keeping continuity at the joints.
3. The result is continuous on all of $\mathbb{R}$ with the exact prescribed values.

So the lemma holds here by an explicit formula; normality guarantees such functions exist in general.

## Worked: the metric formula

Build the separator in any metric space with distance $d$:

1. Set $g(x)=d(x, A)/(d(x, A)+d(x, B))$; the denominator never vanishes since $A,B$ are disjoint and closed.
2. On $A$ the numerator is 0, so $g=0$; on $B$ the second term is 0, so $g=1$.
3. Distance functions are continuous, hence $g$ is continuous on all of $X$.

So metric spaces are normal with room to spare: the distance ratio writes down the Urysohn function directly.

## When separation fails

Try to separate in a non-normal space:

1. Normality means disjoint closed sets admit disjoint open neighbourhoods first.
2. Without those neighbourhoods there is no room to fit the transitioning values between 0 and 1.
3. Hence the lemma can fail outside normal spaces: normality is the sharp hypothesis.

So Urysohn's lemma characterizes the setting: normal spaces are exactly where closed sets can be separated by continuous functions.
