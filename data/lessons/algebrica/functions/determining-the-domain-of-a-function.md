> Content sourced from [Algebrica](https://algebrica.org/determining-the-domain-of-a-function/) — CC BY-NC 4.0

## A systematic method for determining the domain of a function

In the introduction to [functions](<../functions>), we discussed the idea of the domain of a function, that is, the set of input values for which an expression is mathematically meaningful. We also examined how to determine the domain of basic families of functions such as [polynomials](<../polynomials/>), [radicals](<../radicals/>), [logarithms](<../logarithms/>), and [trigonometric](<../sine-and-cosine/>) expressions.

In practice, however, we often encounter more elaborate expressions that combine several types of functions in a single formula. In these situations, identifying the domain is not always immediate, because each internal component may introduce its own restrictions. When dealing with such cases, the most effective strategy can be summarized in the following steps:

  * Start by examining the inner components of the expression and determine the domain required by each of them.
  * Move outward, layer by layer, propagating the restrictions imposed by every new operation or function involved.
  * Combine all the resulting conditions by taking their intersection, since the overall domain consists only of the values that satisfy every constraint simultaneously.


To understand how this procedure works in practice, it is helpful to move directly to a concrete example and apply each step to real expressions.

## Types of intervals

When studying the domain of a function, it is useful to recall how real intervals are defined, since the domain is always expressed as a union of such sets. Intervals describe continuous portions of the real line and provide a compact way to specify which values of \\(x\\) are allowed. A first type is the open interval, which contains all points strictly between two endpoints while excluding the endpoints themselves. In descriptive form it is written  
\\[\lbrace\,x : a < x < b\,\rbrace \\] and in interval notation it appears as: \\[(a,\,b) \\] This notation is particularly convenient when the values \\(a\\) and \\(b\\) correspond to points where the function is not defined or where a discontinuity occurs. Graphically, the illustrations on Algebrica will follow this convention. The endpoints \\(a\\) and \\(b\\) are shown as open circles to indicate that they are not included in the interval, while the segment connecting them represents all values strictly between \\(a\\) and \\(b\\). This visual notation makes an open interval immediately recognizable, clearly emphasizing that its endpoints are excluded.

| \\[a\\]| \\[b\\]|   
---|---|---|---  
| | |   
| | |   
  
The second type is the closed interval, which includes both endpoints. In descriptive form it is written: \\[\lbrace\,x : a \le x \le b\,\rbrace \\] and the corresponding interval notation is  
\\[[\,a,\,b\,] \\] This representation is typical when the function is defined and continuous even at the boundary points \\(a\\) and \\(b\\). In the same spirit as the illustration shown above, the closed interval is represented graphically by filling the endpoints \\(a\\) and \\(b\\). The solid dots indicate that both boundary values are included, while the segment between them depicts every point from \\(a\\) to \\(b\\). This convention makes the structure of a closed interval immediately clear, emphasizing that its endpoints form part of the set.

| \\[a\\]| \\[b\\]|   
---|---|---|---  
| | |   
| | |   
  
Naturally, it is easy to see that other combinations of open and closed endpoints are also possible. For instance, one may have a half-closed interval such as \\([a, +\infty)\\), or more generally intervals that mix inclusion and exclusion at their boundaries such as \\((a, b]\\). These forms are illustrated in the figure below and are commonly used when describing domains that extend indefinitely or start from a specific boundary value.

| \\[a\\]| \\[b\\]|   
---|---|---|---  
| | |   
| | |   
| | |   
  
The graphical notation introduced above is particularly useful because it allows the domain of a function to be identified at a glance. When several conditions of existence arise from different components of the expression, each condition can be represented as its own interval diagram. The overall domain is then obtained by **intersecting these conditions** : in other words, only the intervals that appear simultaneously across all lines of the diagram form the set of admissible values for \\(x\\). This visual approach makes the structure of the domain clear and helps avoid errors in combining multiple constraints.

## Domain of elementary functions

As an additional point to keep in mind, it is useful to summarize the domain restrictions induced by the most common expressions:

  * **Polynomials:** always defined on \\(\mathbb{R}.\\)
  * **Rational functions:** denominator \\(\neq 0.\\)
  * **Even roots:** radicand \\(\gt 0.\\)
  * **Odd roots:** always defined on \\(\mathbb{R}.\\)
  * **Logarithms:** argument \\(> 0\\), base \\(> 0\\), base \\(\neq 1.\\)
  * **Absolute values:** always defined on \\(\mathbb{R}.\\)
  * **Trigonometric functions:** always defined on \\(\mathbb{R}.\\)
  * **Inverse trigonometric functions:** argument within specific bounds


## Example 1

To give an example of how the graphical method is used, let us consider the function:

\\[f(x) = \log(x - 1) + \sqrt{x + 4} \\]

To determine its domain, each component must be examined separately, because the function is defined only for those values of \\(x\\) that satisfy all underlying conditions at the same time.


The logarithmic term \\(\log(x - 1)\\) requires its argument to be strictly positive. This leads to the condition: \\[x - 1 > 0 \\]

which means that only values greater than \\(1\\) are admissible. Any value less than or equal to \\(1\\) immediately violates the definition of the logarithm. Graphically, we can represent the condition in the following way:

| \\[1\\]|   
---|---|---  
| |   
| |   
  
The square root term \\(\sqrt{x + 4}\\) requires its argument to be non-negative. This condition is expressed by: \\[x + 4 \ge 0 \\] so the values \\(x \ge -4\\) are allowed. Compared to the logarithmic constraint, this requirement is less restrictive, since every value greater than \\(1\\) is automatically greater than \\(-4\\). By adding this result to the previous diagram, we obtain:

| \\[-4\\]| \\[1\\]|   
---|---|---|---  
| | |   
| | |   
| | |   
| | |   
  
Putting everything together, the domain of the function is obtained by intersecting the two conditions \\(x > 1\\) and \\(x \ge -4\\). Because \\(x > 1\\) already implies \\(x \ge -4\\), the intersection reduces to the simple inequality \\(x > 1)\\).

