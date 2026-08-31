# Overfitting and Generalization

**Overfitting:** The model performs well on training data but poorly on unseen test data — it memorized noise rather than pattern.

## Diagnosing Overfitting

### Training vs Test Error
Overfit: training error $\ll$ test error. Well-fit: both low and close. Underfit: both high.

### Capacity Curve
As model capacity grows, test error is U-shaped: drops while learning signal, then rises when fitting noise. The minimum marks the best capacity.

## Example

Polynomial fit on 20 noisy points: degree 2 underfits (high train and test error), degree 12 overfits (train error near 0, test error large), degree 4 balances both and generalizes.
