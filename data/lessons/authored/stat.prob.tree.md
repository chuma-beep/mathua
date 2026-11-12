# Probability tree diagrams

A probability tree draws the stages of a random process as branches, with each branch labelled by the conditional probability of that stage given the path so far.

## Branches and paths

The first split covers the outcomes of the first stage. Each later split covers a stage conditional on the branch it grows from. Two draws without replacement from a bag of 5 red and 3 blue counters give first branches \(5/8\) and \(3/8\). After a red counter is gone the next red has probability \(4/7\), while after blue it is \(5/7\).

## Multiplying along a path

The probability of a complete path is the product of its branch probabilities. The path red then red has probability \((5/8)(4/7) = 20/56 = 5/14\). The path red then blue has \((5/8)(3/7) = 15/56\). Blue then red gives the same value, because both orderings select the same pair of counters.

## Complementary paths

The leaf probabilities sum to 1, which is a check and also a shortcut. The probability of at least one red in two draws is \(1 - (3/8)(2/7) = 1 - 6/56 = 50/56\). The complement path, blue then blue, is cheaper than the three paths it replaces.

A tree is worth drawing whenever a process has two or three stages, because the picture keeps each condition visible.
