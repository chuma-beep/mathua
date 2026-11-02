# The kernel trick

## From inner products to kernels

The dual form of the support vector machine depends on the training data only through inner products between pairs of points. A kernel replaces that inner product with a similarity function computed without explicitly building the transformed feature vector.

This allows a linear method to act in a high-dimensional feature space while the computational cost stays at the cost of evaluating the kernel on the original inputs.

## Common kernels

The polynomial kernel adds a constant and raises the inner product to a power. The radial basis function kernel depends on the squared distance between points and a bandwidth, giving large similarity to nearby points and small similarity to distant ones.

A valid kernel corresponds to an inner product in some feature space. The kernel matrix must be symmetric and positive semidefinite, which is the condition checked by Mercer's theorem.

## Kernel matrices

For n training points the kernel matrix has one entry for every ordered pair, so it contains n squared entries. Symmetry means only the diagonal and one triangle need to be computed, giving n times n plus one divided by two distinct values.

Kernel methods scale poorly when n is very large because the matrix grows quadratically, which is why approximations such as random features and the Nystrom method are used in practice.
