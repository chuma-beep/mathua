> Content sourced from [Algebrica](https://algebrica.org/median/) — CC BY-NC 4.0

## The central position in a data distribution

The median is a key measure of central tendency used to describe the typical value within a dataset. While the [mean](<../introduction-to-the-mean>) expresses the numerical balance of all values, the median focuses instead on position, it identifies the value that divides the ordered data into two equally sized groups. Half of the observations fall below the median, and half lie above it. Because it depends only on the relative ordering of the data, the median remains stable in the presence of extreme values. This makes it especially useful when the distribution is skewed or when outliers would otherwise distort the [arithmetic mean](<../arithmetic-mean>).


A formal definition of the median can be expressed as follows: Let \\( n \\) numerical observations be arranged in ascending order:

\\[x_{1} \leq x_{2} \leq \dots \leq x_{n} \\]

If \\( n \\) is odd, the median corresponds to the central element of the sequence:

\\[\tilde{M} = x_{\left(\frac{n+1}{2}\right)} \\]

If \\( n \\) is even, it is defined as the average of the two middle values:

\\[\tilde{M} = \frac{x_{\left(\frac{n}{2}\right)} + x_{\left(\frac{n}{2}+1\right)}}{2} \\]

##### This definition ensures that the number of data points smaller than the median equals the number that are larger, providing a perfectly balanced partition of the dataset.

## Example 1

To illustrate how the median is calculated when the number of observations is odd, let’s look at a simple example showing the monthly salaries of seven employees working in the same company:

Employee \\(x_i\\)| Salary ($)  
---|---  
\\(x_1\\)| 1,200  
\\(x_2\\)| 1,300  
\\(x_3\\)| 1,400  
\\(x_4\\)| 1,500  
\\(x_5\\)| 3,000  
\\(x_6\\)| 3,200  
\\(x_7\\)| 4,000  
  
Since there are \\( n = 7 \\) observations, the median is the value that occupies the central position:

\\[\tilde{M} = x_{\left(\frac{n+1}{2}\right)} = x_{4} \\]

The fourth value in the ordered list is 1,500, so the median is \\(\tilde{M} = 1,500\\).

##### This means that half of the employees earn less than $ 1,500, and the other half earn more.

## A distance-based definition of the median

The median can be described as the point that minimizes the total absolute deviation of all observations:

\\[\tilde{M} = \arg \min_{m} \sum_{i=1}^{n} |x_i - m| \\]

where \\(|x_i - m|\\) is the sum of absolute deviations. This property emphasizes that the median identifies the value that best captures the central location of a distribution when balance is defined in terms of distances instead of averages. Because it is less affected by extreme values, the median is often preferred when data are skewed or contain outliers.


Let’s apply this property to the data from Example 1, where the salaries (in dollars) are:

\\[1200,\ 1300,\ 1400,\ 1500,\ 3000,\ 3200,\ 4000 \\]

We compute the total absolute deviation for several possible values of \\( m \\).

For \\(m= 1200\\) we have \\(\sum_{i=1}^{n} |x_i - m|\\): \\[\begin{aligned} &= |1200 - 1200| + \\\ &\+ |1300 - 1200| \\\ &\+ |1400 - 1200| \\\ &\+ |1500 - 1200| \\\ &\+ |3000 - 1200| \\\ &\+ |3200 - 1200| \\\ &\+ |4000 - 1200| \\\ &= 7,100 \end{aligned} \\]

For \\(m= 1300\\) we have \\(\sum_{i=1}^{n} |x_i - m|\\): \\[\begin{aligned} &= |1200 - 1300| + \\\ &\+ |1300 - 1300| \\\ &\+ |1400 - 1300| \\\ &\+ |1500 - 1300| \\\ &\+ |3000 - 1300| \\\ &\+ |3200 - 1300| \\\ &\+ |4000 - 1300| \\\ &= 6,700 \end{aligned} \\]

Iterating the same procedure for each value of \\(m\\) gives the following total results.

\\(m\\)| \\(S(m) = \sum |x_i - m|\\)  
---|---  
1200| 7,100  
1300| 6,700  
1400| 6,400  
1500| 6,300  
3000| 7,800  
3200| 8,400  
4000| 12,400  
  
The median is therefore \\(1500\\), as this value minimizes the total sum of absolute deviations, with \\(S(m) = 6300\\) being the smallest among all those computed.

## Example 2

Let’s now consider another example, this time with an even number of observations. The dataset shows the selling prices of six houses recently sold in the same neighborhood:

House \\(x_i\\)| Price $  
---|---  
\\(x_1\\)| 180,000  
\\(x_2\\)| 190,000  
\\(x_3\\)| 200,000  
\\(x_4\\)| 220,000  
\\(x_5\\)| 300,000  
\\(x_6\\)| 450,000  
  
Since there are \\( n = 6 \\) observations, the median is given by the average of the two central values:

\\[\tilde{M} = \frac{x_{\left(\frac{n}{2}\right)} + x_{\left(\frac{n}{2}+1\right)}}{2} = \frac{x_{3} + x_{4}}{2} \\]

Substituting the corresponding values:

\\[\tilde{M} = \frac{200,000 + 220,000}{2} = 210,000 \\]

The median house price is $ 210,000.

##### This means that half of the houses were sold for less than $ 210,000 and the other half for more.

## Example 3

Let’s now consider a slightly more complex case, where data are organized into frequency classes. When data are grouped into classes, the median is obtained by interpolation within the median class, providing an estimated measure of the central tendency of the distribution. The following table shows the monthly household electricity consumption (in kWh) for a group of families, grouped into four classes:

Class interval (kWh)| Frequency (\\(f_i\\))  
---|---  
(100, 200]| 5  
(200, 300]| 8  
(300, 400]| 12  
(400, 500]| 5  
  
The total number of observations is:

\\[n = \sum f_i = 5 + 8 + 12 + 5 = 30 \\]

To find the median, we first determine the class that contains the middle observation. Since \\( n = 30 \\), the middle position is:

\\[\frac{n}{2} = \frac{30}{2} = 15 \\]

We now determine the cumulative frequencies, which represent the progressive total of observations up to each class. In this example, the cumulative frequency allows us to locate the class that contains the middle observation (that is, the median class) from which the value of the median will be computed.

Class interval (kWh)| Frequency (\\(f_i\\))| Cumulative frequency  
---|---|---  
(100, 200]| 5| 5  
(200, 300]| 8| 13  
(300, 400]| 12| 25  
(400, 500]| 5| 30  
  
The 15th observation falls into the class (300, 400], which is therefore the median class. The median for grouped data is given by:

\\[\tilde{M} = L + \left( \frac{\frac{n}{2} - F}{f_m} \right) \times c \\]

where:

  * \\(L\\) = lower boundary of the median class
  * \\(F\\) = cumulative frequency before the median class
  * \\(f_m\\) = frequency of the median class
  * \\(c\\) = class width


Substituting the values we obtain:

\\[L = 300, \quad F = 13, \quad f_m = 12, \quad c = 100 \\]

\\[\tilde{M} = 300 + \left( \frac{15 - 13}{12} \right) \times 100 = 300 + \left( \frac{2}{12} \times 100 \right) = 316.7 \\]

The median electricity consumption is approximately 317 kWh.

##### This means that half of the families consume less than 317 kWh per month, and the other half consume more.

## Median and quantiles

The median is a particular case of the **quantiles** , statistical measures that divide a distribution into equal parts. Specifically, the median corresponds to the quantile that separates the lower 50% of observations from the upper 50%, and it is therefore denoted as the quantile \\(x_{0.5}\\).

To introduce the concept of a quantile, we begin by considering a real number \\(p\\) such that \\(0 < p < 1\\). This parameter identifies the proportion of observations lying below a specific threshold within a distribution. The corresponding value, denoted by \\(x_p\\), is called the quantile of order \\(p\\) and represents the point that divides the data in such a way that a fraction \\(p\\) of the observations are less than or equal to it, and the remaining \\(1 - p\\) are greater.

From an analytical perspective, the quantile \\(x_p\\) can also be defined as the value that minimizes a global loss function obtained from the sum of asymmetric deviations:

\\[g(x_i, \bar{x}) = \begin{cases} (1 - p)\, |x_i - \bar{x}| & \text{if } x_i \le \bar{x} \\\\[6pt] p\, |x_i - \bar{x}| & \text{if } x_i > \bar{x} \end{cases} \\]

Commonly used quantiles include:

  * \\(x_{0.25}\\), first quartile, that identifies the value below which lies \\(1/4\\) of the distribution.
  * \\(x_{0.50}\\), median, that divides the data into two equal halves \\(1/2\\).
  * \\(x_{0.75}\\), third quartile, that marks the value below which lies \\(3/4\\) of the observations.


The **interquantile range** is a measure of statistical dispersion that expresses the spread of the central portion of a distribution. It is defined as the difference between two quantiles of order \\(p_1\\) and \\(p_2\\) \\((0 < p_1 < p_2 < 1)\\):

\\[IQR = x_{p_2} - x_{p_1} \\]

When the two quantiles correspond to the first and third quartiles, that is \\(p_1 = 0.25\\) and \\(p_2 = 0.75\\), the interquantile range becomes:

\\[IQR = x_{0.75} - x_{0.25} \\]

This interval contains the central 50% of the observations and is particularly useful for describing the variability of a dataset in a way that is robust to outliers, since it ignores the extreme values located in the tails of the distribution.

## Example 4

Let’s consider a small dataset representing the monthly salaries of ten employees in a company:

\\(i\\)| \\(x_i\\) (Salary in $)  
---|---  
1| 2,200  
2| 2,400  
3| 2,600  
4| 2,700  
5| 2,800  
6| 2,900  
7| 3,100  
8| 3,300  
9| 3,400  
10| 3,700  
  
Since there are \\(n = 10\\) observations, we can compute the positions of the first and third quartiles as:

\\[Q_1 = x_{0.25} = x_{(n+1)\times0.25} = x_{2.75} \\] \\[Q_3 = x_{0.75} = x_{(n+1)\times0.75} = x_{8.25} \\]

These positions correspond to fractional ranks, so we interpolate between the nearest observations.

\\(x_{2.75}\\) lies between the 2nd and 3rd values:  
\\[x_{0.25} = 2400 + 0.75 \times (2600 - 2400) = 2400 + 150 = 2550 \\]

\\(x_{8.25}\\) lies between the 8th and 9th values:  
\\[x_{0.75} = 3300 + 0.25 \times (3400 - 3300) = 3300 + 25 = 3325 \\]

We now calculate the interquantile range between the third and the first quantile. We obtain:

\\[IQR = x_{0.75} - x_{0.25} = 3325 - 2550 = 775 \\]

The interquantile range is 775 dollars.

##### This means that the central 50% of the employees have salaries that fall within a range of 775 dollars.

Probability and Statistics

Probability and statistics provide a framework for analyzing uncertainty, modeling random phenomena.

1.5k

[Introduction to the Mean](https://algebrica.org/introduction-to-the-mean/)

1.2k

[Arithmetic Mean](https://algebrica.org/arithmetic-mean/)

1.6k

[Geometric Mean](https://algebrica.org/geometric-mean/)

1.2k

[Harmonic Mean](https://algebrica.org/harmonic-mean/)

1k

[Root Mean Square](https://algebrica.org/root-mean-square/)

1.1k

[Variance](https://algebrica.org/variance/)

945

[Discrete Random Variables](https://algebrica.org/discrete-random-variables/)

1.2k

[Continuous Random Variables](https://algebrica.org/continuous-random-variables/)

1.3k

[Mean or Expected Value of a Random Variable](https://algebrica.org/mean-or-expected-value-of-a-random-variable/)

1.1k

[Variance and Covariance of a Random Variable](https://algebrica.org/variance-and-covariance-of-a-random-variable/)

1.3k

[Bernoulli Distribution](https://algebrica.org/bernoulli-distribution/)

1.4k

[Binomial Distribution](https://algebrica.org/binomial-distribution/)

1.5k

[Hypergeometric Distribution](https://algebrica.org/hypergeometric-distribution/)

1.8k

[Geometric Distribution](https://algebrica.org/geometric-distribution/)

1.4k

[Poisson Distribution](https://algebrica.org/poisson-distribution/)

1k

[Uniform Distribution](https://algebrica.org/uniform-distribution/)

1.7k

[Beta Distribution](https://algebrica.org/beta-distribution/)

6.4k

[Normal Distribution](https://algebrica.org/normal-distribution/)

1.3k

[Standard Normal Z Table](https://algebrica.org/standard-normal-z-table/)

1.3k

[Gamma Distribution](https://algebrica.org/gamma-distribution/)

1.4k

[Chi-square Distribution](https://algebrica.org/chi-square-distribution/)

1k

[Student’s t Distribution](https://algebrica.org/student-t-distribution/)

1.3k

[Exponential Distribution](https://algebrica.org/exponential-distribution/)

938

[Sampling Distributions](https://algebrica.org/sampling-distributions/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
