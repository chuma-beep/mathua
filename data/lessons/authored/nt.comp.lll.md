# Lattice reduction

The LLL algorithm finds a short basis for a lattice. Short vectors turn number-theoretic problems into linear algebra.

## The goal

A lattice is the integer span of a set of vectors. A reduced basis has short, nearly orthogonal vectors, which makes the lattice's structure visible.

## The approximation

LLL does not find the shortest vector in general. It guarantees a vector within a factor that grows exponentially in the dimension.

## Applications

Lattice reduction breaks knapsack cryptosystems and finds small solutions to integer relations. It is also a building block for stronger reduction algorithms.

To check your understanding, describe the lattice, state the factor and attack a knapsack.

LLL is an approximation algorithm, and its exponential factor is the price of staying polynomial in the dimension.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
