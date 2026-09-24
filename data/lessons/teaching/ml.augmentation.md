# Data Augmentation

**Data augmentation:** Expand training data with label-preserving transforms — flips, crops, rotations, color jitter, cutout, mixup. Augmented copies teach invariances the architecture lacks, acting as regularization through effective data size.

## Worked: doubling CIFAR-10

Take 50000 training images with random crops plus horizontal flips:
1. Each epoch shows each image cropped differently — the model rarely sees the same pixels twice.
2. Flips double the effective poses: roughly 100000 distinct views per epoch from 50000 stored images.
3. Generalization typically gains about 2 points: the model learns flip-invariance instead of memorizing orientations.

So augmentation multiplies data without collecting any: variety comes from transforms, labels ride along unchanged.

## Worked: one mixup pair

Take two samples $x_1, x_2$ with one-hot labels $y_1, y_2$ and mixing weight $\lambda = 0.3$:
1. Blend inputs: $\tilde{x} = 0.3\,x_1 + 0.7\,x_2$.
2. Blend labels identically: $\tilde{y} = 0.3\,y_1 + 0.7\,y_2$.
3. Train on $(\tilde{x}, \tilde{y})$: the model learns linear behavior between samples instead of sharp cliffs at each point.

So mixup fills the gaps between data points with soft targets — smoother boundaries, less memorization.

## When augmentation breaks labels

Transforms must preserve meaning: horizontal flips keep cats as cats, but rotating digit 6 by half a turn makes 9 — the label lies. Augment only with transforms the task is invariant to; verify each transform keeps labels true.
