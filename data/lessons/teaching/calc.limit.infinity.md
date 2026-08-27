> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 4.6 Limits at Infinity and Asymptotes

### Learning Objectives

- 4.6.1 Calculate the limit of a function as $x$ increases or decreases without bound.
- 4.6.2 Recognize a horizontal asymptote on the graph of a function.
- 4.6.3 Estimate the end behavior of a function as $x$ increases or decreases without bound.
- 4.6.4 Recognize an oblique asymptote on the graph of a function.
- 4.6.5 Analyze a function and its derivatives to draw its graph.

We have shown how to use the first and second derivatives of a function to describe the shape of a graph. To graph a function $f$ defined on an unbounded domain, we also need to know the behavior of $f$ as $x\rightarrow\text{\pm}\infty.$ In this section, we define limits at infinity and show how these limits affect the graph of a function. At the end of this section, we outline a strategy for graphing an arbitrary function $f.$

### Limits at Infinity

We begin by examining what it means for a function to have a finite limit at infinity. Then we study the idea of a function with an infinite limit at infinity. Back in Introduction to Functions and Graphs, we looked at vertical asymptotes; in this section we deal with horizontal and oblique asymptotes.

#### Limits at Infinity and Horizontal Asymptotes

Recall that $\underset{x\rightarrow a}{\text{lim}}f(x) = L$ means $f(x)$ becomes arbitrarily close to $L$ as long as $x$ is sufficiently close to $a.$ We can extend this idea to limits at infinity. For example, consider the function $f(x) = 2 + \frac{1}{x}.$ As can be seen graphically in Figure 4.40 and numerically in Table 4.2, as the values of $x$ get larger, the values of $f(x)$ approach $2.$ We say the limit as $x$ approaches $\infty$ of $f(x)$ is $2$ and write $\underset{x\rightarrow\infty}{\text{lim}}f(x) = 2.$ Similarly, for $x < 0,$ as the values $|x|$ get larger, the values of $f(x)$ approaches $2.$ We say the limit as $x$ approaches $\text{-}\infty$ of $f(x)$ is $2$ and write $\underset{x\rightarrow - \infty}{\text{lim}}f(x) = 2.$

*Figure 4.40 The function approaches the asymptote y = 2 y = 2 as x x approaches ± ∞ . ± ∞ .*

|                         |         |          |           |             |
|-------------------------|---------|----------|-----------|-------------|
| **$x$**               | $10$  | $100$  | $1,000$ | $10,000$  |
| **$2 + \frac{1}{x}$** | $2.1$ | $2.01$ | $2.001$ | $2.0001$  |
| **$x$**               | $-10$ | $-100$ | $-1000$ | $-10,000$ |
| **$2 + \frac{1}{x}$** | $1.9$ | $1.99$ | $1.999$ | $1.9999$  |

Table 4.2 Values of a function $f$ as $x\rightarrow\text{\pm}\infty$

More generally, for any function $f,$ we say the limit as $x\rightarrow\infty$ of $f(x)$ is $L$ if $f(x)$ becomes arbitrarily close to $L$ as long as $x$ is sufficiently large. In that case, we write $\underset{x\rightarrow\infty}{\text{lim}}f(x) = L.$ Similarly, we say the limit as $x\rightarrow\text{-}\infty$ of $f(x)$ is $L$ if $f(x)$ becomes arbitrarily close to $L$ as long as $x < 0$ and $|x|$ is sufficiently large. In that case, we write $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x) = L.$ We now look at the definition of a function having a limit at infinity.

### Definition

(Informal) If the values of $f(x)$ become arbitrarily close to $L$ as $x$ becomes sufficiently large, we say the function $f$ has a limit at infinity and write

$$
\underset{x\rightarrow\infty}{\text{lim}}f(x) = L.
$$

If the values of $f(x)$ becomes arbitrarily close to $L$ for $x < 0$ as $|x|$ becomes sufficiently large, we say that the function $f$ has a limit at negative infinity and write

$$
\underset{x\rightarrow{–\infty}}{\text{lim}}f(x) = L.
$$

If the values $f(x)$ are getting arbitrarily close to some finite value $L$ as $x\rightarrow\infty$ or $x\rightarrow\text{-}\infty,$ the graph of $f$ approaches the line $y = L.$ In that case, the line $y = L$ is a horizontal asymptote of $f$ (Figure 4.41). For example, for the function $f(x) = \frac{1}{x},$ since $\underset{x\rightarrow\infty}{\text{lim}}f(x) = 0,$ the line $y = 0$ is a horizontal asymptote of $f(x) = \frac{1}{x}.$

### Definition

If $\underset{x\rightarrow\infty}{\text{lim}}f(x) = L$ or $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x) = L,$ we say the line $y = L$ is a horizontal asymptote of $f.$

*Figure 4.41 (a) As x → ∞ , x → ∞ , the values of f f are getting arbitrarily close to L . L . The line y = L y = L is a horizontal asymptote of f . f . (b) As x → − ∞ , x → − ∞ , the values of f f are getting arbitrarily close to M . M . The line y = M y = M is a horizontal asymptote of f . f .*

A function cannot cross a vertical asymptote because the graph must approach infinity (or $\text{-}\infty)$ from at least one direction as $x$ approaches the vertical asymptote. However, a function may cross a horizontal asymptote. In fact, a function may cross a horizontal asymptote an unlimited number of times. For example, the function ${f(x) = \frac{\left( {\text{cos} x} \right)}{x}} + 1$ shown in Figure 4.42 intersects the horizontal asymptote $y = 1$ an infinite number of times as it oscillates around the asymptote with ever-decreasing amplitude.

*Figure 4.42 The graph of f ( x ) = ( cos x ) / x + 1 f ( x ) = ( cos x ) / x + 1 crosses its horizontal asymptote y = 1 y = 1 an infinite number of times.*

The algebraic limit laws and squeeze theorem we introduced in Introduction to Limits also apply to limits at infinity. We illustrate how to use these laws to compute several limits at infinity.

### Example 4.21

#### Computing Limits at Infinity

For each of the following functions $f,$ evaluate $\underset{x\rightarrow\infty}{\text{lim}}f(x)$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x).$ Determine the horizontal asymptote(s) for $f.$

1.  $f(x) = 5 - \frac{2}{x^{2}}$
2.  $f(x) = \frac{\text{sin} x}{x}$
3.  $f(x) = \text{tan}^{-1}(x)$

#### Solution

1.  Using the algebraic limit laws, we have $\underset{x\rightarrow\infty}{\text{lim}}\left( {5 - \frac{2}{x^{2}}} \right) = \underset{x\rightarrow\infty}{\text{lim}}5 - 2\left( {\underset{x\rightarrow\infty}{\text{lim}}\frac{1}{x}} \right).\left( {\underset{x\rightarrow\infty}{\text{lim}}\frac{1}{x}} \right) = 5 - 2 \cdot 0 = 5.$  
    Similarly, $\underset{x\rightarrow - \infty}{\text{lim}}f(x) = 5.$ Therefore, $f(x) = 5 - \frac{2}{x^{2}}$ has a horizontal asymptote of $y = 5$ and $f$ approaches this horizontal asymptote as $x\rightarrow\text{\pm}\infty$ as shown in the following graph.  

    *Figure 4.43 This function approaches a horizontal asymptote as x → ± ∞ . x → ± ∞ .*

2.  Since $-1 \leq \text{sin} x \leq 1$ for all $x,$ we have  
    ``` math
    \frac{-1}{x} \leq \frac{\text{sin} x}{x} \leq \frac{1}{x}
    ```
      
    for all $x \neq 0.$ Also, since  
    ``` math
    \underset{x\rightarrow\infty}{\text{lim}}\frac{-1}{x} = 0 = \underset{x\rightarrow\infty}{\text{lim}}\frac{1}{x},
    ```
      
    we can apply the squeeze theorem to conclude that  
    ``` math
    \underset{x\rightarrow\infty}{\text{lim}}\frac{\text{sin} x}{x} = 0.
    ```
      
    Similarly,  
    ``` math
    \underset{x\rightarrow\text{−}\infty}{\text{lim}}\frac{\text{sin} x}{x} = 0.
    ```
      
    Thus, $f(x) = \frac{\text{sin} x}{x}$ has a horizontal asymptote of $y = 0$ and $f(x)$ approaches this horizontal asymptote as $x\rightarrow\text{\pm}\infty$ as shown in the following graph.  

    *Figure 4.44 This function crosses its horizontal asymptote multiple times.*

3.  To evaluate $\underset{x\rightarrow\infty}{\text{lim}}\text{tan}^{-1}(x)$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}\text{tan}^{-1}(x),$ we first consider the graph of $y = \text{tan}(x)$ over the interval $\left( {\text{-}\pi\text{/}2,\pi\text{/}2} \right)$ as shown in the following graph.  

    *Figure 4.45 The graph of tan x tan x has vertical asymptotes at x = ± π 2 x = ± π 2*

