> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 3.3 Differentiation Rules

### Learning Objectives

- 3.3.1 State the constant, constant multiple, and power rules.
- 3.3.2 Apply the sum and difference rules to combine derivatives.
- 3.3.3 Use the product rule for finding the derivative of a product of functions.
- 3.3.4 Use the quotient rule for finding the derivative of a quotient of functions.
- 3.3.5 Extend the power rule to functions with negative exponents.
- 3.3.6 Combine the differentiation rules to find the derivative of a polynomial or rational function.

Finding derivatives of functions by using the definition of the derivative can be a lengthy and, for certain functions, a rather challenging process. For example, previously we found that $\frac{d}{dx}\left( \sqrt{x} \right) = \frac{1}{2\sqrt{x}}$ by using a process that involved multiplying an expression by a conjugate prior to evaluating a limit. The process that we could use to evaluate $\frac{d}{dx}\left( \sqrt[3]{x} \right)$ using the definition, while similar, is more complicated. In this section, we develop rules for finding derivatives that allow us to bypass this process. We begin with the basics.

### The Basic Rules

The functions $f(x) = c$ and $g(x) = x^{n}$ where $n$ is a positive integer are the building blocks from which all polynomials and rational functions are constructed. To find derivatives of polynomials and rational functions efficiently without resorting to the limit definition of the derivative, we must first develop formulas for differentiating these basic functions.

#### The Constant Rule

We first apply the limit definition of the derivative to find the derivative of the constant function, $f(x) = c.$ For this function, both $f(x) = c$ and $f\left( {x + h} \right) = c,$ so we obtain the following result:

$$
\begin{array}{cl}
{f'(x)} & {= \underset{h\rightarrow 0}{\text{lim}}\frac{f\left( {x + h} \right) - f(x)}{h}} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{c - c}{h}} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{0}{h}} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}0 = 0.}
\end{array}
$$

The rule for differentiating constant functions is called the constant rule. It states that the derivative of a constant function is zero; that is, since a constant function is a horizontal line, the slope, or the rate of change, of a constant function is $0.$ We restate this rule in the following theorem.

### Theorem 3.2

#### The Constant Rule

Let $c$ be a constant.

If $f(x) = c,$ then $f'(x) = 0.$

Alternatively, we may express this rule as

$$
\frac{d}{dx}(c) = 0.
$$

### Example 3.17

#### Applying the Constant Rule

Find the derivative of $f(x) = 8.$

#### Solution

This is just a one-step application of the rule:

$$
f'(x) = 0.
$$

### Checkpoint 3.11

Find the derivative of $g(x) = -3.$

### The Power Rule

We have shown that

$$
\frac{d}{dx}\left( x^{2} \right) = 2x\ \text{and}\ \frac{d}{dx}\left( x^{1\text{/}2} \right) = \frac{1}{2}x^{\text{−}{1\text{/}2}}.
$$

At this point, you might see a pattern beginning to develop for derivatives of the form $\frac{d}{dx}\left( x^{n} \right).$ We continue our examination of derivative formulas by differentiating power functions of the form $f(x) = x^{n}$ where $n$ is a positive integer. We develop formulas for derivatives of this type of function in stages, beginning with positive integer powers. Before stating and proving the general rule for derivatives of functions of this form, we take a look at a specific case, $\frac{d}{dx}(x^{3}).$ As we go through this derivation, note that the technique used in this case is essentially the same as the technique used to prove the general case.

### Example 3.18

#### Differentiating $x^{3}$

Find $\frac{d}{dx}\left( x^{3} \right).$

#### Solution

$$
\begin{array}{clccc}
{\frac{d}{dx}\left( x^{3} \right)} & {= \underset{h\rightarrow 0}{\text{lim}}\frac{{(x + h)}^{3} - x^{3}}{h}} & & & \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{x^{3} + 3x^{2}h + 3xh^{2} + h^{3} - x^{3}}{h}} & & & \begin{array}{l}
\text{Notice that the first term in the expansion of} \\
{{(x + h)}^{3}\ \text{is}\ x^{3}\ \text{and the second term is}\ 3x^{2}h.\ \text{All}} \\
{\text{other terms contain powers of}\ h\ \text{that are two or}} \\
\text{greater.}
\end{array} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{3x^{2}h + 3xh^{2} + h^{3}}{h}} & & & \begin{array}{l}
{\text{In this step the}\ x^{3}\ \text{terms have been cancelled,}} \\
{\text{leaving only terms containing}\ h.}
\end{array} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{h(3x^{2} + 3xh + h^{2})}{h}} & & & {\text{Factor out the common factor of}\ h.} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}(3x^{2} + 3xh + h^{2})} & & & \begin{array}{l}
{\text{After cancelling the common factor of}\ h,\text{the}} \\
{\text{only term not containing}\ h\ \text{is}\ 3x^{2}.}
\end{array} \\
 & {= 3x^{2}} & & & {\text{Let}\ h\ \text{go to 0.}}
\end{array}
$$

### Checkpoint 3.12

