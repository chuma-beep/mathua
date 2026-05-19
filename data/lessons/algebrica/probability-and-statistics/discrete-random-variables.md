> Content sourced from [Algebrica](https://algebrica.org/discrete-random-variables/) — CC BY-NC 4.0

## Definition of a discrete random variable

A **discrete random variable** is a function that assigns a real number to each element of a discrete sample space. In other words, it maps the possible outcomes of a random experiment to numerical values that can be analyzed statistically. Formally, a discrete random variable is a function:

\\[X : \Omega \rightarrow \mathbb{R} \\]

where \\(\Omega\\) is a discrete sample space.

##### When the sample space is continuous, composed of infinitely many infinitesimally close outcomes, we speak of [continuous random variables](<../continuous-random-variables/>).


To illustrate the concept in a simple way, consider an experiment where a single die is rolled twice, and let the random variable \\(X\\) represent the number of sixes obtained. The possible values of \\(X\\) are 0, 1, and 2 where \\(0\\) means that no six appears in the two rolls, \\(1\\) means that exactly one six appears, and \\(2\\) means that both rolls show a six.

\\(x\\)| 0| 1| 2  
---|---|---|---  
\\[f(x)\\]| \\[\frac{25}{36}\\]| \\[\frac{10}{36}\\]| \\[\frac{1}{36}\\]  
  
where \\(x\\) represents the possible outcomes of the random variable \\(X\\) and \\(f(x)\\) represents the probability associated with each outcome. The probabilities satisfy:

\\[\sum f(x) = 1 \\]

This is consistent with the law of total probability, which states that the sum of the probabilities of all mutually exclusive outcomes of a random variable must equal \\(1\\). It ensures that the probability distribution accounts for every possible event in the experiment.


Since probability calculations can be tricky at first, the following shows how the values of \\(f(x)\\) for 0, 1, and 2 are obtained.

\\[\begin{aligned} P(X = 0) &= \left(\frac{5}{6}\right)^2 = \frac{25}{36} \\\\[1em] P(X = 1) &= 2 \cdot \frac{1}{6} \cdot \frac{5}{6} = \frac{10}{36} \\\\[1em] P(X = 2) &= \left(\frac{1}{6}\right)^2 = \frac{1}{36} \end{aligned} \\]

  * In the case of \\(x = 0\\), both dice show numbers other than six. Since the probability of not getting a six on a single roll is \\(\tfrac{5}{6}\\), the probability that this happens twice in a row is \\(\left(\tfrac{5}{6}\right)^2\\).
  * In the case of \\(x = 1\\), exactly one six appears in the two rolls. There are two possible ways this can happen: the first die shows a six and the second does not,  
or the first does not show a six and the second does. Each event has a probability of \\(\tfrac{1}{6} \cdot \tfrac{5}{6}\\), so the total probability is \\(2 \cdot \tfrac{1}{6} \cdot \tfrac{5}{6}\\).
  * Finally, in the case of \\(x = 2\\), both dice show a six. Since the probability of rolling a six on a single die is \\(\tfrac{1}{6}\\), the probability that this occurs twice in a row is \\(\left(\tfrac{1}{6}\right)^2\\).


## Discrete probability distribution

A discrete random variable has a certain probability of taking each of its possible values.  
This probability is described by a function \\(f(x)\\), called the probability mass function or **discrete probability distribution**. For such a distribution, the following conditions must hold:

\\[\begin{aligned} & f(x) \ge 0 \\\\[8pt] & \sum_x f(x) = 1 \\\\[5pt] & P(X = x) = f(x) \end{aligned} \\]

These conditions ensure that all probabilities are non-negative, that their total equals one,  
and that the probability of a specific value \\(x\\) is exactly given by its corresponding \\(f(x)\\).


When dealing with a discrete random variable \\(X\\), it is often useful to describe the probability that \\(X\\) takes a value up to a certain threshold (x). This leads to the definition of the **cumulative distribution function** , denoted by \\(F(x)\\):

\\[F(x) = P(X \le x) = \sum_{t \le x} f(t) \\]

The function \\(F(x)\\) expresses the total probability accumulated up to \\(x\\). It is defined for all real values of \\(x\\) and increases step by step as new probability mass is added. Being cumulative by nature, \\(F(x)\\) is always non-decreasing and never exceeds \\(1\\). To better illustrate the concept, let us return to the example of rolling two dice and show how the cumulative distribution function is constructed. The random variable \\(X\\) represents the number of sixes obtained. Its probability mass function is:

\\(x\\)| 0| 1| 2  
---|---|---|---  
\\[f(x)\\]| \\[\frac{25}{36}\\]| \\[\frac{10}{36}\\]| \\[\frac{1}{36}\\]  
  
The cumulative distribution function (F(x)) is obtained by adding the probabilities  
up to each value of \\(x\\):

\\(x\\)| 0| 1| 2  
---|---|---|---  
\\[F(x)\\]| \\[\frac{25}{36}\\]| \\[\frac{35}{36}\\]| \\[1\\]  
  
In fact, we have:

\\[\begin{aligned} F(0) &= P(X \le 0) = f(0) = \tfrac{25}{36} \\\\[6pt] F(1) &= P(X \le 1) = f(0) + f(1) = \tfrac{25}{36} + \tfrac{10}{36} = \tfrac{35}{36} \\\\[6pt] F(2) &= P(X \le 2) = f(0) + f(1) + f(2) = 1 \end{aligned} \\]

The function \\(F(x)\\) shows how probability accumulates as \\(x\\) increases. It starts at \\(\tfrac{25}{36}\\) when no sixes are obtained and reaches 1 when all possible outcomes have been included.

## Joint probability distributions

In cases where the sample space is multidimensional, meaning that each outcome depends on two or more random variables, the corresponding probabilities are described by **joint probability distributions** for discrete random variables. In some experiments, two discrete random variables can occur together, each taking specific values within the same outcome. The probability of this combined occurrence is described by a function \\(f(x, y)\\), which assigns a probability to every possible pair \\((x, y)\\). We have:

\\[f(x, y) = P(X = x; \quad Y = y) \\]

This function expresses how likely it is that \\(X\\) takes the value \\(x\\) while, at the same time, \\(Y\\) takes the value \\(y\\). For joint probability distributions, the following conditions must hold:

\\[\begin{aligned} & f(x, y) \ge 0 \quad \forall \ (x, y) \\\\[8pt] & \sum_x \sum_y f(x, y) = 1 \\\\[5pt] & P(X = x; \ Y = y) = f(x, y) \end{aligned} \\]

These conditions state that all probabilities are non-negative, that their total sum over all possible pairs \\((x, y)\\) equals one, and that each joint probability \\(P(X = x, Y = y)\\) is represented by the value of \\(f(x, y)\\).

## Example 1

To better illustrate the concept of a joint probability distribution for discrete random variables, consider the following simple example.Consider a small box containing 4 balls, 2 white and 2 black. Two balls are drawn at random without replacement. Let:

  * \\(X\\) = the number of black balls drawn
  * \\(Y\\) = the number of white balls drawn


The possible pairs \\((x, y)\\) represent all combinations of black and white balls that can be drawn. Since only two balls are extracted, \\(x + y = 2\\), and the possible pairs are:

\\[(0, 2),\ (1, 1),\ (2, 0) \\]

The joint probability distribution \\(f(x, y)\\) is given by:

\\[f(x, y) = \frac{\binom{2}{x}\binom{2}{y}}{\binom{4}{2}} \\]

By representing the values assumed by each pair \\((x, y)\\), we obtain the following table showing the joint probability distribution \\(f(x, y)\\):

\\[\begin{array}{c|ccc|c} f(x, y) & 0 & 1 & 2 & \text{Totals} \\\\[6pt] \hline 0 & \frac{0}{6} & \frac{0}{6} & \frac{1}{6} & \frac{1}{6} \\\\[6pt] 1 & \frac{0}{6} & \frac{4}{6} & \frac{0}{6} & \frac{4}{6} \\\\[6pt] 2 & \frac{1}{6} & \frac{0}{6} & \frac{0}{6} & \frac{1}{6} \\\\[6pt] \hline \text{Totals} & \tfrac{1}{6} & \tfrac{4}{6} & \tfrac{1}{6} & 1 \end{array} \\]

This example helps visualize how probabilities can be distributed across two discrete random variables. Each cell in the table represents the likelihood of a specific combination of black and white balls being drawn. By summing across rows and columns, we obtain the marginal probabilities of \\(X\\) and \\(Y\\), confirming that the total probability of all possible outcomes equals one.

##### It’s a simple yet effective way to understand how joint distributions organize and relate probabilities in a two-variable system.

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
