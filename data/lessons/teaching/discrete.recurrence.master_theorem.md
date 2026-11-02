# The Master Theorem

## The Divide-and-Conquer Recurrence

A divide-and-conquer algorithm splits an input of size n into a smaller subproblems, each of size n/b, and spends f(n) work outside the recursive calls. Its running time obeys a recurrence of the form T(n) = a T(n/b) + f(n).

The master theorem compares the growth of the recursive part with the growth of the extra work. The key quantity is log base b of a, the rate at which the number of subproblems outgrows the shrinking input size.

Write the recurrence in this standard shape before applying any case. Constants outside the recursion belong in f(n), and the base case is handled separately.

## Which Case Applies

When f(n) grows more slowly than n raised to log base b of a, the recursion dominates and the solution has that same exponent. When the two grow at the same rate, the solution gains an extra logarithmic factor.

When f(n) grows faster and satisfies a regularity condition, the extra work dominates and the solution follows f(n). Comparing exponents is enough for the polynomial examples in this course.

A quick check on a small input confirms that the exponent is plausible, but the theorem itself supplies the asymptotic conclusion.

## Reading the Exponent

The exponent of n in the solution is the larger of log base b of a and the polynomial exponent of f(n), with a logarithmic tie-breaker. Many exam questions ask only for that exponent.

For example, two halves with linear extra work gives an exponent of one and an extra log factor, while two halves with quadratic extra work gives an exponent of two. The arithmetic is in the exponents, not the constants.

Practice with several triples of a, b and f until the comparison is automatic.
