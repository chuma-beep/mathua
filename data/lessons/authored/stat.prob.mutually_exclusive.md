# Mutually exclusive events

Two events are mutually exclusive when they cannot happen at the same time. The definition is about the intersection, not about the probabilities of the separate events. A coin cannot land heads and tails on one toss, so those two outcomes are mutually exclusive. A card can be a heart and a face card at once, so those two events are not.

## Disjoint events

The intersection of mutually exclusive events is empty, which gives the probability statement

\[
P(A \text{ and } B) = 0.
\]

That single equation is the test. If the joint probability is zero, the events are disjoint. If it is positive, they overlap and are not mutually exclusive, however small the overlap happens to be.

Disjointness is a property of the events themselves, while the probability of zero is a consequence. In a continuous model a single point can have probability zero without being impossible, so the probability statement is the operational definition in this course. The distinction matters in measure theory, but not in the finite settings used here.

## The addition rule

For mutually exclusive events the general addition rule simplifies to

\[
P(A \text{ or } B) = P(A) + P(B).
\]

The general rule subtracts the intersection, and that term vanishes when the events are disjoint. If a question supplies the two marginals and the union, the test is to compare the union with the sum. Equality means disjointness.

If the union is smaller than the sum, the events overlap. The difference is exactly the probability of the intersection. That is why the addition rule can be rearranged to recover the overlap, and why a union that is too large signals an arithmetic error rather than a new phenomenon.

## When disjointness is possible

Two events with probabilities \(P(A)\) and \(P(B)\) can be mutually exclusive only if

\[
P(A) + P(B) \le 1.
\]

If the sum exceeds one, the events must overlap somewhere, because the total probability cannot exceed one. If the sum equals one, the events can be disjoint and together exhaust the sample space. If the sum is below one, disjointness is possible and leaves room for outcomes outside both events.

This inequality is a quick feasibility check. It does not prove that two particular events are disjoint, because it uses only the marginals. It proves whether disjointness is available at all. A positive answer means the model can place the events in separate regions, while a negative answer rules that arrangement out.
