# Independent events

Two events are independent when knowing that one occurred does not change the probability of the other. Independence is a statement about a probability model, not about whether the events look unrelated in a story. The definition is short, but it has several equivalent forms, and each one is useful in a different question.

## The product rule

Events \(A\) and \(B\) are independent exactly when

\[
P(A \text{ and } B) = P(A)P(B).
\]

The product rule is the definition. If the joint probability equals the product of the marginal probabilities, then the occurrence of one event carries no information about the other. If the joint probability is larger than the product, the events tend to occur together. If it is smaller, they tend to exclude one another.

A common mistake is to check the product rule with rounded decimals. The comparison should be made with exact fractions or with integers cross-multiplied, because two probabilities that differ by a rounding error can look equal. The generator prints fractions for that reason.

## Conditioning does not change the probability

The conditional probability of \(A\) given \(B\) is

\[
P(A \mid B) = \frac{P(A \text{ and } B)}{P(B)}.
\]

When \(A\) and \(B\) are independent, substituting the product rule gives \(P(A \mid B) = P(A)\). This is the interpretation most people carry: the conditioning event does not move the probability.

The equivalence runs in both directions. If \(P(A \mid B) = P(A)\) and \(P(B) > 0\), then the joint probability factors and the events are independent. Comparing a conditional probability with a marginal probability is therefore a complete independence test, not a rough diagnostic.

## Recovering an intersection

Sometimes the question supplies \(P(A)\), \(P(B)\) and \(P(A \text{ or } B)\) instead of the intersection. The addition rule recovers it:

\[
P(A \text{ and } B) = P(A) + P(B) - P(A \text{ or } B).
\]

Once the intersection is known, the product rule decides independence. This route is common when the data arrive as a union, for example when a survey reports the share of people who own a car, the share who own a bicycle and the share who own at least one of the two.

The subtraction must be done exactly. A union larger than the sum of the two marginals would make the intersection negative, which is impossible. A union equal to the sum means the events are disjoint, and disjoint events with positive probabilities are never independent. That connection between disjointness and dependence is worth remembering, because the two words are often confused.
