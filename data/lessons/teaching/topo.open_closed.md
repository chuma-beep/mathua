# Open and Closed Sets

In a metric space $(X, d)$, we define open and closed sets using open balls.

## Open Sets

A set $U \subseteq X$ is **open** if for every $x \in U$, there exists $r > 0$ such that $B_r(x) \subseteq U$.

**Examples:**
- Open intervals $(a, b)$ in $\mathbb{R}$ are open
- The open ball $B_r(x)$ is itself open
- The empty set $\emptyset$ and the whole space $X$ are open

## Closed Sets

A set $F \subseteq X$ is **closed** if its complement $X \setminus F$ is open.

Equivalently, $F$ is closed if it contains all its limit points: whenever $x_n \to x$ with $x_n \in F$, then $x \in F$.

**Examples:**
- Closed intervals $[a, b]$ in $\mathbb{R}$ are closed
- Closed balls $\{y : d(x, y) \le r\}$ are closed
- The set $\{1/n : n \in \mathbb{N}\} \cup \{0\}$ is closed (includes the limit point 0)

## Properties

**Open sets:**
- Arbitrary unions of open sets are open
- Finite intersections of open sets are open

**Closed sets:**
- Arbitrary intersections of closed sets are closed
- Finite unions of closed sets are closed

## Important Distinctions

- A set can be both open and closed ("clopen"): e.g., $\emptyset$ and $X$ in any metric space; more interestingly, a set like $[0, 1]$ in $\mathbb{R}$ with the discrete metric.
- A set can be neither open nor closed: e.g., $(0, 1]$ in $\mathbb{R}$ with the usual metric.
- Open does not mean "not closed" — a door analogy does not apply in topology.

## Interior, Closure, Boundary

- **Interior** $\text{int}(A)$: the largest open set contained in $A$
- **Closure** $\overline{A}$: the smallest closed set containing $A$
- **Boundary** $\partial A$: $\overline{A} \setminus \text{int}(A)$

**Example:** For $A = (0, 1]$ in $\mathbb{R}$:
- $\text{int}(A) = (0, 1)$
- $\overline{A} = [0, 1]$
- $\partial A = \{0, 1\}$
