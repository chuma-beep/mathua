> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 3.9 Derivatives of Exponential and Logarithmic Functions

### Learning Objectives

- 3.9.1 Find the derivative of exponential functions.
- 3.9.2 Find the derivative of logarithmic functions.
- 3.9.3 Use logarithmic differentiation to determine the derivative of a function.

So far, we have learned how to differentiate a variety of functions, including trigonometric, inverse, and implicit functions. In this section, we explore derivatives of exponential and logarithmic functions. As we discussed in Introduction to Functions and Graphs, exponential functions play an important role in modeling population growth and the decay of radioactive materials. Logarithmic functions can help rescale large quantities and are particularly helpful for rewriting complicated expressions.

### Derivative of the Exponential Function

Just as when we found the derivatives of other functions, we can find the derivatives of exponential and logarithmic functions using formulas. As we develop these formulas, we need to make certain basic assumptions. The proofs that these assumptions hold are beyond the scope of this course.

First of all, we begin with the assumption that the function $B(x) = b^{x},b > 0,$ is defined for every real number and is continuous. In previous courses, the values of exponential functions for all rational numbers were defined—beginning with the definition of $b^{n},$ where $n$ is a positive integer—as the product of $b$ multiplied by itself $n$ times. Later, we defined $b^{0} = 1,b^{\text{-}n} = \frac{1}{b^{n}},$ for a positive integer $n,$ and $b^{s\text{/}t} = (\sqrt[t]{b})^{s}$ for positive integers $s$ and $t.$ These definitions leave open the question of the value of $b^{r}$ where $r$ is an arbitrary real number. By assuming the *continuity* of $B(x) = b^{x},b > 0,$ we may interpret $b^{r}$ as $\underset{x\rightarrow r}{\text{lim}}b^{x}$ where the values of $x$ as we take the limit are rational. For example, we may view $4^{\pi}$ as the number satisfying

$$
\begin{array}{l}
{4^{3} < 4^{\pi} < 4^{4},4^{3.1} < 4^{\pi} < 4^{3.2},4^{3.14} < 4^{\pi} < 4^{3.15},} \\
{4^{3.141} < 4^{\pi} < 4^{3.142},4^{3.1415} < 4^{\pi} < 4^{3.1416}\text{,}\text{…}.}
\end{array}
$$

As we see in the following table, $4^{\pi} \approx 77.88.$

| $x$       | $4^{x}$     | $x$        | $4^{x}$     |
|-------------|---------------|--------------|---------------|
| $3$       | 64            | $3.141593$ | 77.8802710486 |
| $3.1$     | 73.5166947198 | $3.1416$   | 77.8810268071 |
| $3.14$    | 77.7084726013 | $3.142$    | 77.9242251944 |
| $3.141$   | 77.8162741237 | $3.15$     | 78.7932424541 |
| $3.1415$  | 77.8702309526 | $3.2$      | 84.4485062895 |
| $3.14159$ | 77.8799471543 | $4$        | 256           |

Table 3.6 Approximating a Value of $4^{\pi}$

We also assume that for $B(x) = b^{x},b > 0,$ the value $B'(0)$ of the derivative exists. In this section, we show that by making this one additional assumption, it is possible to prove that the function $B(x)$ is differentiable everywhere.

We make one final assumption: that there is a unique value of $b > 0$ for which $B'(0) = 1.$ We define $e$ to be this unique value, as we did in Introduction to Functions and Graphs. Figure 3.33 provides graphs of the functions $y = 2^{x},y = 3^{x},y = 2.7^{x},$ and $y = 2.8^{x}.$ A visual estimate of the slopes of the tangent lines to these functions at 0 provides evidence that the value of *e* lies somewhere between 2.7 and 2.8. The function $E(x) = e^{x}$ is called the natural exponential function. Its inverse, $L(x) = \text{log}_{e}x = \text{ln} x$ is called the natural logarithmic function.

