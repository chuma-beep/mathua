> Content sourced from [Algebrica](https://algebrica.org/harmonic-mean/) — CC BY-NC 4.0

## What is the harmonic mean?

The **harmonic mean** belongs to the broader family of power [means](<../introduction-to-the-mean/>) and plays a distinctive role whenever the data being analyzed combine reciprocally rather than additively. Instead of summing the observations, it takes the reciprocal of each value, computes their arithmetic mean, and then takes the reciprocal of that result.

This mean is particularly appropriate when averaging rates, ratios, or speeds, such as velocity, cost per unit, or productivity, cases where smaller values exert a stronger influence on the overall balance. While the [arithmetic mean](<../arithmetic-mean>) emphasizes larger [numbers](<../types-of-numbers/>), the harmonic mean highlights the contribution of smaller ones, providing a more accurate picture when the data vary inversely with respect to a fixed total.


In general form, the harmonic mean is expressed as:

\\[M_{-1} = \frac{n}{\displaystyle\sum_{i=1}^{n} \frac{1}{x_i}} \\]

where \\(x_1, x_2, \ldots, x_n\\) are the observed positive values and \\(n\\) is the total number of elements in the dataset.

  * The harmonic mean gives greater weight to smaller values, making it suitable for datasets based on rates or proportional quantities.
  * It can only be calculated for positive, non-zero values because it involves taking the reciprocal of each observation.
  * It is always less than or equal to the [geometric mean](<../geometric-mean>), which in turn is less than or equal to the arithmetic mean.
  * When data represent uniform measures of work or distance completed at varying speeds, the harmonic mean expresses the true average rate more accurately than other means.


##### The harmonic mean is often denoted as \\( M_{-1} \\) because it represents a specific case within the [Hölder mean family](<../introduction-to-the-mean>) (or power means), corresponding to the exponent ( s = -1 ).

## Example 1

To understand how the harmonic mean works in practice, let’s look at a simple situation involving average speed. Imagine a car that travels a road divided into two equal segments:

  * On the first half, the car moves at 60 km/h.
  * On the second half, it moves faster, at 90 km/h.


Even though the distance is the same, the time spent on each part of the trip is not. Because the slower speed takes more time, it has a greater influence on the overall average. That’s why using the arithmetic mean \\((75 \text{ km/h})\\) would give a misleading result, the correct approach is the harmonic mean.

Substituting the two speed values to the formula we obtain:

\\[M_{-1} = \frac{2}{\frac{1}{60} + \frac{1}{90}} = \frac{2}{\frac{5}{180}} = \frac{360}{5} = 72 \\]

##### The harmonic mean accurately represents the true average rate when distances are equal, because it reflects the additional time spent at lower speeds. Its formulation captures the reciprocal relationship between the variables, recognizing that time varies inversely with velocity. In essence, the harmonic mean describes balance within rate-based or proportional data, offering a precise and unbiased measure whenever the values being averaged represent performance, efficiency, or speed rather than direct quantities.

Therefore, the harmonic mean speed for the trip is:

\\[M_h = 72 \text{ km/h} \\]

## Example 2

Consider a scenario involving a machine that operates at different production rates over five equal time periods. Each period lasts the same amount of time, but the output rate, measured in units per minute, changes due to varying efficiency or workload conditions.

Period| Rate (units/minute)  
---|---  
1| 10  
2| 12  
3| 8  
4| 15  
5| 9  
  
Since each interval has the same duration, the correct way to find the overall average rate is through the harmonic mean, not the arithmetic one. This is because the slower periods have a stronger impact on the final result, reflecting the inverse relationship between time and rate.

Substituting the observed values we obtain:

\\[\begin{align} M_{-1} &= \frac{5}{\frac{1}{10} + \frac{1}{12} + \frac{1}{8} + \frac{1}{15} + \frac{1}{9}} \\\\[3pt] &= \frac{5}{0.1 + 0.0833 + 0.125 + 0.0667 + 0.1111} \\\\[8pt] &= \frac{5}{0.4861} \approx 10.29 \end{align} \\]

Hence, the harmonic mean rate of production is approximately:

\\[M_{-1} \approx 10.3 \text{ units per minute} \\]

Probability and Statistics

Probability and statistics provide a framework for analyzing uncertainty, modeling random phenomena.

1.5k

[Introduction to the Mean](https://algebrica.org/introduction-to-the-mean/)

1.2k

[Arithmetic Mean](https://algebrica.org/arithmetic-mean/)

1.6k

[Geometric Mean](https://algebrica.org/geometric-mean/)

1k

[Root Mean Square](https://algebrica.org/root-mean-square/)

1k

[Median and Quantiles](https://algebrica.org/median/)

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
