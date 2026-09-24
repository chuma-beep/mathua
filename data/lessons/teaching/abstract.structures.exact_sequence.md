# Exact Sequences

**Exact sequence:** A chain $\cdots \to A \to B \to C \to \cdots$ with image equals kernel at every step: everything arriving continues, nothing extra. Short exact $0 \to A \to B \to C \to 0$ says $A$ injects, $C$ is the cokernel, and $|B| = |A|\cdot|C|$ for finite groups.

## Worked: 0 → Z --2→ Z → Z/2 → 0

Verify exactness at each joint:
1. At the first $\mathbb Z$: kernel of $\times 2$ is 0, matching the image of $0 \to \mathbb Z$.
2. At the second $\mathbb Z$: image of $\times 2$ is $2\mathbb Z$; kernel of the projection is the evens $2\mathbb Z$. Match.
3. At $\mathbb Z/2$: the projection is onto, kernel of $\mathbb Z/2 \to 0$ is everything. Match.

So exactness is three local checks: kernels equal incoming images, joint by joint.

## Worked: orders multiply

For $0 \to A \to B \to C \to 0$ with $|A| = 2$, $|C| = 3$:
1. $A$ injects: 2 distinct elements land in $B$.
2. $C \cong B/\mathrm{im}(A)$: each coset has $|A| = 2$ elements, with $|C| = 3$ cosets.
3. So $|B| = 2 \cdot 3 = 6$: the middle order is forced by the ends.

So short exact sequences turn $B$ into $|A|$ copies of $C$ — counting does the rest.

## Splitting

If the sequence splits, $B \cong A \oplus C$: the extension is trivial. Non-split extensions (like $Z_4$ over $0 \to Z_2 \to Z_4 \to Z_2 \to 0$) are where the theory bites.
