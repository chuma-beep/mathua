> Content sourced from [Algebrica](https://algebrica.org/variance-and-covariance-of-a-random-variable/) — CC BY-NC 4.0

## What is the variance of a random variable?

In descriptive statistics, the [variance](<../variance>) expresses how much a set of values differs, on average, from its [mean](<../introduction-to-the-mean>). It is obtained by squaring each deviation from the mean and taking their average and is generally expressed as:

\\[\sigma^2 = \frac{1}{n}\sum_{i=1}^{n}(x_i - M)^2 \\]

An equivalent form emphasizes that variance can also be obtained through the relationship:

\\[\sigma^2 = M(x^2) - [M(x)]^2 \\]

This second formulation connects the average of squared values with the square of the average, offering a simpler way to compute dispersion when both quantities are known.


When we move to random variables, the concept of the variance of a random variable \\(X\\) naturally follows from the general definition introduced above, with a few important refinements. In general terms, the variance of a random variable \\(X\\) expresses how much the possible values of \\(X\\) are expected to deviate from their mean, that is, from their [expected value](<../mean-or-expected-value-of-a-random-variable>) \\(E[X]\\). It quantifies the overall spread of the distribution, providing a measure of how far the outcomes of the variable tend to lie from this theoretical center.

In the case of [discrete random variables](<../discrete-random-variables>), the variance is expressed as:

\\[\sigma^2 = E[(X - \mu)^2] = \sum_x (x - \mu)^2 f(x) \\]

  * \\(x\\) represents each possible value that the random variable \\(X\\) can take
  * \\(\mu = E[X]\\) is the [expected value](<../mean-or-expected-value-of-a-random-variable>) or theoretical mean of the variable
  * \\(f(x)\\) is the probability mass function, which assigns a probability to each possible outcome of \\(X\\).


In the case of [continuous random variables](<../continuous-random-variables>), the variance is given by the following [integral](<../definite-integrals/>):

\\[\sigma^2 = E[(X - \mu)^2] = \int_{-\infty}^{+\infty} (x - \mu)^2 f(x)\,dx \\]

Here, \\(f(x)\\) is the probability density function, which describes how the probability is distributed over the possible values of \\(X\\).

##### It should be recalled that when moving from variables that take discrete values to those defined over a continuous range, the summation is replaced by an integral.

## Deviation and Standard Deviation

The term \\(x - \mu\\) is called the deviation of an observation from its mean. It indicates how far a specific value \\(x\\) lies from the central reference point \\(\mu = E[X]\\), showing whether that value is above or below the mean and by how much.

The positive square root of the variance, denoted by \\(\sigma\\), is known as the standard deviation. It represents the typical amount by which the values of a random variable tend to differ from their mean, providing a more intuitive measure of spread expressed in the same units as the data.

## Example 1

Let us consider a random variable \\(X\\) representing the number of defective pieces found in a batch of ten components produced by a manufacturing line. Based on previous observations, the probability distribution of \\(X\\) is given as follows:

\\(x\\)| \\(0\\)| \\(1\\)| \\(2\\)| \\(3\\)| \\(4\\)  
---|---|---|---|---|---  
\\(f(x)\\)| \\(0.1\\)| \\(0.3\\)| \\(0.4\\)| \\(0.1\\)| \\(0.1\\)  
  
Let us now try to calculate the variance and the standard deviation of this random variable.


We start by finding the expected value, which represents the theoretical mean of the distribution. For a discrete random variable, the expected value is defined as:

\\[E[X] = \sum_x x f(x) \\]

By substituting the numerical values, we obtain:

\\[\begin{align} E[X] &= (0)(0.1) + (1)(0.3) + (2)(0.4) + (3)(0.1) + (4)(0.1) \\\\[6pt] &= 0 + 0.3 + 0.8 + 0.3 + 0.4 = 1.8 \end{align} \\]

Therefore, the mean value of (X) is \\(mu = 1.8\\) This means that, on average, about 1.8 defective pieces are expected in each production batch.


We can now compute the variance, which for a discrete random variable is defined as:

\\[\mathrm{Var}(X) = \sum_x (x - \mu)^2 f(x) \\]

Substituting the values from the table gives:

\\[\begin{align} \mathrm{Var}(X) &= (0 - 1.8)^2(0.1) + (1 - 1.8)^2(0.3) + (2 - 1.8)^2(0.4) + (3 - 1.8)^2(0.1) + (4 - 1.8)^2(0.1) \\\\[6pt] &= (3.24)(0.1) + (0.64)(0.3) + (0.04)(0.4) + (1.44)(0.1) + (4.84)(0.1)\\\\[6pt] &= 0.324 + 0.192 + 0.016 + 0.144 + 0.484 = 1.16 \end{align} \\]

Finally, the standard deviation is the positive square root of the variance:

\\[\sigma = \sqrt{\mathrm{Var}(X)} = \sqrt{1.16} \approx 1.08 \\]

In conclusion, the variance of \\(X\\) is \\( \sigma^2 = 1.16 \\), and the standard deviation is \\( \sigma \approx 1.08 \\).

##### This result shows that the number of defective pieces per batch typically deviates by about one unit from the average value of 1.8, indicating a moderate level of variability in the production process.

## Covariance of a random variable

When we want to examine how two random variables are related, we use the covariance. It quantifies the extent to which two variables, \\(X\\) and \\(Y\\), vary together. In the case of discrete random variables, the covariance is given by:

\\[\sigma_{XY} = E[(X - \mu_X)(Y - \mu_Y)] = \sum_x \sum_y (x - \mu_X)(y - \mu_Y) f(x, y) \\]

  * \\((x - \mu_X)\\) represents the deviation of each value of \\(X\\) from its mean.
  * \\((y - \mu_Y)\\) represents the deviation of each value of \\(Y\\) from its mean.
  * \\(f(x, y)\\) denotes the [joint probability distribution](<../discrete-random-variables/>) of the two variables, which assigns a probability to each possible pair of outcomes \\((x, y)\\).


In the case of continuous random variables, the covariance is defined as:

\\[\sigma_{XY} = E[(X - \mu_X)(Y - \mu_Y)] = \int_{-\infty}^{+\infty} \int_{-\infty}^{+\infty} (x - \mu_X)(y - \mu_Y) f(x, y), dx, dy \\]

In this expression, \\(f(x, y)\\) denotes the [joint probability density function](<../continuous-random-variables>) of the two variables.

## Correlation coefficient

To understand how strongly two random variables \\(X\\) and \\(Y\\) are linearly related, we use a measure known as the correlation coefficient. Unlike covariance, which only indicates whether two quantities move in the same or in opposite directions, the correlation coefficient provides a standardized measure that is independent of the scale of the variables.

It is obtained by dividing the covariance by the product of the standard deviations of \\(X\\) and \\(Y\\). This normalization step removes the influence of the measurement units and allows the relationship between the two variables to be expressed on a fixed numerical scale.

Formally, the correlation coefficient, denoted by \\(\rho_{XY}\\), is defined as:

\\[\rho_{XY} = \frac{\sigma_{XY}}{\sigma_X \sigma_Y} = \frac{\mathrm{Cov}(X, Y)}{\sqrt{\mathrm{Var}(X)} \sqrt{\mathrm{Var}(Y)}} \\]

  * \\(\sigma_{XY}\\) represents the covariance between \\(X\\) and \\(Y\\)
  * \\(\sigma_X\\), \\(\sigma_Y\\) are the respective standard deviations of the two variables.
  * The value of \\(\rho_{XY}\\) always lies between \\(−1\\) and \\(+1\\).


A value close to \\(1\\) indicates a strong positive linear relationship and both variables tend to increase together. A value close to \\(−1\\) indicates a strong negative relationship, so when one variable increases, the other tends to decrease. A value near \\(0\\) suggests that there is no linear association between them.

## Example 2

Let us consider a situation in which two random variables are strongly, but not perfectly, related. The variable \\(X\\) may represent the number of hours spent studying during a week, and \\(Y\\) the corresponding result obtained in a short assessment. Their joint probability distribution is defined as follows:

\\(X\\)| \\(Y\\)| \\(f(x, y)\\)  
---|---|---  
1| 2| 0.2  
2| 3| 0.3  
3| 5| 0.3  
4| 6| 0.2  
  
Let us compute the correlation coefficient.

##### The table shows a clear positive relationship: larger values of \\(X\\) are generally associated with higher values of \\(Y\\), though not in a perfectly proportional way.


We begin by calculating the expected values of the two variables, which are obtained as follows:

\\[\begin{align} E[X] &= (1)(0.2) + (2)(0.3) + (3)(0.3) + (4)(0.2) = 2.5 \\\\[6pt] E[Y] &= (2)(0.2) + (3)(0.3) + (5)(0.3) + (6)(0.2)\ = 4.0 \end{align} \\]

The expected value of the product \\(XY\\) is: \\[\begin{align} E[XY] &= (1 \cdot 2)(0.2) + (2 \cdot 3)(0.3) + (3 \cdot 5)(0.3) + (4 \cdot 6)(0.2) \\\\[6pt] &= 0.4 + 1.8 + 4.5 + 4.8 = 11.5 \end{align} \\]


We can now compute the covariance:

\\[\begin{align} \mathrm{Cov}(X, Y) &= E[XY] - E[X]E[Y] \\\\[6pt] &= 11.5 - (2.5)(4.0) = 11.5 - 10 = 1.5 \end{align} \\]


The variance of \\(X\\) is: \\[\mathrm{Var}(X) = E[X^2] - [E[X]]^2 = 7.3 - (2.5)^2 = 7.3 - 6.25 = 1.05 \\]

The variance of \\(Y\\) is: \\[\mathrm{Var}(Y) = E[Y^2] - [E[Y]]^2 = 18.2 - (4.0)^2 = 18.2 - 16 = 2.2 \\]


The correlation coefficient is obtained as: \\[\begin{align} \rho_{XY} &= \frac{\mathrm{Cov}(X, Y)}{\sqrt{\mathrm{Var}(X)} \sqrt{\mathrm{Var}(Y)}} \\\\[6pt] &= \frac{1.5}{\sqrt{1.05}\sqrt{2.2}} \\\\[6pt] &= \frac{1.5}{1.075 \times 1.483} = \frac{1.5}{1.594} \approx 0.94 \end{align} \\]

The correlation between \\(X\\) and \\(Y\\) is approximately \\(0.94\\). This value indicates a very strong positive linear relationship: when \\(X\\) increases, \\(Y\\) tends to increase almost proportionally, though not in a perfectly exact way.
