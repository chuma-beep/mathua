# Fixed-parameter tractability

A problem is fixed-parameter tractable when it can be solved in time \(f(k)\cdot n^{c}\) for a parameter \(k\) and a constant \(c\). The exponential cost is confined to the parameter.

## The definition

FPT requires the input size to appear only polynomially. A running time of \(n^k\) is not FPT because the exponent depends on the parameter.

## Kernels

A kernelisation reduces an instance to one of size bounded by a function of \(k\) in polynomial time. A problem with a kernel is FPT.

## Examples

Vertex cover parameterised by solution size is FPT. Independent set is not known to be, which is why the choice of parameter matters.

To check your understanding, separate parameter from size, reduce in polynomial time and parameterise vertex cover.

Parameterised complexity explains why some NP-hard problems are easy on real instances: the hardness sits in a small parameter that can be isolated.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
