# Fill-in

Elimination on a sparse matrix creates nonzeros where the original had zeros. That fill-in, not the original sparsity, decides the cost of a sparse factorisation.

## Where it comes from

Eliminating a pivot couples every pair of its remaining neighbours. A pivot with many neighbours therefore turns a sparse block dense.

## Counting it

The factors hold more entries than the original matrix. Fill is the factor count minus the original nonzero count, and it depends on the elimination order.

## Control

Ordering heuristics and symbolic factorisation predict fill before any floating-point work. A good order can change the cost by orders of magnitude.

To check your understanding, eliminate a pivot, count the factor entries and predict fill symbolically.

Fill is the hidden cost that separates a sparse method from a dense one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
