# Correlation

Correlation summarises the direction and strength of a straight-line association between two quantitative variables.

## Association and sign

A positive association means that larger values of $x$ tend to go with larger values of $y$. A negative association means that larger values of $x$ tend to go with smaller values of $y$.

The correlation coefficient $r$ lies between -1 and 1. Its sign gives the direction. A value near 0 means little straight-line association.

A correlation of $r=0.9$ is a strong positive association. A correlation of $r=-0.9$ is a strong negative association. Both are equally strong, because strength depends on the distance from 0, not on the sign.

## The correlation coefficient

For paired data,

$$r = \frac{1}{n-1}\sum_{i=1}^{n}\left(\frac{x_i-\bar{x}}{s_x}\right)\left(\frac{y_i-\bar{y}}{s_y}\right).$$

Each value is converted to a standardised distance from its mean. The products are large and positive when the points lie in the same direction from their means. They are large and negative when one variable is above its mean while the other is below.

If every point lies exactly on a line with positive slope, then $r=1$. If every point lies exactly on a line with negative slope, then $r=-1$.

## Correlation is not causation

A strong correlation does not prove that one variable causes the other. Ice-cream sales and drowning deaths correlate because both rise in hot weather. The heat is a lurking variable.

Correlation also ignores curved relationships. A perfect parabola can have $r=0$ even though $y$ is completely determined by $x$. Always plot the data before trusting a single number.
