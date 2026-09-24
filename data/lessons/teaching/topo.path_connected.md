# Path-Connected Spaces

**Path-connected:** $X$ is path-connected if any two points $x,y$ join by a continuous path $p:[0, 1]\to X$ with $p(0)=x$ and $p(1)=y$. Every path-connected space is connected; the converse fails.

## Worked: intervals are path-connected

Join $x$ and $y$ in $(0, 1)$:

1. Set $p(t)=(1-t)x+ty$, the straight-line segment between them.
2. Each $p(t)$ stays strictly between $x$ and $y$, hence inside $(0, 1)$.
3. $p$ is continuous with $p(0)=x$ and $p(1)=y$.

So intervals are path-connected via linear paths, and the same formula works in $\mathbb{R}^n$.

## Worked: the punctured plane stays path-connected

Join two points in $\mathbb{R}^2$ minus the origin:

1. If the straight segment misses the origin, take it.
2. Otherwise detour: travel from $x$ outward to a circle clearing the origin, arc around, then head into $y$.
3. The detour is a concatenation of continuous paths, hence a path in the punctured plane.

So deleting one point does not disconnect by paths: go around the hole.

## When connected does not imply path-connected

Inspect the topologist's sine curve $S=\{(x,\sin(1/x)):x>0\}$ plus its limit segment $L=\{0\}\times[-1, 1]$:

1. $S$ union $L$ is connected: $L$ consists of limit points of $S$.
2. Any path from $S$ to $L$ would have to oscillate infinitely fast approaching $L$, breaking continuity at the endpoint.
3. Hence no path joins the curve to the segment, though no open split separates them.

So path-connected is strictly stronger: this space is connected but not path-connected, and $\mathbb{Q}$ is neither.
