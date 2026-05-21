> Content sourced from [Algebrica](https://algebrica.org/bayes-theorem/) — CC BY-NC 4.0

## What it is and what Bayes’ Theorem is used for

Bayes’ Theorem is a fundamental result in probability theory that describes how to compute the conditional probability of a hypothesis given observed evidence. It provides a formal mechanism for updating prior beliefs in light of new data, by relating the posterior probability to the prior probability and the likelihood of the observed evidence.


In formal terms, given two events \\(A\\) and \\(B\\), Bayes’ Theorem states that the posterior probability of event \\(A\\) given event \\(B\\) is equal to the likelihood of event \\(B\\) given event \\(A\\) multiplied by the prior probability of event \\(A\\), all divided by the marginal probability (or evidence) of event \\(B\\).

\\[P(A|B) = \frac{P(B|A)P(A)}{P(B)}\\]

  * \\(P(A|B)\\), posterior probability: the probability of event \\(A\\) occurring given that event \\(B\\) has already occurred. It represents our updated belief about \\(A\\) after observing \\(B\\).

  * \\(P(B|A)\\), likelihood: the probability of observing event $B$ if event $A$ were true. It measures how compatible the observed data (\\(B\\)) is with the hypothesis (\\(A\\)).

  * \\(P(A)\\), prior probability, the initial probability of event \\(A\\) occurring before observing any evidence (\\(B\\)). It represents our initial belief about \\(A\\).

  * \\(P(B)\\), marginal probability, or evidence: the overall probability of event \\(B\\) occurring. It can be calculated as the sum (or integral in the continuous case) of the probabilities of \\(B\\) conditioned on all possible states of \\(A\\), weighted by their prior probabilities.


Formally, the marginal probability of an event \\(B\\) is calculated using the law of total probability, which expresses \\(P(B)\\) as the sum of the probabilities of \\(B\\) conditioned on all possible events of another complete and exclusive set of events (such as \\(A\\) and its complement \\(¬A\\), weighted by the probabilities of these events:

\\[P(B) = P(B|A)P(A) + P(B|\neg A)P(\neg A)\\]

##### This formula is crucial because it has significant practical implications for calculating the marginal probability of an event (B) when solving problems that utilize Bayes’ Theorem.

## How to derive Bayes’ Theorem

To derive the Bayes’ Theorem, we consider the joint probability of the two events \\(A\\) and \\(B\\), which can be expressed as:

\\[(A \cap B) = P(A|B)P(B)\\]

Since the intersection of two sets is commutative, the order does not change the result, we also have that the joint probability of \\(B\\) and \\(A\\) can be expressed as:

\\[(B \cap A) = P(B|A)P(A)\\]

From this, it follows that:

\\[(A \cap B) = P(A|B)P(B) = P(B|A)P(A) = (B \cap A)\\]

Therefore, we obtain:

\\[( P(A|B) = \frac{P(B|A)P(A)}{P(B)}\\]

## Example

An email filtering system tries to classify emails into two categories: spam and not spam. The filter uses the presence of certain keywords to make this decision. Let’s consider the word “discount”.

First, let’s define the events that occur in this problem:

  * \\(S\\): the email is spam.
  * \\(\neg S\\): the email is not spam.
  * \\(D\\): the email contains the word “discount”.


Now, let’s consider the following probabilities, assuming that the word “discount” is present in 95% of spam emails, while it is present in only 2% of non-spam emails.

  * Prior probability of spam: \\(P(S) = 0.40\\).
  * Prior probability of not spam: \\(P(\neg S) = 0.60\\).
  * Likelihood of “discount” given spam: \\(P(D|S) = 0.95\\).
  * Likelihood of “discount” given not spam: \\(P(D|\neg S) = 0.02\\).


We want to find the probability that the email is spam given that it contains the word “discount”, which is \\(P(S|D)\\).


Let’s apply Bayes’ Theorem to our problem, and we get

\\[P(S|D) = \frac{P(D|S)P(S)}{P(D)}\\]

We are missing the marginal probability of finding the word “discount” in any email, \\(P(D)\\). We can calculate it using the law of total probability:

\\[P(D) = P(D|S)P(S) + P(D|\neg S)P(\neg S)\\]

Calculating \\(P(D)\\) gives us:

\begin{align*} P(D) &= (0.95 \times 0.40) + (0.02 \times 0.60) \\\\[0.5em] P(D) &= 0.38 + 0.012 \\\\[0.5em] P(D) &= 0.392 \end{align*}

Substituting the value \\(P(B)\\) into Bayes’ Theorem formula, we get:

\begin{align*} P(S|D) &= \frac{0.95 \times 0.40}{0.392} \\\\[1em] P(S|D) &= \frac{0.38}{0.392} \\\\[1.5em] P(S|D) &\approx 0.969 \end{align*}

Therefore, it is concluded that given that the email contains the word “discount”, the probability that it is spam is approximately 96.94%.

## Glossary

  * Bayes’ Theorem: a fundamental result in probability theory that describes how to compute the conditional probability of a hypothesis given observed evidence, providing a framework for updating prior beliefs.

  * Posterior probability \\(P(A|B)\\): the updated probability of an event \\(A\\) occurring after observing new evidence \\(B\\).

  * Prior probability \\(P(A)\\): the initial probability of an event \\(A\\) occurring before any evidence is observed.

  * Likelihood \\(P(B|A)\\): the probability of observing the evidence \\(B\\) if the hypothesis \\(A\\) were true.

  * Marginal probability \\(P(B)\\): the overall probability of the evidence \\(B\\) occurring, regardless of the truth of the hypothesis \\(A\\). Also referred to as the evidence.

  * Conditional probability: the probability of an event occurring given that another event has already occurred.

  * Joint probability \\(P(A \cap B)\\): the probability that two events \\(A\\) and \\(B\\) both occur.