Find $\frac{d}{dx}\left( x^{4} \right).$

As we shall see, the procedure for finding the derivative of the general form $f(x) = x^{n}$ is very similar. Although it is often unwise to draw general conclusions from specific examples, we note that when we differentiate $f(x) = x^{3},$ the power on $x$ becomes the coefficient of $x^{2}$ in the derivative and the power on $x$ in the derivative decreases by 1. The following theorem states that the power rule holds for all positive integer powers of $x.$ We will eventually extend this result to negative integer powers. Later, we will see that this rule may also be extended first to rational powers of $x$ and then to arbitrary powers of $x.$ Be aware, however, that this rule does not apply to functions in which a constant is raised to a variable power, such as $f(x) = 3^{x}.$

### Theorem 3.3

#### The Power Rule

Let $n$ be a positive integer. If $f(x) = x^{n},$ then

$$
f'(x) = nx^{n - 1}.
$$

Alternatively, we may express this rule as

$$
\frac{d}{dx}x^{n} = nx^{n - 1}.
$$

#### Proof

For $f(x) = x^{n}$ where $n$ is a positive integer, we have

$$
f'(x) = \underset{h\rightarrow 0}{\text{lim}}\frac{{(x + h)}^{n} - x^{n}}{h}.
$$

$$
\text{Since}\ {(x + h)}^{n} = x^{n} + nx^{n - 1}h + \begin{pmatrix}
n \\
2
\end{pmatrix}\ x^{n - 2}h^{2} + \begin{pmatrix}
n \\
3
\end{pmatrix}\ x^{n - 3}h^{3} + \text{…} + nxh^{n - 1} + h^{n},
$$

we see that

$$
{(x + h)}^{n} - x^{n} = nx^{n - 1}h + \begin{pmatrix}
n \\
2
\end{pmatrix}\ x^{n - 2}h^{2} + \begin{pmatrix}
n \\
3
\end{pmatrix}\ x^{n - 3}h^{3} + \text{…} + nxh^{n - 1} + h^{n}.
$$

Next, divide both sides by *h*:

$$
\frac{\left( {x + h} \right)^{n} - x^{n}}{h} = \frac{nx^{n - 1}h + \begin{pmatrix}
n \\
2
\end{pmatrix}\ x^{n - 2}h^{2} + \begin{pmatrix}
n \\
3
\end{pmatrix}\ x^{n - 3}h^{3} + \text{…} + nxh^{n - 1} + h^{n}}{h}.
$$

Thus,

$$
\frac{\left( {x + h} \right)^{n} - x^{n}}{h} = nx^{n - 1} + \begin{pmatrix}
n \\
2
\end{pmatrix}\ x^{n - 2}h + \begin{pmatrix}
n \\
3
\end{pmatrix}\ x^{n - 3}h^{2} + \text{…} + nxh^{n - 2} + h^{n - 1}.
$$

Finally,

$$
\begin{array}{cl}
{f'(x)} & {= \underset{h\rightarrow 0}{\text{lim}}\left( {nx^{n - 1} + \begin{pmatrix}
n \\
2
\end{pmatrix}\ x^{n - 2}h + \begin{pmatrix}
n \\
3
\end{pmatrix}\ x^{n - 3}h^{2} + \text{…} + nxh^{n - 2} + h^{n - 1}} \right)} \\
 & {= nx^{n - 1}.}
\end{array}
$$

□

### Example 3.19

#### Applying the Power Rule

Find the derivative of the function $f(x) = x^{10}$ by applying the power rule.

#### Solution

Using the power rule with $n = 10,$ we obtain

$$
f'(x) = 10x^{10 - 1} = 10x^{9}.
$$

### Checkpoint 3.13

Find the derivative of $f(x) = x^{7}.$

### The Sum, Difference, and Constant Multiple Rules

We find our next differentiation rules by looking at derivatives of sums, differences, and constant multiples of functions. Just as when we work with functions, there are rules that make it easier to find derivatives of functions that we add, subtract, or multiply by a constant. These rules are summarized in the following theorem.

### Theorem 3.4

#### Sum, Difference, and Constant Multiple Rules

Let $f(x)$ and $g(x)$ be differentiable functions and $k$ be a constant. Then each of the following equations holds.

Sum Rule. The derivative of the sum of a function $f$ and a function $g$ is the same as the sum of the derivative of $f$ and the derivative of $g.$

$$
\frac{d}{dx}\left( {f(x) + g(x)} \right) = \frac{d}{dx}\left( {f(x)} \right) + \frac{d}{dx}\left( {g(x)} \right);
$$

that is,

$$
\text{for}\ j(x) = f(x) + g(x),j'(x) = f'(x) + g'(x).
$$

Difference Rule. The derivative of the difference of a function *f* and a function *g* is the same as the difference of the derivative of *f* and the derivative of $g\text{:}$

$$
\frac{d}{dx}\left( {f(x) - g(x)} \right) = \frac{d}{dx}\left( {f(x)} \right) - \frac{d}{dx}\left( {g(x)} \right);
$$

that is,

