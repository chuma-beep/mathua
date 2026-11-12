# Sobolev spaces

A Sobolev space collects functions whose weak derivatives up to a given order lie in \(L^{2}\).

## The Sobolev norm

The norm on \(H^{k}\) is the square root of the sum of the squared \(L^{2}\) norms of all derivatives of order at most \(k\).

Order zero is plain \(L^{2}\). The space is complete in that norm. The count of terms grows with the number of variables. On the unit interval the constant function 1 has norm 1 at every order.

## Embedding into continuous functions

A function in \(H^{s}\) on \(n\)-dimensional space is continuous when \(s\) exceeds \(n/2\).

In one dimension \(H^{1}\) sits inside the continuous functions, since 1 exceeds 0.5. In three dimensions \(H^{2}\) is enough, since 2 exceeds 1.5, while \(H^{1}\) is not.

At the borderline the embedding fails, which is why the inequality is strict.

## A worked norm

Take \(u(x)=x\) on the interval from 0 to 1.

The squared \(L^{2}\) norm is the integral of \(x^{2}\), which is \(1/3\), while the derivative 1 contributes 1. The \(H^{1}\) norm is then \(2/\sqrt{3}\), about 1.1547.

For \(u(x)=x^{2}\) the contributions are \(1/5\) and \(4/3\), so the norm squared is \(23/15\), about 1.5333.

Higher order means more control. That control buys continuity.