*Figure 3.33 The graph of E ( x ) = e x E ( x ) = e x is between y = 2 x y = 2 x and y = 3 x . y = 3 x .*

For a better estimate of $e,$ we may construct a table of estimates of $B'(0)$ for functions of the form $B(x) = b^{x}.$ Before doing this, recall that

$$
B'(0) = \underset{x\rightarrow 0}{\text{lim}}\frac{b^{x} - b^{0}}{x - 0} = \underset{x\rightarrow 0}{\text{lim}}\frac{b^{x} - 1}{x} \approx \frac{b^{x} - 1}{x}
$$

for values of $x$ very close to zero. For our estimates, we choose $x = 0.00001$ and $x = -0.00001$ to obtain the estimate

$$
\frac{b^{-0.00001} - 1}{-0.00001} < B'(0) < \frac{b^{0.00001} - 1}{0.00001}.
$$

See the following table.

| $b$ | $\frac{b^{-0.00001} - 1}{-0.00001} < B'(0) < \frac{b^{0.00001} - 1}{0.00001}$ | $b$ | $\frac{b^{-0.00001} - 1}{-0.00001} < B'(0) < \frac{b^{0.00001} - 1}{0.00001}$ |
|----|----|----|----|
| $2$ | $0.693145 < B'(0) < 0.69315$ | $2.7183$ | $1.000002 < B'(0) < 1.000012$ |
| $2.7$ | $0.993247 < B'(0) < 0.993257$ | $2.719$ | $1.000259 < B'(0) < 1.000269$ |
| $2.71$ | $0.996944 < B'(0) < 0.996954$ | $2.72$ | $1.000627 < B'(0) < 1.000637$ |
| $2.718$ | $0.999891 < B'(0) < 0.999901$ | $2.8$ | $1.029614 < B'(0) < 1.029625$ |
| $2.7182$ | $0.999965 < B'(0) < 0.999975$ | $3$ | $1.098606 < B'(0) < 1.098618$ |

Table 3.7 Estimating a Value of $e$

The evidence from the table suggests that $2.7182 < e < 2.7183.$

The graph of $E(x) = e^{x}$ together with the line $y = x + 1$ are shown in Figure 3.34. This line is tangent to the graph of $E(x) = e^{x}$ at $x = 0.$

*Figure 3.34 The tangent line to E ( x ) = e x E ( x ) = e x at x = 0 x = 0 has slope 1.*

Now that we have laid out our basic assumptions, we begin our investigation by exploring the derivative of $B(x) = b^{x},b > 0.$ Recall that we have assumed that $B'(0)$ exists. By applying the limit definition to the derivative we conclude that

$$
B'(0) = \underset{h\rightarrow 0}{\text{lim}}\frac{b^{0 + h} - b^{0}}{h} = \underset{h\rightarrow 0}{\text{lim}}\frac{b^{h} - 1}{h}.
$$

(3.28)

Turning to $B'(x),$ we obtain the following.