Since

$$
\underset{x\rightarrow{(\pi\text{/}2)}^{-}}{\text{lim}}\text{tan} x = \infty,
$$

it follows that

$$
\underset{x\rightarrow\infty}{\text{lim}}\text{tan}^{-1}(x) = \frac{\pi}{2}.
$$

Similarly, since

$$
\underset{x\rightarrow{({–\pi}\text{/}2)}^{+}}{\text{lim}}\text{tan} x = \text{−}\infty,
$$

it follows that

$$
\underset{x\rightarrow\text{−}\infty}{\text{lim}}\text{tan}^{-1}(x) = - \frac{\pi}{2}.
$$

As a result, $y = \frac{\pi}{2}$ and $y = - \frac{\pi}{2}$ are horizontal asymptotes of $f(x) = \text{tan}^{-1}(x)$ as shown in the following graph.

*Figure 4.46 This function has two horizontal asymptotes.*

### Checkpoint 4.20

Evaluate $\underset{x\rightarrow\text{-}\infty}{\text{lim}}\left( {3 + \frac{4}{x}} \right)$ and $\underset{x\rightarrow\infty}{\text{lim}}\left( {3 + \frac{4}{x}} \right).$ Determine the horizontal asymptotes of $f(x) = 3 + \frac{4}{x},$ if any.

#### Infinite Limits at Infinity

Sometimes the values of a function $f$ become arbitrarily large as $x\rightarrow\infty$ (or as $x\rightarrow\text{-}\infty).$ In this case, we write $\underset{x\rightarrow\infty}{\text{lim}}f(x) = \infty$ (or $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x) = \infty).$ On the other hand, if the values of $f$ are negative but become arbitrarily large in magnitude as $x\rightarrow\infty$ (or as $x\rightarrow\text{-}\infty),$ we write $\underset{x\rightarrow\infty}{\text{lim}}f(x) = \text{-}\infty$ (or $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x) = \text{-}\infty).$

For example, consider the function $f(x) = x^{3}.$ As seen in Table 4.3 and Figure 4.47, as $x\rightarrow\infty$ the values $f(x)$ become arbitrarily large. Therefore, $\underset{x\rightarrow\infty}{\text{lim}}x^{3} = \infty.$ On the other hand, as $x\rightarrow\text{-}\infty,$ the values of $f(x) = x^{3}$ are negative but become arbitrarily large in magnitude. Consequently, $\underset{x\rightarrow\text{-}\infty}{\text{lim}}x^{3} = \text{-}\infty.$

|  |  |  |  |  |  |
|----|----|----|----|----|----|
| **$x$** | $10$ | $20$ | $50$ | $100$ | $1000$ |
| **$x^{3}$** | $1000$ | $8000$ | $125,000$ | $1,000,000$ | $1,000,000,000$ |
| **$x$** | $-10$ | $-20$ | $-50$ | $-100$ | $-1000$ |
| **$x^{3}$** | $-1000$ | $-8000$ | $-125,000$ | $-1,000,000$ | $-1,000,000,000$ |

Table 4.3 Values of a power function as $x\rightarrow\text{\pm}\infty$

*Figure 4.47 For this function, the functional values approach infinity as x → ± ∞ . x → ± ∞ .*

### Definition

(Informal) We say a function $f$ has an infinite limit at infinity and write

$$
\underset{x\rightarrow\infty}{\text{lim}}f(x) = \infty.
$$

if $f(x)$ becomes arbitrarily large for $x$ sufficiently large. We say a function has a negative infinite limit at infinity and write

$$
\underset{x\rightarrow\infty}{\text{lim}}f(x) = \text{−}\infty.
$$

if $f(x) < 0$ and $\left| {f(x)} \right|$ becomes arbitrarily large for $x$ sufficiently large. Similarly, we can define infinite limits as $x\rightarrow\text{-}\infty.$

#### Formal Definitions

Earlier, we used the terms *arbitrarily close*, *arbitrarily large*, and *sufficiently large* to define limits at infinity informally. Although these terms provide accurate descriptions of limits at infinity, they are not precise mathematically. Here are more formal definitions of limits at infinity. We then look at how to use these definitions to prove results involving limits at infinity.

### Definition

(Formal) We say a function $f$ has a limit at infinity, if there exists a real number $L$ such that for all $\varepsilon > 0,$ there exists $N > 0$ such that

$$
\left| {f(x) - L} \right| < \varepsilon
$$

for all $x > N.$ In that case, we write

$$
\underset{x\rightarrow\infty}{\text{lim}}f(x) = L
$$

(see Figure 4.48).

We say a function $f$ has a limit at negative infinity if there exists a real number $L$ such that for all $\varepsilon > 0,$ there exists $N < 0$ such that

$$
\left| {f(x) - L} \right| < \varepsilon
$$

for all $x < N.$ In that case, we write

$$
\underset{x\rightarrow\text{−}\infty}{\text{lim}}f(x) = L.
$$

*Figure 4.48 For a function with a limit at infinity, for all x \> N , x \> N , \| f ( x ) − L \| \< ε . \| f ( x ) − L \| \< ε .*

Earlier in this section, we used graphical evidence in Figure 4.40 and numerical evidence in Table 4.2 to conclude that $\underset{x\rightarrow\infty}{\text{lim}}\left( 2 + \frac{1}{x} \right) = 2.$ Here we use the formal definition of limit at infinity to prove this result rigorously.

### Example 4.22

#### A Finite Limit at Infinity Example

Use the formal definition of limit at infinity to prove that $\underset{x\rightarrow\infty}{\text{lim}}\left( 2 + \frac{1}{x} \right) = 2.$

#### Solution

Let $\varepsilon > 0.$ Let $N = \frac{1}{\varepsilon}.$ Therefore, for all $x > N,$ we have

$$
{\left| {2 + \frac{1}{x} - 2} \right| = \left| \frac{1}{x} \right| = \frac{1}{x} < \frac{1}{N} = \varepsilon}\text{.}
$$

### Checkpoint 4.21

Use the formal definition of limit at infinity to prove that $\underset{x\rightarrow\infty}{\text{lim}}\left( 3 - \frac{1}{x^{2}} \right) = 3.$

We now turn our attention to a more precise definition for an infinite limit at infinity.

### Definition

(Formal) We say a function $f$ has an infinite limit at infinity and write

$$
\underset{x\rightarrow\infty}{\text{lim}}f(x) = \infty
$$

if for all $M > 0,$ there exists an $N > 0$ such that

$$
f(x) > M
$$

for all $x > N$ (see Figure 4.49).

We say a function has a negative infinite limit at infinity and write

$$
\underset{x\rightarrow\infty}{\text{lim}}f(x) = \text{−}\infty
$$

if for all $M < 0,$ there exists an $N > 0$ such that

$$
f(x) < M
$$

for all $x > N.$

Similarly we can define limits as $x\rightarrow\text{-}\infty.$

*Figure 4.49 For a function with an infinite limit at infinity, for all x \> N , x \> N , f ( x ) \> M . f ( x ) \> M .*

Earlier, we used graphical evidence (Figure 4.47) and numerical evidence (Table 4.3) to conclude that $\underset{x\rightarrow\infty}{\text{lim}}x^{3} = \infty.$ Here we use the formal definition of infinite limit at infinity to prove that result.

### Example 4.23

#### An Infinite Limit at Infinity

Use the formal definition of infinite limit at infinity to prove that $\underset{x\rightarrow\infty}{\text{lim}}x^{3} = \infty.$

#### Solution

Let $M > 0.$ Let $N = \sqrt[3]{M}.$ Then, for all $x > N,$ we have

$$
x^{3} > N^{3} = \left( \sqrt[3]{M} \right)^{3} = M.
$$

Therefore, $\underset{x\rightarrow\infty}{\text{lim}}x^{3} = \infty.$

### Checkpoint 4.22

Use the formal definition of infinite limit at infinity to prove that $\underset{x\rightarrow\infty}{\text{lim}}3x^{2} = \infty.$

### End Behavior

The behavior of a function as $x\rightarrow\text{\pm}\infty$ is called the function’s end behavior. At each of the function’s ends, the function could exhibit one of the following types of behavior:

1.  The function $f(x)$ approaches a horizontal asymptote $y = L.$
2.  The function $f(x)\rightarrow\infty$ or $f(x)\rightarrow\text{-}\infty.$
3.  The function does not approach a finite limit, nor does it approach $\infty$ or $\text{-}\infty.$ In this case, the function may have some oscillatory behavior.

Let’s consider several classes of functions here and look at the different types of end behaviors for these functions.

#### End Behavior for Polynomial Functions

Consider the power function $f(x) = x^{n}$ where $n$ is a positive integer. From Figure 4.50 and Figure 4.51, we see that

