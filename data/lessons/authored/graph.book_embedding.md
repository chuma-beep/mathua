# Book embeddings

A book embedding places the vertices on a line and draws each edge in one page without crossings, so the pages are half-planes sharing the spine.

## The definition

All vertices lie on a common line called the spine, and every edge lies in one of several half-planes called pages. Edges on the same page must not cross, and the book thickness is the least number of pages.

## Outerplanar graphs

A graph has book thickness one exactly when it is outerplanar, because one page is a half-plane with every vertex on the boundary. Trees are therefore one-page graphs.

## Complete graphs

The complete graph \(K_n\) has book thickness \(\lceil n/2\rceil\). The lower bound comes from the fact that at most one edge of a triangle can sit on each page between a fixed pair of spine positions.

To check your understanding, place vertices on the spine, draw one page and state the page count.

A book embedding is a planar drawing split across pages, and the page count measures how far the graph is from outerplanar.
