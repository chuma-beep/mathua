# The Kaplan-Meier estimator

The Kaplan-Meier estimator estimates a survival curve from censored data.

## Compute a survival step

At each event time multiply the current survival probability by one minus events over at-risk. Censored observations leave the curve unchanged at their time. The result is a step function.

## Account for censoring

A censored observation is removed from the risk set after its censoring time. It does not count as an event. This is why censoring reduces the at-risk denominator without lowering the curve directly.

## Find the median survival

The median is the first time the curve reaches one half. If the curve never reaches one half the median is not estimable. The median is read from the step plot.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
