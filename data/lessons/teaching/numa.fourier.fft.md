# The Fast Fourier Transform

## Roots of Unity

The discrete Fourier transform evaluates a polynomial at the complex roots of unity. Those roots have special symmetry: squaring a root gives another root of half the order.

The symmetry allows a transform of length n to be split into two transforms of length n/2. One handles the even-indexed samples and the other handles the odd-indexed samples.

The split is repeated recursively until the transforms have length one. At that point the transform is just the sample itself, and the results are combined on the way back up.

## Divide and Conquer

The combining step multiplies the odd transform by a twiddle factor and then adds or subtracts the two halves. This is the butterfly operation.

Each stage of the recursion performs n/2 butterflies. The total number of stages is the base-two logarithm of n, so the work is n log n instead of n squared.

The transform is inverted by reversing the sign of the exponent and dividing by n. The same butterfly structure serves both directions.

## Counting Operations

A direct transform of length n needs about n squared complex multiplications. The fast transform reduces that to about n log base two of n multiplications plus a similar number of additions.

For a length of one thousand twenty-four the saving is roughly a factor of one hundred. This is why the fast transform made spectral methods and digital signal processing practical.

The transform is usually implemented iteratively with a bit-reversal permutation at the start. The iterative form avoids recursion overhead and uses the same small butterfly routine at every stage.
