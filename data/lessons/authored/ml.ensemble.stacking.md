# Stacking ensembles

## Base models

A stacking ensemble trains several base models and combines their predictions with a meta-model. The base models should make different kinds of errors so the combination adds information.

The meta-model is trained on out-of-fold predictions. Training it on predictions the base models saw during fitting would leak information.

## Meta-features

The meta-features are the base predictions. A linear meta-model combines them with weights and a bias. The weighted sum is the meta-prediction.

The meta-weights measure how much each base model contributes. A large positive weight means the base prediction is useful in the combination.

## Comparing models

A stacking ensemble is useful when it predicts better than the best base model. Compare absolute errors on held-out data, not on the training data.

A more complex ensemble is not automatically better. If the meta-model overfits, a simple base model can win.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
