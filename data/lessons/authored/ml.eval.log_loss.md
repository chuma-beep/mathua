# Log loss and proper scoring rules

Log loss scores a predicted probability for a binary label. It is a proper scoring rule, so honest probabilities are rewarded over confident guesses.

## Log loss for one example

For a positive label and predicted probability p, log loss is minus the logarithm of p. In base 2 the answer is an exact number of bits. A smaller predicted probability for the true label gives a larger loss.

## Proper scoring rules

A scoring rule is proper when the expected score is optimized by the true probability. Log loss has that property. Accuracy alone does not reward calibrated probabilities.

## Log loss and accuracy

Accuracy counts correct labels and ignores confidence. Log loss penalizes confident mistakes heavily. Two models with equal accuracy can have very different log loss.

Practice connects the three sections: state the relation, carry out the computation and interpret the result. A second method is a useful check when one is available.

A complete answer states the governing relation, substitutes the given values and reports the result with the correct sign and units.
