> Content sourced from [Algebrica](https://algebrica.org/standard-normal-z-table/) — CC BY-NC 4.0

## From the normal to the standard normal distribution

A generic [normal distribution](<../normal-distribution>) \\( \mathcal{N}(x; \mu, \sigma) \\) can always be transformed into its standardized form \\( N(x; 0, 1) \\) by introducing the standard variable \\( Z \\), defined as

\\[Z = \frac{X - \mu}{\sigma} \\]

This transformation, known as standardization, expresses each value of the [continuous random variable](<../continuous-random-variables/>) \\( X \\) in terms of the number of [standard deviations](<../variance-and-covariance-of-a-random-variable/>) it lies away from the [mean](<../mean-or-expected-value-of-a-random-variable/>) \\( \mu \\). As a result, the new variable \\( Z \\) follows a standard normal distribution with mean \\( 0 \\) and standard deviation \\( 1 \\).


The process of standardization is particularly useful because it places different normal distributions on a common scale. Once the variable \\( X \\) has been transformed into the standard variable \\( Z \\), probabilities and critical values can be derived directly from the standard normal distribution \\( N(x; 0, 1) \\).

This allows statistical problems to be solved using a single universal reference table (Z table), simplifying calculations and comparisons across different contexts.

## Standard Z Table

**The standard Z table** lists, for each cell, the cumulative probability to the left of a given \\(Z\\) value. In other words, it shows the probability that a standard normal variable takes on a value less than or equal to a specified Z-score. These values correspond to the area under the standard normal curve from negative infinity up to the chosen point on the horizontal axis, represented by the value of \\( z \\).

![](https://algebrica.org/wp-content/uploads/resources/images/normal-distribution-standard-1.png)

The complete Z table is extensive, as it contains values for a wide range of Z-scores. For brevity, only a portion of the table is shown below as an illustrative extract.

z| .00| .01| .02| .03| .04| .05| .06| .07| .08| …  
---|---|---|---|---|---|---|---|---|---|---  
0.0| .5000| .5040| .5080| .5120| .5160| .5199| .5239| .5279| .5319| …  
0.1| .5398| .5438| .5478| .5517| .5557| .5596| .5636| .5675| .5714| …  
0.2| .5793| .5832| .5871| .5910| .5948| .5987| .6026| .6064| .6103| …  
0.3| .6179| .6217| .6255| .6293| .6331| .6368| .6406| .6443| .6480| …  
0.4| .6554| .6591| .6628| .6664| .6700| .6736| .6772| .6808| .6844| …  
0.5| .6915| .6950| .6985| .7019| .7054| .7088| .7123| .7157| .7190| …  
0.6| .7257| .7291| .7324| .7357| .7389| .7422| .7454| .7486| .7517| …  
0.7| .7580| .7611| .7642| .7673| .7704| .7734| .7764| .7794| .7823| …  
0.8| .7881| .7910| .7939| .7967| .7995| .8023| .8051| .8078| .8106| …  
0.9| .8159| .8186| .8212| .8238| .8264| .8289| .8315| .8340| .8365| …  
1.0| .8413| .8438| .8461| .8485| .8508| .8531| .8554| .8577| .8599| …  
…| …| …| …| …| …| …| …| …| …| …  
  
##### Several online resources, both static and [interactive](https://ztable.io), allow users to compute cumulative probabilities for different values of \\( z \\).


Using the Z table is straightforward.

  * The row corresponding to the \\(Z\\) value contains the [integer](<../integers/>) part and the first decimal place
  * The column represents the remaining decimal places starting from the second one.
  * The cell at the intersection of the selected row and column gives the cumulative probability up to the specified value of \\( z \\).


## Example 1

Let us consider a simple example to illustrate how the Z table is used to determine the cumulative probability corresponding to a given value of the standardized variable \\( Z \\). Let us consider the case of a standardized variable \\( Z \\) such that:

\\[z = 0.16\\]

To find the cumulative probability to the left of the standardized variable \\( z \\), we locate the value 0.1 in the row and 0.06 in the column of the Z table. The intersection of these two entries gives the corresponding probability value, which in this case is \\( 0.5636 \\).

z| .00| .01| .02| .03| .04| .05| **.06**| .07| .08| …  
---|---|---|---|---|---|---|---|---|---|---  
0.0| .5000| .5040| .5080| .5120| .5160| .5199| .5239| .5279| .5319| …  
**0.1**| .5398| .5438| .5478| .5517| .5557| .5596| **.5636**| .5675| .5714| …  
…| …| …| …| …| …| …| …| …| …| …  
  

This means that approximately 56.36% of the observations in a standard normal distribution fall below this value of \\( z \\). We can therefore state that the probability of the standardized variable \\( Z \\) being less than \\( z = 0.16 \\) is equal to \\( 0.5636 \\), that is:

\\[P(Z < 0.16) = 0.5636 \\]

![](https://algebrica.org/wp-content/uploads/resources/images/normal-distribution-standard-2-1.png)

Since the total area under the standard normal curve equals \\(1\\), and the distribution is symmetric with respect to its [mean](<../mean-or-expected-value-of-a-random-variable/>), we can immediately deduce that

\\[P(Z > 0.16) = 1 - 0.5636 \\]

Moreover, the probability that \\( Z \\) lies between \\(0\\) and \\(0.16\\) is obtained by subtracting the cumulative probability up to \\( Z = 0 \\) from that up to \\( Z = 0.16 \\):

\\[P(0 < Z < 0.16) = P(Z < 0.16) - P(Z < 0) \\]

Substituting the corresponding values gives

\\[P(0 < Z < 0.16) = 0.5636 - 0.5 = 0.0636 \\]

Finally, due to the symmetry of the normal distribution about its mean, the probability that \\( Z \\) lies within the [interval](<../intervals/>) \\( -0.16 < Z < 0.16 \\) is twice the probability of being between \\(0\\) and \\(0.16\\):

\\[P(-0.16 < Z < 0.16) = 2 \times P(0 < Z < 0.16) = 0.1272 \\]

## Example 2

Suppose that the average lifetime of a rechargeable battery is \\( \mu = 8.0 \\) hours, with a standard deviation of \\( \sigma = 1.2 \\) hours, and that battery life follows a normal distribution. We want to calculate the probability that a randomly selected battery lasts less than \\( 6.5 \\) hours.


To find the probability \\( P(X < 6.5) \\), we need to determine the area under the normal curve to the left of \\( 6.5 \\). For this purpose, we apply the standardization formula:

\\[z = \frac{x - \mu}{\sigma} \\]

where \\( x \\) represents the observed value of the random variable, which in this case corresponds to 6.5 hours. We have:

\\[z = \frac{6.5 - 8.0}{1.2} = \frac{-1.5}{1.2} = -1.25 \\]

We can then use the Z table to find \\( P(Z < -1.25) \\), which gives the cumulative probability associated with this value of \\( z \\).

z| .00| .01| .02| .03| .04| **.05**| .06| .07| .08| …  
---|---|---|---|---|---|---|---|---|---|---  
-1.0| .1587| .1562| .1539| .1515| .1492| .1469| .1446| .1423| .1401| …  
**-1.2**| .1151| .1131| .1112| .1093| .1075| **.1056**| .1038| .1020| .1003| …  
-1.3| .0968| .0951| .0934| .0918| .0901| .0885| .0869| .0853| .0838| …  
…| …| …| …| …| …| …| …| …| …| …  
  

From the intersection of the selected row and column, we obtain a probability value equal to \\( 0.1056 \\). The figure below illustrates where this probability is located under the standard normal curve, corresponding to the cumulative area to the left of \\( z = -1.25 \\).

![](https://algebrica.org/wp-content/uploads/resources/images/normal-distribution-standard-4.png)

From the standard normal table, we obtain

\\[P(Z < -1.25) = 0.1056 \\]

##### This means that approximately 10.56% of all batteries are expected to last less than 6.5 hours.

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

1k

[Median and Quantiles](https://algebrica.org/median/)

1.1k

[Variance](https://algebrica.org/variance/)

946

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
