> Content sourced from [Algebrica](https://algebrica.org/variance/) — CC BY-NC 4.0

## Introduction to mean deviations

One of the key principles in statistics is that the [mean](<../introduction-to-the-mean>) alone cannot fully describe a dataset. What truly matters is understanding how the individual observations are spread around the mean, that is, how far each value deviates from the central tendency. For this reason, it is essential to introduce the concept of **mean deviations** , expressed by the following general formula:

\\[^sS_{\bar{x}} = \left[ \frac{ \sum_{i=1}^{n} \left| x_i - \bar{x} \right|^s }{n} \right]^{1/s} \\]

  * The symbol \\(\bar{x}\\) represents the reference mean value around which deviations are measured.
  * The superscript \\(s\\) denotes the order of the powered mean being considered.


Within the family of mean deviations, the **quadratic mean deviation** (or root mean square deviation) measures the average magnitude of deviations from the mean. It is defined as the square root of the ratio between the sum of squared deviations (also known as deviance) and the total number of observations:

\\[^2S = \sqrt{ \frac{\sum_{i=1}^{n}(x_i - M)^2}{n} } \\]

Here, \\(x_i\\) represents the \\(i-th\\) observation in the dataset, while \\(M\\) denotes the reference mean value, typically the [arithmetic mean](<../arithmetic-mean>).

## Example 1

To illustrate how the quadratic mean deviation is computed, let us consider a simple dataset composed of five observed values. This measure will allow us to understand how far, on average, each observation lies from the mean when larger deviations are given proportionally greater weight. The dataset is as follows:

\\(X_i\\)| Observation  
---|---  
\\(x_{1}\\)| 3  
\\(x_{2}\\)| 5  
\\(x_{3}\\)| 7  
\\(x_{4}\\)| 10  
\\(x_{5}\\)| 15  
  
The [arithmetic mean](<../arithmetic-mean>) of the observations is:

\\[\bar{x} = \frac{3 + 5 + 7 + 10 + 15}{5} = \frac{40}{5} = 8 \\]


Next, we calculate the squared deviations of each observation from the mean:

\\(x_i\\)| \\(x_i - \bar{x}\\)| \\((x_i - \bar{x})^2\\)  
---|---|---  
3| -5| 25  
5| -3| 9  
7| -1| 1  
10| 2| 4  
15| 7| 49  
  
The sum of squared deviations is:

\\[\sum (x_i - \bar{x})^2 = 25 + 9 + 1 + 4 + 49 = 88 \\]


Dividing this value by the total number of observations and taking the square root, we obtain:

\\[^2S = \sqrt{ \frac{ \sum (x_i - \bar{x})^2 }{n} } = \sqrt{ \frac{88}{5} } = \sqrt{17.6} \approx 4.195 \\]

The quadratic mean deviation of the dataset is therefore \\(^2S \approx 4.20\\). This means that, on average, the observations differ from the mean by about 4.2 units.

## Variance

The square of the quadratic mean deviation is called the **variance**. It represents the average of the squared differences between each observation and a reference value \\(M\\), providing a measure of how widely the data are dispersed around that point. Unlike the quadratic mean deviation, which retains the same unit of measurement as the data, the variance is expressed in squared units. Its expression is:

\\[\sigma^2 = \frac{\sum_{i=1}^{n}(x_i - M)^2}{n} \\]

  * When \\(\sigma^2 = 0\\), all the observed values coincide with the reference value \\(M\\). In this case, there is no dispersion at all, since every observation is identical and perfectly aligned with the mean.
  * When \\(\sigma^2 = 1\\), the average squared deviation of the observations from the reference value (M) is exactly one unit. This means that, on average, each observation differs from (M) by one unit in squared terms. In practical terms, it indicates a moderate level of dispersion: the data are not identical, but their deviations from the mean are relatively small and balanced.
  * If a constant value is added to all observations, the variance remains unchanged.
  * If each observation is multiplied by a constant \\(a\\), the variance is multiplied by \\(a^2\\).
  * The larger the variance, the greater the spread of the data around \\(M\\), indicating higher heterogeneity within the dataset.


The variance can also be expressed in an alternative and more compact form by expanding the squared term in its definition. We start from the general formula:

\\[\sigma^2 = \frac{\sum_{i=1}^{n}(x_i - M)^2}{n} \\]

[Expanding the square](<../notable-products/>) gives:

\\[(x_i - M)^2 = x_i^2 - 2Mx_i + M^2 \\]

By substituting this expression into the original formula, we obtain:

\\[\sigma^2 = \frac{\sum_{i=1}^{n}x_i^2 - 2M\sum_{i=1}^{n}x_i + nM^2}{n} \\]

At this point, we recall that the mean \\(M\\) is defined as \\(M = \frac{\sum_{i=1}^{n}x_i}{n}\\). Replacing this relationship inside the equation allows us to simplify the expression step by step:

\\[\sigma^2 = \frac{\sum_{i=1}^{n}x_i^2}{n} - 2M^2 + M^2 \\]

Simplifying further, the two terms in \\(M^2\\) combine as follows:

\\[\sigma^2 = \frac{\sum_{i=1}^{n}x_i^2}{n} - M^2 \\]

We can therefore express the variance in a more compact and elegant way:

\\[\sigma^2 = M(x^2) - M^2 \\]

Here, \\(M(x^2)\\) represents the mean of the squared observations, while \\(M\\) denotes the mean of the original data. This identity shows that the variance can be computed simply as the mean of the squares minus the square of the mean.

## Example 2

Let us now compute the variance for the following dataset, which consists of five observed values:

\\(X_i\\)| Observation  
---|---  
\\(x_{(1)}\\)| 5  
\\(x_{(2)}\\)| 7  
\\(x_{(3)}\\)| 9  
\\(x_{(4)}\\)| 12  
\\(x_{(5)}\\)| 17  
  
The first step is to determine the mean of the observations. By summing all the values and dividing by the number of cases, we obtain:

\\[M = \frac{5 + 7 + 9 + 12 + 17}{5} = \frac{50}{5} = 10 \\]


Each observation can now be compared with the mean to find out how far it deviates from it.  
We then square these deviations to avoid sign cancellation and to give greater weight to larger differences:

\\(x_i\\)| \\(x_i - M\\)| \\((x_i - M)^2\\)  
---|---|---  
5| -5| 25  
7| -3| 9  
9| -1| 1  
12| 2| 4  
17| 7| 49  
  
Adding up all the squared deviations gives:

\\[\sum (x_i - M)^2 = 25 + 9 + 1 + 4 + 49 = 88 \\]


To obtain the variance, we divide this total by the number of observations:

\\[\sigma^2 = \frac{88}{5} = 17.6 \\]

The variance of this dataset is therefore equal to:

\\[\sigma^2 = 17.6 \\]

##### This result means that, on average, the squared distance of the observations from the mean is 17.6 units. If we take the square root of this value, we find the corresponding quadratic mean deviation \\(\sigma = \sqrt{17.6} \approx 4.20\\). In other words, the data values differ from the mean by about 4.2 units on average, indicating a moderate level of dispersion around the central value.

## Example 3

In some cases, the variance can be computed more efficiently by using an alternative expression that involves the mean of the squared values and the square of the mean. This formulation provides the same result as the standard definition but simplifies the calculation, especially when the necessary sums are already known. Let’s see how it works through a simple example. Consider the following dataset:

\\(i\\)| \\(x_i\\)  
---|---  
1| 3  
2| 9  
3| 11  
4| 14  
  
First, we calculate the mean of the observations:

\\[M = \frac{3 + 9 + 11 + 14}{4} = \frac{37}{4} = 9.25 \\]


Next, we find the mean of the squared values:

\\[\begin{align} M(x^2) &= \frac{3^2 + 9^2 + 11^2 + 14^2}{4} \\\\[3pt] &= \frac{9 + 81 + 121 + 196}{4} \\\\[3pt] &= \frac{407}{4} = 101.75 \end{align} \\]


Now we apply the simplified formula for the variance:

\\[\sigma^2 = M(x^2) - M^2 \\]

Substituting the values gives:

\\[\sigma^2 = 101.75 - (9.25)^2 = 101.75 - 85.56 = 16.19 \\]

Therefore, the variance of the dataset is \\(\sigma^2 = 16.19\\).

##### This approach shows that the variance can be derived directly from the mean of the squared values and the square of the mean, offering a more streamlined way to measure data dispersion when working with aggregated information.

## Variance of discrete and continuous random variables

In the case of [discrete random variables](<../discrete-random-variables>), the [variance](<../variance-and-covariance-of-a-random-variable/>) is defined as:

\\[\sigma^2 = E[(X - \mu)^2] = \sum_x (x - \mu)^2 f(x) \\]

where \\( x \\) denotes each possible value that the random variable \\( X\\) can assume, \\( \mu = E[X] \\) represents the [expected value](<../mean-or-expected-value-of-a-random-variable/>) or theoretical mean of the variable and \\( f(x) \\) is the probability mass function (PMF), which assigns a probability to every possible outcome of \\( X \\).


In the case of [continuous random variables](<../continuous-random-variables>), the variance is defined as:

\\[\sigma^2 = E[(X - \mu)^2] = \int_{-\infty}^{+\infty} (x - \mu)^2 f(x),dx \\]

Here, \\( f(x) \\) is the probability density function, which describes how the probability is distributed over the possible values of \\( X \\).

## Variance of a sampling distribution

The [sample variance](<../sampling-distributions/>) describes how the data in a sample are spread around the sample mean. A larger variance means that the observations are more widely dispersed, while a smaller variance indicates that they are closer to the mean. Formally, given \\( X_1, X_2, \ldots, X_n \\) random variables, the sample variance is defined as:

\\[S^2 = \frac{1}{n - 1} \sum_{i=1}^{n} (X_i - \overline{X})^2 \\]

where \\( S^2 \\) is the sample variance, \\( \overline{X} \\) is the sample mean, and \\( n \\) is the sample size.
