# Box plots and the five-number summary

A box plot displays a data set through five numbers: the minimum, the lower quartile, the median, the upper quartile and the maximum.

## The five-number summary

Order the data first. The median splits the data into two halves. The lower quartile $Q_1$ is the median of the lower half. The upper quartile $Q_3$ is the median of the upper half.

For $2, 4, 6, 8, 10, 12, 14$, the median is 8. The lower half is $2, 4, 6$, whose median is 4. The upper half is $10, 12, 14$, whose median is 12. The five-number summary is $2, 4, 8, 12, 14$.

The box runs from $Q_1$ to $Q_3$. The line inside the box marks the median. The whiskers reach the smallest and largest values that are not outliers.

## Quartiles and the interquartile range

The interquartile range is

$$IQR = Q_3 - Q_1.$$

It measures the spread of the middle half of the data. In the example, $IQR = 12 - 4 = 8$.

The IQR is resistant to extreme values because it ignores the lowest quarter and the highest quarter of the data.

## Outlier fences

A value is a potential outlier when it lies below

$$Q_1 - 1.5\,IQR$$

or above

$$Q_3 + 1.5\,IQR.$$

For the example, the fences are $4 - 12 = -8$ and $12 + 12 = 24$. A value of 30 would lie above the upper fence and would be plotted as a separate point.

Always report the rule you used. Different software packages sometimes compute quartiles in slightly different ways.
