# Combining functions

Two functions can be combined point by point. At every input where both functions are defined, their values can be added, subtracted or multiplied, and the result is a new function. The notation \((f+g)(x)\) means evaluate each function separately and then combine the two numbers.

## The sum of two functions

The sum is defined by \((f+g)(x) = f(x) + g(x)\). The definition is pointwise, so the work at a particular input is short: find \(f(x)\), find \(g(x)\), then add.

For \(f(x) = 3x + 1\) and \(g(x) = x - 4\), at \(x = 2\) the values are \(f(2) = 7\) and \(g(2) = -2\). The sum is \(7 + (-2) = 5\). Checking against the combined rule gives \((f+g)(x) = 4x - 3\), which is also 5 at \(x = 2\).

The domain of the sum is the set of inputs where both functions are defined. Addition does not create new inputs, and it does not remove any that both functions already share.

## The difference of two functions

The difference is defined by \((f-g)(x) = f(x) - g(x)\). Order matters here, because subtraction is not commutative. \((f-g)(2)\) and \((g-f)(2)\) are opposites of each other.

With the same functions, \(f(2) - g(2) = 7 - (-2) = 9\). The combined rule is \((f-g)(x) = 2x + 5\), which gives 9 at \(x = 2\). Reversing the order gives -9.

A common slip is to subtract the functions in the order the names appear rather than the order in the notation. Reading \((f-g)\) as "f first" prevents it.

## The product of two functions

The product is defined by \((f \cdot g)(x) = f(x) \cdot g(x)\). Multiply the two values at the chosen input, not the two formulas unless a general rule is wanted.

At \(x = 2\) the product is \(7 \times (-2) = -14\). Expanding the general rule gives \((f \cdot g)(x) = (3x+1)(x-4) = 3x^2 - 11x - 4\), and substituting \(x = 2\) gives \(12 - 22 - 4 = -14\).

Combining functions is the foundation for composition, which applies one function to the output of another. Sums, differences and products combine outputs side by side. Composition chains them.
