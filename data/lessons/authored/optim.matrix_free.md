# Matrix-free methods

A matrix-free method never forms the matrix. It asks for products with vectors and builds the solution from those products.

## Products only

Conjugate gradients need one product per iteration and a few vector updates. The matrix can therefore be a function that computes the product on demand.

## Iterations

The method converges in at most the dimension of the problem in exact arithmetic. In practice the iteration count depends on the condition number.

## Memory

Only a handful of vectors are stored. A factorisation would store a matrix that may not fit, so the memory saving is the point.

To check your understanding, count products, bound exact iterations and count vectors.

Matrix-free methods replace storage with repeated multiplication.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
