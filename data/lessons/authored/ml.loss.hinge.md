# Hinge loss

## Margin

The margin of a classifier is the label times the score. A positive margin means the prediction is on the correct side. A margin of at least one is required for zero hinge loss.

The margin combines the true label and the score into one number. It is not the same as the score alone.

## Loss

Hinge loss is the maximum of zero and one minus the margin. A margin above one contributes no loss. A margin below one contributes a positive penalty that grows as the margin becomes more negative.

The loss is piecewise linear. It punishes confident mistakes more than uncertain correct predictions.

## Subgradient

The subgradient of hinge loss with respect to the score is negative label when the margin is below one, and zero otherwise. At the margin equal to one the subgradient is set-valued, and either endpoint can be used.

The subgradient drives the update. A zero subgradient means the example is already correct with a sufficient margin and does not move the classifier.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
