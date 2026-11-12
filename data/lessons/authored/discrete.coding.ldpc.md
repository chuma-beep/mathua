# LDPC codes

A low-density parity-check code has a sparse parity-check matrix. Sparsity makes iterative belief propagation cheap, and well-designed codes approach the Shannon limit.

## The parity-check matrix

Codewords are the vectors orthogonal to every row of a sparse matrix \(H\). Each row is a parity check on a small set of bits.

## Iterative decoding

Variable nodes and check nodes exchange messages along the edges of a bipartite graph. Each round improves the posterior estimates.

## The rate

For an \(m\times n\) check matrix of rank \(m\), the code rate is \((n-m)/n\). Sparsity keeps both memory and per-round work linear.

To check your understanding, check a codeword, send variable messages and subtract the rank.

LDPC codes trade algebraic structure for graph structure. Their performance comes from message passing on a sparse graph rather than from a closed-form decoder.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
