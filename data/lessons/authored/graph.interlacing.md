# Eigenvalue interlacing

Deleting a vertex interlaces the adjacency eigenvalues, which bounds how far one deletion can move the spectrum.

## Cauchy interlacing

If \(H\) is obtained from \(G\) by deleting a vertex, the adjacency eigenvalues interlace: \(\lambda_1(G)\ge\lambda_1(H)\ge\lambda_2(G)\ge\cdots\). Each eigenvalue of the subgraph is trapped between two consecutive eigenvalues of the larger graph.

## Repeated deletion

Deleting \(k\) vertices leaves the eigenvalues interlaced between the eigenvalues of the original graph shifted by at most \(k\). The \(i\)-th eigenvalue of the subgraph lies between \(\lambda_i(G)\) and \(\lambda_{i+k}(G)\).

## Consequences

Interlacing bounds the spectral radius of an induced subgraph by that of the whole graph. It also shows that the largest eigenvalue of a proper induced subgraph is strictly smaller when the graph is connected and the eigenvector has no zero entries.

To check your understanding, delete one vertex, delete \(k\) vertices and bound a subgraph radius.

A subgraph cannot have an eigenvalue outside the interval set by its parent, and the ordering is rigid.
