# The Four Color Theorem

## Colouring a Map

A proper coloring assigns colors so that adjacent vertices receive different colors. Planar graphs model maps, with regions as vertices and shared borders as edges. The four color theorem says that every planar graph can be colored with at most four colors.

The theorem is a bound, not an algorithm that always needs four colors. Many planar graphs need fewer, and some need exactly four. The statement is celebrated because it was the first major theorem proved with substantial computer assistance.

## The Four Color Bound

For a planar graph with n vertices, the guaranteed number of colors is the smaller of n and four. When n is at most three the complete graph K_n already needs n colors, so the bound cannot be improved. For n at least four, four colors always suffice.

The bound applies to simple planar graphs. Adding multiple edges or loops does not change the chromatic number in the usual vertex-coloring model, because those features do not create new adjacency between distinct vertices.

## Sharpness and Examples

The complete graph K_4 is planar and needs four colors, so the bound is sharp. Odd cycles need three colors, while even cycles need two. Complete graphs K_n with n at least five are not planar, which is why they do not contradict the theorem.

A quick way to practice the theorem is to compare the chromatic number of a family with the planar bound. The comparison shows both when the bound is tight and when a particular graph is far from needing all four colors.
