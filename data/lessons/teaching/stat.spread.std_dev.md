# Standard deviation

Standard deviation measures how far a data set typically sits from its mean. It is the square root of variance, so it is in the same units as the data.

## Standard deviation as a distance

For a population of $N$ values with mean $\mu$, the standard deviation is

$$\sigma = \sqrt{\frac{1}{N}\sum_{i=1}^{N}(x_i-\mu)^2}.$$

Each term is a squared distance from the mean. Dividing by $N$ averages those squared distances. The square root returns the answer to the original units.

A small standard deviation means the values cluster near the mean. A large one means they spread farther away.

## Computing population and sample spread

The population variance of $2, 2, 8, 8$ is 9, because each value is 3 from the mean 5. The population standard deviation is $\sqrt{9}=3$.

For a sample of $n$ values, the sample variance divides the squared deviations by $n-1$ rather than $n$. The sample standard deviation is the square root of that quantity.

The denominator $n-1$ corrects the tendency of a sample to underestimate the spread of the population it came from.

## Comparing distributions

Two data sets can have the same mean and different standard deviations. The set $48, 49, 50, 51, 52$ is tight around 50. The set $30, 40, 50, 60, 70$ has the same centre but a much larger standard deviation.

When you compare spread, compare standard deviations only when the units and the means are comparable. A standard deviation of 5 means something different for incomes in dollars than for heights in centimetres.
