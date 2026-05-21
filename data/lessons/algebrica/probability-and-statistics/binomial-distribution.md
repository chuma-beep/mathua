> Content sourced from [Algebrica](https://algebrica.org/binomial-distribution/) — CC BY-NC 4.0

## Introduction to the binomial distribution

The **binomial distribution** is a discrete probability distribution that models the number of successes in a sequence of independent experiments, each one following a [Bernoulli distribution](<../bernoulli-distribution/>) with the same probability of success. In this sense, the binomial model is derived directly from repeated Bernoulli trials. In each trial, only two outcomes are possible: success or failure. The binomial model therefore quantifies, for every possible value of \\( x \\), the probability of observing exactly \\( x \\) successes in a total of \\( n \\) trials, assuming that the probability of success in each trial remains constant and equal to \\( p \\). This type of experiment must satisfy a series of properties:

  * There are only two possible outcomes for each trial: success, with probability \\( p \\), and failure, with probability \\( 1 - p \\).
  * The trials are independent.
  * The probability of success \\( p \\) remains constant across all trials.
  * The total number of trials \\( n \\) is fixed in advance.
  * The random variable \\( X \\) represents the number of successes obtained in the sequence of trials.
  * Each trial produces a single, mutually exclusive outcome, either success or failure.


Formally, the binomial distribution is expressed as

\\[P(X = x) = b(x; n, p) = \binom{n}{x} p^{x} q^{n - x} \\]

where:

  * \\( q = 1 - p \\).
  * \\( n \\) represents the total number of independent trials.
  * \\( x \\) is the number of observed successes.
  * \\( p \\) is the probability of success in each individual trial.
  * \\( q \\) is the probability of failure, equal to ( 1 - p ).
  * \\( \dbinom{n}{x} \\) is the [binomial coefficient](<../binomial-coefficient>), which counts the number of distinct ways to obtain \\( x \\) successes out of \\( n \\) trials.


Considering that \\( p + q = 1 \\), the binomial model satisfies the fundamental condition required of any probability distribution:

\\[\sum_{x = 0}^{n} b(x; n, p) = 1 \\]

This relationship ensures that the total probability across all possible outcomes is equal to one, meaning that the distribution fully describes every potential result of the \\( n \\) Bernoulli trials.

##### The binomial distribution is directly related to the [binomial theorem](<../binomial-theorem>), which states that the expansion of \\( (p + q)^n \\) yields all possible combinations of successes and failures across \\( n \\) trials. Each term in this expansion corresponds to a possible number of successes \\( x \\), with its associated probability given precisely by the binomial formula above.

## Key features

  * \\[\text{1. } \quad P(X = x) = b(x; n, p) = \binom{n}{x}\, p^{x}\, q^{\,n - x} \quad x = 0,1,\dots,n \\]

  * \\[\text{2. } \quad \mu = E(X) = np \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = np(1 - p) \\]

  * \\[\text{4. } \quad \sigma = \sqrt{\,np(1 - p)\,} \\]


##### Each expression highlights a key property of the binomial distribution, summarizing how it models the number of successes, where its average behavior lies, and how its variability increases with the number of trials.

## Mean of the binomial distribution

The [mean](<../introduction-to-the-mean/>), or [expected value](<../mean-or-expected-value-of-a-random-variable/>), of a binomial distribution represents the average number of successes that can be expected over a large number of identical experiments. To compute the mean formally, we start from the definition of the expected value:

\\[\mu = E(X) = \sum_{x=0}^{n} x \, P(X = x) = n p \\]

Substituting the probability mass function of the binomial distribution we have:

\\[E(X) = \sum_{x=0}^{n} x \, \binom{n}{x} p^{x} (1 - p)^{n - x} \\]

Using the identity, which connects two related binomial coefficients by reducing both \\( n \\) and \\( x \\) by one:

\\[x \binom{n}{x} = n \binom{n-1}{x-1} \\]

the expression becomes:

\\[E(X) = n p \sum_{x=1}^{n} \binom{n-1}{x-1} p^{x-1} (1 - p)^{(n-1)-(x-1)} \\]

The summation term equals 1, because it corresponds to the total probability of a binomial distribution with parameters \\( n - 1 \\) and \\( p \\). Therefore, we obtain:

\\[\mu = E(X) = n p \\]

This shows that the mean of a binomial distribution depends linearly on both the number of trials and the probability of success in each trial. On average, we expect \\( n p \\) of the \\( n \\) experiments to produce a successful outcome.


This result can also be derived by noting that the binomial random variable \\( X \\) can be expressed as the sum of ( n ) independent Bernoulli random variables \\( X_1, X_2, \ldots, X_n \\), each taking the value 1 (success) with probability \\( p \\) and 0 (failure) with probability \\( 1 - p \\):

\\[X = X_1 + X_2 + \cdots + X_n \\]

By the linearity of expectation:

\\[E(X) = E(X_1) + E(X_2) + \cdots + E(X_n) = n p \\]

## Variance of the binomial distribution

The [variance](<../variance-and-covariance-of-a-random-variable/>) of a binomial distribution measures how much the number of observed successes is expected to vary around the mean value \\( \mu = n p \\). While the mean describes the central tendency of the distribution, the variance quantifies its spread that is, how concentrated or dispersed the outcomes are across repeated experiments. Formally, the variance is defined as:

\\[\sigma^2 = \mathrm{Var}(X) = E(X^2) - [E(X)]^2 = n p q \\]

To compute it, we recall that the binomial variable \\( X \\) can be expressed as the sum of \\( n \\) independent Bernoulli variables \\( X_1, X_2, \ldots, X_n \\), where each trial has probability of success \\( p \\):

\\[X = X_1 + X_2 + \cdots + X_n \\]

Since the variance of a Bernoulli variable is \\( \mathrm{Var}(X_i) = p(1 - p) \\), and the trials are independent, the variance of their sum is simply the sum of the individual variances:

\\[\mathrm{Var}(X) = \mathrm{Var}(X_1) + \mathrm{Var}(X_2) + \cdots + \mathrm{Var}(X_n) \\]

Therefore, the variance of the binomial distribution is:

\\[\sigma^2 = n p q \\]

##### This result shows that the variability of the distribution increases linearly with the number of trials \\( n \\), and depends on both the probability of success \\( p \\) and the probability of failure \\( (1 - p) \\).

## Example 1

Consider a factory that produces electronic sensors. Each sensor is tested to verify whether it operates correctly under specific temperature conditions. Suppose 6 sensors are selected at random from a production batch, and the probability that a single sensor passes the test is ( p = 0.7 ). We want to find the probability that exactly 4 of the 6 sensors will function properly during the test.


Assuming that each test is independent, the probability is given by:

\\[\begin{align} b(4; 6, 0.7) &= \binom{6}{4} (0.7)^4 (0.3)^2 \\\\[5pt] &= \frac{6!}{4! \, 2!} (0.7)^4 (0.3)^2 \\\\[10pt] &= 15 \times 0.2401 \times 0.09 = 0.324135 \end{align} \\]

##### The [factorial](<../factorial>) symbol \\((!)\\) indicates the product of all positive integers up to a given number.

Therefore, the probability that exactly four sensors out of six will pass the test is approximately \\( 0.324 \\), or 32.4%.

## Cumulative binomial distribution

In some contexts, the goal is not to find the probability of getting exactly \\( x \\) successes, but rather the probability of getting no more than a certain number of successes \\( r \\) in \\( n \\) Bernoulli trials. This type of probability is obtained by summing all individual terms of the binomial distribution from \\( x = 0 \\) up to \\( x = r \\):

\\[B(r; n, p) = \sum_{x = 0}^{r} b(x; n, p) \\]

where \\( b(x; n, p) \\) represents the probability mass function of the binomial distribution. The values of the cumulative binomial distribution are often provided in dedicated tables, similar in purpose to the [standard normal Z table](<../standard-normal-z-table/>). Each entry in these tables corresponds to the cumulative probability of obtaining up to a given number of successes, for specific combinations of \\( n \\) and \\( p \\).

In more formal terms, the cumulative probability of the binomial distribution can be expressed as:

\\[P[X \le c] = \sum_{x = 0}^{c} \binom{n}{x} p^{x} (1 - p)^{n - x} \\]

This probability is often computed using cumulative binomial tables, which list pre-calculated values of \\( P[X \le c] \\) for selected values of \\( n \\) and \\( p \\). A simplified portion of such a table is shown below.

n| c| p = 0.05| p = 0.10| p = 0.20| p = 0.30| p = 0.40| p = 0.50| …  
---|---|---|---|---|---|---|---|---  
1| 0| 0.950| 0.900| 0.800| 0.700| 0.600| 0.500| …  
1| 1| 1.000| 1.000| 1.000| 1.000| 1.000| 1.000| …  
2| 0| 0.903| 0.810| 0.640| 0.490| 0.360| 0.250| …  
2| 1| 0.998| 0.990| 0.960| 0.910| 0.840| 0.750| …  
2| 2| 1.000| 1.000| 1.000| 1.000| 1.000| 1.000| …  
…| …| …| …| …| …| …| …| …  
  
These tables allow users to quickly find cumulative probabilities without computing each term of the summation, in the same way that the [standard normal Z table](<../standard-normal-z-table/>) is used to determine cumulative probabilities under the normal distribution. The intersection between a row and a column gives the cumulative probability \\( P[X \le c] \\), that is, the probability of obtaining up to \\( c \\) successes for the corresponding value of \\( n \\) and \\( p \\).

##### Such tabulated values allow for quick estimation of binomial probabilities without performing the full summation manually, just as the Z table simplifies the computation of probabilities under the [normal distribution](<../normal-distribution>).

## Normal approximation of the binomial distribution

In certain cases, the binomial distribution, which is inherently discrete, can be effectively approximated by a continuous [normal distribution](<../normal-distribution>). This approximation is appropriate when the binomial distribution exhibits a roughly symmetric and bell-shaped form, closely resembling the profile of a normal curve. Such situations typically occur when the number of trials \\( n \\) is large and the probability of success \\( p \\) is not too close to 0 or 1.

In these conditions, the random variable \\( X \\), distributed as \\( \mathrm{Bin}(n, p) \\), can be approximated by a normal distribution with the same mean and variance:

\\[X \approx \mathcal{N}(x; n p, n p q) \\]

This method simplifies calculations and allows binomial probabilities to be estimated using the tools and [Z tables](https://algebrica.org/standard-normal-z-table/) associated with the normal distribution, providing accurate results in most practical applications.


If \\( X \\) is a binomial random variable with mean \\( \mu = n p \\) and variance \\( \sigma^2 = n p q \\), then the limiting distribution of the standardized variable \\( Z \\), defined as:

\\[\frac{X - n p}{\sqrt{n p q}} \xrightarrow{d} \mathcal{N}(x; 0, 1) \quad \text{as } n \to \infty \\]

is the [standard normal distribution](<../normal-distribution/>):

\\[Z \sim \mathcal{N}(x; 0, 1) \\]

##### The notation \\( \xrightarrow{d} \\) denotes convergence in distribution, meaning that the probability distribution of the standardized variable gradually approaches the standard normal distribution as \\( n \\) increases.

## Example 2

To better illustrate the normal approximation to the binomial distribution, consider a company that produces small electronic components. During quality control, each item is tested to check whether it meets the required electrical specifications. Suppose that 100 components are tested and that the probability a single component passes the test is \\( p = 0.9 \\). We want to find the probability that between 85 and 92 components pass the inspection successfully.


Let \\( X \\) be the discrete random variable representing the number of components that pass the test. Then the probability of interest can be written as:

\\[P(85 \le X \le 92) = \sum_{x=85}^{92} b(x; 100, 0.9) \\]

Since \\( n \\) is large and \\( p \\) is not too close to 0 or 1, we can use the normal approximation with:

\\[\mu = n p = 100 \times 0.9 = 90 \\] \\[\sigma = \sqrt{n p q} = \sqrt{100 \times 0.9 \times 0.1} = 3 \\]

Because the binomial variable \\( X \\) is discrete (it only takes [integer](<../integers/>) values) and the normal distribution is continuous, the probability of \\( X \\) lying between 85 and 92 is better approximated by the area under the normal curve from half a unit below 85 to half a unit above 92. For this reason, we replace the discrete bounds with

\\[x_1 = 84.5 \quad \text{and} \quad x_2 = 92.5 \\]


We then compute the corresponding standardized values using the [transformation](<../normal-distribution>):

\\[z = \frac{x - \mu}{\sigma} \\]

which gives:

\\[z_1 = \frac{84.5 - 90}{3} = -1.83, \quad z_2 = \frac{92.5 - 90}{3} = 0.83 \\]

We can write:

\\[P(85 \le X \le 92) = \sum_{x=85}^{92} b(x; 100, 0.9) \approx P(-1.83 \le Z \le 0.83) \\]


According to the standard normal distribution, we have:

\\[P(-1.83 \le Z \le 0.83) = P(Z \le 0.83) - P(Z \le -1.83) \\]

From the [standard normal Z table](<../standard-normal-z-table/>):

\\[P(Z \le 0.83) = 0.7967, \quad P(Z \le -1.83) = 0.0336 \\]


Thus we obtain:

\\[P(-1.83 \le Z \le 0.83) = 0.7967 - 0.0336 = 0.7631 \\]

Therefore, the probability that between 85 and 92 components pass the quality test is approximately 0.763, or 76.3%:

## Comparison between the binomial and hypergeometric distributions

The binomial distribution applies when each trial has the same probability of success and when one trial does not influence the next. This setting is appropriate for sampling with replacement or for situations in which the population is large enough that removing a single item does not change its overall composition. In many real-world problems, however (such as quality control), the selected items are not replaced. In these cases, the composition of the population changes after each draw, and the binomial assumptions are no longer valid. The correct model becomes the [hypergeometric distribution](<../hypergeometric-distribution>).

When the population size \\(N\\) is not large compared with the sample size \\(n\\), the probability of success can no longer be treated as constant across trials. The binomial model must then be replaced by its finite-population counterpart:

\\[X \sim \text{Hyp}(N, K, n) \\]

where \\(K\\) is the number of successes in the population. Moving from the binomial to the hypergeometric distribution reflects the shift from independent trials with fixed probabilities to the more realistic setting of sampling without replacement.

## Connection between the binomial and the Poisson distribution

When the number of trials \\( n \\) is large and the probability of success \\( p \\) is small, while keeping \\( n p = \lambda \\) constant, the binomial distribution converges to the [Poisson Distribution](<../poisson-distribution/>):

\\[P(X = x) = \frac{e^{-\lambda}\lambda^x}{x!} \\]

##### The binomial distribution thus serves as the foundation from which the Poisson model arises in the limit, offering a simplified representation of rare and independent events occurring over time or space.
