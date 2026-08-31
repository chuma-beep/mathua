# Convolutional Neural Networks

**CNN:** Stack of convolutional layers that slide learnable filters over the input, sharing weights spatially to detect local patterns, followed by activations (ReLU) and pooling.

## Local Structure

### Convolutions and Sharing
A $3×3$ filter on a $5×5$ image (stride $1$, no padding) yields $3×3$ output: $(5-3)/1+1=3$. Same filter reused at each position — far fewer parameters than fully connected.

### Pooling and Invariance
Max-pooling (e.g. $2×2$) downsamples and provides translation invariance; features survive small shifts.

## Example

Grayscale $5×5$, filter $3×3$, stride $1$: $9$ multiply-adds per position, $9$ positions → $81$ ops, vs $25·9=225$ for fully connected to $9$ outputs with sharing saving parameters.
