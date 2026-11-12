# Total probability from a partition

The law of total probability computes the probability of an event by splitting the sample space into non-overlapping parts and adding its probability within each part, weighted by the size of that part.

## Partitioning the sample space

A partition is a collection of events that are mutually exclusive and together exhaust the sample space. Two machines partition the output if every component comes from exactly one of them. If machine A makes 60 percent of the output and machine B the rest, the parts carry probabilities 0.6 and 0.4.

## Weighting each part

The rule is \(P(E) = P(E \mid A) P(A) + P(E \mid B) P(B)\), with one term per part. Each term is the probability of reaching the part and then observing the event, the same product a tree's path uses. The weights sum to 1, which keeps the total a probability.

## A worked total

Machine A produces 60 percent of the components and 2 percent of its output is defective. Machine B produces 40 percent and 5 percent is defective. The overall defect probability is \((0.02)(0.6) + (0.05)(0.4) = 0.012 + 0.020 = 0.032\). Machine A supplies 0.012 of that total, 37.5 percent of all defects, despite its lower rate.

The law converts a question about the whole population into one question per part, which is usually the only way the answer is available.
