# Radial basis function kernels

The radial basis function kernel depends only on the distance between points. A length scale sets how quickly similarity decays, and the kernel is positive definite for every positive scale.

## The formula

The kernel is \(\exp(-\gamma\|x-y\|^2)\) for a positive \(\gamma\). At zero distance its value is one.

## The length scale

Writing \(\gamma=1/(2\sigma^2)\) makes \(\sigma\) the length scale. Doubling \(\sigma\) makes distant points more similar.

## The feature space

The kernel expands into an infinite-dimensional feature space. That expansion is why a finite data set can still be separated by a smooth boundary.

To check your understanding, evaluate at zero distance, convert gamma to sigma and name the dimension.

The RBF kernel measures similarity with a Gaussian bump around each point.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
