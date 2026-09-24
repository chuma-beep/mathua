# Transformers and Attention

**Attention:** For queries $Q$, keys $K$, values $V$, output $\text{softmax}(QK^T/\sqrt{d_k})V$: score every query against every key, normalize to weights, return weighted values. Multi-head attention runs $h$ such heads in parallel and concatenates them.

## Worked: one attention computation

Take one query with scores $[3, 1, 0]$ against three keys, values $[10, 20, 30]$:
1. Exponentiate: $e^3 \approx 20.1$, $e^1 \approx 2.7$, $e^0 = 1$; total $\approx 23.8$.
2. Normalize: weights $\approx [0.84, 0.11, 0.04]$ — they sum to 1, spending most weight on the best match.
3. Mix values: $0.84\cdot 10 + 0.11\cdot 20 + 0.04\cdot 30 \approx 11.8$, pulled toward the top-matching value 10.

So attention is soft lookup: similarity scores become mixture weights, and the output is a weighted average of values.

## Worked: dimensions and cost

Take sequence length $n = 4$, model width 512, $h = 8$ heads:
1. The score matrix $QK^T$ is $4\times 4$: every query meets every key, so cost grows as $O(n^2)$.
2. Each head works in $512/8 = 64$ dimensions — the split keeps total compute flat while adding subspaces.
3. Without positional encoding the operation is permutation equivariant (reorder inputs, outputs reorder identically), so sinusoidal encodings add order back.

So the transformer's power and price share one source: all-pairs comparison, quadratic in length, order-free without encodings.

## When attention is too expensive

At $n = 100000$ the $n\times n$ matrix is infeasible — long documents and genomes need sparse, linear, or recurrent alternatives. The decoder's causal mask adds a further constraint: generation attends only to past tokens, staying autoregressive.