Therefore, the domain of the function is represented by the highlighted black segment, which corresponds to the interval: \\[(1, +\infty) \\]

## Example 2

Let us consider a relatively simple example that already combines a square root, a logarithm, and a trigonometric function. The function is given by:

\\[f(x) = \sqrt{\frac{\log(2 + \sin x)}{\sinh(x)}} \\]

At first glance, determining the domain of this function is not entirely straightforward. One initial observation we can make, however, is that the presence of a trigonometric term suggests that any conditions we obtain are likely to repeat periodically. We now apply the method outlined above, examining step by step the domain restrictions introduced by each individual component of the expression, starting from the innermost elements.


The first element we encounter is the [logarithm](<../logarithms>): \\[\log(2 + \sin x)\\]

Since the logarithm is defined only for strictly positive values, we require:

\\[2 + \sin x > 0\\]

Because \\(\sin x\\) ranges between \\(-1\\) and \\(1\\), the smallest value this expression can take is \\(2 - 1 = 1\\), which is already positive. This means that the logarithm introduces no restriction on the domain, as its argument is positive for every real value of \\(x\\).


Moving a step outward, the logarithm itself appears in the numerator of a fraction whose denominator is \\(\sinh(x)\\). A denominator cannot be zero because division by zero is not defined in the [real numbers](<../types-of-numbers>), so we must exclude the values of \\(x\\) for which \\(\sinh(x) = 0\\). The hyperbolic sine vanishes only at \\(x = 0\\), so this point must be removed from the domain.


Finally, the entire fraction appears under a square root, which requires its argument to be non-negative.

  * The numerator \\(\log(2 + \sin x)\\) is always non-negative and becomes zero precisely when \\(\sin x = -1\\), that is, at the points \\(x = -\frac{\pi}{2} + 2k\pi\\). At those points the fraction is equal to zero, and thus acceptable, provided the denominator is non-zero, which is always the case except at \\(x = 0\\), already excluded.

  * For all other values, \\(\log(2 + \sin x)\\) is strictly positive, so the sign of the entire fraction depends solely on the sign of \\(\sinh(x)\\). The hyperbolic sine is positive for \\(x > 0\\) and negative for \\(x < 0\\). Since the radicand must be non-negative, the fraction is admissible exactly when \\(x > 0\\), or when the numerator is zero.


Putting these observations together, the domain consists of all \\(x\\) greater than zero, along with the isolated points where \\(\sin x = -1\\). Therefore, the domain is:

\\[(0, +\infty) \;\cup\; \left\\{ -\frac{\pi}{2} + 2k\pi \,\middle|\, k \in \mathbb{Z} \right\\} \\]

## Example 3

Let us consider another example and determine the domain of the function:

\\[f(x) = \sqrt{\cos^{2}(3x - 1) - \log\\!\bigl(5 - |2x|\bigr)} \\]

We determine the domain by examining the expression from the inside out, identifying the restrictions introduced by each component and combining them at the end. The first element that imposes a constraint is the logarithm. Since the natural logarithm accepts only strictly positive arguments, we must require

\\[5 - |2x| > 0 \\]

Solving this [inequality](<../linear-inequalities/>), which involves an [absolute value](<../absolute-value/>), gives:

\\[|2x| < 5 \quad \rightarrow \quad -\frac{5}{2} < x < \frac{5}{2} \\]


The trigonometric part, \\(\cos^{2}(3x - 1)\\), does not introduce additional restrictions, because the [cosine function](<../cosine-function/>) and its square are defined for all real values of \\(x\\).


We now move one level outward and consider the entire expression inside the square root. A square root is defined only when its argument is non-negative, so we must impose

\\[\cos^{2}(3x - 1) - \log(5 - |2x|) \ge 0 \\]

This inequality links together two terms with very different behaviors. The squared cosine oscillates between \\(0\\) and \\(1\\), while the logarithm varies over the interval \\((-5/2,\; 5/2)\\) and becomes arbitrarily large as \\(|2x|\\) approaches 5 from below. The inequality therefore holds only for those values of \\(x\\) for which the logarithmic term does not exceed the value of the squared cosine. In other words, we must have:

\\[\log(5 - |2x|) \le \cos^{2}(3x - 1) \\]

The set of solutions of this inequality is contained within the interval \\((-5/2, 5/2)\\), because outside this interval the logarithm would no longer be defined. Finally, we combine all restrictions by intersecting them, since the function is defined only when every condition is satisfied simultaneously. The logarithm gives the interval

\\[-\frac{5}{2} < x < \frac{5}{2} \\]

and within this interval we must additionally select only those values of \\(x\\) for which

\\[\cos^{2}(3x - 1) - \log(5 - |2x|) \ge 0 \\]

In conclusion, the domain of the function is given by all real values \\(x\\) such that:  
\\[x \in \left(-\frac{5}{2},\, \frac{5}{2}\right) \\]
