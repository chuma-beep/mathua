# Weighted-sum scalarisation

A weighted sum turns several objectives into one by giving each a weight. Minimising the sum returns a point on the front.

## The scalarisation

The combined objective is the sum of each weight times its objective. Minimising it for positive weights returns a Pareto optimal point.

## Supported points

Every weighted-sum optimum is Pareto optimal, but not every Pareto point is the optimum of any weighted sum. Points hidden in a nonconvex part of the front are unsupported.

## Weights and trade-offs

Along a smooth front the ratio of two weights equals the negative slope between the objectives. Weights therefore encode a local exchange rate.

To check your understanding, form a weighted sum, recognise a supported point and read a slope.

A weighted sum is simple and fast, and it cannot see the nonconvex part of a front.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
