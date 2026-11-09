# Chernoff bounds

Chernoff bounds turn a moment generating function into an exponential tail bound. The exponent is then optimised to make the bound as small as possible.

## The moment generating function

The moment generating function of X is the expectation of e to the tX. For a Bernoulli variable with success probability p, it is one minus p plus p e to the t. Its derivatives at zero give the moments.

## Optimising the exponent

The Chernoff bound is e to the minus t a times the moment generating function at t. Differentiating the logarithm with respect to t gives the optimal exponent. For a one-sided event the boundary value often gives the best simple choice.

## The tail bound

Substituting the chosen exponent into the bound gives a numerical upper bound on the tail probability. The bound is exact as an inequality, though it is usually not tight. It is most useful when the moment generating function has a simple closed form.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
