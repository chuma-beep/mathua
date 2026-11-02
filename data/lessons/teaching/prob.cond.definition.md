# Conditional Probability

**Conditional probability** is the chance of an event after the sample space has been restricted to another event.

## Conditioning on an event

For events $A$ and $B$ with $P(B)>0$, the conditional probability of $A$ given $B$ is
$$P(A \mid B) = \frac{P(A \cap B)}{P(B)}.$$
The denominator rescales the probability so that the outcomes outside $B$ no longer count. Conditioning is therefore not a new probability rule. It is the ordinary rule applied to the smaller sample space $B$.

## The conditional probability formula

A school has 40 students in the band, and 6 of them also play in the orchestra. If a band member is chosen at random, the chance that the student also plays in the orchestra is
$$P(\text{orchestra} \mid \text{band}) = \frac{6}{40} \Big/ \frac{40}{40} = \frac{6}{40} = 0.15.$$
The numerator counts the students in both groups. The denominator counts only the students in the conditioning group, so students outside the band never enter the calculation.

## Updating a sample space

A two-way table makes the restricted sample space visible. Put the conditioning event in one row or column, then divide the joint count by that row or column total. If 25 students play a sport, and 10 of those also play an instrument, then
$$P(\text{instrument} \mid \text{sport}) = \frac{10}{25} = 0.4.$$
The same table also gives $P(\text{sport} \mid \text{instrument})$ with a different denominator. Reversing the condition reverses which total is used, which is why $P(A \mid B)$ and $P(B \mid A)$ are usually different.
