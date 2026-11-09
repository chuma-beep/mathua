# Gradient boosting

Gradient boosting builds an additive model by fitting each new tree to the residual of the current model.

## Residuals

The residual is the target minus the current prediction. A new weak learner is fitted to that residual. The residual measures the part of the target the current model has not explained.

## Learning rate

The new tree is added with a learning rate. A rate of one half adds half of the fitted residual. A smaller rate takes a shorter step and usually needs more trees.

## The additive model

After several trees the prediction is the base value plus the sum of the scaled tree outputs. The final prediction is therefore an additive combination. The order of the trees matters only through the residuals they fit.

The three sections form one chain: the first names the object, the second carries out the computation and the third reads the result in context.

A complete answer names the relation, substitutes the quantities in the question and states the result in the language of the concept. Where a second computation is available, compare the two results before moving on.