$$
\text{for}\ j(x) = f(x) - g(x),j'(x) = f'(x) - g'(x).
$$

Constant Multiple Rule. The derivative of a constant *k* multiplied by a function *f* is the same as the constant multiplied by the derivative:

$$
\frac{d}{dx}\left( {kf(x)} \right) = k\frac{d}{dx}\left( {f(x)} \right);
$$

that is,

$$
\text{for}\ j(x) = kf(x),j'(x) = kf'(x).
$$

#### Proof

We provide only the proof of the sum rule here. The rest follow in a similar manner.

For differentiable functions $f(x)$ and $g(x),$ we set $j(x) = f(x) + g(x).$ Using the limit definition of the derivative we have

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\frac{j\left( {x + h} \right) - j(x)}{h}.
$$

By substituting $j\left( {x + h} \right) = f\left( {x + h} \right) + g\left( {x + h} \right)$ and $j(x) = f(x) + g(x),$ we obtain

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\frac{\left( {f\left( {x + h} \right) + g\left( {x + h} \right)} \right) - \left( {f(x) + g(x)} \right)}{h}.
$$

Rearranging and regrouping the terms, we have

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\left( {\frac{f\left( {x + h} \right) - f(x)}{h} + \frac{g\left( {x + h} \right) - g(x)}{h}} \right).
$$

We now apply the sum law for limits and the definition of the derivative to obtain

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\left( \frac{f\left( {x + h} \right) - f(x)}{h} \right) + \underset{h\rightarrow 0}{\text{lim}}\left( \frac{g\left( {x + h} \right) - g(x)}{h} \right) = f'(x) + g'(x).
$$

□

### Example 3.20

#### Applying the Constant Multiple Rule

Find the derivative of $g(x) = 3x^{2}$ and compare it to the derivative of $f(x) = x^{2}.$

#### Solution

We use the power rule directly:

$$
g'(x) = \frac{d}{dx}\left( {3x^{2}} \right) = 3\frac{d}{dx}\left( x^{2} \right) = 3\left( {2x} \right) = 6x.
$$

Since $f(x) = x^{2}$ has derivative $f'(x) = 2x,$ we see that the derivative of $g(x)$ is 3 times the derivative of $f(x).$ This relationship is illustrated in Figure 3.18.

*Figure 3.18 The derivative of g ( x ) g ( x ) is 3 times the derivative of f ( x ) . f ( x ) .*

### Example 3.21

#### Applying Basic Derivative Rules

Find the derivative of $f(x) = 2x^{5} + 7.$

#### Solution

We begin by applying the rule for differentiating the sum of two functions, followed by the rules for differentiating constant multiples of functions and the rule for differentiating powers. To better understand the sequence in which the differentiation rules are applied, we use Leibniz notation throughout the solution:

