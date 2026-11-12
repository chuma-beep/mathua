# Portfolio selection as a cone program

A portfolio problem trades expected return against a risk measured by a norm. When the risk is a standard deviation the model is a second-order cone program.

## Risk as a norm

The covariance matrix supplies a square root, and the risk of a weight vector is the norm of its transformed version. Fixing a risk budget therefore fixes a norm bound.

## The return constraint

The expected return is a linear function of the weights. Requiring a minimum return adds one linear row to the same cone program.

## The efficient frontier

Sweeping the risk budget traces the frontier. Each point on it is the best return available at that risk, which is what makes the curve useful.

To check your understanding, multiply the weights by the covariance square root, multiply the weights by expected returns and compare risk at equal return.

The weights must sum to one, so the model carries an equality row as well as the cone. The frontier starts at the minimum-risk portfolio and rises from there.
