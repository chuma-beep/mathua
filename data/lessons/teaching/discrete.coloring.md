# Graph Coloring and Chromatic Number

**Chromatic number:** For a graph $G$, $\chi(G)$ is the fewest colors needed so adjacent vertices get different colors. Complete graphs need the most: $\chi(K_n)=n$, since every pair is adjacent. Bipartite graphs need only 2, and odd cycles need 3.

## Coloring complete graphs

Take $K_3$, the triangle: each vertex touches the other two, so no two share a color and $\chi(K_3)=3$. In $K_n$ every vertex is adjacent to all $n-1$ others, forcing $n$ distinct colors. At the other extreme, a bipartite graph splits into two independent sides, so 2 colors always suffice and $\chi=2$ for any bipartite graph with at least one edge.

## Worked: C_5 needs 3 colors

Try 2-coloring the 5-cycle $C_5$ with vertices $v_1,\dots,v_5$ in order:
1. Alternate colors: $v_1,v_3$ red and $v_2,v_4$ blue satisfies every edge among the first four.
2. $v_5$ is adjacent to both $v_4$ (blue) and $v_1$ (red), so neither color fits $v_5$.
3. A third color for $v_5$ completes a valid coloring, so $\chi(C_5)=3$.

So an even cycle needs only 2 colors but an odd cycle $C_{2k+1}$ always needs 3: the alternation fails exactly at the closing edge.

## Worked: the greedy bound

Order the vertices arbitrarily and color each with the first available color:
1. When a vertex $v$ is colored, at most $\Delta$ of its neighbors are already colored, where $\Delta$ is the maximum degree.
2. Those neighbors block at most $\Delta$ colors, so among colors $1,\dots,\Delta+1$ at least one is free.
3. Giving $v$ the first free color keeps the coloring valid at every step.

So greedy coloring never needs more than $\Delta+1$ colors, and $\chi(G)\le\Delta+1$ for every graph $G$.
