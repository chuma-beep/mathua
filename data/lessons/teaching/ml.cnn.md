# Convolutional Neural Networks

**CNN:** Layers that slide small learnable filters over the input, reusing the same weights at every position to detect local patterns, followed by ReLU activations and pooling. Weight sharing makes them the default for image data.

## Worked: output size

Take a $5\times 5$ image, a $3\times 3$ filter, stride 1, no padding:
1. Slide the filter: it fits at rows 1–3, 2–4, 3–5 — 3 positions vertically, 3 horizontally.
2. In general, $(5-3)/1 + 1 = 3$ positions per side, so the output is $3\times 3$.
3. Each output cell is one filter placement: 9 multiply-adds, 9 placements, 81 ops total.

So convolution arithmetic is placement counting: output size follows from input size, filter size, stride, and padding.

## Worked: why sharing saves parameters

Take the same $5\times 5$ input mapped to 9 outputs:
1. Fully connected needs $25\cdot 9 = 225$ weights plus 9 biases — every input-output pair is independent.
2. A $3\times 3$ convolutional filter needs 9 weights plus 1 bias, reused at all 9 positions.
3. Fewer parameters means less memorization and translation robustness: the same edge detector fires wherever the edge appears.

So sharing trades position-specific freedom for sample efficiency — the right trade wherever location should not matter.

## When convolutions are wrong

On data with no spatial locality (shuffled pixels, tabular features) the locality assumption buys nothing and the constraint only hurts. Convolve where neighbors correlate; use fully connected or attention layers where they do not.
