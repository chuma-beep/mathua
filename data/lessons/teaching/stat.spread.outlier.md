# Detecting outliers

An outlier is a value that lies far from the rest of the data. It may be a genuine extreme observation or a recording error, so detection is not the same as deletion.

## What an outlier is

There is no universal cutoff. A value can be unusual in one data set and ordinary in another. Always compare a suspect value with the spread and context of the other observations.

An outlier can change the mean and standard deviation substantially. The median and interquartile range are more resistant to extreme values.

## The interquartile range fence

Order the data and find the first quartile $Q_1$ and third quartile $Q_3$. The interquartile range is

$$IQR = Q_3 - Q_1.$$

The lower fence is $Q_1 - 1.5\,IQR$ and the upper fence is $Q_3 + 1.5\,IQR$. Values outside the fences are flagged as potential outliers.

## Deciding what to do

Do not delete a value just because it is outside a fence. First check the measurement and the context. A genuine extreme observation may be the most informative part of the data.

If a value is removed or corrected, report the decision and show the analysis with and without it when the conclusion changes.
