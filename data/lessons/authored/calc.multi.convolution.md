# Multivariable convolution

Multivariable convolution combines two functions by sliding one over the other. On a finite grid it is a circular sum over all shifts.

## Definition

The convolution of f and g at a point sums f at one position times g at the reflected displacement. Each output point uses all input pairs. On a 2 by 2 grid the wrap-around makes the sum finite.

## Fourier product

The Fourier transform turns convolution into pointwise multiplication. This is why convolution is central in signal processing and partial differential equations. The total mass of the convolution is the product of the total masses.

## A point value

To compute one convolution value, list the four contributing pairs and add their products. The wrap-around indices are taken modulo two. The result is an integer when the inputs are integers.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
