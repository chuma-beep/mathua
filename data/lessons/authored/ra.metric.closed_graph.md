# The closed graph theorem

The closed graph theorem gives continuity from a closedness condition on the graph. It applies to operators between Banach spaces.

## Closed graphs

The graph of \(T\) is the set of pairs \((x, Tx)\). It is closed when \(x_n \to x\) and \(T x_n \to y\) together force \(y = Tx\). For a bounded operator the graph is closed, since continuity gives \(T x_n \to Tx\).

## Continuity from closedness

If \(X\) and \(Y\) are Banach and the graph of a linear \(T\) is closed, then \(T\) is bounded. The graph is a Banach space under the norm \(\|x\| + \|Tx\|\), and the projection onto \(X\) is a bounded bijection. The open mapping theorem then bounds the other projection.

## Applications

Differentiation on \(C^1[0,1]\) with the supremum norm has a closed graph and is unbounded, since \(\|D f_n\|/\|f_n\| = n\) for \(f_n(t) = \sin(nt)/n\). The space \(C^1\) with the supremum norm is not complete, so the theorem does not apply. That failure is exactly the point of the example.

Closedness is cheap to check and continuity is expensive, and completeness is what makes the exchange legal.
