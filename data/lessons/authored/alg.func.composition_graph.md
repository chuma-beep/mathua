# Composition from graphs

Composition chains two functions. The output of the inner function becomes the input of the outer one. When both functions are shown as graphs or tables, the work is a sequence of two readings rather than one formula.

## Reading the inner function first

For \((f \circ g)(x)\), the inner function is \(g\). Read \(g(x)\) first, then use that number as the input to \(f\). The notation is read from the inside out, even though the symbol \(f\) appears first.

Suppose a graph of \(g\) gives \(g(2) = 3\), and a graph of \(f\) gives \(f(3) = 7\). Then \((f \circ g)(2) = 7\). The intermediate value 3 is the bridge between the two graphs.

The order is not optional. Using \(f(2)\) and then \(g\) of that value computes a different composition, and the two results usually disagree.

## Composing in the other order

The composition \((g \circ f)(x)\) reverses the roles. Now \(f\) is inner and \(g\) is outer. Read \(f(x)\) first and feed the result to \(g\).

With the same values, \(f(2) = 3\) and \(g(3) = 7\) would make \((g \circ f)(2) = 7\) as well, but that is a coincidence of the chosen numbers. In general \((f \circ g)(x)\) and \((g \circ f)(x)\) are different functions.

Composition is not commutative. It is associative, so three functions can be grouped either way, but swapping the order of two changes the result.

## Composing a function with itself

The composition \((f \circ f)(x)\) applies the same rule twice. First compute \(f(x)\), then substitute that output back into \(f\).

For \(f(x) = 2x + 1\), \(f(3) = 7\) and \(f(7) = 15\), so \((f \circ f)(3) = 15\). The general rule is \((f \circ f)(x) = 2(2x+1)+1 = 4x+3\), which gives 15 at \(x = 3\).

Graph reading makes the repeated step visible. Start at the input on the horizontal axis, move to the graph of \(f\), and use the height as the next input. Repeating the move is composition with itself.
