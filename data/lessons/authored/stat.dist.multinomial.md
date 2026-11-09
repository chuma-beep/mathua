# Multinomial distribution

The multinomial distribution generalizes the binomial from two outcomes to several categories. A trial still has one outcome, but that outcome now falls into one of \(k\) categories with probabilities \(p_1, p_2, \dots, p_k\) that sum to one. After \(n\) independent trials the data are the counts in each category, and the distribution describes how those counts vary.

## Counting sequences

The number of distinct sequences with counts \(n_1, n_2, \dots, n_k\) is the multinomial coefficient

\[
\binom{n}{n_1, n_2, \dots, n_k} = \frac{n!}{n_1! n_2! \cdots n_k!}.
\]

The numerator counts every ordering of the \(n\) trials. The denominator removes the orderings inside each category, because two sequences that differ only by swapping two trials in the same category are the same sequence of outcomes.

For three categories with counts 3, 2 and 1, the coefficient is \(6!/(3!2!1!) = 60\). That number is the size of the sample space when only the counts are recorded. It is not a probability, and it can be much larger than the number of trials.

The coefficient also appears in the binomial as a special case. Setting \(k = 2\) and writing \(n_1 = j\), \(n_2 = n-j\) recovers the familiar \(\binom{n}{j}\).

## Expected counts

Each trial contributes one count to exactly one category, so the count \(N_i\) in category \(i\) has expectation

\[
E[N_i] = n p_i.
\]

Linearity of expectation does not require the categories to be independent, and it does not require the trials to be independent either. The expected count is simply the number of trials multiplied by the probability of landing in that category.

For a multinomial experiment with 200 trials and a category probability of 0.35, the expected count is \(200 \times 0.35 = 70\). The variance of that count is \(n p_i (1-p_i)\), the same as a binomial variance, because the question of whether a trial lands in category \(i\) is a Bernoulli question.

Expected counts are the bridge to the chi-square goodness-of-fit test. That test compares the observed count in each category with \(n p_i\), so the multinomial mean is the quantity that test is built around.

## The multinomial probability

The probability of a particular count vector is the coefficient multiplied by the probability of one sequence:

\[
P(N_1 = n_1, \dots, N_k = n_k) = \frac{n!}{n_1! \cdots n_k!} p_1^{n_1} \cdots p_k^{n_k}.
\]

Each sequence with the stated counts has probability \(p_1^{n_1} \cdots p_k^{n_k}\), and the coefficient counts how many such sequences exist. When the categories are equally likely, each sequence has probability \((1/k)^n\), so the count probability is the coefficient divided by \(k^n\).

With three equally likely categories and counts 2, 1 and 0 in three trials, the coefficient is \(3!/(2!1!0!) = 3\), and the probability is \(3/27 = 1/9\). Keeping the answer as a fraction shows the exact value before any decimal approximation.
