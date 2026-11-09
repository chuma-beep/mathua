# Coefficient of variation

The coefficient of variation expresses the standard deviation as a fraction of the mean. A standard deviation of 5 points means something different for a test scored out of 20 and for a salary measured in thousands, so dividing by the mean produces a unitless number that can be compared across measurements. The result is usually written as a percentage.

## Spread relative to the mean

For a sample with mean \(\bar{x}\) and standard deviation \(s\), the coefficient of variation is

\[
CV = \frac{s}{\bar{x}} \times 100\%.
\]

The mean must be positive for the ratio to be interpreted in the usual way, because a negative mean changes the sign without changing the spread. A coefficient of 20 percent says the typical distance from the mean is one fifth of the mean itself.

The measure is scale-free. Multiplying every observation by a constant multiplies both the mean and the standard deviation by that constant, so the ratio is unchanged. That is exactly what a comparison between centimetres and inches needs, and it is why the coefficient of variation appears whenever the units of two data sets differ.

## Recovering the standard deviation

The formula can be read in reverse. Given a mean and a coefficient of variation, the standard deviation is

\[
s = \frac{CV \cdot \bar{x}}{100}.
\]

If a sample has mean 80 and a coefficient of variation of 60 percent, then \(s = 60 \times 80 / 100 = 48\). The same rearrangement recovers the mean when the standard deviation and the coefficient are known, provided the coefficient is not zero.

This direction is common in reporting. A journal may publish a mean and a coefficient of variation because the original standard deviation was measured on a different scale. Recovering the standard deviation puts the result back into the units of the data.

## Comparing two samples

Two samples are compared by forming the ratio of their coefficients of variation. A ratio above one says the first sample is more variable relative to its own mean, and a ratio below one says the opposite.

The comparison is only meaningful when both means are positive and measured on scales that can be rescaled. A coefficient of variation for a temperature in degrees Celsius is not comparable with one for a temperature in kelvin, because the zero point matters. For ratio-scale measurements such as length, mass or money, the comparison is safe.

Reporting the ratio rather than the two coefficients also removes the percent sign and makes the conclusion a single number. A ratio of 1.4 says the first sample is 40 percent more variable than the second, relative to their own means.
