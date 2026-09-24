# Order of an Element

**Order:** The order of $g$ is the smallest $n > 0$ with $g^n = e$ (or infinite if none exists). In additive notation, the smallest $n > 0$ with $ng = 0$.

## Worked: order of 2 in Z_6

Add 2 to itself mod 6 until hitting 0:
1. $2 + 2 = 4$, still nonzero.
2. $2 + 2 + 2 = 6 = 0$ mod 6.
3. Three copies suffice and fewer do not, so the order of 2 is 3.

So order is found by repeated operation, stopping at the first return to identity.

## Worked: a transposition has order 2

Apply $(1\,2)$ twice:
1. First application swaps 1 and 2.
2. Second application swaps them back: identity.
3. So $(1\,2)$ has order 2, while the 3-cycle $(1\,2\,3)$ needs three applications: order 3.

So in $S_n$, cycle type reads off the order: disjoint cycle lengths, take the lcm.

## Infinite order

In $\mathbb Z$ under addition, $-1$ never returns to 0: $n \cdot (-1) = 0$ only for $n = 0$. Elements of infinite order generate infinite cyclic subgroups.
