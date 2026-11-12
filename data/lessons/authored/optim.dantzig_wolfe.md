# Dantzig-Wolfe decomposition

Dantzig-Wolfe writes each block as a convex combination of its extreme points. The master then chooses the weights and the pricing problem is a block optimisation.

## Extreme points as columns

Any point of a bounded block is a convex combination of its extreme points. The master's variables are the weights of those combinations.

## The master and pricing

The master enforces the coupling through the weights. Pricing maximises the reduced profit over one block, which is a small linear program.

## Reading the bound

The restricted master bounds the full optimum from the appropriate side. The difference between that bound and the best known value is the gap.

To check your understanding, multiply the extreme points by their weights, solve the restricted master and subtract the known value from the bound.

The master sees only the weights. Block detail stays inside pricing, which is why the method handles huge block-structured models.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
