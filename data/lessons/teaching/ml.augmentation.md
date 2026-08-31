# Data Augmentation

**Data augmentation:** Artificially expand training data via label-preserving transforms: flips, crops, rotations, color jitter, cutout, mixup — approximates desired invariances.

## Regularization via Augmentation

### Invariance
Augmenting with flips teaches model to ignore horizontal reflection; mixup $λx_1+(1-λ)x_2$, $λy_1+(1-λ)y_2$ encourages linear behaviour between samples.

### Train vs Test
Augmentation only at training; test uses clean data. Effectively increases data, reduces overfitting like regularization.

## Example

CIFAR-10: random crop $32×32$ from $40×40$ padded image plus horizontal flip doubles effective data and improves generalization by $≈2\%$.