$$
\begin{array}{clccc}
{f'(x)} & {= \frac{d}{dx}\left( {2x^{5} + 7} \right)} & & & \\
 & {= \frac{d}{dx}\left( {2x^{5}} \right) + \frac{d}{dx}(7)} & & & \text{Apply the sum rule.} \\
 & {= 2\frac{d}{dx}\left( x^{5} \right) + \frac{d}{dx}(7)} & & & \text{Apply the constant multiple rule.} \\
 & {= 2\left( {5x^{4}} \right) + 0} & & & \text{Apply the power rule and the constant rule.} \\
 & {= 10x^{4}.} & & & \text{Simplify.}
\end{array}
$$

### Checkpoint 3.14

Find the derivative of $f(x) = 2x^{3} - 6x^{2} + 3.$

### Example 3.22

#### Finding the Equation of a Tangent Line

Find an equation of the line tangent to the graph of $f(x) = x^{2} - 4x + 6$ at $x = 1.$

#### Solution

To find an equation of the tangent line, we need a point and a slope. To find the point, compute

$$
f(1) = 1^{2} - 4(1) + 6 = 3.
$$

This gives us the point $\left( {1,3} \right).$ Since the slope of the tangent line at 1 is $f'(1),$ we must first find $f'(x).$ Using the definition of a derivative, we have

$$
f'(x) = 2x - 4
$$

so the slope of the tangent line is $f'(1) = -2.$ Using the point-slope formula, we see that the equation of the tangent line is

$$
y - 3 = -2\left( {x - 1} \right).
$$

Putting the equation of the line in slope-intercept form, we obtain

$$
y = -2x + 5.
$$

### Checkpoint 3.15

Find an equation of the line tangent to the graph of $f(x) = 3x^{2} - 11$ at $x = 2.$ Use the point-slope form.

### The Product Rule

Now that we have examined the basic rules, we can begin looking at some of the more advanced rules. The first one examines the derivative of the product of two functions. Although it might be tempting to assume that the derivative of the product is the product of the derivatives, similar to the sum and difference rules, the product rule does not follow this pattern. To see why we cannot use this pattern, consider the function $f(x) = x^{2},$ whose derivative is $f'(x) = 2x$ and not $\frac{d}{dx}(x) \cdot \frac{d}{dx}(x) = 1 \cdot 1 = 1.$

### Theorem 3.5

#### Product Rule

Let $f(x)$ and $g(x)$ be differentiable functions. Then

$$
\frac{d}{dx}\left( {f(x)g(x)} \right) = \frac{d}{dx}\left( {f(x)} \right) \cdot g(x) + \frac{d}{dx}\left( {g(x)} \right) \cdot f(x).
$$

That is,

$$
\text{if}\ j(x) = f(x)g(x),\text{then}\ j'(x) = f'(x)g(x) + g'(x)f(x).
$$

This means that the derivative of a product of two functions is the derivative of the first function times the second function plus the derivative of the second function times the first function.

#### Proof

We begin by assuming that $f(x)$ and $g(x)$ are differentiable functions. At a key point in this proof we need to use the fact that, since $g(x)$ is differentiable, it is also continuous. In particular, we use the fact that since $g(x)$ is continuous, $\underset{h\rightarrow 0}{\text{lim}}g\left( {x + h} \right) = g(x).$

By applying the limit definition of the derivative to $j(x) = f(x)g(x),$ we obtain

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\frac{f\left( {x + h} \right)g\left( {x + h} \right) - f(x)g(x)}{h}.
$$

By adding and subtracting $f(x)g(x + h)$ in the numerator, we have

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\frac{f\left( {x + h} \right)g\left( {x + h} \right) - f(x)g\left( {x + h} \right) + f(x)g\left( {x + h} \right) - f(x)g(x)}{h}.
$$

After breaking apart this quotient and applying the sum law for limits, the derivative becomes

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\left( \frac{f\left( {x + h} \right)g\left( {x + h} \right) - f(x)g\left( {x + h} \right)}{h} \right) + \underset{h\rightarrow 0}{\text{lim}}\left( \frac{f(x)g\left( {x + h} \right) - f(x)g(x)}{h} \right).
$$

Rearranging, we obtain

$$
j'(x) = \underset{h\rightarrow 0}{\text{lim}}\left( {\frac{f\left( {x + h} \right) - f(x)}{h} \cdot g(x + h)} \right) + \underset{h\rightarrow 0}{\text{lim}}\left( {\frac{g\left( {x + h} \right) - g(x)}{h} \cdot f(x)} \right).
$$

By using the continuity of $g(x),$ the definition of the derivatives of $f(x)$ and $g(x),$ and applying the limit laws, we arrive at the product rule,

$$
j'(x) = f'(x)g(x) + g'(x)f(x).
$$

□

### Example 3.23

#### Applying the Product Rule to Functions at a Point

For $j(x) = f(x)g(x),$ use the product rule to find $j'(2)$ if $f(2) = 3,f'(2) = -4,g(2) = 1,$ and $g'(2) = 6.$

#### Solution

Since $j(x) = f(x)g(x),j'(x) = f'(x)g(x) + g'(x)f(x),$ and hence

$$
j'(2) = f'(2)g(2) + g'(2)f(2) = (-4)(1) + (6)(3) = 14.
$$

### Example 3.24

#### Applying the Product Rule to Binomials

For $j(x) = (x^{2} + 2)(3x^{3} - 5x),$ find $j'(x)$ by applying the product rule. Check the result by first finding the product and then differentiating.

#### Solution

If we set $f(x) = x^{2} + 2$ and $g(x) = 3x^{3} - 5x,$ then $f'(x) = 2x$ and $g'(x) = 9x^{2} - 5.$ Thus,

$$
j'(x) = f'(x)g(x) + g'(x)f(x) = \left( {2x} \right)\left( {3x^{3} - 5x} \right) + (9x^{2} - 5)(x^{2} + 2).
$$

Simplifying, we have

$$
j'(x) = 15x^{4} + 3x^{2} - 10.
$$

To check, we see that $j(x) = 3x^{5} + x^{3} - 10x$ and, consequently, $j'(x) = 15x^{4} + 3x^{2} - 10.$

### Checkpoint 3.16

Use the product rule to obtain the derivative of $j(x) = 2x^{5}\left( {4x^{2} + x} \right).$

### The Quotient Rule

Having developed and practiced the product rule, we now consider differentiating quotients of functions. As we see in the following theorem, the derivative of the quotient is not the quotient of the derivatives; rather, it is the derivative of the function in the numerator times the function in the denominator minus the derivative of the function in the denominator times the function in the numerator, all divided by the square of the function in the denominator. In order to better grasp why we cannot simply take the quotient of the derivatives, keep in mind that

$$
\frac{d}{dx}\left( x^{2} \right) = 2x,\text{not}\ \frac{\frac{d}{dx}\left( x^{3} \right)}{\frac{d}{dx}(x)} = \frac{3x^{2}}{1} = 3x^{2}.
$$

### Theorem 3.6

#### The Quotient Rule

Let $f(x)$ and $g(x)$ be differentiable functions. Then

$$
\frac{d}{dx}\left( \frac{f(x)}{g(x)} \right) = \frac{\frac{d}{dx}(f(x)) \cdot g(x) - \frac{d}{dx}(g(x)) \cdot f(x)}{{(g(x))}^{2}}.
$$

That is,

$$
\text{if}\ j(x) = \frac{f(x)}{g(x)},\text{then}\ j'(x) = \frac{f'(x)g(x) - g'(x)f(x)}{{(g(x))}^{2}}.
$$

The proof of the quotient rule is very similar to the proof of the product rule, so it is omitted here. Instead, we apply this new rule for finding derivatives in the next example.

### Example 3.25

#### Applying the Quotient Rule

Use the quotient rule to find the derivative of $k(x) = \frac{5x^{2}}{4x + 3}.$

#### Solution

Let $f(x) = 5x^{2}$ and $g(x) = 4x + 3.$ Thus, $f'(x) = 10x$ and $g'(x) = 4.$ Substituting into the quotient rule, we have

$$
k'(x) = \frac{f'(x)g(x) - g'(x)f(x)}{{(g(x))}^{2}} = \frac{10x\left( {4x + 3} \right) - 4(5x^{2})}{{(4x + 3)}^{2}}.
$$

Simplifying, we obtain

$$
k'(x) = \frac{20x^{2} + 30x}{{(4x + 3)}^{2}}.
$$

### Checkpoint 3.17

Find the derivative of $h(x) = \frac{3x + 1}{4x - 3}.$

It is now possible to use the quotient rule to extend the power rule to find derivatives of functions of the form $x^{k}$ where $k$ is a negative integer.

### Theorem 3.7

#### Extended Power Rule

If $k$ is a negative integer, then

$$
\frac{d}{dx}\left( x^{k} \right) = kx^{k - 1}.
$$

#### Proof

If $k$ is a negative integer, we may set $n = \text{-}k,$ so that *n* is a positive integer with $k = \text{-}n.$ Since for each positive integer $n,x^{\text{-}n} = \frac{1}{x^{n}},$ we may now apply the quotient rule by setting $f(x) = 1$ and $g(x) = x^{n}.$ In this case, $f'(x) = 0$ and $g'(x) = nx^{n - 1}.$ Thus,

$$
\frac{d}{dx}\left( x^{\text{−}n} \right) = \frac{0\left( x^{n} \right) - 1\left( {nx^{n - 1}} \right)}{\left( x^{n} \right)^{2}}.
$$

Simplifying, we see that

$$
\frac{d}{dx}\left( x^{\text{−}n} \right) = \frac{\text{−}nx^{n - 1}}{x^{2n}} = \text{−}nx^{{({n - 1})} - 2n} = \text{−}nx^{\text{−}n - 1}.
$$

Finally, observe that since $k = \text{-}n,$ by substituting we have

$$
\frac{d}{dx}\left( x^{k} \right) = kx^{k - 1}.
$$

□

### Example 3.26

#### Using the Extended Power Rule

Find $\frac{d}{dx}\left( x^{-4} \right).$

#### Solution

By applying the extended power rule with $k = -4,$ we obtain

$$
\frac{d}{dx}\left( x^{-4} \right) = -4x^{-4 - 1} = -4x^{-5}.
$$

### Example 3.27

#### Using the Extended Power Rule and the Constant Multiple Rule

Use the extended power rule and the constant multiple rule to find the derivative of $f(x) = \frac{6}{x^{2}}.$

#### Solution

It may seem tempting to use the quotient rule to find this derivative, and it would certainly not be incorrect to do so. However, it is far easier to differentiate this function by first rewriting it as $f(x) = 6x^{-2}.$

$$
\begin{array}{clccl}
{f'(x)} & {= \frac{d}{dx}\left( \frac{6}{x^{2}} \right) = \frac{d}{dx}\left( {6x^{-2}} \right)} & & & {\text{Rewrite}\ \frac{6}{x^{2}}\ \text{as}\ 6x^{-2}.} \\
 & {= 6\frac{d}{dx}(x^{-2})} & & & \text{Apply the constant multiple rule.} \\
 & {= 6(-2x^{-3})} & & & {\text{Use the extended power rule to differentiate}\ x^{-2}.} \\
 & {= -12x^{-3}} & & & \text{Simplify.}
\end{array}
$$

### Checkpoint 3.18

Find the derivative of $g(x) = \frac{1}{x^{7}}$ using the extended power rule.

### Combining Differentiation Rules

As we have seen throughout the examples in this section, it seldom happens that we are called on to apply just one differentiation rule to find the derivative of a given function. At this point, by combining the differentiation rules, we may find the derivatives of any polynomial or rational function. Later on we will encounter more complex combinations of differentiation rules. A good rule of thumb to use when applying several rules is to apply the rules in reverse of the order in which we would evaluate the function.

### Example 3.28

#### Combining Differentiation Rules

For $k(x) = 3h(x) + x^{2}g(x),$ find $k'(x).$

#### Solution

Finding this derivative requires the sum rule, the constant multiple rule, and the product rule.

$$
\begin{array}{clccc}
{k'(x)} & {= \frac{d}{dx}\left( {3h(x) + x^{2}g(x)} \right) = \frac{d}{dx}\left( {3h(x)} \right) + \frac{d}{dx}\left( {x^{2}g(x)} \right)} & & & \text{Apply the sum rule.} \\
 & {= 3\frac{d}{dx}\left( {h(x)} \right) + \left( {\frac{d}{dx}\left( x^{2} \right)g(x) + \frac{d}{dx}\left( {g(x)} \right)x^{2}} \right)} & & & \begin{array}{l}
\text{Apply the constant multiple rule to} \\
{\text{differentiate}\ 3h(x)\ \text{and the product}} \\
{\text{rule to differentiate}\ x^{2}g(x).}
\end{array} \\
 & {= 3h^{'}(x) + 2xg(x) + x^{2}g^{'}(x)} & & & 
\end{array}
$$

### Example 3.29

#### Extending the Product Rule

For $k(x) = f(x)g(x)h(x),$ express $k'(x)$ in terms of $f(x),g(x),h(x),$ and their derivatives.

#### Solution

We can think of the function $k(x)$ as the product of the function $f(x)g(x)$ and the function $h(x).$ That is, $k(x) = \left( {f(x)g(x)} \right) \cdot h(x).$ Thus,

$$
\begin{array}{clccl}
{k'(x)} & {= \frac{d}{dx}\left( {f(x)g(x)} \right) \cdot h(x) + \frac{d}{dx}\left( {h(x)} \right) \cdot \left( {f(x)g(x)} \right)} & & & \begin{array}{l}
\text{Apply the product rule to the product} \\
{\text{of}\ f(x)g(x)\ \text{and}\ h(x).}
\end{array} \\
 & {= \left( {f'(x)g(x) + g'(x)f{(x))}h} \right.(x) + h'(x)f(x)g(x)} & & & {\text{Apply the product rule to}\ f(x)g(x).} \\
 & {= f'(x)g(x)h(x) + f(x)g'(x)h(x) + f(x)g(x)h'\left( x\text{).} \right.} & & & \text{Simplify.}
\end{array}
$$

### Example 3.30

#### Combining the Quotient Rule and the Product Rule

For $h(x) = \frac{2x^{3}k(x)}{3x + 2},$ find $h'(x).$

#### Solution

This procedure is typical for finding the derivative of a rational function.

$$
\begin{matrix}
{h^{'}(x)} & {= \frac{\frac{d}{dx}\left( 2x^{3}k(x) \right) \cdot (3x + 2) - \frac{d}{dx}(3x + 2) \cdot \left( 2x^{3}k(x) \right)}{(3x + 2)^{2}}} & & & \text{Apply the quotient rule.} & \\
 & {= \frac{\left( 6x^{2}k(x) + k^{'}(x) \cdot 2x^{3} \right)(3x + 2) - 3\left( 2x^{3}k(x) \right)}{(3x + 2)^{2}}} & & & \begin{array}{l}
\text{Apply the product rule to find} \\
{\frac{d}{dx}\left( 2x^{3}k(x) \right).\ \text{Use}\ \frac{d}{dx}(3x + 2) = 3.}
\end{array} & \\
 & {= \frac{-6x^{3}k(x) + 18x^{3}k(x) + 12x^{2}k(x) + 6x^{4}k^{'}(x) + 4x^{3}k^{'}(x)}{(3x + 2)^{2}}} & & & \text{Simplify.} & \\
 & {= \frac{12k(x)\left( {x^{3} + x^{2}} \right) + 2k'(x)\left( {3x^{4} + 2x^{3}} \right)}{\left( {3x + 2} \right)^{2}}} & & & & 
\end{matrix}
$$

### Checkpoint 3.19

Find $\frac{d}{dx}\left( {3f(x) - 2g(x)} \right).$

### Example 3.31

#### Determining Where a Function Has a Horizontal Tangent

Determine the values of $x$ for which $f(x) = x^{3} - 7x^{2} + 8x + 1$ has a horizontal tangent line.

#### Solution

To find the values of $x$ for which $f(x)$ has a horizontal tangent line, we must solve $f'(x) = 0.$ Since

$$
f'(x) = 3x^{2} - 14x + 8 = \left( {3x - 2} \right)\left( {x - 4} \right),
$$

we must solve $\left( {3x - 2} \right)\left( {x - 4} \right) = 0.$ Thus we see that the function has horizontal tangent lines at $x = \frac{2}{3}$ and $x = 4$ as shown in the following graph.

*Figure 3.19 This function has horizontal tangent lines at x = 2/3 and x = 4.*

### Example 3.32

#### Finding a Velocity

The position of an object on a coordinate axis at time $t$ is given by $s(t) = \frac{t}{t^{2} + 1}.$ What is the initial velocity of the object?

#### Solution

Since the initial velocity is $v(0) = s'(0),$ begin by finding $s'(t)$ by applying the quotient rule:

$$
s'(t) = \frac{1\left( {t^{2} + 1} \right) - 2t(t)}{\left( {t^{2} + 1} \right)^{2}} = \frac{1 - t^{2}}{\left( {t^{2} + 1} \right)^{2}}.
$$

After evaluating, we see that $v(0) = 1.$

### Checkpoint 3.20

Find the values of $x$ for which the graph of $f(x) = 4x^{2} - 3x + 2$ has a tangent line parallel to the line $y = 2x + 3.$

### Student Project

#### Formula One Grandstands

Formula One car races can be very exciting to watch and attract a lot of spectators. Formula One track designers have to ensure sufficient grandstand space is available around the track to accommodate these viewers. However, car racing can be dangerous, and safety considerations are paramount. The grandstands must be placed where spectators will not be in danger should a driver lose control of a car (Figure 3.20).

*Figure 3.20 The grandstand next to a straightaway of the Circuit de Barcelona-Catalunya race track, located where the spectators are not in danger.*

**********

Safety is especially a concern on turns. If a driver does not slow down enough before entering the turn, the car may slide off the racetrack. Normally, this just results in a wider turn, which slows the driver down. But if the driver loses control completely, the car may fly off the track entirely, on a path tangent to the curve of the racetrack.

Suppose you are designing a new Formula One track. One section of the track can be modeled by the function $f(x) = x^{3} + 3x^{2} + x$ (Figure 3.21). The current plan calls for grandstands to be built along the first straightaway and around a portion of the first curve. The plans call for the front corner of the grandstand to be located at the point $\left( {-1.9,2.8} \right).$ We want to determine whether this location puts the spectators in danger if a driver loses control of the car.

*Figure 3.21 (a) One section of the racetrack can be modeled by the function f ( x ) = x 3 + 3 x 2 + x . f ( x ) = x 3 + 3 x 2 + x . (b) The front corner of the grandstand is located at ( −1.9 , 2.8 ) . ( −1.9 , 2.8 ) .*

1.  Physicists have determined that drivers are most likely to lose control of their cars as they are coming into a turn, at the point where the slope of the tangent line is 1. Find the $\left( {x,y} \right)$ coordinates of this point near the turn.
2.  Find an equation of the tangent line to the curve at this point.
3.  To determine whether the spectators are in danger in this scenario, find the *x*-coordinate of the point where the tangent line crosses the line $y = 2.8.$ Is this point safely to the right of the grandstand? Or are the spectators in danger?
4.  What if a driver loses control earlier than the physicists project? Suppose a driver loses control at the point $\left( {-2.5,0.625} \right).$ What is the slope of the tangent line at this point?
5.  If a driver loses control as described in part 4, are the spectators safe?
6.  Should you proceed with the current design for the grandstand, or should the grandstands be moved?

### Section 3.3 Exercises

For the following exercises, find $f'(x)$ for each function.

106\.

$f(x) = x^{7} + 10$

107\.

$f(x) = 5x^{3} - x + 1$

108\.

$f(x) = 4x^{2} - 7x$

109\.

$f(x) = 8x^{4} + 9x^{2} - 1$

110\.

$f(x) = x^{4} + \frac{2}{x}$

111\.

$f(x) = 3x\left( {18x^{4} + \frac{13}{x + 1}} \right)$

112\.

$f(x) = \left( {x + 2} \right)\left( {2x^{2} - 3} \right)$

113\.

$f(x) = x^{2}\left( {\frac{2}{x^{2}} + \frac{5}{x^{3}}} \right)$

114\.

$f(x) = \frac{x^{3} + 2x^{2} - 4}{3}$

115\.

$f(x) = \frac{4x^{3} - 2x + 1}{x^{2}}$

116\.

$f(x) = \frac{x^{2} + 4}{x^{2} - 4}$

117\.

$f(x) = \frac{x + 9}{x^{2} - 7x + 1}$

For the following exercises, find an equation of the tangent line $T(x)$ to the graph of the given function at the indicated point. Use a graphing calculator to graph the function and the tangent line.

118\.

**\[T\]** $y = 3x^{2} + 4x + 1$ at $\left( {0,1} \right)$

119\.

**\[T\]** $y = \frac{2}{x^{2}} + 1$ at $\left( {1,3} \right)$

120\.

**\[T\]** $y = \frac{2x}{x - 1}$ at $\left( {-1,1} \right)$

121\.

**\[T\]** $y = \frac{2}{x} - \frac{3}{x^{2}}$ at $\left( {1,-1} \right)$

For the following exercises, assume that $f(x)$ and $g(x)$ are both differentiable functions for all $x.$ Find the derivative of each of the functions $h(x).$

122\.

$h(x) = 4f(x) + \frac{g(x)}{7}$

123\.

$h(x) = x^{3}f(x)$

124\.

$h(x) = \frac{f(x)g(x)}{2}$

125\.

$h(x) = \frac{3f(x)}{g(x) + 2}$

For the following exercises, assume that $f(x)$ and $g(x)$ are both differentiable functions with values as given in the following table. Use the following table to calculate the following derivatives.

|               |        |       |        |        |
|---------------|--------|-------|--------|--------|
| **$x$**     | $1$  | $2$ | $3$  | $4$  |
| **$f(x)$**  | $3$  | $5$ | $-2$ | $0$  |
| **$g(x)$**  | $2$  | $3$ | $-4$ | $6$  |
| **$f'(x)$** | $-1$ | $7$ | $8$  | $-3$ |
| **$g'(x)$** | $4$  | $1$ | $2$  | $9$  |

126\.

Find $h'(1)$ if $h(x) = xf(x) + 4g(x).$

127\.

Find $h'(2)$ if $h(x) = \frac{f(x)}{g(x)}.$

128\.

Find $h'(3)$ if $h(x) = 2x + f(x)g(x).$

129\.

Find $h'(4)$ if $h(x) = \frac{1}{x} + \frac{g(x)}{f(x)}.$

For the following exercises, use the following figure to find the indicated derivatives, if they exist.

130\.

Let $h(x) = f(x) + g(x).$ Find

1.  $h'(1),$
2.  $h'(3),$ and
3.  $h'(4).$

131\.

Let $h(x) = f(x)g(x).$ Find

1.  $h'(1),$
2.  $h'(3),$ and
3.  $h'(4).$

132\.

Let $h(x) = \frac{f(x)}{g(x)}.$ Find

1.  $h'(1),$
2.  $h'(3),$ and
3.  $h'(4).$

For the following exercises,

1.  evaluate $f'(a),$ and
2.  graph the function $f(x)$ and the tangent line at $x = a.$

133\.

**\[T\]** $f(x) = 2x^{3} + 3x - x^{2},a = 2$

134\.

**\[T\]** $f(x) = \frac{1}{x} - x^{2},a = 1$

135\.

**\[T\]** $f(x) = x^{2} - x^{12} + 3x + 2,a = 0$

136\.

**\[T\]** $f(x) = \frac{1}{x} - x^{2},a = -1$

137\.

Find an equation of the tangent line to the graph of $f(x) = 2x^{3} + 4x^{2} - 5x - 3$ at $x = -1.$

138\.

Find an equation of the tangent line to the graph of $f(x) = x^{2} + \frac{4}{x} - 10$ at $x = 8.$

139\.

Find an equation of the tangent line to the graph of $f(x) = (3x - x^{2})(3 - x - x^{2})$ at $x = 1.$

140\.

Find the point on the graph of $f(x) = x^{3}$ such that the tangent line at that point has an $x$ intercept of 6.

141\.

Find an equation of the line passing through the point $P(3,3)$ and tangent to the graph of $f(x) = \frac{6}{x - 1}.$

142\.

Determine all points on the graph of $f(x) = x^{3} + x^{2} - x - 1$ for which

1.  the tangent line is horizontal
2.  the tangent line has a slope of $-1.$

143\.

Find a quadratic polynomial such that $f(1) = 5,f'(1) = 3$ and $f^{''}(1) = -6.$

144\.

A car driving along a freeway with traffic has traveled $s(t) = t^{3} - 6t^{2} + 9t$ meters in $t$ seconds.

1.  Determine the time in seconds when the velocity of the car is 0.
2.  Determine the acceleration of the car when the velocity is 0.

145\.

**\[T\]** A herring swimming along a straight line has traveled $s(t) = \frac{t^{2}}{t^{2} + 2}$ feet in $t$ seconds.

Determine the velocity of the herring when it has traveled 3 seconds.

146\.

The population in millions of arctic flounder in the Atlantic Ocean is modeled by the function $P(t) = \frac{8t + 3}{0.2t^{2} + 1},$ where $t$ is measured in years.

1.  Determine the initial flounder population.
2.  Determine $P'(10)$ and briefly interpret the result.

147\.

**\[T\]** The concentration of antibiotic in the bloodstream $t$ hours after being injected is given by the function $C(t) = \frac{2t^{2} + t}{t^{3} + 50},$ where $C$ is measured in milligrams per liter of blood.

1.  Find the rate of change of $C(t).$
2.  Determine the rate of change for $t = 8,12,24,$ and $36.$
3.  Briefly describe what seems to be occurring as the number of hours increases.

148\.

A book publisher has a cost function given by $C(x) = \frac{x^{3} + 2x + 3}{x^{2}},$ where *x* is the number of copies of a book in thousands and *C* is the cost, per book, measured in dollars. Evaluate $C'(2)$ and explain its meaning.

149\.

**\[T\]** According to Newton’s law of universal gravitation, the force $F$ between two bodies of constant mass $m_{1}$ and $m_{2}$ is given by the formula $F = \frac{Gm_{1}m_{2}}{d^{2}},$ where $G$ is the gravitational constant and $d$ is the distance between the bodies.

1.  Suppose that $G,m_{1},\text{and}\ m_{2}$ are constants. Find the rate of change of force $F$ with respect to distance $d.$
2.  Find the rate of change of force $F$ with gravitational constant $G = 6.67\  \times \ 10^{-11}$ $\text{Nm}^{2}\text{/}\text{kg}^{2},$ on two bodies 10 meters apart, each with a mass of 1000 kilograms.
