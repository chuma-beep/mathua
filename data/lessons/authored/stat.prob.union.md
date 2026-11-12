# The addition rule for probability

The addition rule gives the probability that at least one of two events occurs, by adding their probabilities and correcting for the outcomes counted twice.

## The union of two events

The union of A and B holds the outcomes in A, in B or in both. Its probability is not the sum of the two probabilities unless the events cannot both happen. One card drawn from a pack: a heart has probability \(13/52\) and a picture card has probability \(12/52\). Adding gives 25 over 52, but the three picture hearts have each been counted twice.

## Subtracting the overlap

The corrected rule is \(P(A \cup B) = P(A) + P(B) - P(A \cap B)\). The subtracted term is the intersection. For the cards the intersection holds 3, so the union has probability \((13 + 12 - 3)/52 = 22/52 = 11/26\). When the events are mutually exclusive the intersection is empty and plain addition is correct.

## A worked union

In a survey of 200 households, 120 own a car, 70 own a bicycle and 40 own both. The union is \(120 + 70 - 40 = 150\) households, so a random household owns at least one of the two with probability \(150/200 = 0.75\). Without the subtraction the count would be 190, implying only 10 own neither, when the true count is 50.

The rule is a statement about sets; drawing the two overlapping regions keeps the correction in the right place.
