# Conditional Densities

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Conditioning on a continuous value

For jointly continuous variables, the conditional density of $X$ given $Y=y$ is

$$f_{X\mid Y}(x\mid y)=\frac{f_{X,Y}(x,y)}{f_Y(y)},$$

provided the marginal density $f_Y(y)$ is positive. The definition mirrors conditional probability, but the denominator is a density rather than a probability.

Because the conditioning event has probability zero, the conditional density is defined as a limit. The formula is the practical form: divide the joint density by the marginal density at the fixed value of $Y$.

## A conditional mean

The conditional expectation given $Y=y$ is

$$E[X\mid Y=y]=\int x\,f_{X\mid Y}(x\mid y)\,dx.$$

For a uniform conditional distribution on an interval, the mean is the midpoint. For example, if $X$ is uniform on the interval from zero to $y$, the conditional mean is $y/2$. The answer depends on $y$, so the conditional mean is a function of the conditioning value.

A conditional probability is found the same way, by integrating the conditional density over the event of interest.

## Independence through densities

Two continuous variables are independent exactly when their joint density factors:

$$f_{X,Y}(x,y)=f_X(x)f_Y(y)$$

for almost every pair. Equivalently, the conditional density of $X$ given $Y=y$ does not depend on $y$.

A rectangular support with a constant joint density is the simplest independent case. A triangular support usually signals dependence, because the allowable range of one variable changes with the other.

