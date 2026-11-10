# Spectral clustering

## The Fiedler vector

Spectral clustering uses the eigenvectors of a graph Laplacian. The Fiedler vector is the eigenvector for the second-smallest eigenvalue. Its entries give a one-dimensional embedding of the vertices.

The sign of each entry suggests a cluster. Entries with the same sign tend to belong together. The magnitude measures how strongly a vertex belongs to its side.

## Signs and clusters

A simple spectral clustering assigns negative entries to one cluster and positive entries to the other. The number of sign changes across the graph indicates how mixed the partition is.

The sign rule is a heuristic. For more than two clusters, several eigenvectors are used and the signs are combined.

## Normalized cuts

The normalized cut compares the edges crossing a partition with the total edge volume of each side. A small normalized cut indicates a good separation.

The Fiedler vector approximates the minimum normalized cut. The approximation is why spectral clustering is fast but not always exact.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