$$
\underset{x\rightarrow\infty}{\text{lim}}x^{n} = \infty;n = 1,2,3\text{,…}
$$

and

$$
\underset{x\rightarrow\text{−}\infty}{\text{lim}}x^{n} = \left\{ {\begin{array}{l}
{\infty;n = 2,4,6\text{,…}} \\
{\text{−}\infty;n = 1,3,5\text{,…}}
\end{array}.} \right.
$$

*Figure 4.50 For power functions with an even exponent of n , n , lim x → ∞ x n = ∞ = lim x → − ∞ x n . lim x → ∞ x n = ∞ = lim x → − ∞ x n .*

*Figure 4.51 For power functions with an odd exponent of n , n , lim x → ∞ x n = ∞ lim x → ∞ x n = ∞ and lim x → − ∞ x n = − ∞ . lim x → − ∞ x n = − ∞ .*

Using these facts, it is not difficult to evaluate $\underset{x\rightarrow\infty}{\text{lim}}cx^{n}$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}cx^{n},$ where $c$ is any constant and $n$ is a positive integer. If $c > 0,$ the graph of $y = cx^{n}$ is a vertical stretch or compression of $y = x^{n},$ and therefore

$$
\underset{x\rightarrow\infty}{\text{lim}}cx^{n} = \underset{x\rightarrow\infty}{\text{lim}}x^{n}\ \text{and}\ \underset{x\rightarrow\text{−}\infty}{\text{lim}}cx^{n} = \underset{x\rightarrow\text{−}\infty}{\text{lim}}x^{n}\ \text{if}\ c > 0.
$$

If $c < 0,$ the graph of $y = cx^{n}$ is a vertical stretch or compression combined with a reflection about the $x$-axis, and therefore

$$
\underset{x\rightarrow\infty}{\text{lim}}cx^{n} = \text{−}\underset{x\rightarrow\infty}{\text{lim}}x^{n}\ \text{and}\ \underset{x\rightarrow\text{−}\infty}{\text{lim}}cx^{n} = \text{−}\underset{x\rightarrow\text{−}\infty}{\text{lim}}x^{n}\ \text{if}\ c < 0.
$$

If $c = 0,y = cx^{n} = 0,$ in which case $\underset{x\rightarrow\infty}{\text{lim}}cx^{n} = 0 = \underset{x\rightarrow\text{-}\infty}{\text{lim}}cx^{n}.$

### Example 4.24

#### Limits at Infinity for Power Functions

For each function $f,$ evaluate $\underset{x\rightarrow\infty}{\text{lim}}f(x)$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x).$

1.  $f(x) = -5x^{3}$
2.  $f(x) = 2x^{4}$

#### Solution

1.  Since the coefficient of $x^{3}$ is $-5,$ the graph of $f(x) = -5x^{3}$ involves a vertical stretch and reflection of the graph of $y = x^{3}$ about the $x$-axis. Therefore, $\underset{x\rightarrow\infty}{\text{lim}}\left( {-5x^{3}} \right) = \text{-}\infty$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}\left( {-5x^{3}} \right) = \infty.$
2.  Since the coefficient of $x^{4}$ is $2,$ the graph of $f(x) = 2x^{4}$ is a vertical stretch of the graph of $y = x^{4}.$ Therefore, $\underset{x\rightarrow\infty}{\text{lim}}2x^{4} = \infty$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}2x^{4} = \infty.$

### Checkpoint 4.23

Let $f(x) = -3x^{4}.$ Find $\underset{x\rightarrow\infty}{\text{lim}}f(x).$

We now look at how the limits at infinity for power functions can be used to determine $\underset{x\rightarrow\text{\pm}\infty}{\text{lim}}f(x)$ for any polynomial function $f.$ Consider a polynomial function

$$
f(x) = a_{n}x^{n} + a_{n - 1}x^{n - 1} + \text{…} + a_{1}x + a_{0}
$$

of degree $n \geq 1$ so that $a_{n} \neq 0.$ Factoring, we see that

$$
f(x) = a_{n}x^{n}\left( {1 + \frac{a_{n - 1}}{a_{n}}\ \frac{1}{x} + \text{…} + \frac{a_{1}}{a_{n}}\ \frac{1}{x^{n - 1}} + \frac{a_{0}}{a_{n}}}\frac{1}{x^{n}} \right).
$$

As $x\rightarrow\text{\pm}\infty,$ all the terms inside the parentheses approach zero except the first term. We conclude that

$$
\underset{x\rightarrow\text{±}\infty}{\text{lim}}f(x) = \underset{x\rightarrow\text{±}\infty}{\text{lim}}a_{n}x^{n}.
$$

For example, the function $f(x) = 5x^{3} - 3x^{2} + 4$ behaves like $g(x) = 5x^{3}$ as $x\rightarrow\text{\pm}\infty$ as shown in Figure 4.52 and Table 4.4.

*Figure 4.52 The end behavior of a polynomial is determined by the behavior of the term with the largest exponent.*

|  |  |  |  |
|----|----|----|----|
| **$x$** | $10$ | $100$ | $1000$ |
| **$f(x) = 5x^{3} - 3x^{2} + 4$** | $4704$ | $4,970,004$ | $4,997,000,004$ |
| **$g(x) = 5x^{3}$** | $5000$ | $5,000,000$ | $5,000,000,000$ |
| **$x$** | $-10$ | $-100$ | $-1000$ |
| **$f(x) = 5x^{3} - 3x^{2} + 4$** | $-5296$ | $-5,029,996$ | $-5,002,999,996$ |
| **$g(x) = 5x^{3}$** | $-5000$ | $-5,000,000$ | $-5,000,000,000$ |

Table 4.4 A polynomial’s end behavior is determined by the term with the largest exponent.

#### End Behavior for Algebraic Functions

The end behavior for rational functions and functions involving radicals is a little more complicated than for polynomials. In Example 4.25, we show that the limits at infinity of a rational function $f(x) = \frac{p(x)}{q(x)}$ depend on the relationship between the degree of the numerator and the degree of the denominator. To evaluate the limits at infinity for a rational function, we divide the numerator and denominator by the highest power of $x$ appearing in the denominator. This determines which term in the overall expression dominates the behavior of the function at large values of $x.$

### Example 4.25

#### Determining End Behavior for Rational Functions

For each of the following functions, determine the limits as $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty.$ Then, use this information to describe the end behavior of the function.

1.  $f(x) = \frac{3x - 1}{2x + 5}$ (*Note:* The degree of the numerator and the denominator are the same.)
2.  $f(x) = \frac{3x^{2} + 2x}{4x^{3} - 5x + 7}$ (*Note:* The degree of numerator is less than the degree of the denominator.)
3.  $f(x) = \frac{3x^{2} + 4x}{x + 2}$ (*Note:* The degree of numerator is greater than the degree of the denominator.)

#### Solution

1.  The highest power of $x$ in the denominator is $x.$ Therefore, dividing the numerator and denominator by $x$ and applying the algebraic limit laws, we see that  
    ``` math
    \begin{array}{cl}
    {\underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{3x - 1}{2x + 5}} & {= \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{3 - 1\text{/}x}{2 + 5\text{/}x}} \\
     & {= \frac{\underset{x\rightarrow\text{±}\infty}{\text{lim}}\left( {3 - 1\text{/}x} \right)}{\underset{x\rightarrow\text{±}\infty}{\text{lim}}\left( {2 + 5\text{/}x} \right)}} \\
     & {= \frac{\underset{x\rightarrow\text{±}\infty}{\text{lim}}3 - \underset{x\rightarrow\text{±}\infty}{\text{lim}}1\text{/}x}{\underset{x\rightarrow\text{±}\infty}{\text{lim}}2 + \underset{x\rightarrow\text{±}\infty}{\text{lim}}5\text{/}x}} \\
     & {= \frac{3 - 0}{2 + 0} = \frac{3}{2}.}
    \end{array}
    ```
      
    Since $\underset{x\rightarrow\text{\pm}\infty}{\text{lim}}f(x) = \frac{3}{2},$ we know that $y = \frac{3}{2}$ is a horizontal asymptote for this function as shown in the following graph.  

    *Figure 4.53 The graph of this rational function approaches a horizontal asymptote as x → ± ∞ . x → ± ∞ .*

2.  Since the largest power of $x$ appearing in the denominator is $x^{3},$ divide the numerator and denominator by $x^{3}.$ After doing so and applying algebraic limit laws, we obtain  
    ``` math
    \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{3x^{2} + 2x}{4x^{3} - 5x + 7} = \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{3\text{/}x + 2\text{/}x^{2}}{4 - 5\text{/}x^{2} + 7\text{/}x^{3}} = \frac{3(0) + 2(0)}{4 - 5(0) + 7(0)} = 0.
    ```
      
    Therefore $f$ has a horizontal asymptote of $y = 0$ as shown in the following graph.  

    *Figure 4.54 The graph of this rational function approaches the horizontal asymptote y = 0 y = 0 as x → ± ∞ . x → ± ∞ .*

