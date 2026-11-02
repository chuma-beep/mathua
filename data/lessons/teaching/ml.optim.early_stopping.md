# Early stopping

## Validation monitoring

Early stopping monitors a validation metric during training and ends the run when the metric has not improved for a fixed number of epochs. That number is the patience.

If the best validation epoch is b and the patience is p, training stops at epoch b plus p. The final model is usually the checkpoint from epoch b rather than the last epoch.

## Why it works

Training error usually keeps falling while validation error eventually rises, which is the signature of overfitting. Stopping near the minimum of the validation curve keeps the model before that turn.

Early stopping acts as a form of regularisation because it limits how far the parameters can move from their initial values. It is cheap because it needs no change to the training objective.

## Bookkeeping

If training is allowed to run to a fixed budget of T epochs, stopping at epoch b plus p saves T minus b minus p epochs compared with using the whole budget.

The best checkpoint is restored after stopping. Resuming from that point with a lower learning rate is a common alternative to a single stopping decision.
