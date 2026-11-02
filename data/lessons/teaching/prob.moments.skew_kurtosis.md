# Skewness and Kurtosis

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Third and fourth moments

The mean and variance describe location and spread. Skewness uses the third central moment and kurtosis uses the fourth:

$$\mu_3=E[(X-\mu)^3],\qquad \mu_4=E[(X-\mu)^4].$$

Dividing by the correct power of the standard deviation makes each one scale-free. Standardised skewness is $\mu_3/\sigma^3$, and standardised kurtosis is $\mu_4/\sigma^4$.

A symmetric distribution has skewness zero. A long right tail gives positive skewness and a long left tail gives negative skewness.

## Reading skewness

When a distribution is right-skewed, the mean is usually larger than the median because the tail pulls the average upward. Income, city sizes and waiting times are common examples. When it is left-skewed, the mean is usually smaller than the median.

Skewness is a property of the distribution, not of a single sample. A sample estimate can be noisy, especially when the sample is small or the tail is heavy, so a histogram and a box plot should accompany any single number.

## Excess kurtosis

Kurtosis measures how much probability sits in the tails relative to a normal distribution. The normal distribution has standardised kurtosis 3, so excess kurtosis subtracts 3 and gives zero for a normal.

A positive excess kurtosis means heavier tails and a sharper peak than the normal. A negative value means lighter tails and a flatter shape. The uniform distribution has negative excess kurtosis, while the exponential distribution has positive excess kurtosis.