3.  Dividing the numerator and denominator by $x,$ we have  
    ``` math
    \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{3x^{2} + 4x}{x + 2} = \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{3x + 4}{1 + 2\text{/}x}.
    ```
      
    As $x\rightarrow\text{\pm}\infty,$ the denominator approaches $1.$ As $x\rightarrow\infty,$ the numerator approaches $+ \infty.$ As $x\rightarrow\text{-}\infty,$ the numerator approaches $\text{-}\infty.$ Therefore $\underset{x\rightarrow\infty}{\text{lim}}f(x) = \infty,$ whereas $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x) = \text{-}\infty$ as shown in the following figure.  

    *Figure 4.55 As x → ∞ , x → ∞ , the values f ( x ) → ∞ . f ( x ) → ∞ . As x → − ∞ , x → − ∞ , the values f ( x ) → − ∞ . f ( x ) → − ∞ .*

### Checkpoint 4.24

Evaluate $\underset{x\rightarrow\text{\pm}\infty}{\text{lim}}\frac{3x^{2} + 2x - 1}{5x^{2} - 4x + 7}$ and use these limits to determine the end behavior of $f(x) = \frac{3x^{2} + 2x - 1}{5x^{2} - 4x + 7}.$

Before proceeding, consider the graph of $f(x) = \frac{\left( {3x^{2} + 4x} \right)}{\left( {x + 2} \right)}$ shown in Figure 4.56. As $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty,$ the graph of $f$ appears almost linear. Although $f$ is certainly not a linear function, we now investigate why the graph of $f$ seems to be approaching a linear function. First, using long division of polynomials, we can write

$$
f(x) = \frac{3x^{2} + 4x}{x + 2} = 3x - 2 + \frac{4}{x + 2}.
$$

Since $\frac{4}{\left( {x + 2} \right)}\rightarrow 0$ as $x\rightarrow\text{\pm}\infty,$ we conclude that

$$
\underset{x\rightarrow\text{±}\infty}{\text{lim}}\left( {f(x) - \left( {3x - 2} \right)} \right) = \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{4}{x + 2} = 0.
$$

Therefore, the graph of $f$ approaches the line $y = 3x - 2$ as $x\rightarrow\text{\pm}\infty.$ This line is known as an oblique asymptote for $f$ (Figure 4.56).

*Figure 4.56 The graph of the rational function f ( x ) = ( 3 x 2 + 4 x ) / ( x + 2 ) f ( x ) = ( 3 x 2 + 4 x ) / ( x + 2 ) approaches the oblique asymptote y = 3 x − 2 as x → ± ∞ . y = 3 x − 2 as x → ± ∞ .*

We can summarize the results of Example 4.25 to make the following conclusion regarding end behavior for rational functions. Consider a rational function

$$
f(x) = \frac{p(x)}{q(x)} = \frac{a_{n}x^{n} + a_{n - 1}x^{n - 1} + \text{…} + a_{1}x + a_{0}}{b_{m}x^{m} + b_{m - 1}x^{m - 1} + \text{…} + b_{1}x + b_{0}},
$$

where $a_{n} \neq 0\ \text{and}\ b_{m} \neq 0.$

1.  If the degree of the numerator is the same as the degree of the denominator $\left( {n = m} \right),$ then $f$ has a horizontal asymptote of $y = a_{n}\text{/}b_{m}$ as $x\rightarrow\text{\pm}\infty.$
2.  If the degree of the numerator is less than the degree of the denominator $\left( {n < m} \right),$ then $f$ has a horizontal asymptote of $y = 0$ as $x\rightarrow\text{\pm}\infty.$
3.  If the degree of the numerator is greater than the degree of the denominator $\left( {n > m} \right),$ then $f$ does not have a horizontal asymptote. The limits at infinity are either positive or negative infinity, depending on the signs of the leading terms. In addition, using long division, the function can be rewritten as  
    ``` math
    f(x) = \frac{p(x)}{q(x)} = g(x) + \frac{r(x)}{q(x)},
    ```
      
    where the degree of $r(x)$ is less than the degree of $q(x).$ As a result, $\underset{x\rightarrow\text{\pm}\infty}{\text{lim}}r(x)\text{/}q(x) = 0.$ Therefore, the values of $\left\lbrack {f(x) - g(x)} \right\rbrack$ approach zero as $x\rightarrow\text{\pm}\infty.$ If the degree of $p(x)$ is exactly one more than the degree of $q(x)$ $\left( {n = m + 1} \right),$ the function $g(x)$ is a linear function. In this case, we call $g(x)$ an oblique asymptote.  
    Now let’s consider the end behavior for functions involving a radical.

### Example 4.26

#### Determining End Behavior for a Function Involving a Radical

Find the limits as $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty$ for $f(x) = \frac{3x - 2}{\sqrt{4x^{2} + 5}}$ and describe the end behavior of $f.$

#### Solution

Let’s use the same strategy as we did for rational functions: divide the numerator and denominator by a power of $x.$ To determine the appropriate power of $x,$ consider the expression $\sqrt{4x^{2} + 5}$ in the denominator. Since

$$
\sqrt{4x^{2} + 5} \approx \sqrt{4x^{2}} = 2|x|
$$

for large values of $x$ in effect $x$ appears just to the first power in the denominator. Therefore, we divide the numerator and denominator by $|x|.$ Then, using the fact that $|x| = x$ for $x > 0,$ $|x| = \text{-}x$ for $x < 0,$ and $|x| = \sqrt{x^{2}}$ for all $x,$ we calculate the limits as follows:

$$
\begin{array}{cll}
{\underset{x\rightarrow\infty}{\text{lim}}\frac{3x - 2}{\sqrt{4x^{2} + 5}}} & = & {\underset{x\rightarrow\infty}{\text{lim}}\frac{\left( {1\text{/}|x|} \right)\left( {3x - 2} \right)}{\left( {1\text{/}|x|} \right)\sqrt{4x^{2} + 5}}} \\
 & = & {\underset{x\rightarrow\infty}{\text{lim}}\frac{\left( {1\text{/}x} \right)\left( {3x - 2} \right)}{\sqrt{\left( {1\text{/}x^{2}} \right)\left( {4x^{2} + 5} \right)}}} \\
 & = & {\underset{x\rightarrow\infty}{\text{lim}}\frac{3 - 2\text{/}x}{\sqrt{4 + 5\text{/}x^{2}}} = \frac{3}{\sqrt{4}} = \frac{3}{2}} \\
{\underset{x\rightarrow\text{−}\infty}{\text{lim}}\frac{3x - 2}{\sqrt{4x^{2} + 5}}} & = & {\underset{x\rightarrow\text{−}\infty}{\text{lim}}\frac{\left( {1\text{/}|x|} \right)\left( {3x - 2} \right)}{\left( {1\text{/}|x|} \right)\sqrt{4x^{2} + 5}}} \\
 & = & {\underset{x\rightarrow\text{−}\infty}{\text{lim}}\frac{\left( {-1\text{/}x} \right)\left( {3x - 2} \right)}{\sqrt{\left( {1\text{/}x^{2}} \right)\left( {4x^{2} + 5} \right)}}} \\
 & = & {\underset{x\rightarrow\text{−}\infty}{\text{lim}}\frac{-3 + 2\text{/}x}{\sqrt{4 + 5\text{/}x^{2}}} = \frac{-3}{\sqrt{4}} = \frac{-3}{2}.}
\end{array}
$$

Therefore, $f(x)$ approaches the horizontal asymptote $y = \frac{3}{2}$ as $x\rightarrow\infty$ and the horizontal asymptote $y = - \frac{3}{2}$ as $x\rightarrow\text{-}\infty$ as shown in the following graph.

*Figure 4.57 This function has two horizontal asymptotes and it crosses one of the asymptotes.*

### Checkpoint 4.25

Evaluate $\underset{x\rightarrow\infty}{\text{lim}}\frac{\sqrt{3x^{2} + 4}}{x + 6}.$

#### Determining End Behavior for Transcendental Functions

The six basic trigonometric functions are periodic and do not approach a finite limit as $x\rightarrow\text{\pm}\infty.$ For example, $\text{sin} x$ oscillates between $1\ \text{and}\ -1$ (Figure 4.58). The tangent function, $\text{tan}(x)$, has an infinite number of vertical asymptotes as $x\rightarrow\text{\pm}\infty;$ therefore, it does not approach a finite limit nor does it approach $\text{\pm}\infty$ as $x\rightarrow\text{\pm}\infty$ as shown in Figure 4.59.

