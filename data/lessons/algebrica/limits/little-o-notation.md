> Content sourced from [Algebrica](https://algebrica.org/little-o-notation/) — CC BY-NC 4.0

## What is little-o notation

The symbol \\( o(x) \\), referred to as little-o of \\(x\\), belongs to the Landau symbol family, which is used to characterise asymptotic relationships between functions. This notation indicates that one function is negligible compared to another as the input approaches a given limit. Thus, \\( o(x) \\) formalises [asymptotic](<../asymptotes/>) control by signifying that the growth rate of one function is insignificant relative to the other in the [limit](<../limits/>).

Let \\( f, g : A \to \mathbb{R} \\) (or \\( \mathbb{C} \\)) be two [functions](<../functions/>), and let \\( x_0 \\) be a limit point of \\( A \\). We say that \\( f(x) \\) is little-o of \\( g(x) \\) as \\( x \to x_0 \\) if \\( g(x) \neq 0 \\) in a neighborhood of \\( x_0 \\) (except possibly at \\( x_0 \\) itself) and:

\\[\lim_{x \to x_0} \frac{f(x)}{g(x)} = 0 \\]

Equivalently, for every \\( \varepsilon > 0 \\) there exists \\( \delta > 0 \\) such that whenever \\( 0 < |x - x_0| < \delta \\), we have \\( |f(x)| \le \varepsilon \cdot |g(x)| \\). This definition means that \\( f(x) \\) grows asymptotically slower than \\( g(x) \\) near \\( x_0 \\).

> This notation applies to limits at infinity by replacing \\( x \to x_0 \\) with \\( x \to \infty \\). It also applies to [sequences](<../sequences/>), where the independent variable \\( x \\) is replaced by \\( n \to \infty \\).

## Example 1

To make the concept clearer, let us consider a simple example given by the following limit:

\\[\lim_{x \to 0} \frac{x^2}{x} = \lim_{x \to 0} x = 0 \\]

This limit demonstrates that \\( x^2 \\) grows asymptotically slower than \\( x \\) as \\( x \\) approaches 0. Therefore, we can write:

\\[x^2 = o(x) \quad \text{as} \quad x \to 0 \\]

Let’s explore why this is the case. If we compare \\(x\\) and \\(x^2\\) as \\(x\\) approaches zero, we observe that both functions tend to zero, but at different rates. The graph clearly shows that, near zero, \\(x^2\\) is much smaller than \\(x\\).

![](https://algebrica.org/wp-content/uploads/resources/images/little-o-1.png)

Here are a few examples:

  * If \\(x = 0.1\\), then \\(x = 0.1\\) and \\(x^2 = 0.01\\).
  * If \\(x = 0.01\\), then \\(x = 0.01\\) and \\(x^2 = 0.0001\\).
  * If \\(x = 0.001\\), then \\(x = 0.001\\) and \\(x^2 = 0.000001\\).


As these examples show, \\(x^2\\) becomes much smaller than \\(x\\) as \\(x\\) approaches zero. This is the key idea behind the little-o notation: it captures the fact that one function can become asymptotically negligible compared to another as the input approaches a particular value.

## Example 2

Little-o notation remains applicable when the input increases without bound. The following limit illustrates this as \\( x \to \infty \\):

\\[\lim_{x \to \infty} \frac{x}{x^2} = \lim_{x \to \infty} \frac{1}{x} = 0 \\]

Because the ratio approaches zero, the following expression holds:

\\[x = o(x^2) \quad \text{as} \quad x \to \infty \\]

This result indicates that \\( x \\) grows asymptotically more slowly than \\( x^2 \\) as \\( x \\) increases without bound. More generally, for any two polynomials \\( x^a \\) and \\( x^b \\) with \\( a < b \\), the following relationship holds:

\\[x^a = o(x^b) \quad \text{as} \quad x \to \infty \\]

Another example compares a logarithmic function with a power function. Since:

\\[\lim_{x \to \infty} \frac{\log x}{x} = 0 \\]

it follows that \\( \log x = o(x) \\) as \\( x \to \infty \\). This result is particularly relevant in algorithm analysis, as it confirms that logarithmic growth is strictly dominated by linear growth.

> The little-o condition at infinity parallels the condition at a finite point: the ratio of the two functions must approach zero as the input increases without bound, rather than as it approaches a fixed value.

## The meaning of \\( o(1) \\)

The symbol \\( o(1) \\) represents the class of functions that tend to zero as \\( x \\) approaches a specific point \\( x_0 \\). In other words, a function \\( f(x) \\) is said to belong to \\( o(1) \\) when it becomes infinitesimally small compared to a constant, specifically 1, in the limit \\( x \to x_0 \\). Formally, we write \\(f(x) = o(1)\\) as \\(x \to x_0\\) if and only if:

\\[\lim_{x \to x_0} \frac{f(x)}{1} = \lim_{x \to x_0} f(x) = 0 \\]

The set of all functions that belong to \\( o(1) \\) can be described as follows:

\\[o_{x_0}(1) = \lbrace f : B(x_0, \delta) \setminus \\{x_0\\} \to \mathbb{R} \,\Big|\, \lim_{x \to x_0} f(x) = 0 \rbrace \\]

In practice, this expression is used to describe the concept of \\( o(1) \\) by specifying that:

  * The functions considered must be defined on a neighborhood of \\( x_0 \\), excluding the point \\( x_0 \\) itself.
  * The function must tend to zero as \\( x \\) approaches \\( x_0 \\).
  * The notation \\( o(1) \\) represents the set of all functions that are infinitesimal compared to a constant, specifically to \\(1\\).
  * The symbol \\( B(x_0, \delta) \\) denotes an open neighborhood of \\( x_0 \\) with radius \\( \delta \\), where the function is defined and the limit is taken.


## Example 3

Let us consider the limit:

\\[\lim_{x \to 0} \frac{\sin(x)}{x} \\]

Using the Taylor expansion of \\( \sin(x) \\) near zero, we have:

\\[\sin(x) = x - \frac{x^3}{6} + o(x^3) \quad \text{as} \quad x \to 0 \\]


Dividing both sides by \\( x \\), we get:

\\[\frac{\sin(x)}{x} = 1 - \frac{x^2}{6} + o(x^2) \quad \text{as} \quad x \to 0 \\]

Now observe that both the term \\( \frac{x^2}{6} \\) and the remainder \\( o(x^2) \\) tend to zero as \\( x \to 0 \\).

Therefore, as \\(x \to 0\\) we can write:

\\[\frac{\sin(x)}{x} = 1 + o(1) \\]

> This expression shows that the difference between \\( \frac{\sin(x)}{x} \\) and the constant \\(1\\) tends to zero in the limit, and the correction terms are asymptotically smaller than 1.

## Properties

One fundamental property of the little-o notation is that, by definition, if \\( g(x) = o(f(x)) \\) as \\( x \to x_0 \\), then the ratio of the two functions tends to zero. In formal terms:

\\[\lim_{x \to x_0} \frac{o(f(x))}{f(x)} = 0 \\]


Multiplying a function by a nonzero constant does not change its asymptotic behavior in little-o notation. For any constant \\( c \in \mathbb{R} \\) and function \\( g(x) \\), as \\(x \to x_0\\), we have:

\\[o(c \cdot g(x)) = o(g(x)) \\]

\\[c \cdot o(g(x)) = o(g(x)) \\]

> This shows that little-o notation is about relative growth: scaling by a constant factor does not affect the asymptotic behavior near \\( x_0 \\).


Little-o terms also behave predictably under addition. The sum of two little-o terms of the same function remains a little-o term of that function. Formally as \\(x \to x_0\\):

\\[o(f(x)) + o(f(x)) = o(f(x)) \\]

> This means that adding two functions that are each asymptotically smaller than \\( f(x) \\) does not affect the fact that the sum is still asymptotically smaller than \\( f(x) \\).


When multiplying a little-o term by a function, the result is a new little-o term where the asymptotic behavior scales accordingly. Specifically, for functions \\( f(x) \\) and \\( g(x) \\), as \\(x \to x_0\\):

\\[f(x) \, o(g(x)) = o(f(x)g(x)) \\]

For example, if \\( g(x) = x \\) and \\( o(g(x)) = o(x) \\), then multiplying by \\( f(x) = x^2 \\) gives:

\\[x^2 \cdot o(x) = o(x^3) \\]


Another important property of the little-o notation involves powers of functions. If a function \\( f(x) \\) is asymptotically smaller than \\( g(x) \\) as \\( x \to x_0 \\), then raising both functions to the same positive power preserves the little-o relationship. Formally, for \\( a > 0 \\) if \\(f(x) = o(g(x)) \\) as \\(x \to x_0\\) then \\([f(x)]^a = o([g(x)]^a)\\) as \\(x \to x_0\\).

This property demonstrates that the asymptotic behaviour remains consistent when scaled by positive powers. if \\( f(x) \\) becomes negligible compared to \\( g(x) \\), then \\( [f(x)]^a \\) is also negligible compared to \\( [g(x)]^a \\) as \\( x \to x_0 \\). For example, if \\( f(x) = o(x) \\) as \\( x \to 0 \\), then as \\(x \to x_0\\) we have:

\\[[f(x)]^2 = o(x^2) \\]


Little-o notation exhibits transitivity. Specifically, if \\( f(x) = o(g(x)) \\) and \\( g(x) = o(h(x)) \\) as \\( x \to x_0 \\), then:

\\[f(x) = o(h(x)) \quad \text{as} \quad x \to x_0 \\]

This result follows directly from the definition. Since both ratios approach zero, their product also approaches zero, and therefore \\( f(x)/h(x) \to 0 \\). For example, since \\( x^3 = o(x^2) \\) and \\( x^2 = o(x) \\) as \\( x \to 0 \\), it follows that \\( x^3 = o(x) \\) as \\( x \to 0 \\).

> This property enables the composition of chains of asymptotic comparisons: if \\(f\\) grows more slowly than \\(g\\), and \\(g\\) grows more slowly than \\(h\\), then \\(f\\) also grows more slowly than \\(h\\).


The composition of two little-o terms reduces to a single term. Specifically, if \\( h(x) = o(g(x)) \\) and \\( g(x) = o(f(x)) \\) as \\( x \to x_0 \\), then any function that is little-o of \\( g \\) is also little-o of \\( f \\). In compact notation:

\\[o(o(f(x))) = o(f(x)) \quad \text{as} \quad x \to x_0 \\]

This result follows directly from the transitivity property: if \\( h = o(g) \\) and \\( g = o(f) \\), then \\( h = o(f) \\). For example, since \\( x^2 = o(x) \\) as \\( x \to 0 \\), any function that is \\( o(x^2) \\) is also \\( o(x) \\).

> This property is especially useful for simplifying nested asymptotic expressions, as it ensures that iterated little-o terms can be represented by a single term.

## Distinction Between Little-o and Big-O Notation

Little-o and [Big-O notation](<../big-o-notation/>) are both members of the Landau symbol family, but they describe distinct asymptotic behaviours. Big-O notation provides an upper bound for a function up to a constant multiple, whereas little-o notation imposes a stricter requirement: the ratio of the two functions must approach zero in the limit.

Formally, \\( f(x) = O(g(x)) \\) as \\( x \to x_0 \\) if there exist constants \\( M > 0 \\) and \\( \delta > 0 \\) such that: \\(|f(x)| \leq M |g(x)| \\) whenever \\(0 < |x - x_0| < \delta\\). This distinction is illustrated by the following example.

As \\( x \to 0 \\), \\( x^2 = o(x) \\), which also implies \\( x^2 = O(x) \\). However, \\( x = O(x) \\) does not imply \\( x = o(x) \\), as demonstrated below:

\\[\lim_{x \to 0} \frac{x}{x} = 1 \neq 0 \\]

The limit fails to vanish, so the little-o condition is not satisfied. This asymmetry is the heart of the distinction: little-o requires the ratio to go to zero, while Big-O only requires it to stay bounded.

This relationship can be visualised in terms of set inclusion: the set of functions satisfying \\( f = o(g) \\) is strictly contained within the set of functions satisfying \\( f = O(g) \\). Every little-o relationship is also a Big-O relationship, but the converse does not hold.

> Little-o notation excludes functions that merely keep pace with \\( g \\), admitting only those that fall strictly behind in the limit.

## Little-o Notation in Taylor Expansions

In [Taylor expansions](<../taylor-series/>), little-o notation provides a rigorous framework for quantifying the error introduced by truncating an infinite series at a finite order. Instead of enumerating each remaining term, the remainder is represented by a single symbol that specifies its precise asymptotic order. Given a function \\( f(x) \\) that is \\( n \\) times differentiable at \\( x_0 \\), its Taylor expansion to order \\( n \\) takes the form:

\\[f(x) = f(x_0) + f’(x_0)(x - x_0) + \frac{f’'(x_0)}{2!}(x - x_0)^2 + \cdots + \frac{f^{(n)}(x_0)}{n!}(x - x_0)^n + o\bigl((x - x_0)^n\bigr) \\]

The remainder \\( o((x - x_0)^n) \\) conveys precise asymptotic information. Specifically, the error decreases more rapidly than \\( (x - x_0)^n \\) as \\( x \to x_0 \\), rendering it negligible compared to the last explicit term in the expansion.


The following table presents the Taylor expansions of commonly encountered functions near \\( x = 0 \\), each expressed with an explicit little-o remainder:

|   
---|---  
\\[e^x \\]| \\[1 + x + \dfrac{x^2}{2!} + \dfrac{x^3}{3!} + o(x^3) \\]  
\\[\sin x \\]| \\[x - \dfrac{x^3}{6} + o(x^3) \\]  
\\[\cos x \\]| \\[1 - \dfrac{x^2}{2} + \dfrac{x^4}{24} + o(x^4) \\]  
\\[\ln(1+x) \\]| \\[x - \dfrac{x^2}{2} + \dfrac{x^3}{3} + o(x^3) \\]  
\\[(1+x)^\alpha \\]| \\[1 + \alpha x + \dfrac{\alpha(\alpha-1)}{2}x^2 + o(x^2) \\]  
  
These expansions are particularly effective for evaluating limits involving indeterminate forms. Replacing the original function with its Taylor expansion transforms the problem into an algebraic manipulation, where the little-o remainder vanishes in the limit. The following example demonstrates this approach:

\\[\lim_{x \to 0} \frac{e^x - 1 - x}{x^2} = \lim_{x \to 0} \frac{\dfrac{x^2}{2} + o(x^2)}{x^2} = \frac{1}{2} \\]

As \\( x \to 0 \\), the little-o term becomes negligible, allowing the limit to be determined by the leading coefficient.

> In a [Taylor expansion](<../taylor-series/>), the little-o remainder does not simply indicate omitted terms; it specifies the rate at which the approximation improves, thereby anchoring the truncation to a precise asymptotic scale.

Definition

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Advanced

2

Requires

0

Enables

The following concepts, [Functions](https://algebrica.org/functions/), [Limits](https://algebrica.org/limits/), are required as prerequisites for this entry.

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

1.8k

[Big O Notation](https://algebrica.org/big-o-notation/)
