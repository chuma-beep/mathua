> Content sourced from [Algebrica](https://algebrica.org/functions/) — CC BY-NC 4.0

## What is a function

A function is a mathematical rule that connects two non-empty [subsets](<../sets/>) of the [real numbers](<../real-numbers/>), typically denoted as \\( A \subseteq \mathbb{R} \\) and \\( B \subseteq \mathbb{R} \\). A function \\( f \\) from \\( A \\) to \\( B \\) assigns exactly one real number in \\( B \\) to each real number in \\( A \\). This relationship is written as: \\[f \colon A \to B \\]

  * The set \\( A \\) is called the [domain](<#domain>) of the function.
  * The set \\( B \\) is the codomain.
  * For every \\( x \in A \\), the function produces a unique value \\( f(x) \in B \\).
  * The variable \\( x \\) is called the independent variable, while \\( y \\) is the dependent variable.


![](/diagrams/algebrica/functions-5.png)

> If such a rule holds, we say that the function is well-defined. Otherwise, the relation is not a function, since it either assigns no value or more than one value to an element of the domain.


A function from a set \\( A \\) (domain) to a set \\( B \\) (codomain), \\( f : A \to B, \\) can be classified as:

  * Injective: if every element of \\( B \\) is the image of at most one element of \\( A \\), that is, if for any \\(x_1, x_2 \in A\\) with \\(x_1 \neq x_2\\), we have \\(f(x_1) \neq f(x_2)\\); equivalently, if for every \\(y \in B\\) there is at most one \\(x \in A\\) such that \\(f(x) = y\\).
  * Surjective: if every element of \\( B \\) is the image of at least one element of \\( A \\), that is, if \\(\forall \, y \in B\\) there exists at least one \\(x \in A\\) such that \\(f(x) = y.\\)
  * Bijective: if the function is both injective and surjective, that is, if for every \\(y \in B\\) there exists a unique \\(x \in A\\) such that \\(f(x) = y\\).


An alternative definition states that a function \\(f : A \to B\\) is bijective if and only if there exists a function \\(g : B \to A\\) such that: \\[(g \circ f)(x) = x, \; \forall \, x \in A\\] \\[(f \circ g)(y) = y, \; \forall y \in B\\] Whenever such a function \\(g\\) exists, it is uniquely determined. In this case \\(g\\) is called the [inverse](<../inverse-function/>) of \\(f\\) and is denoted by \\(f^{-1}\\).

> An example of an inverse function is the [logarithmic function](<../logarithmic-function/>), which serves as the inverse of the [exponential function](<../exponential-function/>), and conversely.

## What is not a function

For greater clarity, let us illustrate the case in which we face a situation similar to the one shown in the previous figure, but where it is not possible to define a function. In formal terms, this occurs when there exist at least two distinct points \\( (x, y_1) \\) and \\( (x, y_2) \\) such that \\( y_1 \ne y_2 \\). In this case, the relation does not satisfy the definition of a function, since a single value of \\( x \\) is associated with more than one value of \\( y \\).

![What is not a function.](/diagrams/algebrica/functions-4-1.png)

In the image, a single element of the domain, denoted by \\( x_0 \\), corresponds to two distinct values in the codomain. This situation is not admissible according to the very definition of a function. Formally, a relation \\( R \subseteq A \times B \\) is a function if and only if:

\\[\forall \, x \in A, \; \exists! \, y \in B \; : \; (x, y) \in R \\]

In this case, there exist \\( y_1, y_2 \in B \\) with \\( y_1 \ne y_2 \\) such that \\((x_0, y_1) \in R\\) and \\((x_0, y_2) \in R\\) which violates the condition of uniqueness required for \\( R \\) to be a function. In simpler terms, imagine a situation where certain values of \\(x\\) are associated with corresponding values of \\(y\\), as shown in the following table.

X| -3| 1| -3| 5| 2  
---|---|---|---|---|---  
Y| 7| 4| 10| -2| 8  
  
In this case, the relation does not represent a function because for the same value \\( x = -3 \\), there are two different corresponding values of \\( y \\), namely \\( y = 7 \\) and \\( y = 10 \\). This violates the fundamental definition of a function, which requires that each element of the domain be associated with one and only one element of the codomain.


A practical graphical criterion for deciding whether a curve in the plane represents a function is the vertical line test. Draw any vertical line and count how many times it intersects the curve. If every vertical line meets the curve in at most one point, then the curve represents a function.

  * For each value of \\( x \\) there is at most one corresponding value of \\( y \\).
  * If some vertical line crosses the curve in two or more points, the curve does not represent a function, because that value of \\( x \\) would be associated with multiple outputs.


![](/diagrams/algebrica/functions-10-1.png)

The figure shows that the curve on the left (a [parabola](<../parabola/>)) is a function, since each \\( x\ ) corresponds to exactly one \\( y \\), whereas the curve on the right is not, since for \\( x_2 \\) there are multiple possible values of \\( y \\).

## Difference between codomain and range

  * The codomain is the set that we declare as the potential target of the outputs of a function. It is explicitly stated in the function’s definition. For instance, in a function written as \\( f: A \to B \\), the set \\( B \\) is the codomain.

  * The range (or image) is the actual set of outputs that the function attains when evaluated over its domain. It includes all values \\( f(x) \\) for \\( x \in A \\). The range is always a subset of the codomain.


## Function equality and zeros

Two functions \\( y = f(x) \\) and \\( y = g(x) \\) are considered equal if they share the same domain \\( D \\) and satisfy the condition:

\\[f(x) = g(x) \quad \forall \, x \in D \\]


A real number \\( a \in \mathbb{R} \\) is called a zero of the function \\( y = f(x) \\) if the function evaluates to zero at that point, that is:

\\[f(a) = 0 \\]

This means that the graph of the function intersects the x-axis at the point \\( (a, 0) \\). Identifying the zeros of a function is crucial for analyzing its behavior, solving [equations](<../equations/>), and determining where the function changes sign.

## Symmetric functions

When studying real-valued functions, it is often useful to understand how a function behaves with respect to reflections around the origin. This leads to the distinction between [even and odd functions](<../even-and-odd-functions/>). Let \\(A \subseteq \mathbb{R}\\) be a domain that is symmetric with respect to the origin, meaning that \\(x \in A \Rightarrow -x \in A\\). A function \\(f : A \to \mathbb{R}\\) is said to be:

  * Even if \\(f(-x) = f(x) \quad \forall \; x \in A. \\)
  * Odd if \\(f(-x) = -f(x) \quad \forall \; x \in A. \\)


More generally, for the function \\(f(x) = x^{n}\\) with \\(n \in \mathbb{N}\\), the function is even exactly when the exponent \\(n\\) is even, and it is odd exactly when \\(n\\) is odd.

## Bounded functions

Before introducing the formal definitions, it is useful to reflect on how a real-valued function behaves across its entire domain. In many situations we want to know whether the function stays within certain limits, never rising above a fixed threshold or never dropping below one. A function \\(f : A \subseteq \mathbb{R} \to \mathbb{R}\\) is said to be:

  * Upper bounded if there exists a real number \\(M \in \mathbb{R}\\) such that \\(f(x) \leq M\\) \\( \forall \; x \in A.\\)
  * Lower bounded if there exists a real number \\(m \in \mathbb{R}\\) such that \\(m \leq f(x) \forall \; x \in A.\\)
  * Bounded if it satisfies both conditions above, meaning that there exist \\(m, M \in \mathbb{R}\\) for which \\(m \leq f(x) \leq M\\) holds \\( \forall \; x \in A\\).


It is important to distinguish between the notion of a bounded function and the existence of a [maximum or minimum](<../maximum-minimum-and-inflection-points/>). Saying that a function is bounded does not imply that it actually attains a maximal or minimal value, either globally or locally. The requirement of boundedness simply means that the values of the function remain confined within a finite interval \\([m, M]\\).

Of course, if a function has a global maximum, then it is automatically bounded from above. Similarly, the presence of a global minimum ensures that it is bounded from below. However, the converse is not true: a function may be bounded without attaining any maxima or minima. A typical example is:

\\[f(x) = \sin\\!\left(\frac{1}{x}\right) \quad x \neq 0  
\\]

This function always stays between \\(-1\\) and \\(1\\), so it is indeed bounded. Nevertheless, it has neither a global maximum nor a global minimum, because as \\(x\\) approaches \\(0\\) the oscillations become arbitrarily rapid and the function never reaches its extreme values.

## Monotone functions

In the analysis of real-valued functions, understanding whether a function increases, decreases, or maintains a consistent directional behavior is essential. These ideas are captured by the concepts of [increasing, decreasing, and monotone functions](<../increasing-and-decreasing-functions/>), which describe how the output values evolve as the input grows. Let \\(A \subseteq \mathbb{R}\\) and let \\(x_1, x_2 \in A\\) with \\(x_1 < x_2\\). A function \\(f : A \to \mathbb{R}\\) is said to be:

  * Increasing if \\(f(x_1) \leq f(x_2).\\)
  * Strictly increasing if \\(f(x_1) < f(x_2).\\)
  * Decreasing if \\(f(x_1) \geq f(x_2).\\)
  * Strictly decreasing if \\(f(x_1) > f(x_2).\\)
  * Monotone if it satisfies one of the conditions above throughout its domain.


## Periodic functions

A function is called periodic when its values repeat after a fixed horizontal shift. More formally, if \\(X \subseteq \mathbb{R}\\) and \\(T > 0\\) is such that \\(x + T \in X\\) whenever \\(x \in X\\), a function: \\[f : X \to \mathbb{R} \\] is periodic with period \\(T\\) when \\(T\\) is the smallest positive number for which  
\\[f(x + T) = f(x) \\] holds \\(\forall \; x\\) in the domain. As a simple illustration, the [sine](<../sine-function/>) and [cosine functions](<../cosine-function/>) are periodic: both repeat their entire pattern after an interval of length \\(2\pi\\).

## Classification of functions

Functions are classified as algebraic or transcendental. A function is said to be algebraic if its analytical expression \\( y = f(x) \\) involves only a finite number of operations such as addition, subtraction, multiplication, division, exponentiation to a rational [power](<../powers/>), and [root extraction](<../radicals/>). Algebraic functions can be further categorized based on the structure of their expressions:

  * [Polynomial functions](<../polynomial-function/>) are defined by a [polynomial](<../polynomials/>) expression involving powers of \\( x \\) with constant coefficients.
  * Rational functions are expressed as the [ratio of two polynomials](<../rational-equations/>).
  * Irrational functions contain the variable \\( x \\) under a [root symbol](<../radicals/>), such as \\( \sqrt{x} \\).


Transcendental functions go beyond algebraic operations and include expressions that cannot be written using a finite combination of addition, subtraction, multiplication, division, and root extraction. Common examples are [exponential functions](<../exponential-function/>), [logarithmic functions](<../logarithmic-function/>), and [trigonometric functions](<../sine-function/>).

## Domain of the main functions

Below we will review the domains of the most common elementary functions. In practice, however, many functions encountered in real problems are considerably more intricate, and their domains cannot be identified at a glance. For all such situations in which the structure of the expression makes the analysis less immediate, we refer to [this method](<../determining-the-domain-of-a-function/>), which provides a systematic way to determine the domain of more complex functions.


[Polynomial functions](<../polynomial-function/>) are defined by expressions of the form:

\\[y = a_0 x^n + a_1 x^{n-1} + \dots + a_n \\]

where \\( a_0, a_1, \dots, a_n \\) are real coefficients and \\( n \in \mathbb{N} \\). The domain of a polynomial function is the entire set of real numbers, \\( \mathbb{R} \\), since it does not involve any operations that could restrict its definition. An example of a polynomial function is:

\\[y = 2x^3 - 5x^2 + 3x - 1 \\]

This expression can be viewed as a linear combination of powers of \\( x \\), where each term \\( a_i x^i \\) has a real coefficient \\( a_i \\). Polynomial functions of this type are defined for every real number \\( x \in \mathbb{R} \\), since no operation in their form imposes any restriction on the domain.


[Rational functions](<../rational-functions/>) are expressed in the form:

\\[y = \frac{N(x)}{D(x)} \\]

where \\( N(x) \\) and \\( D(x) \\) are polynomials. These functions are defined for all real numbers \\( x \\) such that \\( D(x) \neq 0 \\). The domain is therefore \\( \mathbb{R} \\) excluding the values that make the denominator zero. An example of a rational function is:

\\[y = \frac{x^2 - 4}{x - 2} \\]

This expression represents the ratio of two polynomials and can be interpreted as a linear combination of powers of \\( x \\) in both the numerator and denominator. Rational functions of this type are defined for all real values of \\( x \\) except those that make the denominator equal to zero. In this example, the function is undefined at \\( x = 2 \\), since it would make the denominator vanish.


Irrational functions are expressions of the form:

\\[y = \sqrt[n]{f(x)} \\]

The domain depends on the parity of the index \\( n \\). If \\( n \\) is even, the function is defined only when \\( f(x) \geq 0 \\), so the domain is:

\\[{ x \in \mathbb{R} \mid f(x) \geq 0 } \\]

An example of an irrational function with an even index is:

\\[y = \sqrt{x - 2} \\]

In this case, the expression under the radical must be non-negative, so the domain is given by \\( x \ge 2 \\). The function is therefore defined only for values of \\( x \\) that make the radicand \\( x - 2 \\) greater than or equal to zero. If \\( n \\) is odd, the function is defined for every value in the domain of \\( f(x) \\). An example with an odd index is:

\\[y = \sqrt[3]{x - 2} \\]

Here, the cube root is defined for every real value of \\( x \\), since odd roots can also take negative arguments. Thus, the domain of the function is the entire set of real numbers, \\( \mathbb{R} \\).


[Logarithmic functions](<../logarithmic-function/>) are defined as:

\\[y = \log_a{f(x)} \quad \text{with} \quad a > 0,\; a \ne 1 \\]

They are defined only when the argument of the [logarithm](<../logarithms/>) is strictly positive. Therefore, their domain is:

\\[{x \in \mathbb{R} \mid f(x) > 0 } \\]

An example of a logarithmic function is:

\\[y = \log_2(x - 1) \\]

This function is defined only when the argument of the logarithm, \\( x - 1 \\), is strictly positive. Therefore, the domain is \\( x > 1 \\). For any value of \\( x \\) less than or equal to \\(1\\), the expression becomes undefined because the [logarithm](<../logarithms/>) of a non-positive number does not exist in the real domain. Another example is:

\\[y = \ln(3x + 6) \\]

In this case, the argument \\( 3x + 6 \\) must also be positive, so the domain is \\( x > -2 \\). The same principle applies to all logarithmic functions: the argument of the logarithm must always be greater than zero for the function to be defined.


[Exponential functions](<../exponential-function/>) of the form:

\\[y = a^{f(x)} \quad \text{with} \quad a > 0,\; a \ne 1 \\]

are defined for all values in the domain of \\( f(x) \\). An example of an exponential function is:

\\[y = 2^x \\]

This function is defined for every real value of \\( x \\), since the base \\( a = 2 \\) is positive and different from \\(1\\). Its domain is therefore the entire set of real numbers, \\( \mathbb{R} \\), while the range is strictly positive, \\( y > 0 \\).


Exponential functions of the form:

\\[y = [f(x)]^{g(x)} \\]

are defined only when the base \\( f(x) > 0 \\), since real-valued exponentiation requires a positive base. Therefore, their domain is the intersection of:

\\[\\{ x \in \mathbb{R} \mid f(x) > 0 \\} \, \cap \, \text{domain of } g(x) \\]


Exponential functions of the form:

\\[f(x)^{\alpha} \\]

where \\( \alpha \in \mathbb{R} \setminus \mathbb{Q} \\) are defined under the following conditions:

\\[\\{ x \in \mathbb{R} \mid f(x) \geq 0 \\}, \quad \text{if } \alpha > 0 \\] \\[\\{ x \in \mathbb{R} \mid f(x) > 0 \\}, \quad \text{if } \alpha < 0 \\]

This ensures that the result of the expression remains within the set of real numbers.


For trigonometric functions, the following domains apply:

\\( y = \sin x, \quad y = \cos x \\)  
Domain: \\( \mathbb{R} \\)

\\( y = \tan x \\)  
Domain: \\( \mathbb{R} \setminus \left\lbrace \dfrac{\pi}{2} + k\pi \right\rbrace, \quad k \in \mathbb{Z} \\)

\\( y = \cot x \\)  
Domain: \\( \mathbb{R} \setminus \left\lbrace k\pi \right\rbrace, \quad k \in \mathbb{Z} \\)

\\( y = \arcsin x, \quad y = \arccos x \\)  
Domain: \\( [-1, 1] \\)

\\( y = \arctan x, \quad y = \text{arccot}\,x \\)  
Domain: \\( \mathbb{R} \\)

> Understanding the domain of a function is a crucial step in solving equations, inequalities, and in [analyzing functions](<../analyzing-the-graphs-of-functions/>). In most cases, identifying the domain involves recognizing a few fundamental situations, summarized below.

## Operations between functions

When two real functions are defined on the same domain, it is possible to combine them through the usual algebraic operations, obtaining new functions derived from the original ones. This generalizes the basic arithmetic of real numbers to the context of functions, where each operation is applied point by point for every value of \\( x \\) in the domain. In formal terms, if

\\[f : X_1 \subseteq \mathbb{R} \to \mathbb{R} \quad \text{and} \quad g : X_2 \subseteq \mathbb{R} \to \mathbb{R} \\]

are two functions, we define: \\(X = X_1 \cap X_2\\) as the common domain on which the following operations can be performed.


The sum of two functions \\( f \\) and \\( g \\) is defined by:

\\[(f + g)(x) = f(x) + g(x) \\]

The resulting function assigns to each \\( x \in X \\) the sum of the corresponding values of \\( f \\) and \\( g \\).


The difference of two functions is given by:

\\[(f - g)(x) = f(x) - g(x) \\]


The product of two functions is defined as:

\\[(f \cdot g)(x) = f(x) \, g(x) \\]

The resulting function represents the pointwise multiplication of the two values.


The quotient of two functions is expressed by

\\[\left(\frac{f}{g}\right)(x) = \frac{f(x)}{g(x)} \\]

but it is defined only for those values of \\( x \in X \\) where \\( g(x) \neq 0 \\). The domain of this function is therefore obtained by excluding from \\( X \\) all points that make the denominator vanish.


The composition of two functions, written \\( (g \circ f)(x) = g(f(x)) \\), is another fundamental operation and is treated in detail on the dedicated page on [composite functions](<../composite-functions/>).

Concept
