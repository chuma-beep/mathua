# Stochastic variance-reduced gradients

Full gradients are accurate and expensive. A sample gradient is cheap and noisy. SVRG keeps a snapshot of the parameters and a snapshot gradient, then corrects each sample gradient against the snapshot.

## The snapshot correction

At a snapshot point the method computes one full gradient. A later sample gradient is corrected by subtracting the snapshot gradient of the same sample and adding the full snapshot gradient. The correction leaves the expected value unchanged and removes most of the noise near the snapshot.

## The inner loop

Each inner step draws one sample, forms the corrected gradient and takes a gradient step. The snapshot is refreshed after a fixed number of inner steps, usually once per pass over the data. The refresh cost is one full gradient for every pass.

## The variance argument

The corrected gradient is unbiased because the two snapshot terms cancel in expectation. Its variance shrinks as the iterate approaches the snapshot, so the method behaves like a full gradient method near the solution while paying a sample cost per step.

To check your understanding, subtract the sample snapshot gradient, draw one sample and cancel in expectation.

A variance-reduced method buys a lower noise floor with memory for two vectors and one periodic full pass, which is a good trade when the sample count is large.