*Figure 4.58 The function f ( x ) = sin x f ( x ) = sin x oscillates between 1 and −1 1 and −1 as x → ± ∞ x → ± ∞*

*Figure 4.59 The function f ( x ) = tan x f ( x ) = tan x does not approach a limit and does not approach ± ∞ ± ∞ as x → ± ∞ x → ± ∞*

Recall that for any base $b > 0,b \neq 1,$ the function $y = b^{x}$ is an exponential function with domain $\left( {\text{-}\infty,\infty} \right)$ and range $\left( {0,\infty} \right).$ If $b > 1,y = b^{x}$ is increasing over $`\left( {\text{-}\infty,\infty} \right).$ If $0 < b < 1,$ $y = b^{x}$ is decreasing over $\left( {\text{-}\infty,\infty} \right).$ For the natural exponential function $f(x) = e^{x},$ $e \approx 2.718 > 1.$ Therefore, $f(x) = e^{x}$ is increasing on $`\left( {\text{-}\infty,\infty} \right)$ and the range is $`\left( {0,\infty} \right).$ The exponential function $f(x) = e^{x}$ approaches $\infty$ as $x\rightarrow\infty$ and approaches $0$ as $x\rightarrow\text{-}\infty$ as shown in Table 4.5 and Figure 4.60.

|               |             |           |       |           |             |
|---------------|-------------|-----------|-------|-----------|-------------|
| **$x$**     | $-5$      | $-2$    | $0$ | $2$     | $5$       |
| **$e^{x}$** | $0.00674$ | $0.135$ | $1$ | $7.389$ | $148.413$ |

Table 4.5 End behavior of the natural exponential function

*Figure 4.60 The exponential function approaches zero as x → − ∞ x → − ∞ and approaches ∞ ∞ as x → ∞ . x → ∞ .*

Recall that the natural logarithm function $f(x) = \text{ln}(x)$ is the inverse of the natural exponential function $y = e^{x}.$ Therefore, the domain of $f(x) = \text{ln}(x)$ is $\left( {0,\infty} \right)$ and the range is $\left( {\text{-}\infty,\infty} \right).$ The graph of $f(x) = \text{ln}(x)$ is the reflection of the graph of $y = e^{x}$ about the line $y = x.$ Therefore, $\text{ln}(x)\rightarrow\text{-}\infty$ as $x\rightarrow 0^{+}$ and $\text{ln}(x)\rightarrow\infty$ as $x\rightarrow\infty$ as shown in Figure 4.61 and Table 4.6.

|                      |            |            |       |           |           |
|----------------------|------------|------------|-------|-----------|-----------|
| **$x$**            | $0.01$   | $0.1$    | $1$ | $10$    | $100$   |
| **$\text{ln}(x)$** | $-4.605$ | $-2.303$ | $0$ | $2.303$ | $4.605$ |

Table 4.6 End behavior of the natural logarithm function

*Figure 4.61 The natural logarithm function approaches ∞ ∞ as x → ∞ . x → ∞ .*

### Example 4.27

#### Determining End Behavior for a Transcendental Function

Find the limits as $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty$ for $f(x) = \frac{\left( {2 + 3e^{x}} \right)}{\left( {7 - 5e^{x}} \right)}$ and describe the end behavior of $f.$

#### Solution

To find the limit as $x\rightarrow\infty,$ divide the numerator and denominator by $e^{x}\text{:}$

$$
\begin{array}{cl}
{\underset{x\rightarrow\infty}{\text{lim}}f(x)} & {= \underset{x\rightarrow\infty}{\text{lim}}\frac{2 + 3e^{x}}{7 - 5e^{x}}} \\
 & {= \underset{x\rightarrow\infty}{\text{lim}}\frac{\left( {2\text{/}e^{x}} \right) + 3}{\left( {7\text{/}e^{x}} \right) - 5}.}
\end{array}
$$

As shown in Figure 4.60, $e^{x}\rightarrow\infty$ as $x\rightarrow\infty.$ Therefore,

$$
\underset{x\rightarrow\infty}{\text{lim}}\frac{2}{e^{x}} = 0 = \underset{x\rightarrow\infty}{\text{lim}}\frac{7}{e^{x}}.
$$

We conclude that $\underset{x\rightarrow\infty}{\text{lim}}f(x) = - \frac{3}{5},$ and the graph of $f$ approaches the horizontal asymptote $y = - \frac{3}{5}$ as $x\rightarrow\infty.$ To find the limit as $x\rightarrow\text{-}\infty,$ use the fact that $e^{x}\rightarrow 0$ as $x\rightarrow\text{-}\infty$ to conclude that $\underset{x\rightarrow\infty}{\text{lim}}f(x) = \frac{2}{7},$ and therefore the graph of approaches the horizontal asymptote $y = \frac{2}{7}$ as $x\rightarrow\text{-}\infty.$

### Checkpoint 4.26

Find the limits as $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty$ for $f(x) = \frac{\left( {3e^{x} - 4} \right)}{\left( {5e^{x} + 2} \right)}.$

### Guidelines for Drawing the Graph of a Function

We now have enough analytical tools to draw graphs of a wide variety of algebraic and transcendental functions. Before showing how to graph specific functions, let’s look at a general strategy to use when graphing any function.

### Problem-Solving Strategy

#### Drawing the Graph of a Function

Given a function $f,$ use the following steps to sketch a graph of $f\text{:}$

1.  Determine the domain of the function.
2.  Locate the $x$- and $y$-intercepts.
3.  Evaluate $\underset{x\rightarrow\infty}{\text{lim}}f(x)$ and $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x)$ to determine the end behavior. If either of these limits is a finite number $L,$ then $y = L$ is a horizontal asymptote. If either of these limits is $\infty$ or $\text{-}\infty,$ determine whether $f$ has an oblique asymptote. If $f$ is a rational function such that $f(x) = \frac{p(x)}{q(x)},$ where the degree of the numerator is greater than the degree of the denominator, then $f$ can be written as  
    ``` math
    f(x) = \frac{p(x)}{q(x)} = g(x) + \frac{r(x)}{q(x)},
    ```
      
    where the degree of $r(x)$ is less than the degree of $q(x).$ The values of $f(x)$ approach the values of $g(x)$ as $x\rightarrow\text{\pm}\infty.$ If $g(x)$ is a linear function, it is known as an *oblique asymptote*.
4.  Determine whether $f$ has any vertical asymptotes.
5.  Calculate $f'.$ Find all critical points and determine the intervals where $f$ is increasing and where $f$ is decreasing. Determine whether $f$ has any local extrema.
6.  Calculate $f^{''}.$ Determine the intervals where $f$ is concave up and where $f$ is concave down. Use this information to determine whether $f$ has any inflection points. The second derivative can also be used as an alternate means to determine or verify that $f$ has a local extremum at a critical point.

Now let’s use this strategy to graph several different functions. We start by graphing a polynomial function.

### Example 4.28

#### Sketching a Graph of a Polynomial

Sketch a graph of $f(x) = \left( {x - 1} \right)^{2}\left( {x + 2} \right).$

#### Solution

Step 1. Since $f$ is a polynomial, the domain is the set of all real numbers.

Step 2. When $x = 0,f(x) = 2.$ Therefore, the $y$-intercept is $\left( {0,2} \right).$ To find the $x$-intercepts, we need to solve the equation $\left( {x - 1} \right)^{2}\left( {x + 2} \right) = 0,$ gives us the $x$-intercepts $\left( {1,0} \right)$ and $\left( {-2,0} \right)$

Step 3. We need to evaluate the end behavior of $f.$ As $x\rightarrow\infty,$ $\left( {x - 1} \right)^{2}\rightarrow\infty$ and $\left( {x + 2} \right)\rightarrow\infty.$ Therefore, $\underset{x\rightarrow\infty}{\text{lim}}f(x) = \infty.$ As $x\rightarrow\text{-}\infty,$ $\left( {x - 1} \right)^{2}\rightarrow\infty$ and $\left( {x + 2} \right)\rightarrow\text{-}\infty.$ Therefore, $\underset{x\rightarrow\text{-}\infty}{\text{lim}}f(x) = \text{-}\infty.$ To get even more information about the end behavior of $f,$ we can multiply the factors of $f.$ When doing so, we see that

$$
f(x) = \left( {x - 1} \right)^{2}\left( {x + 2} \right) = x^{3} - 3x + 2.
$$

Since the leading term of $f$ is $x^{3},$ we conclude that $f$ behaves like $y = x^{3}$ as $x\rightarrow\text{\pm}\infty.$

Step 4. Since $f$ is a polynomial function, it does not have any vertical asymptotes.

Step 5. The first derivative of $f$ is

$$
f'(x) = 3x^{2} - 3.
$$

Therefore, $f$ has two critical points: $x = 1,-1.$ Divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the three smaller intervals: $\left( {\text{-}\infty,-1} \right),$ $\left( {-1,1} \right),$ and $\left( {1,\infty} \right).$ Then, choose test points $x = -2,$ $x = 0,$ and $x = 2$ from these intervals and evaluate the sign of $f'(x)$ at each of these test points, as shown in the following table.

| Interval | Test Point | Sign of Derivative $f'(x) = 3x^{2} - 3 = 3\left( {x - 1} \right)\left( {x + 1} \right)$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,-1} \right)$ | $x = -2$ | $\left( \text{+} \right)\left( \text{-} \right)\left( \text{-} \right) = +$ | $f$ is increasing. |
| $\left( {-1,1} \right)$ | $x = 0$ | $\left( \text{+} \right)\left( \text{-} \right)\left( \text{+} \right) = \text{-}$ | $f$ is decreasing. |
| $\left( {1,\infty} \right)$ | $x = 2$ | $\left( \text{+} \right)\left( \text{+} \right)\left( \text{+} \right) = +$ | $f$ is increasing. |

From the table, we see that $f$ has a local maximum at $x = -1$ and a local minimum at $x = 1.$ Evaluating $f(x)$ at those two points, we find that the local maximum value is $f(-1) = 4$ and the local minimum value is $f(1) = 0.$

Step 6. The second derivative of $f$ is

$$
f^{''}(x) = 6x.
$$

The second derivative is zero at $x = 0.$ Therefore, to determine the concavity of $f,$ divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the smaller intervals $\left( {\text{-}\infty,0} \right)$ and $\left( {0,\infty} \right),$ and choose test points $x = -1$ and $x = 1$ to determine the concavity of $f$ on each of these smaller intervals as shown in the following table.

| Interval | Test Point | Sign of $f^{''}(x) = 6x$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,0} \right)$ | $x = -1$ | $-$ | $f$ is concave down. |
| $\left( {0,\infty} \right)$ | $x = 1$ | $+$ | $f$ is concave up. |

We note that the information in the preceding table confirms the fact, found in step $5,$ that $f$ has a local maximum at $x = -1$ and a local minimum at $x = 1.$ In addition, the information found in step $5$—namely, $f$ has a local maximum at $x = -1$ and a local minimum at $x = 1,$ and $f'(x) = 0$ at those points—combined with the fact that $f^{''}$ changes sign only at $x = 0$ confirms the results found in step $6$ on the concavity of $f.$

Combining this information, we arrive at the graph of $f(x) = \left( {x - 1} \right)^{2}\left( {x + 2} \right)$ shown in the following graph.

### Checkpoint 4.27

Sketch a graph of $f(x) = \left( {x - 1} \right)^{3}\left( {x + 2} \right).$

### Example 4.29

#### Sketching a Rational Function

Sketch the graph of ${f(x) = \frac{x^{2}}{\left( {1 - x^{2}} \right)}}\text{.}$

#### Solution

Step 1. The function $f$ is defined as long as the denominator is not zero. Therefore, the domain is the set of all real numbers $x$ except $x = \text{\pm}1.$

Step 2. Find the intercepts. If $x = 0,$ then $f(x) = 0,$ so $0$ is an intercept. If $y = 0,$ then $\frac{x^{2}}{\left( {1 - x^{2}} \right)} = 0,$ which implies $x = 0.$ Therefore, $\left( {0,0} \right)$ is the only intercept.

Step 3. Evaluate the limits at infinity. Since $f$ is a rational function, divide the numerator and denominator by the highest power in the denominator: $x^{2}.$ We obtain

$$
\underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{x^{2}}{1 - x^{2}} = \underset{x\rightarrow\text{±}\infty}{\text{lim}}\frac{1}{\frac{1}{x^{2}} - 1} = -1.
$$

Therefore, $f$ has a horizontal asymptote of $y = -1$ as $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty.$

Step 4. To determine whether $f$ has any vertical asymptotes, first check to see whether the denominator has any zeroes. We find the denominator is zero when $x = \text{\pm}1.$ To determine whether the lines $x = 1$ or $x = -1$ are vertical asymptotes of $f,$ evaluate $\underset{x\rightarrow 1}{\text{lim}}f(x)$ and $\underset{x\rightarrow\text{-}1}{\text{lim}}f(x).$ By looking at each one-sided limit as $x\rightarrow 1,$ we see that

$$
\underset{x\rightarrow 1^{+}}{\text{lim}}\frac{x^{2}}{1 - x^{2}} = \text{−}\infty\ \text{and}\ \underset{x\rightarrow 1^{-}}{\text{lim}}\frac{x^{2}}{1 - x^{2}} = \infty.
$$

In addition, by looking at each one-sided limit as $x\rightarrow\text{-}1,$ we find that

$$
\underset{x\rightarrow\text{−}1^{+}}{\text{lim}}\frac{x^{2}}{1 - x^{2}} = \infty\ \text{and}\ \underset{x\rightarrow\text{−}1^{-}}{\text{lim}}\frac{x^{2}}{1 - x^{2}} = \text{−}\infty.
$$

Step 5. Calculate the first derivative:

$$
f'(x) = \frac{\left( {1 - x^{2}} \right)\left( {2x} \right) - x^{2}\left( {-2x} \right)}{\left( {1 - x^{2}} \right)^{2}} = \frac{2x}{\left( {1 - x^{2}} \right)^{2}}.
$$

Critical points occur at points $x$ where $f'(x) = 0$ or $f'(x)$ is undefined. We see that $f'(x) = 0$ when $x = 0.$ The derivative $f'$ is not undefined at any point in the domain of $f.$ However, $x = \text{\pm}1$ are not in the domain of $f.$ Therefore, to determine where $f$ is increasing and where $f$ is decreasing, divide the interval $\left( {\text{-}\infty,\infty} \right)$ into four smaller intervals: $\left( {\text{-}\infty,-1} \right),$ $\left( {-1,0} \right),$ $\left( {0,1} \right),$ and $\left( {1,\infty} \right),$ and choose a test point in each interval to determine the sign of $f'(x)$ in each of these intervals. The values $x = -2,$ $x = - \frac{1}{2},$ $x = \frac{1}{2},$ and $x = 2$ are good choices for test points as shown in the following table.

| Interval | Test Point | Sign of $f'(x) = \frac{2x}{\left( {1 - x^{2}} \right)^{2}}$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,-1} \right)$ | $x = -2$ | $\text{-}\text{/} + = \text{-}$ | $f$ is decreasing. |
| $\left( {-1,0} \right)$ | $x = -1\text{/}2$ | $\text{-}\text{/} + = \text{-}$ | $f$ is decreasing. |
| $\left( {0,1} \right)$ | $x = 1\text{/}2$ | $+ \text{/} + = +$ | $f$ is increasing. |
| $\left( {1,\infty} \right)$ | $x = 2$ | $+ \text{/} + = +$ | $f$ is increasing. |

From this analysis, we conclude that $f$ has a local minimum at $x = 0$ but no local maximum.

Step 6. Calculate the second derivative:

$$
\begin{matrix}
{f^{''}(x)} & {= \frac{\left( {1 - x^{2}} \right)^{2}(2) - 2x\left( {2\left( {1 - x^{2}} \right)\left( {-2x} \right)} \right)}{\left( {1 - x^{2}} \right)^{4}}} \\
 & {= \frac{\left( {1 - x^{2}} \right)\left\lbrack {2\left( {1 - x^{2}} \right) + 8x^{2}} \right\rbrack}{\left( {1 - x^{2}} \right)^{4}}} \\
 & {= \frac{2\left( {1 - x^{2}} \right) + 8x^{2}}{\left( {1 - x^{2}} \right)^{3}}} \\
 & {= \frac{6x^{2} + 2}{\left( {1 - x^{2}} \right)^{3}}.}
\end{matrix}
$$

To determine the intervals where $f$ is concave up and where $f$ is concave down, we first need to find all points $x$ where $f^{''}(x) = 0$ or $f^{''}(x)$ is undefined. Since the numerator $6x^{2} + 2 \neq 0$ for any $x,$ $f^{''}(x)$ is never zero. Furthermore, $f^{''}$ is not undefined for any $x$ in the domain of $f.$ However, as discussed earlier, $x = \text{\pm}1$ are not in the domain of $f.$ Therefore, to determine the concavity of $f,$ we divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the three smaller intervals $\left( {\text{-}\infty,-1} \right),$ $\left( {-1,-1} \right),$ and $\left( {1,\infty} \right),$ and choose a test point in each of these intervals to evaluate the sign of $f^{''}(x).$ in each of these intervals. The values $x = -2,$ $x = 0,$ and $x = 2$ are possible test points as shown in the following table.

| Interval | Test Point | Sign of $f^{''}(x) = \frac{6x^{2} + 2}{\left( {1 - x^{2}} \right)^{3}}$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,-1} \right)$ | $x = -2$ | $+ \text{/} - = \text{-}$ | $f$ is concave down. |
| $\left( {-1,-1} \right)$ | $x = 0$ | $+ \text{/} + = +$ | $f$ is concave up. |
| $\left( {1,\infty} \right)$ | $x = 2$ | $+ \text{/} - = \text{-}$ | $f$ is concave down. |

Combining all this information, we arrive at the graph of $f$ shown below. Note that, although $f$ changes concavity at $x = -1$ and $x = 1,$ there are no inflection points at either of these places because $f$ is not continuous at $x = -1$ or $x = 1.$

### Checkpoint 4.28

Sketch a graph of $f(x) = \frac{\left( {3x + 5} \right)}{\left( {8 + 4x} \right)}.$

### Example 4.30

#### Sketching a Rational Function with an Oblique Asymptote

Sketch the graph of $f(x) = \frac{x^{2}}{\left( {x - 1} \right)}$

#### Solution

Step 1. The domain of $f$ is the set of all real numbers $x$ except $x = 1.$

Step 2. Find the intercepts. We can see that when $x = 0,$ $f(x) = 0,$ so $\left( {0,0} \right)$ is the only intercept.

Step 3. Evaluate the limits at infinity. Since the degree of the numerator is one more than the degree of the denominator, $f$ must have an oblique asymptote. To find the oblique asymptote, use long division of polynomials to write

$$
f(x) = \frac{x^{2}}{x - 1} = x + 1 + \frac{1}{x - 1}.
$$

Since $1\text{/}\left( {x - 1} \right)\rightarrow 0$ as $x\rightarrow\text{\pm}\infty,$ $f(x)$ approaches the line $y = x + 1$ as $x\rightarrow\text{\pm}\infty.$ The line $y = x + 1$ is an oblique asymptote for $f.$

Step 4. To check for vertical asymptotes, look at where the denominator is zero. Here the denominator is zero at $x = 1.$ Looking at both one-sided limits as $x\rightarrow 1,$ we find

$$
\underset{x\rightarrow 1^{+}}{\text{lim}}\frac{x^{2}}{x - 1} = \infty\ \text{and}\ \underset{x\rightarrow 1^{-}}{\text{lim}}\frac{x^{2}}{x - 1} = \text{−}\infty.
$$

Therefore, $x = 1$ is a vertical asymptote, and we have determined the behavior of $f$ as $x$ approaches $1$ from the right and the left.

Step 5. Calculate the first derivative:

$$
f'(x) = \frac{\left( {x - 1} \right)\left( {2x} \right) - x^{2}(1)}{\left( {x - 1} \right)^{2}} = \frac{x^{2} - 2x}{\left( {x - 1} \right)^{2}}.
$$

We have $f'(x) = 0$ when $x^{2} - 2x = x\left( {x - 2} \right) = 0.$ Therefore, $x = 0$ and $x = 2$ are critical points. Since $f$ is undefined at $x = 1,$ we need to divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the smaller intervals $\left( {\text{-}\infty,0} \right),$ $\left( {0,1} \right),$ $\left( {1,2} \right),$ and $\left( {2,\infty} \right),$ and choose a test point from each interval to evaluate the sign of $f'(x)$ in each of these smaller intervals. For example, let $x = -1,$ $x = \frac{1}{2},$ $x = \frac{3}{2},$ and $x = 3$ be the test points as shown in the following table.

| Interval | Test Point | Sign of $f'(x) = \frac{x^{2} - 2x}{\left( {x - 1} \right)^{2}} = \frac{x\left( {x - 2} \right)}{\left( {x - 1} \right)^{2}}$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,0} \right)$ | $x = -1$ | $\left( \text{-} \right)\left( \text{-} \right)\text{/} + = +$ | $f$ is increasing. |
| $\left( {0,1} \right)$ | $x = 1\text{/}2$ | $\left( \text{+} \right)\left( \text{-} \right)\text{/} + = \text{-}$ | $f$ is decreasing. |
| $\left( {1,2} \right)$ | $x = 3\text{/}2$ | $\left( \text{+} \right)\left( \text{-} \right)\text{/} + = \text{-}$ | $f$ is decreasing. |
| $\left( {2,\infty} \right)$ | $x = 3$ | $\left( \text{+} \right)\left( \text{+} \right)\text{/} + = +$ | $f$ is increasing. |

From this table, we see that $f$ has a local maximum at $x = 0$ and a local minimum at $x = 2.$ The value of $f$ at the local maximum is $f(0) = 0$ and the value of $f$ at the local minimum is $f(2) = 4.$ Therefore, $\left( {0,0} \right)$ and $\left( {2,4} \right)$ are important points on the graph.

Step 6. Calculate the second derivative:

$$
\begin{array}{cl}
{f^{''}(x)} & {= \frac{\left( {x - 1} \right)^{2}\left( {2x - 2} \right) - \left( {x^{2} - 2x} \right)\left( {2\left( {x - 1} \right)} \right)}{\left( {x - 1} \right)^{4}}} \\
 & {= \frac{\left( {x - 1} \right)\left\lbrack {\left( {x - 1} \right)\left( {2x - 2} \right) - 2\left( {x^{2} - 2x} \right)} \right\rbrack}{\left( {x - 1} \right)^{4}}} \\
 & {= \frac{\left( {x - 1} \right)\left( {2x - 2} \right) - 2\left( {x^{2} - 2x} \right)}{\left( {x - 1} \right)^{3}}} \\
 & {= \frac{2x^{2} - 4x + 2 - \left( {2x^{2} - 4x} \right)}{\left( {x - 1} \right)^{3}}} \\
 & {= \frac{2}{\left( {x - 1} \right)^{3}}.}
\end{array}
$$

We see that $f^{''}(x)$ is never zero or undefined for $x$ in the domain of $f.$ Since $f$ is undefined at $x = 1,$ to check concavity we just divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the two smaller intervals $\left( {\text{-}\infty,1} \right)$ and $\left( {1,\infty} \right),$ and choose a test point from each interval to evaluate the sign of $f^{''}(x)$ in each of these intervals. The values $x = 0$ and $x = 2$ are possible test points as shown in the following table.

| Interval | Test Point | Sign of $f^{''}(x) = \frac{2}{\left( {x - 1} \right)^{3}}$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,1} \right)$ | $x = 0$ | $+ \text{/} - = \text{-}$ | $f$ is concave down. |
| $\left( {1,\infty} \right)$ | $x = 2$ | $+ \text{/} + = +$ | $f$ is concave up. |

From the information gathered, we arrive at the following graph for $f.$

### Checkpoint 4.29

Find the oblique asymptote for $f(x) = \frac{\left( {3x^{3} - 2x + 1} \right)}{\left( {2x^{2} - 4} \right)}.$

### Example 4.31

#### Sketching the Graph of a Function with a Cusp

Sketch a graph of $f(x) = \left( {x - 1} \right)^{2\text{/}3}.$

#### Solution

Step 1. Since the cube-root function is defined for all real numbers $x$ and $\left( {x - 1} \right)^{2\text{/}3} = \left( \sqrt[3]{x - 1} \right)^{2},$ the domain of $f$ is all real numbers.

Step 2: To find the $y$-intercept, evaluate $f(0).$ Since $f(0) = 1,$ the $y$-intercept is $\left( {0,1} \right).$ To find the $x$-intercept, solve $\left( {x - 1} \right)^{2\text{/}3} = 0.$ The solution of this equation is $x = 1,$ so the $x$-intercept is $\left( {1,0} \right).$

Step 3: Since $\underset{x\rightarrow\text{\pm}\infty}{\text{lim}}\left( {x - 1} \right)^{2\text{/}3} = \infty,$ the function continues to grow without bound as $x\rightarrow\infty$ and $x\rightarrow\text{-}\infty.$

Step 4: The function has no vertical asymptotes.

Step 5: To determine where $f$ is increasing or decreasing, calculate $f'.$ We find

$$
f'(x) = \frac{2}{3}\left( {x - 1} \right)^{-1\text{/}3} = \frac{2}{3\left( {x - 1} \right)^{1\text{/}3}}.
$$

This function is not zero anywhere, but it is undefined when $x = 1.$ Therefore, the only critical point is $x = 1.$ Divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the smaller intervals $\left( {\text{-}\infty,1} \right)$ and $\left( {1,\infty} \right),$ and choose test points in each of these intervals to determine the sign of $f'(x)$ in each of these smaller intervals. Let $x = 0$ and $x = 2$ be the test points as shown in the following table.

| Interval | Test Point | Sign of $f'(x) = \frac{2}{3\left( {x - 1} \right)^{1\text{/}3}}$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,1} \right)$ | $x = 0$ | $+ \text{/} - = \text{-}$ | $f$ is decreasing. |
| $\left( {1,\infty} \right)$ | $x = 2$ | $+ \text{/} + = +$ | $f$ is increasing. |

We conclude that $f$ has a local minimum at $x = 1.$ Evaluating $f$ at $x = 1,$ we find that the value of $f$ at the local minimum is zero. Note that $f'(1)$ is undefined, so to determine the behavior of the function at this critical point, we need to examine $\underset{x\rightarrow 1}{\text{lim}}f'(x).$ Looking at the one-sided limits, we have

$$
\underset{x\rightarrow 1^{+}}{\text{lim}}\frac{2}{3\left( {x - 1} \right)^{1\text{/}3}} = \infty\ \text{and}\ \underset{x\rightarrow 1^{-}}{\text{lim}}\frac{2}{3\left( {x - 1} \right)^{1\text{/}3}} = \text{−}\infty.
$$

Therefore, $f$ has a cusp at $x = 1.$

Step 6: To determine concavity, we calculate the second derivative of $f\text{:}$

$$
f^{''}(x) = - \frac{2}{9}\left( {x - 1} \right)^{-4\text{/}3} = \frac{-2}{9\left( {x - 1} \right)^{4\text{/}3}}.
$$

We find that $f^{''}(x)$ is defined for all $x,$ but is undefined when $x = 1.$ Therefore, divide the interval $\left( {\text{-}\infty,\infty} \right)$ into the smaller intervals $\left( {\text{-}\infty,1} \right)$ and $\left( {1,\infty} \right),$ and choose test points to evaluate the sign of $f^{''}(x)$ in each of these intervals. As we did earlier, let $x = 0$ and $x = 2$ be test points as shown in the following table.

| Interval | Test Point | Sign of $f^{''}(x) = \frac{-2}{9\left( {x - 1} \right)^{4\text{/}3}}$ | Conclusion |
|----|----|----|----|
| $\left( {\text{-}\infty,1} \right)$ | $x = 0$ | $\text{-}\text{/} + = \text{-}$ | $f$ is concave down. |
| $\left( {1,\infty} \right)$ | $x = 2$ | $\text{-}\text{/} + = \text{-}$ | $f$ is concave down. |

From this table, we conclude that $f$ is concave down everywhere. Combining all of this information, we arrive at the following graph for $f.$

### Checkpoint 4.30

Consider the function $f(x) = 5 - x^{2\text{/}3}.$ Determine the point on the graph where a cusp is located. Determine the end behavior of $f.$

### Section 4.6 Exercises

For the following exercises, examine the graphs. Identify where the vertical asymptotes are located.

251\. 252. 253. 254. 255.

For the following functions $f(x),$ determine whether there is an asymptote at $x = a.$ Justify your answer without graphing on a calculator.

256\.

$f(x) = \frac{x + 1}{x^{2} + 5x + 4},a = -1$

257\.

$f(x) = \frac{x}{x - 2},a = 2$

258\.

$f(x) = \left( {x + 2} \right)^{3\text{/}2},a = -2$

259\.

$f(x) = \left( {x - 1} \right)^{-1\text{/}3},a = 1$

260\.

$f(x) = 1 + x^{-2\text{/}5},a = 1$

For the following exercises, evaluate the limit.

261\.

$\underset{x\rightarrow\infty}{\text{lim}}\frac{1}{3x + 6}$

262\.

$\underset{x\rightarrow\infty}{\text{lim}}\frac{2x - 5}{4x}$

263\.

$\underset{x\rightarrow\infty}{\text{lim}}\frac{x^{2} - 2x + 5}{x + 2}$

264\.

$\underset{x\rightarrow\text{-}\infty}{\text{lim}}\frac{3x^{3} - 2x}{x^{2} + 2x + 8}$

265\.

$\underset{x\rightarrow\text{-}\infty}{\text{lim}}\frac{x^{4} - 4x^{3} + 1}{2 - 2x^{2} - 7x^{4}}$

266\.

$\underset{x\rightarrow\infty}{\text{lim}}\frac{3x}{\sqrt{x^{2} + 1}}$

267\.

$\underset{x\rightarrow\text{-}\infty}{\text{lim}}\frac{\sqrt{4x^{2} - 1}}{x + 2}$

268\.

$\underset{x\rightarrow\infty}{\text{lim}}\frac{4x}{\sqrt{x^{2} - 1}}$

269\.

$\underset{x\rightarrow\text{-}\infty}{\text{lim}}\frac{4x}{\sqrt{x^{2} - 1}}$

270\.

$\underset{x\rightarrow\infty}{\text{lim}}\frac{2\sqrt{x}}{x - \sqrt{x} + 1}$

For the following exercises, find the horizontal and vertical asymptotes.

271\.

$f(x) = x - \frac{9}{x}$

272\.

$f(x) = \frac{1}{1 - x^{2}}$

273\.

$f(x) = \frac{x^{3}}{4 - x^{2}}$

274\.

$f(x) = \frac{x^{2} + 3}{x^{2} + 1}$

275\.

$f(x) = \text{sin}(x)\text{sin}\left( {2x} \right)$

276\.

$f(x) = \text{cos} x + \text{cos}\left( {3x} \right) + \text{cos}\left( {5x} \right)$

277\.

$f(x) = \frac{x\text{sin}(x)}{x^{2} - 1}$

278\.

$f(x) = \frac{x}{\text{sin}(x)}$

279\.

$f(x) = \frac{1}{x^{3} + x^{2}}$

280\.

$f(x) = \frac{1}{x - 1} - 2x$

281\.

$f(x) = \frac{x^{3} + 1}{x^{3} - 1}$

282\.

$f(x) = \frac{\text{sin} x + \text{cos} x}{\text{sin} x - \text{cos} x}$

283\.

$f(x) = x - \text{sin} x$

284\.

$f(x) = \frac{1}{x} - \sqrt{x}$

For the following exercises, construct a function $f(x)$ that has the given asymptotes.

285\.

$x = 1$ and $y = 2$

286\.

$x = 1$ and $y = 0$

287\.

$y = 4,$ $x = -1$

288\.

$x = 0$

For the following exercises, graph the function on a graphing calculator on the window $x = \left\lbrack {-5,5} \right\rbrack$ and estimate the horizontal asymptote or limit. Then, calculate the actual horizontal asymptote or limit.

289\.

**\[T\]** $f(x) = \frac{1}{x + 10}$

290\.

**\[T\]** $f(x) = \frac{x + 1}{x^{2} + 7x + 6}$

291\.

**\[T\]** $\underset{x\rightarrow\text{-}\infty}{\text{lim}}x^{2} + 10x + 25$

292\.

**\[T\]** $\underset{x\rightarrow\text{-}\infty}{\text{lim}}\frac{x + 2}{x^{2} + 7x + 6}$

293\.

**\[T\]** $\underset{x\rightarrow\infty}{\text{lim}}\frac{3x + 2}{x + 5}$

For the following exercises, draw a graph of the functions without using a calculator. Be sure to notice all important features of the graph: local maxima and minima, inflection points, and asymptotic behavior.

294\.

$y = 3x^{2} + 2x + 4$

295\.

$y = x^{3} - 3x^{2} + 4$

296\.

$y = \frac{2x + 1}{x^{2} + 6x + 5}$

297\.

$y = \frac{x^{3} + 4x^{2} + 3x}{3x + 9}$

298\.

$y = \frac{x^{2} + x - 2}{x^{2} - 3x - 4}$

299\.

$y = \sqrt{x^{2} - 5x + 4}$

300\.

$y = 2x\sqrt{16 - x^{2}}$

301\.

$y = \frac{\text{cos} x}{x},$ on $x = \left\lbrack {-2\pi,2\pi} \right\rbrack$

302\.

$y = \frac{\sqrt{x^{2} + 2}}{x + 1}$

303\.

$y = x\text{tan} x,x = \left\lbrack {\text{-}\pi,\pi} \right\rbrack$

304\.

$y = x\text{ln}(x),x > 0$

305\.

$y = x^{2}\text{sin}(x),x = \left\lbrack {-2\pi,2\pi} \right\rbrack$

306\.

For $f(x) = \frac{P(x)}{Q(x)}$ to have an asymptote at $y = 2$ then the polynomials $P(x)$ and $Q(x)$ must have what relation?

307\.

For $f(x) = \frac{P(x)}{Q(x)}$ to have an asymptote at $x = 0,$ then the polynomials $P(x)$ and $Q(x).$ must have what relation?

308\.

If $f'(x)$ has asymptotes at $y = 3$ and $x = 1,$ then $f(x)$ has what asymptotes?

309\.

Both $f(x) = \frac{1}{\left( {x - 1} \right)}$ and $g(x) = \frac{1}{\left( {x - 1} \right)^{2}}$ have asymptotes at $x = 1$ and $y = 0.$ What is the most obvious difference between these two functions?

310\.

True or false: Every ratio of polynomials has vertical asymptotes.
