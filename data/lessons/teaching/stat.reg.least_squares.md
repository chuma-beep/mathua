# The least-squares regression line

A regression line predicts a response $y$ from an explanatory variable $x$. The least-squares line is the line that makes the sum of squared vertical residuals as small as possible.

## The least-squares line

Write the fitted line as

$$\hat{y} = b_0 + b_1 x.$$

The slope is

$$b_1 = r\frac{s_y}{s_x},$$

and the intercept is

$$b_0 = \bar{y} - b_1\bar{x}.$$

The line always passes through the point $(\bar{x}, \bar{y})$. The slope says how many units $\hat{y}$ changes when $x$ increases by one unit.

## Slope and intercept from data

Suppose a line passes through $(2, 5)$ and $(6, 13)$. The slope is

$$b_1 = \frac{13-5}{6-2} = 2.$$

Using the first point, $5 = 2(2) + b_0$, so $b_0 = 1$. The fitted line is $\hat{y} = 2x + 1$.

With real data the points rarely lie exactly on a line. The least-squares formulas choose the line that balances the vertical distances.

## Prediction and residuals

A residual is

$$e = y - \hat{y}.$$

It is the vertical gap between an observed value and the prediction. A positive residual means the observation lies above the line.

Prediction is safest inside the range of the observed $x$ values. Extrapolating far beyond that range assumes the straight-line pattern continues, which the data may not support.
