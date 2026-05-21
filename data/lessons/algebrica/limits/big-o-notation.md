> Content sourced from [Algebrica](https://algebrica.org/big-o-notation/) — CC BY-NC 4.0

## What is Big O notation

The symbol \\( O(x) \\), commonly known as big O of \\( x \\), is part of the Landau symbol family. It represents the concept of [asymptotic](<../asymptotes/>) upper bounds, signifying that a [function](<../functions>) does not grow faster than another function as the input approaches a specified [limit](<../limits>). Formally, it indicates that the ratio of the two functions remains bounded as the input approaches the limit.

Let \\(f(x)\\) and \\(g(x)\\) be two functions defined on a set \\(A\\), and let \\(x_0\\) be a limit point for \\(A\\). If there exist positive constants \\(M\\) and \\(\delta\\) such that for all \\(x\\) in a neighborhood of \\(x_0\\) (excluding \\(x_0\\) itself), \\(|f(x)| \leq M|g(x)|\\), then we say that \\(f(x)\\) is big O of \\(g(x)\\) as \\(x\\) approaches \\(x_0\\). In symbols:

\\[\exists \, M > 0, \, \delta > 0 : |x - x_0| < \delta \rightarrow |f(x)| \leq M|g(x)| \quad \rightarrow \quad f(x) = O_{x_0}(g(x))\\]

This notation indicates that \\(f(x)\\) grows asymptotically no faster than \\(g(x)\\) near \\(x_0\\). To illustrate this definition, let \\( f(x) = 3x^2 + 2x + 1 \\) and \\( g(x) = x^2 \\). The graph displays \\( f(x) \\) together with \\( 6x^2 = M \cdot g(x) \\), which serves as the asymptotic upper bound.

![](/diagrams/algebrica/big-o-notation-1.png)

For all \\( x \geq 1 \\), the curve \\( f(x) \\) remains entirely below \\( 6x^2 \\). Although both curves exhibit the same growth rate, \\( f(x) \\) does not exceed its bound. This relationship represents the geometric interpretation of \\( f(x) = O(x^2) \\).

## Example

To make the concept clearer, let us consider a simple example. Consider the function \\(f(x) = 3x^2 + 2x + 1\\) as \\(x \to \infty\\). We want to show that \\(f(x) = O(x^2)\\).

For large values of \\(x\\), we can analyze:

\\[\frac{3x^2 + 2x + 1}{x^2} = 3 + \frac{2}{x} + \frac{1}{x^2}\\]

As \\(x \to \infty\\), this ratio approaches 3. More precisely, for \\(x \geq 1\\):

\\[\frac{3x^2 + 2x + 1}{x^2} = 3 + \frac{2}{x} + \frac{1}{x^2} \leq 3 + 2 + 1 = 6\\]

Therefore, we can choose \\(M = 6\\) and conclude:

\\[3x^2 + 2x + 1 = O(x^2) \quad \text{as} \quad x \to \infty\\]

Let’s explore why this is the case. If we compare \\(f(x) = 3x^2 + 2x + 1\\) and \\(g(x) = x^2\\) as \\(x\\) approaches infinity, we observe that \\(f(x)\\) grows at most as fast as a constant multiple of \\(x^2\\).

Here are a few examples:

  * If \\(x = 10\\), then \\(f(x) = 300 + 20 + 1 = 321\\) and \\(6x^2 = 600\\).
  * If \\(x = 100\\), then \\(f(x) = 30000 + 200 + 1 = 30201\\) and \\(6x^2 = 60000\\).
  * If \\(x = 1000\\), then \\(f(x) = 3000000 + 2000 + 1 = 3002001\\) and \\(6x^2 = 6000000\\).


As these examples show, \\(f(x)\\) is always bounded above by \\(6x^2\\) for \\(x \geq 1\\). This is the key idea behind the big O notation: it captures the fact that one function grows at most as fast as another (up to a constant factor) as the input approaches a particular value.

## The meaning of \\( O(1) \\)

The symbol \\(O(1)\\) represents the class of functions that remain bounded as \\(x\\) approaches a specific point \\(x_0\\). In other words, a function \\(f(x)\\) is said to belong to \\(O(1)\\) when it remains bounded compared to a constant as \\(x \to x_0\\). Formally, we write \\(f(x) = O(1)\\) as \\(x \to x_0\\) if and only if:

\\[\exists M > 0, \delta > 0 : |x - x_0| < \delta \Rightarrow |f(x)| \leq M\\]

The set of all functions that belong to \\(O(1)\\) can be described as follows:

\\[O_{x_0}(1) = \left\lbrace f : B(x_0, \delta) \setminus {x_0} \to \mathbb{R} ,\Big|, \exists M > 0 : |f(x)| \leq M \text{ for } x \text{ near } x_0 \right\rbrace \\]

In practice, this expression is used to describe the concept of \\(O(1)\\) by specifying that:

  * The functions considered must be defined on a neighborhood of \\(x_0\\), excluding the point \\(x_0\\) itself.
  * The function must remain bounded as \\(x\\) approaches \\(x_0\\).
  * The notation \\(O(1)\\) represents the set of all functions that are bounded compared to a constant, specifically to 1.
  * The symbol \\(B(x_0, \delta)\\) denotes an open neighborhood of \\(x_0\\) with radius \\(\delta\\), where the function is defined.


## Example 1

Let us consider the function:

\\[f(x) = \frac{\sin(x)}{x} \quad \text{as} \quad x \to 0\\]

Using the Taylor expansion of \\(\sin(x)\\) near zero, as \\(x \to 0\\), we have:

\\[\sin(x) = x - \frac{x^3}{6} + O(x^5) \\]

Dividing both sides by \\(x\\), as \\(x \to 0\\), we get:

\\[\frac{\sin(x)}{x} = 1 - \frac{x^2}{6} + O(x^4) \\]

Now observe that the term \\(\dfrac{x^2}{6}\\) and the remainder \\(O(x^4)\\) both remain bounded as \\(x \to 0\\). Therefore, as \\(x \to 0\\) we can write:

\\[\frac{\sin(x)}{x} = O(1)\\]

This expression shows that the function remains bounded near \\(x = 0\\), even though the function has a removable singularity at that point

## Example 2

Big O notation frequently arises in the context of Taylor expansions to characterise the remainder term. For example, consider the Taylor expansion of \\( \cos(x) \\) near zero:

\\[\cos(x) = 1 - \frac{x^2}{2} + \frac{x^4}{24} - \cdots \\]

Truncating the series after the second term yields a remainder that is bounded by a constant multiple of \\( x^4 \\) for \\( x \\) near zero. Therefore, it can be expressed as:

\\[\cos(x) = 1 - \frac{x^2}{2} + O(x^4) \quad \text{as} \quad x \to 0 \\]

This means there exists a constant \\( M > 0 \\) such that:

\\[\left| \cos(x) - 1 + \frac{x^2}{2} \right| \leq M x^4 \\]

for all \\( x \\) in a neighborhood of zero. Since the next term in the series is \\( x^4/24 \\), it follows that \\( M = 1/24 \\) is sufficient for sufficiently small \\( x \\). This application of big O notation is common in analysis: instead of retaining the entire infinite series, only the terms of interest are kept, and all higher-order contributions are incorporated into a single remainder term \\( O(x^n) \\).

## Properties

One fundamental property of the big O notation is that, by definition, if \\(g(x) = O(f(x))\\) as \\(x \to x_0\\), then the ratio of the two functions remains bounded. In formal terms:

\\[\limsup_{x \to x_0} \left|\frac{O(f(x))}{f(x)}\right| < \infty\\]


Multiplying a function by a nonzero constant does not change its asymptotic behavior in big O notation. For any constant \\(c \neq 0\\) and function \\(g(x)\\), as \\(x \to x_0\\), we have:

\\[O(c \cdot g(x)) = O(g(x))\\] \\[c \cdot O(g(x)) = O(g(x))\\]

##### This shows that big O notation absorbs constant factors: scaling by a nonzero constant factor does not affect the asymptotic upper bound near \\(x_0\\).


Big O terms also behave predictably under addition. The sum of two big O terms of the same function remains a big O term of that function. Formally as \\(x \to x_0\\):

\\[O(f(x)) + O(f(x)) = O(f(x))\\]

##### This means that adding two functions that are each asymptotically bounded by \\(f(x)\\) results in a sum that is still asymptotically bounded by \\(f(x)\\) (possibly with a larger constant).


When multiplying a big O term by a function, the result is a new big O term where the asymptotic behavior scales accordingly. Specifically, for functions \\(f(x)\\) and \\(g(x)\\), as (x \to x_0):

\\[f(x) , O(g(x)) = O(f(x)g(x))\\]

For example, if \\(g(x) = x\\) and \\(O(g(x)) = O(x)\\), then multiplying by \\(f(x) = x^2\\) gives:

\\[x^2 \cdot O(x) = O(x^3)\\]


Another important property of the big O notation involves powers of functions. If a function \\(f(x)\\) is asymptotically bounded by \\(g(x)\\) as \\(x \to x_0\\), then raising both functions to the same positive power preserves the big O relationship. Formally, for \\(a > 0\\), if \\(f(x) = O(g(x)\\) as \\(x \to x_0\\), then:

\\[[f(x)]^a = O([g(x)]^a)\\]

This property shows that the asymptotic upper bound scales consistently under positive powers: if \\(f(x)\\) is bounded by a multiple of \\(g(x)\\), then \\([f(x)]^a\\) is also bounded by a multiple of \\([g(x)]^a\\) as \\(x \to x_0\\). For example, if \\(f(x) = O(x)\\) as \\(x \to \infty\\), then as \\(x \to \infty\\) we have:

\\[[f(x)]^2 = O(x^2)\\]


Big O notation demonstrates the property of transitivity. Specifically, if \\( f(x) = O(g(x)) \\) and \\( g(x) = O(h(x)) \\) as \\( x \to x_0 \\), then:

\\[f(x) = O(h(x)) \quad \text{as} \quad x \to x_0 \\]

The proof is straightforward. By assumption, there exist positive constants \\( M_1 \\) and \\( M_2 \\), as well as a neighborhood of \\( x_0 \\), such that \\( |f(x)| \leq M_1|g(x)| \\) and \\( |g(x)| \leq M_2|h(x)| \\) both hold.

Combining these inequalities gives \\( |f(x)| \leq M_1 M_2 |h(x)| \\), so \\( f(x) = O(h(x)) \\) with constant \\( M = M_1 M_2 \\). For example, \\( x^3 = O(x^4) \\) and \\( x^4 = O(x^5) \\) as \\( x \to \infty \\), so by transitivity, \\( x^3 = O(x^5) \\) as \\( x \to \infty \\).

##### Transitivity enables chaining of asymptotic bounds: if \\( f \\) is bounded by \\( g \\) and \\( g \\) is bounded by \\( h \\), then \\( f \\) is also bounded by \\( h \\).

## Relationship with little-o notation

It’s important to note the relationship between big O and [little-o notation](<../little-o-notation/>).

  * If \\(f(x) = o(g(x))\\), then \\(f(x) = O(g(x))\\): little-o implies big O.
  * The converse is not necessarily true: \\(f(x) = O(g(x))\\) does not imply \\(f(x) = o(g(x))\\) For example, \\(f(x) = x\\) and \\(g(x) = x\\) satisfy \\(f(x) = O(g(x))\\) but not \\(f(x) = o(g(x))\\) as \\(x \to \infty\\), since \\(\lim_{x \to \infty} \frac{x}{x} = 1 \neq 0.\\)


## Selected references

  * **MIT**. [Big O Notation](https://web.mit.edu/16.070/www/lecture/big_o.pdf)

  * **University of Illinois, A. J. Hildebrand**. [Asymptotic Notations](https://faculty.math.illinois.edu/~hildebr/595ama/ama-ch2.pdf)

  * **Rice University, J. A. Dobelman**. [Big O and Little o](https://www.stat.rice.edu/~dobelman/notes_papers/math/big_O.little_o.pdf)

  * **Simon Fraser University, R. Lockhart**. [Landau Notation: Big O and Little o](https://www.sfu.ca/~lockhart/richard/830/20_3/lectures/Landau/web.pdf)


Limits

Limits describe a function’s behavior near a point.

5.9k

[Limits](https://algebrica.org/limits/)

942

[Algebra of Limits](https://algebrica.org/algebra-of-limits/)

1.6k

[Squeeze Theorem](https://algebrica.org/squeeze-theorem/)

4.1k

[Remarkable Limits](https://algebrica.org/remarkable-limits/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/limits/remarkable-limits.md?plain=1)

3.8k

[Asymptotes](https://algebrica.org/asymptotes/)

1.9k

[Indeterminate Forms of Limits](https://algebrica.org/indeterminate-forms/)

5.1k

[Little-o Notation](https://algebrica.org/little-o-notation/)
