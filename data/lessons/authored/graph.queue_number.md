# The queue number

A queue layout places vertices on a line and puts edges in queues so that no two edges of a queue nest.

## The definition

Vertices lie on a line and each edge is assigned to a queue. Two edges in the same queue may not nest, meaning their endpoints cannot alternate along the line. The queue number is the least number of queues.

## Nesting versus crossing

A book embedding forbids crossings within a page, while a queue layout forbids nesting within a queue. The two parameters are dual, and a graph can have small page number and large queue number or the reverse.

## Values

A path has queue number one, a tree has queue number one and a cycle has queue number two. The complete graph \(K_n\) has queue number \(\lceil n/2\rceil\).

To check your understanding, order the vertices, compare the two patterns and layout a path.

Queues stack edges instead of drawing them on pages, and the forbidden pattern is nesting rather than crossing.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
