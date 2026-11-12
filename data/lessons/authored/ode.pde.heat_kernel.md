# The heat kernel

The heat kernel is the response of the heat equation to a point source. Every solution on the whole line is a weighted average of its initial data against that kernel.

## The Gaussian kernel

The kernel is a normal density whose variance grows in proportion to time. It integrates to one for every positive time, so the total heat is conserved while the profile spreads.

## Convolution

The solution is the convolution of the initial data with the kernel. Each point of the initial profile contributes an independent Gaussian that widens as time advances.

## Scaling

The width grows as the square root of time and the height falls as its reciprocal. The product of height and width stays constant, which is the statement that mass is preserved.

To check your understanding, write the kernel at a given time, form the convolution integral and track the width in time.

The kernel makes the smoothing property explicit, because a convolution with a Gaussian averages the data over a neighbourhood that grows with time.
