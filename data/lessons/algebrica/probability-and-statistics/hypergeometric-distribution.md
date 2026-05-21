> Content sourced from [Algebrica](https://algebrica.org/hypergeometric-distribution/) — CC BY-NC 4.0

## Introduction to the hypergeometric distribution

The hypergeometric distribution is a discrete probability distribution that describes the number of successes drawn from a finite population without replacement. Unlike the [binomial distribution](<../binomial-distribution/>), where each trial is independent and the probability of success stays constant, the hypergeometric setting involves draws that change the composition of the population at every step. As a consequence, the probability of selecting a success varies after each extraction.

To formalize this scenario, we consider a finite population divided into two categories: successes and failures. A sample of fixed size is taken without replacement, and the random variable represents the number of successes observed in that sample. This model relies on the following assumptions:

  * The population has a fixed size, with a known number of successes and failures.
  * A sample of predetermined size is drawn from the population.
  * The draws are made without replacement.
  * Each draw results in selecting either a success or a failure.
  * The [discrete random variable](<../discrete-random-variables/>) \\(X\\) counts how many successes are observed in the sample.


Formally, the hypergeometric distribution is defined as:

\\[P(X = x) = \frac{\binom{K}{x}\,\binom{N-K}{n-x}}{\binom{N}{n}} \\]

where:

  * \\(N\\) is the total size of the population.
  * \\(K\\) is the number of successes in the population.
  * \\(N - K\\) is the number of failures in the population.
  * \\(n\\) is the sample size drawn without replacement.
  * \\(x\\) is the number of observed successes.
  * \\(\binom{K}{x}\\) is the [binomial coefficient](<../binomial-coefficient/>), which counts the number of ways to choose \\(x\\) successes from the \\(K\\) available.
  * \\(\binom{N-K}{n-x}\\) counts the number of ways to choose the remaining items from the failures.
  * \\(\binom{N}{n}\\) represents the total number of distinct samples of size \\(n\\) that can be drawn from a population of size \\(N\\).


##### This distribution is used when independence does not hold and the probability of success changes after each draw. It provides a reliable model for sampling from finite populations, as in quality control, where items are inspected without replacement from a batch with known numbers of defective and non-defective units.

## Key features

  * \\[\text{1. } \quad P(X = x) = \frac{\binom{K}{x}\,\binom{N-K}{\,n-x,}}{\binom{N}{n}} \quad x = 0,1,\dots,n \\]

  * \\[\text{2. } \quad \mu = E(X) = n\,\frac{K}{N} \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = n\,\frac{K}{N}\left(1 - \frac{K}{N}\right)\frac{N-n}{N-1} \\]

  * \\[\text{4. } \quad \sigma = \sqrt{\,n\,\frac{K}{N}\left(1 - \frac{K}{N}\right)\frac{N-n}{N-1}\,} \\]


##### Each expression summarizes a fundamental aspect of the hypergeometric distribution, capturing how it models the number of successes drawn without replacement, where its expected value lies, and how its variability is shaped by the sample size and the finite nature of the population.

## Mean of the hypergeometric distribution

The [mean](<../introduction-to-the-mean/>), or [expected value](<../mean-or-expected-value-of-a-random-variable/>), of a hypergeometric distribution represents the average number of successes that can be expected when drawing a sample without replacement from a finite population. To compute the mean formally, we begin with the definition of the expected value:

\\[\mu = E(X) = \sum_{x=0}^{n} x \, P(X = x) \\]

Substituting the probability mass function of the hypergeometric distribution gives:

\\[E(X) = \sum_{x=0}^{n} x \,\frac{\binom{K}{x}\,\binom{N-K}{\,n-x\,}}{\binom{N}{n}} \\]

To simplify this expression, we use a combinatorial identity that connects two related binomial coefficients by reducing both the number of available successes and the sample size by one:

\\[x\,\binom{K}{x} = K\,\binom{K-1}{\,x-1\,} \\]

Applying this identity to the summation yields:

\\[E(X) = \frac{K}{\binom{N}{n}} \sum_{x=1}^{n} \binom{K-1}{\,x-1\,}\binom{N-K}{\,n-x\,} \\]

We now recognize that the summation corresponds to the total probability of a hypergeometric distribution with parameters \\(N-1\\), \\(K-1\\), and sample size \\(n-1\\). Therefore, the summation equals:

\\[\binom{N-1}{\,n-1\,} \\]

Substituting this into the expression above gives:

\\[E(X) = \frac{K}{\binom{N}{n}} \, \binom{N-1}{\,n-1\,} \\]

Using the identity:

\\[\frac{\binom{N-1}{\,n-1\,}}{\binom{N}{n}} = \frac{n}{N} \\]

we obtain the final expression for the mean:

\\[\mu = E(X) = n\,\frac{K}{N} \\]

##### This result shows that the mean of a hypergeometric distribution depends on the sample size \\(n\\) and on the proportion of successes in the population \\(K/N\\). On average, we expect to observe a fraction \\(K/N\\) of successes in any sample of size \\(n\\), even though the draws are made without replacement.

## Variance of the hypergeometric distribution

The [variance](<../variance-and-covariance-of-a-random-variable/>) of a hypergeometric distribution measures how much the number of observed successes is expected to vary around the mean value \\( \mu = n\,K/N \\). While the mean describes the central tendency of the distribution, the variance quantifies its spread, that is, how concentrated or dispersed the outcomes are when sampling without replacement from a finite population. Formally, the variance is defined as:

\\[\sigma^{2} = \mathrm{Var}(X) = E(X^{2}) - [E(X)]^{2} \\]

To compute it, we recall that the hypergeometric experiment consists of drawing \\(n\\) items without replacement from a finite population of size \\(N\\) containing \\(K\\) successes and \\(N-K\\) failures. Although the draws are not independent, the variance can be derived by considering indicator variables for each draw. Let \\(X\\) be the total number of successes in the sample, and let each draw be represented by an indicator variable:

\\[X = X_{1} + X_{2} + \cdots + X_{n} \\]

where \\(X_{i} = 1\\) if the \\(I\\)-th draw is a success and \\(X_{i} = 0\\) otherwise. Each indicator has expectation:

\\[E(X_{i}) = \frac{K}{N} \\]

and variance:

\\[\mathrm{Var}(X_{i}) = \frac{K}{N}\left(1 - \frac{K}{N}\right) \\]

However, because the sampling is done without replacement, each draw slightly changes the composition of the population. After a success is drawn, fewer successes remain, and after a failure is drawn, fewer failures remain. As a result, the draws influence one another, and the total variability is reduced compared with the binomial case. Taking this effect into account, the variance becomes:

\\[\mathrm{Var}(X) = n\,\frac{K}{N}\left(1 - \frac{K}{N}\right)\frac{N-n}{\,N-1\,} \\]

Therefore, the variance of the hypergeometric distribution is:

\\[\sigma^{2} = n\,\frac{K}{N}\left(1 - \frac{K}{N}\right)\frac{N-n}{\,N-1\,} \\]

##### This expression shows how the spread of the distribution depends not only on the proportion of successes \\(K/N\\), but also on the fact that sampling is done without replacement.

## Example 1

A batch contains 800 items, of which 12% are defective. An inspector selects a sample of 25 items for quality control. Determine the distribution of the random variable \\(X\\) that counts the number of defective items found in the sample. Although this may look similar to a model with independent trials, the situation is different: once an item is selected, it is not placed back into the batch. The probability of drawing a defective item changes after each draw because the composition of the batch changes. For this reason, the draws are not independent.

The problem can be handled using basic combinatorial reasoning. The number of possible samples of 25 items that can be drawn from a batch of 800 is:

\\[\binom{800}{25} \\]

In the original batch, there are \\(0.12 \times 800 = 96\\) defective items and \\(800 - 96 = 704\\) non-defective items. The probability that the sample contains exactly \\(x\\) defective items is:

\\[P(X = x) = \frac{ \binom{96}{x}\, \binom{704}{25 - x} }{ \binom{800}{25} } \\]

Thus, \\(X\\) follows a hypergeometric distribution with parameters \\(N = 800\\), \\(K = 96\\), and \\(n = 25\\).

## Comparison with the binomial distribution

The hypergeometric distribution is often compared to the [binomial distribution](<../binomial-distribution/>) because both describe the number of successes observed in a fixed number of trials. The essential difference lies in the sampling scheme.

  * The binomial distribution assumes independent trials with a constant probability of success, as if each draw were taken from an infinite population or as if the sampled item were replaced before the next draw.
  * The hypergeometric distribution, instead, models sampling without replacement from a finite population, so each draw slightly changes the composition of the remaining items. As a consequence, the probability of success varies across draws and the outcomes are not independent.


Despite these differences, the two distributions are closely related. When the population size \\(N\\) is large compared with the sample size \\(n\\), the effect of removing a few items becomes negligible. In this case, the hypergeometric distribution is well approximated by a binomial distribution with parameter \\(p = K/N\\):

\\[\text{Hyp}(N, K, n) \;\approx\; \text{Bin}\\!\left(n, \frac{K}{N}\right) \\]

##### This approximation highlights how the two models describe similar situations from different perspectives: the binomial focuses on idealized independent trials, while the hypergeometric captures the more realistic behavior of sampling from a finite population.


  * **Sampling**  
Hypergeometric: without replacement  
Binomial: with replacement or independent trials

  * **Probability of success**  
Hypergeometric: changes after each draw  
Binomial: remains constant

  * **Independence**  
Hypergeometric: dependent draws  
Binomial: independent trials

  * **Population**  
Hypergeometric: finite population explicitly matters  
Binomial: population treated as infinite or irrelevant

  * **When it is used**  
Hypergeometric: batch sampling, quality control, card problems  
Binomial: repeated [Bernoulli experiments](<../bernoulli-distribution/>)

  * **Conceptual idea**  
Hypergeometric: removing items alters future probabilities  
Binomial: each draw leaves probabilities unchanged
