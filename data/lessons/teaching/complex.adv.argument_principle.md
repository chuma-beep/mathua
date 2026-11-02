# The Argument Principle

## Zeros and poles

The argument principle relates the change in the argument of an analytic function along a closed contour to the zeros and poles inside. For a function meromorphic inside a simple closed contour with no zeros or poles on the contour, the change in argument divided by two pi equals the number of zeros minus the number of poles, counted with multiplicity.

The statement is a bookkeeping device for the winding number. As the contour is traversed, the image curve winds around the origin once for each zero and once in the opposite direction for each pole. The total winding is the difference. This is why the principle is also called the argument principle or the principle of the argument.

To apply it, factor the function near each zero or pole. A zero of order m contributes m, and a pole of order m contributes minus m. The contour must avoid all of these points, so small circles around them are excluded if necessary. The change in argument is then computed along the remaining curve.

## The winding contribution

The change in argument is a real number, but it is only defined modulo two pi until we choose a continuous branch along the contour. The total change is an integer multiple of two pi precisely because the image curve is closed. That integer is the winding number of the image around the origin.

For a rational function, the argument principle gives the difference between the number of zeros and the number of poles. If the function is a polynomial, there are no poles, so the change in argument counts the zeros. This is the usual way to count roots inside a contour without solving the equation.

Computing the change directly is often harder than counting zeros and poles. The principle is most useful in the other direction: known zeros and poles predict the winding, which can then be used to analyze stability or to prove that a curve surrounds the origin. In control theory this becomes the Nyquist criterion.

## Rouche as a consequence

Rouche's theorem compares two analytic functions on a contour. If one function is strictly larger than the other on the contour, then the two functions have the same number of zeros inside. The proof applies the argument principle to the family that interpolates between them and notes that the winding number cannot change without a zero on the contour.

The usual application counts zeros of a complicated polynomial by comparing it with its dominant term on a large circle. For example, a polynomial of degree n is dominated by its leading term on a sufficiently large circle, so it has n zeros counting multiplicity. This proves the fundamental theorem of algebra without solving the equation.

Rouche is also useful for locating zeros in a region. Choose a comparison function whose zeros are known and verify the strict inequality on the boundary. The conclusion is local to the region, so different regions can be handled separately. The strict inequality is essential; equality on the contour allows the count to change.