$$
\begin{array}{clccl}
{B'(x)} & {= \underset{h\rightarrow 0}{\text{lim}}\frac{b^{x + h} - b^{x}}{h}} & & & \text{Apply the limit definition of the derivative.} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{b^{x}b^{h} - b^{x}}{h}} & & & {\text{Note that}\ b^{x + h} = b^{x}b^{h}.} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{b^{x}(b^{h} - 1)}{h}} & & & {\text{Factor out}\ b^{x}.} \\
 & {= b^{x}\underset{h\rightarrow 0}{\text{lim}}\frac{b^{h} - 1}{h}} & & & \text{Apply a property of limits.} \\
 & {= b^{x}B'(0)} & & & {\text{Use}\ B'(0) = \underset{h\rightarrow 0}{\text{lim}}\frac{b^{0 + h} - b^{0}}{h} = \underset{h\rightarrow 0}{\text{lim}}\frac{b^{h} - 1}{h}.}
\end{array}
$$

We see that on the basis of the assumption that $B(x) = b^{x}$ is differentiable at $0,B(x)$ is not only differentiable everywhere, but its derivative is

$$
B'(x) = b^{x}B'(0).
$$

(3.29)

For $E(x) = e^{x},E'(0) = 1.$ Thus, we have $E'(x) = e^{x}.$ (The value of $B'(0)$ for an arbitrary function of the form $B(x) = b^{x},b > 0,$ will be derived later.)

### Theorem 3.14

#### Derivative of the Natural Exponential Function

Let $E(x) = e^{x}$ be the natural exponential function. Then

$$
E'(x) = e^{x}.
$$

In general,

$$
\frac{d}{dx}\left( e^{g(x)} \right) = e^{g(x)}g'(x).
$$

### Example 3.74

#### Derivative of an Exponential Function

Find the derivative of $f(x) = e^{\text{tan}(2x)}.$

#### Solution

Using the derivative formula and the chain rule,

$$
\begin{array}{cl}
{f'(x)} & {= e^{\text{tan}{({2x})}}\frac{d}{dx}\left( {\text{tan}\left( {2x} \right)} \right)} \\
 & {= e^{\text{tan}(2x)}\text{sec}^{2}\left( {2x} \right) \cdot 2.}
\end{array}
$$

### Example 3.75

#### Combining Differentiation Rules

Find the derivative of $y = \frac{e^{x^{2}}}{x}.$

#### Solution

Use the derivative of the natural exponential function, the quotient rule, and the chain rule.

$$
\begin{array}{clccl}
y' & {= \frac{\left( {e^{x^{2}} \cdot 2} \right)x \cdot x - 1 \cdot e^{x^{2}}}{x^{2}}} & & & \text{Apply the quotient rule.} \\
 & {= \frac{e^{x^{2}}\left( {2x^{2} - 1} \right)}{x^{2}}} & & & \text{Simplify.}
\end{array}
$$

### Checkpoint 3.50

Find the derivative of $h(x) = xe^{2x}.$

### Example 3.76

#### Applying the Natural Exponential Function

A colony of mosquitoes has an initial population of 1000. After $t$ days, the population is given by $A(t) = 1000e^{0.3t}.$ Show that the ratio of the rate of change of the population, $A'(t),$ to the population, $A(t)$ is constant.

#### Solution

First find $A'(t).$ By using the chain rule, we have $A'(t) = 300e^{0.3t}.$ Thus, the ratio of the rate of change of the population to the population is given by

$$
\frac{A^{'}(t)}{A(t)} = \frac{300e^{0.3t}}{1000e^{0.3t}} = 0.3.
$$

The ratio of the rate of change of the population to the population is the constant 0.3.

### Checkpoint 3.51

If $A(t) = 1000e^{0.3t}$ describes the mosquito population after $t$ days, as in the preceding example, what is the rate of change of $A(t)$ after 4 days?

### Derivative of the Logarithmic Function

Now that we have the derivative of the natural exponential function, we can use implicit differentiation to find the derivative of its inverse, the natural logarithmic function.

### Theorem 3.15

#### The Derivative of the Natural Logarithmic Function

If $x > 0$ and $y = \text{ln} x,$ then

$$
\frac{dy}{dx} = \frac{1}{x}.
$$

(3.30)

More generally, let $g(x)$ be a differentiable function. For all values of $x$ for which $g(x) > 0,$ the derivative of $h(x) = \text{ln}\left( {g(x)} \right)$ is given by

$$
h'(x) = \frac{1}{g(x)}g'(x).
$$

(3.31)

#### Proof

If $x > 0$ and $y = \text{ln} x,$ then $e^{y} = x.$ Differentiating both sides of this equation results in the equation

$$
e^{y}\frac{dy}{dx} = 1.
$$

Solving for $\frac{dy}{dx}$ yields

$$
\frac{dy}{dx} = \frac{1}{e^{y}}.
$$

Finally, we substitute $x = e^{y}$ to obtain

$$
\frac{dy}{dx} = \frac{1}{x}.
$$

We may also derive this result by applying the inverse function theorem, as follows. Since $y = g(x) = \text{ln} x$ is the inverse of $f(x) = e^{x},$ by applying the inverse function theorem we have

$$
\frac{dy}{dx} = \frac{1}{f'\left( {g(x)} \right)} = \frac{1}{e^{\text{ln} x}} = \frac{1}{x}.
$$

Using this result and applying the chain rule to $h(x) = \text{ln}\left( {g(x)} \right)$ yields

$$
h'(x) = \frac{1}{g(x)}g'(x).
$$

□

The graph of $y = \text{ln} x$ and its derivative $\frac{dy}{dx} = \frac{1}{x}$ are shown in Figure 3.35.

*Figure 3.35 The function y = ln x The function y = ln x is increasing on ( 0 , + ∞ ) . ( 0 , + ∞ ) . Its derivative y ′ = 1 x y ′ = 1 x is greater than zero on ( 0 , + ∞ ) . ( 0 , + ∞ ) .*

### Example 3.77

#### Taking a Derivative of a Natural Logarithm

Find the derivative of $f(x) = \text{ln}\left( {x^{3} + 3x - 4} \right).$

#### Solution

Use Equation 3.31 directly.

$$
\begin{array}{clccl}
{f'(x)} & {= \frac{1}{x^{3} + 3x - 4} \cdot \left( {3x^{2} + 3} \right)} & & & {\text{Use}\ g(x) = x^{3} + 3x - 4\ \text{in}\ h'(x) = \frac{1}{g(x)}g'(x).} \\
 & {= \frac{3x^{2} + 3}{x^{3} + 3x - 4}} & & & \text{Rewrite.}
\end{array}
$$

### Example 3.78

#### Using Properties of Logarithms in a Derivative

Find the derivative of $f(x) = \text{ln}\left( \frac{x^{2}\text{sin} x}{2x + 1} \right).$

#### Solution

At first glance, taking this derivative appears rather complicated. However, by using the properties of logarithms prior to finding the derivative, we can make the problem much simpler.

$$
\begin{array}{rllccl}
{f(x)} & = & {\text{ln}\left( \frac{x^{2}\text{sin} x}{2x + 1} \right) = 2\text{ln} x + \text{ln}\left( {\text{sin} x} \right) - \text{ln}\left( {2x + 1} \right)} & & & \text{Apply properties of logarithms.} \\
{f'(x)} & = & {\frac{2}{x} + \text{cot} x - \frac{2}{2x + 1}} & & & {\text{Apply sum rule and}\ h'(x) = \frac{1}{g(x)}g'(x).}
\end{array}
$$

### Checkpoint 3.52

Differentiate: $f(x) = \text{ln}\left( {3x + 2} \right)^{5}.$

Now that we can differentiate the natural logarithmic function, we can use this result to find the derivatives of $y = log_{b}x$ and $y = b^{x}$ for $b > 0,b \neq 1.$

### Theorem 3.16

#### Derivatives of General Exponential and Logarithmic Functions

Let $b > 0,b \neq 1,$ and let $g(x)$ be a differentiable function.

1.  If, $y = \text{log}_{b}x,$ then  
    ``` math
    \frac{dy}{dx} = \frac{1}{x\text{ln} b}.
    ```
    (3.32)  
    More generally, if $h(x) = \text{log}_{b}\left( {g(x)} \right),$ then for all values of *x* for which $g(x) > 0,$  
    ``` math
    h'(x) = \frac{g'(x)}{g(x)\text{ln} b}.
    ```
    (3.33)
2.  If $y = b^{x},$ then  
    ``` math
    \frac{dy}{dx} = b^{x}\text{ln} b.
    ```
    (3.34)  
    More generally, if $h(x) = b^{g(x)},$ then  
    ``` math
    h'(x) = b^{g(x)}g\text{'}(x)\text{ln} b.
    ```
    (3.35)

#### Proof

If $y = \text{log}_{b}x,$ then $b^{y} = x.$ It follows that $\text{ln}\left( b^{y} \right) = \text{ln}\ x.$ Thus $y\ \text{ln}\ b = \text{ln}\ x.$ Solving for $y,$ we have $y = \frac{\text{ln} x}{\text{ln} b}.$ Differentiating and keeping in mind that $\text{ln} b$ is a constant, we see that

$$
\frac{dy}{dx} = \frac{1}{x\text{ln} b}.
$$

The derivative in Equation 3.32 now follows from the chain rule.

If $y = b^{x},$ then $\text{ln}\ y = x\text{ln} b.$ Using implicit differentiation, again keeping in mind that $\text{ln} b$ is constant, it follows that $\frac{1}{y}\ \frac{dy}{dx} = \text{ln} b.$ Solving for $\frac{dy}{dx}$ and substituting $y = b^{x},$ we see that

$$
\frac{dy}{dx} = y\text{ln} b = b^{x}\text{ln} b.
$$

The more general derivative (Equation 3.35) follows from the chain rule.

□

### Example 3.79

#### Applying Derivative Formulas

Find the derivative of $h(x) = \frac{3^{x}}{3^{x} + 2}.$

#### Solution

Use the quotient rule and Derivatives of General Exponential and Logarithmic Functions.

$$
\begin{array}{clccl}
{h'(x)} & {= \frac{3^{x}\text{ln} 3\left( {3^{x} + 2} \right) - 3^{x}\text{ln} 3\left( 3^{x} \right)}{\left( {3^{x} + 2} \right)^{2}}} & & & \text{Apply the quotient rule.} \\
 & {= \frac{2 \cdot 3^{x}\text{ln} 3}{\left( {3^{x} + 2} \right)^{2}}} & & & \text{Simplify.}
\end{array}
$$

### Example 3.80

#### Finding the Slope of a Tangent Line

Find the slope of the line tangent to the graph of $y = \text{log}_{2}\left( {3x + 1} \right)$ at $x = 1.$

#### Solution

To find the slope, we must evaluate $\frac{dy}{dx}$ at $x = 1.$ Using Equation 3.33, we see that

$$
\frac{dy}{dx} = \frac{3}{\left( {3x + 1} \right)\text{ln} 2}.
$$

By evaluating the derivative at $x = 1,$ we see that the tangent line has slope

$$
\frac{dy}{dx}{\left| \begin{array}{l}
 \\
{}_{x = 1}
\end{array} \right. = \frac{3}{4\text{ln} 2} = \frac{3}{\text{ln} 16}.}
$$

### Checkpoint 3.53

Find the slope for the line tangent to $y = 3^{x}$ at $x = 2.$

### Logarithmic Differentiation

At this point, we can take derivatives of functions of the form $y = \left( {g(x)} \right)^{n}$ for certain values of $n,$ as well as functions of the form $y = b^{g(x)},$ where $b > 0$ and $b \neq 1.$ Unfortunately, we still do not know the derivatives of functions such as $y = x^{x}$ or $y = x^{\pi}.$ These functions require a technique called logarithmic differentiation, which allows us to differentiate any function of the form $h(x) = g(x)^{f(x)}.$ It can also be used to convert a very complex differentiation problem into a simpler one, such as finding the derivative of $y = \frac{x\sqrt{2x + 1}}{e^{x}\text{sin}^{3}x}.$ We outline this technique in the following problem-solving strategy.

### Problem-Solving Strategy

#### Using Logarithmic Differentiation

1.  To differentiate $y = h(x)$ using logarithmic differentiation, take the natural logarithm of both sides of the equation to obtain $\text{ln}\ y = \text{ln}\left( {h(x)} \right).$
2.  Use properties of logarithms to expand $\text{ln}\left( {h(x)} \right)$ as much as possible.
3.  Differentiate both sides of the equation. On the left we will have $\frac{1}{y}\ \frac{dy}{dx}.$
4.  Multiply both sides of the equation by $y$ to solve for $\frac{dy}{dx}.$
5.  Replace $y$ by $h(x).$

### Example 3.81

#### Using Logarithmic Differentiation

Find the derivative of $y = \left( {2x^{4} + 1} \right)^{\text{tan} x}.$

#### Solution

Use logarithmic differentiation to find this derivative.

$$
\begin{array}{rllccl}
{\text{ln} y} & = & {\text{ln}\left( {2x^{4} + 1} \right)^{\text{tan} x}} & & & \text{Step 1. Take the natural logarithm of both sides.} \\
{\text{ln} y} & = & {\text{tan} x\text{ln}\left( {2x^{4} + 1} \right)} & & & \text{Step 2. Expand using properties of logarithms.} \\
{\frac{1}{y}\ \frac{dy}{dx}} & = & {\text{sec}^{2}x\text{ln}\left( {2x^{4} + 1} \right) + \frac{8x^{3}}{2x^{4} + 1} \cdot \text{tan} x} & & & \begin{array}{l}
\text{Step 3. Differentiate both sides. Use the} \\
\text{product rule on the right.}
\end{array} \\
\frac{dy}{dx} & = & {y \cdot \left( {\text{sec}^{2}x\text{ln}\left( {2x^{4} + 1} \right) + \frac{8x^{3}}{2x^{4} + 1} \cdot \text{tan} x} \right)} & & & {\text{Step 4. Multiply by}\ y\ \text{on both sides.}} \\
\frac{dy}{dx} & = & {\left( {2x^{4} + 1} \right)^{\text{tan} x}\left( {\text{sec}^{2}x\text{ln}\left( {2x^{4} + 1} \right) + \frac{8x^{3}}{2x^{4} + 1} \cdot \text{tan} x} \right)} & & & {\text{Step 5. Substitute}\ y = \left( {2x^{4} + 1} \right)^{\text{tan} x}.}
\end{array}
$$

### Example 3.82

#### Using Logarithmic Differentiation

Find the derivative of $y = \frac{x\sqrt{2x + 1}}{e^{x}\text{sin}^{3}x}.$

#### Solution

This problem really makes use of the properties of logarithms and the differentiation rules given in this chapter.

$$
\begin{array}{rllccl}
{\text{ln} y} & = & {\text{ln}\ \frac{x\sqrt{2x + 1}}{e^{x}\text{sin}^{3}x}} & & & \text{Step 1. Take the natural logarithm of both sides.} \\
{\text{ln} y} & = & {\text{ln} x + \frac{1}{2}\text{ln}\left( {2x + 1} \right) - x\text{ln} e - 3\text{ln}\text{sin} x} & & & \text{Step 2. Expand using properties of logarithms.} \\
{\frac{1}{y}\ \frac{dy}{dx}} & = & {\frac{1}{x} + \frac{1}{2x + 1} - 1 - 3\frac{\text{cos} x}{\text{sin} x}} & & & \text{Step 3. Differentiate both sides.} \\
\frac{dy}{dx} & = & {y\left( {\frac{1}{x} + \frac{1}{2x + 1} - 1 - 3\text{cot} x} \right)} & & & {\text{Step 4. Multiply by}\ y\ \text{on both sides.}} \\
\frac{dy}{dx} & = & {\frac{x\sqrt{2x + 1}}{e^{x}\text{sin}^{3}x}\left( {\frac{1}{x} + \frac{1}{2x + 1} - 1 - 3\text{cot} x} \right)} & & & {\text{Step 5. Substitute}\ y = \frac{x\sqrt{2x + 1}}{e^{x}\text{sin}^{3}x}.}
\end{array}
$$

### Example 3.83

#### Extending the Power Rule

Find the derivative of $y = x^{r}$ where $r$ is an arbitrary real number.

#### Solution

The process is the same as in Example 3.82, though with fewer complications.

$$
\begin{array}{rllccl}
{\text{ln} y} & = & {\text{ln}x^{r}} & & & \text{Step 1. Take the natural logarithm of both sides.} \\
{\text{ln} y} & = & {r\text{ln} x} & & & \text{Step 2. Expand using properties of logarithms.} \\
{\frac{1}{y}\ \frac{dy}{dx}} & = & {r\frac{1}{x}} & & & \text{Step 3. Differentiate both sides.} \\
\frac{dy}{dx} & = & {y\frac{r}{x}} & & & {\text{Step 4. Multiply by}\ y\ \text{on both sides.}} \\
\frac{dy}{dx} & = & {x^{r}\frac{r}{x}} & & & {\text{Step 5. Substitute}\ y = x^{r}.} \\
\frac{dy}{dx} & = & {rx^{r - 1}} & & & \text{Simplify.}
\end{array}
$$

### Checkpoint 3.54

Use logarithmic differentiation to find the derivative of $y = x^{x}.$

### Checkpoint 3.55

Find the derivative of $y = \left( {\text{tan} x} \right)^{\pi}.$

### Section 3.9 Exercises

For the following exercises, find $f'(x)$ for each function.

331\.

$f(x) = x^{2}e^{x}$

332\.

$f(x) = \frac{e^{\text{-}x}}{x}$

333\.

$f(x) = e^{x^{3}\text{ln} x}$

334\.

$f(x) = \sqrt{e^{2x} + 2x}$

335\.

$f(x) = \frac{e^{x} - e^{\text{-}x}}{e^{x} + e^{\text{-}x}}$

336\.

$f(x) = \frac{10^{x}}{\text{ln} 10}$

337\.

$f(x) = 2^{4x} + 4x^{2}$

338\.

$f(x) = 3^{\text{sin} 3\text{x}}$

339\.

$f(x) = x^{\pi} \cdot \pi^{x}$

340\.

$f(x) = \text{ln}\left( {4x^{3} + x} \right)$

341\.

$f(x) = \text{ln}\sqrt{5x - 7}$

342\.

$f(x) = x^{2}\text{ln} 9x$

343\.

$f(x) = \text{log}\left( {\text{sec} x} \right)$

344\.

$f(x) = \text{log}_{7}\left( {6x^{4} + 3} \right)^{5}$

345\.

$f(x) = 2^{x} \cdot \text{log}_{3}7^{x^{2} - 4}$

For the following exercises, use logarithmic differentiation to find $\frac{dy}{dx}.$

346\.

$y = x^{\sqrt{x}}$

347\.

$y = \left( {\text{sin} 2x} \right)^{4x}$

348\.

$y = \left( {\text{ln} x} \right)^{\text{ln} x}$

349\.

$y = x^{\text{log}_{2}x}$

350\.

$y = \left( {x^{2} - 1} \right)^{\text{ln} x}$

351\.

$y = x^{\text{cot} x}$

352\.

$y = \frac{x + 11}{\sqrt[3]{x^{2} - 4}}$

353\.

$y = x^{-1\text{/}2}\left( {x^{2} + 3} \right)^{2\text{/}3}\left( {3x - 4} \right)^{4}$

354\.

**\[T\]** Find an equation of the tangent line to the graph of $f(x) = 4xe^{({x^{2} - 1})}$ at the point where

$x = -1.$ Graph both the function and the tangent line.

355\.

**\[T\]** Find an equation of the line that is normal to the graph of $f(x) = x \cdot 5^{x}$ at the point where $x = 1.$ Graph both the function and the normal line.

356\.

**\[T\]** Find an equation of the tangent line to the graph of $x^{3} - x\text{ln} y + y^{3} = 2x + 5$ at the point (2, 1). (*Hint*: Use implicit differentiation to find $\frac{dy}{dx}.)$ Graph both the curve and the tangent line.

357\.

Consider the function $y = x^{1\text{/}x}$ for $x > 0.$

1.  Determine the points on the graph where the tangent line is horizontal.
2.  Determine the points on the graph where $y' > 0$ and those where $y' < 0.$

358\.

The formula $I(t) = \frac{\text{sin} t}{e^{t}}$ is the formula for a decaying alternating current.

1.  Complete the following table with the appropriate values.  
    | $t$              | $\frac{\text{sin} t}{e^{t}}$ |
    |--------------------|--------------------------------------------|
    | 0                  | \(i\)                                      |
    | $\frac{\pi}{2}$  | \(ii\)                                     |
    | $\pi$            | \(iii\)                                    |
    | $\frac{3\pi}{2}$ | \(iv\)                                     |
    | $2\pi$           | \(v\)                                      |
    | $\frac{5\pi}{2}$ | \(vi\)                                     |
    | $3\pi$           | \(vii\)                                    |
    | $\frac{7\pi}{2}$ | \(viii\)                                   |
    | $4\pi$           | \(ix\)                                     |
2.  Using only the values in the table, determine where the tangent line to the graph of $I(t)$ is horizontal.

359\.

**\[T\]** The population of Toledo, Ohio, in 2000 was approximately 500,000. Assume the population is increasing at a rate of 5% per year.

1.  Write the exponential function that relates the total population as a function of $t.$
2.  Use a. to determine the rate at which the population is increasing in $t$ years.
3.  Use b. to determine the rate at which the population is increasing in 10 years.

360\.

**\[T\]** An isotope of the element erbium has a half-life of approximately 12 hours. Initially there are 9 grams of the isotope present.

1.  Write the exponential function that relates the amount of substance remaining as a function of $t,$ measured in hours.
2.  Use a. to determine the rate at which the substance is decaying in $t$ hours.
3.  Use b. to determine the rate of decay at $t = 4$ hours.

361\.

**\[T\]** The number of cases of influenza in New York City from the beginning of 1960 to the beginning of 1964 is modeled by the function

$N(t) = 5.3e^{0.093t^{2} - 0.87t},(0 \leq t \leq 4),$ where $N(t)$ gives the number of cases (in thousands) and *t* is measured in years, with $t = 0$ corresponding to the beginning of 1960.

1.  Show work that evaluates $N(0)$ and $N(4).$ Briefly describe what these values indicate about the disease in New York City.
2.  Show work that evaluates $N'(0)$ and $N'(3).$ Briefly describe what these values indicate about the disease in New York City.

362\.

**\[T\]** The *relative rate of change* of a differentiable function $y = f(x)$ is given by $\frac{100 \cdot f'(x)}{f(x)}\text{\%}.$ One model for population growth is a Gompertz growth function, given by $P(x) = ae^{\text{-}b \cdot e^{\text{-}cx}}$ where $a,b,$ and $c$ are constants.

1.  Find the relative rate of change formula for the generic Gompertz function.
2.  Use a. to find the relative rate of change of a population in $x = 20$ months when $a = 204,b = 0.0198,$ and $c = 0.15.$
3.  Briefly interpret what the result of b. means.

For the following exercises, use the population of New York City from 1790 to 1860, given in the following table.

| Years since 1790 | Population |
|------------------|------------|
| 0                | 33,131     |
| 10               | 60,515     |
| 20               | 96,373     |
| 30               | 123,706    |
| 40               | 202,300    |
| 50               | 312,710    |
| 60               | 515,547    |
| 70               | 813,669    |

Table 3.8 New York City Population Over Time Source: http://en.wikipedia.org/wiki/Largest_cities_in_the_United_States  
\_by_population_by_decade. 363.

**\[T\]** Using a computer program or a calculator, fit a growth curve to the data of the form $p = ab^{t}.$

364\.

**\[T\]** Using the exponential best fit for the data, write a table containing the derivatives evaluated at each year.

365\.

**\[T\]** Using the exponential best fit for the data, write a table containing the second derivatives evaluated at each year.

366\.

**\[T\]** Using the tables of first and second derivatives and the best fit, answer the following questions:

1.  Will the model be accurate in predicting the future population of New York City? Why or why not?
2.  Estimate the population in 2010. Was the prediction correct from a.?
