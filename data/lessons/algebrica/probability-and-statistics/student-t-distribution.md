> Content sourced from [Algebrica](https://algebrica.org/student-t-distribution/) — CC BY-NC 4.0

## Introduction to the student’s \\(t\\) distribution

In statistical inference, when the goal is to draw conclusions about the mean of a normally distributed population but the [variance](<../variance/>) is unknown and must be estimated from the sample, the [standard normal distribution](<../normal-distribution/>) \\(\mathcal{N}(x; 0, 1) \\) is no longer the correct reference model. In such situations, the **student’s \\(t\\) distribution** is used because it incorporates the additional variability introduced by estimating the population variance. The student’s \\(t\\) [continuous random variable](<../continuous-random-variables/>) is defined as:

\\[T = \frac{Z}{\sqrt{V/k}} \\]

  * \\( Z \\) is a standard normal random variable.
  * \\( V \\) is a [chi-square random variable](<../chi-square-distribution/>) with \\( k \\) degrees of freedom, independent of \\( Z \\).
  * \\( k \\) denotes the degrees of freedom and determines how heavy the tails of the distribution are.


For a student’s \\(t\\) random variable \\( T \\) with \\( k \\) degrees of freedom, the probability density function is:

\\[f(t,k) = \frac{\Gamma\left(\frac{k+1}{2}\right)}{\Gamma\left(\frac{k}{2}\right)\sqrt{k\pi}} \left(1 + \frac{t^{2}}{k}\right)^{-\frac{k+1}{2}} \\]

This expression shows how the shape of the distribution depends on the parameter \\( k \\). Smaller values of \\( k \\) produce heavier tails, reflecting the additional uncertainty associated with estimating the population variance. As \\( k \to \infty \\), the density converges to the standard normal distribution \\(\mathcal{N}(x; 0, 1). \\)

![](/diagrams/algebrica/student-t-distribution.png)

##### The plot visually highlights how the \\(t\\) curves gradually tighten around the center as the degrees of freedom grow, illustrating the smooth transition from a heavier-tailed distribution to the familiar bell shape of the normal model.

  * The total area under the curve equals \\(1\\). This means that the [integral](<../definite-integrals/>) of its probability density function over the entire real line, from \\(-\infty\\) to \\(+\infty\\), is equal to \\(1\\).

  * The curve is [symmetric](<../even-and-odd-functions/>) around the [mean](<../introduction-to-the-mean/>) \\(0\\). Because the student’s \\(t\\) distribution is centered at zero and is symmetric, half of the total probability lies on each side of the origin.

  * The curve has two [inflection points](<../maximum-minimum-and-inflection-points/>), whose location depends on the degrees of freedom \\(k\\). For small \\(k\\), the inflection points lie farther from the center, reflecting the heavier tails; as \\(k\\) increases, they move closer to \\(0\\), approaching those of the standard normal curve.

  * The curve is [asymptotic](<../algebrica.org/asymptotes/>) to the x-axis. As \\(t\\) moves farther away from the center, the probability density approaches (0), but more slowly than the normal distribution when \\(k\\) is small, due to its heavier tails.


## Key features

  * \\[\text{1. } \quad f(t; k) = \frac{\Gamma\left(\frac{k+1}{2}\right)} {\Gamma\left(\frac{k}{2}\right)\sqrt{k\pi}} \left(1 + \frac{t^{2}}{k}\right)^{-\frac{k+1}{2}} \\]

  * \\[\text{2. } \quad \mu = E(T) = 0 \quad \text{for } k > 1 \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(T) = \frac{k}{k-2} \quad \text{for } k > 2 \\]


##### Each expression summarizes a fundamental property of the student’s \\(t\\) distribution, whose center is fixed at zero while its spread and tail behavior depend entirely on the degrees of freedom \\(k\\).

The mean and variance of the student’s \\(t\\) distribution exist only when the integrals that define them converge. For \\(k \le 1\\) the integral for the mean diverges, and for \\(k \le 2\\) the integral for the variance diverges because the distribution has heavy tails, meaning that the density decreases so slowly that observations far from the center contribute enough to make the integral infinite.

## Symmetry of the student’s \\(t\\) distribution

The student’s \\(t\\) distribution is symmetric around zero, and this structural property has important consequences for how its critical values are interpreted. Because the left and right sides of the curve are mirror images of one another, every probability located on one tail corresponds to an equal probability on the opposite tail. This symmetry can be expressed through the relationship

\\[t_{1-\alpha} = -\,t_{\alpha} \\]

which states that the \\(t\\) value leaving an area of \\(1 - \alpha\\) in the right tail is simply the negative of the \\(t\\) value leaving an area of \\(\alpha\\) in that same tail. Equivalently, the [quantile](<../median/>) on the left that captures probability \\(\alpha\\) is positioned at the same height as its right-tail counterpart but reflected across the vertical axis.

![](/diagrams/algebrica/student-t-distribution-2.png)

Assuming that the area to the left of \\( t_{1-\alpha} \\) is roughly 12%, the area to the right of \\( t_{\alpha} \\) is approximately the same, about 12%. Taken together, the two tails account for about 24% of the total area under the curve.

##### Visually, the graph of the \\(t\\) density makes this behavior immediate: the curve descends from its peak at zero in a perfectly balanced way, and the shaded tail areas match in both shape and magnitude.


When computing areas under the student’s \\(t\\) distribution, the procedure is similar to that used for other continuous distributions. In the same way that the standard normal distribution relies on [z-tables](<../standard-normal-z-table/>) to report cumulative probabilities and critical values, the \\(t\\) distribution uses dedicated t-tables. These tables list the critical values associated with specific tail areas, but they also account for the degrees of freedom \\( k \\), which determine the exact shape of the distribution.

The idea is the same as in the normal case: z-tables allow you to find either the probability \\( P(Z < z) \\) or the z-value corresponding to a chosen area. T-tables serve the same purpose for the \\(t\\) distribution, with the only difference being that each row corresponds to a different value of \\( k \\). Because the distribution changes with \\( k \\), the critical values do as well.

##### In practice, both tables provide a convenient way to obtain probabilities and quantiles without performing the integral of the density function, with t-tables extending the logic of z-tables to situations where the variance must be estimated from the sample.

## Example 1

Suppose we want to determine the value of the \\(t\\) statistic with \\( k = 12 \\) degrees of freedom that leaves an area of \\( 0.01 \\) in the left tail of the distribution. Because the student’s \\(t\\) distribution is symmetric around zero, the \\(t\\) value that leaves an area of \\( 0.01 \\) in the left tail corresponds to the negative of the \\(t\\) value that leaves an area of \\( 0.01 \\) in the right tail. In terms of quantiles, this means:

\\[t_{0.99} = -\,t_{0.01} \\]

To find the needed value, we use te T-table and look up the row corresponding to \\( k = 12 \\) and the column for a tail probability of \\( 0.01 \\). The intersection in the is shown below:

\\(k\\)| 0.10| 0.05| 0.025| 0.01| …  
---|---|---|---|---|---  
10| 1.372| 1.812| 2.228| 2.764| …  
11| 1.363| 1.796| 2.201| 2.718| …  
12| 1.356| 1.782| 2.179| **2.681**|  …  
13| 1.350| 1.771| 2.160| 2.650| …  
…| …| …| …| …| …  
  
From the table we obtain:

\\[-t_{0.01} = -2.681 \\]

Thus, the \\(t\\) value that leaves 1% of the total probability on the left side of the distribution is: \\(t_{0.99} = 2.681\\).

## Note on tail probabilities

Since the student’s \\(t\\) distribution is symmetric around zero, the probability contained in the two tails is simply twice the probability in one tail. In the example 1, the right-tail area is \\(0.01\\), and the left-tail area is the same. Therefore, the total probability outside the [interval](<../intervals/>) \\([-t_{0.99},1, t_{0.99}]\\) is:

\\[2 \times 0.01 = 0.02 \\]

Consequently, the central area between \\(-t_{0.99}\\) and \\(t_{0.99}\\) equals:

\\[1 - 0.02 = 0.98 \\]

This illustrates how single-tail critical values allow you to determine both the tail probability and the corresponding central probability under the \\(t\\) distribution.
