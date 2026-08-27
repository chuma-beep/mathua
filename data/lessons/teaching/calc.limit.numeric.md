> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 2.2 The Limit of a Function

### Learning Objectives

- 2.2.1 Using correct notation, describe the limit of a function.
- 2.2.2 Use a table of values to estimate the limit of a function or to identify when the limit does not exist.
- 2.2.3 Use a graph to estimate the limit of a function or to identify when the limit does not exist.
- 2.2.4 Define one-sided limits and provide examples.
- 2.2.5 Explain the relationship between one-sided and two-sided limits.
- 2.2.6 Using correct notation, describe an infinite limit.
- 2.2.7 Define a vertical asymptote.

The concept of a limit or limiting process, essential to the understanding of calculus, has been around for thousands of years. In fact, early mathematicians used a limiting process to obtain better and better approximations of areas of circles. Yet, the formal definition of a limit—as we know and understand it today—did not appear until the late 19th century. We therefore begin our quest to understand limits, as our mathematical ancestors did, by using an intuitive approach. At the end of this chapter, armed with a conceptual understanding of limits, we examine the formal definition of a limit.

We begin our exploration of limits by taking a look at the graphs of the functions

$$
f(x) = \frac{x^{2} - 4}{x - 2},\ g(x) = \frac{\left| {x - 2} \right|}{x - 2},\ \text{and}\ h(x) = \frac{1}{\left( {x - 2} \right)^{2}},
$$

which are shown in Figure 2.12. In particular, let’s focus our attention on the behavior of each graph at and around $x = 2.$

*Figure 2.12 These graphs show the behavior of three different functions around x = 2 . x = 2 .*

Each of the three functions is undefined at $x = 2,$ but if we make this statement and no other, we give a very incomplete picture of how each function behaves in the vicinity of $x = 2.$ To express the behavior of each graph in the vicinity of 2 more completely, we need to introduce the concept of a limit.

### Intuitive Definition of a Limit

Let’s first take a closer look at how the function $f(x) = {{(x^{2} - 4)}\text{/}{(x - 2)}}$ behaves around $x = 2$ in Figure 2.12. As the values of *x* approach 2 from either side of 2, the values of $y = f(x)$ approach 4. Mathematically, we say that the limit of $f(x)$ as *x* approaches 2 is 4. Symbolically, we express this limit as

$$
\underset{x\rightarrow 2}{\text{lim}}f(x) = 4.
$$

From this very brief informal look at one limit, let’s start to develop an intuitive definition of the limit. We can think of the limit of a function at a number *a* as being the one real number *L* that the functional values approach as the *x*-values approach *a,* provided such a real number *L* exists. Stated more carefully, we have the following definition:

### Definition

Let $f(x)$ be a function defined at all values in an open interval containing *a*, with the possible exception of *a* itself, and let *L* be a real number. If *all* values of the function $f(x)$ approach the real number *L* as the values of $x\left( {\neq a} \right)$ approach the number *a*, then we say that the limit of $f(x)$ as *x* approaches *a* is *L*. (More succinct, as *x* gets closer to *a*, $f(x)$ gets closer and stays close to *L*.) Symbolically, we express this idea as

$$
\underset{x\rightarrow a}{\text{lim}}f(x) = L.
$$

(2.3)

We can estimate limits by constructing tables of functional values and by looking at their graphs. This process is described in the following Problem-Solving Strategy.

### Problem-Solving Strategy

#### Evaluating a Limit Using a Table of Functional Values

1.  To evaluate $\underset{x\rightarrow a}{\text{lim}}f(x),$ we begin by completing a table of functional values. We should choose two sets of *x*-values—one set of values approaching *a* and less than *a*, and another set of values approaching *a* and greater than *a*. Table 2.1 demonstrates what your tables might look like.  
    | *x* | $f(x)$ |  | *x* | $f(x)$ |
    |----|----|----|----|----|
    | $a - 0.1$ | $f\left( {a - 0.1} \right)$ |  | $a + 0.1$ | $f\left( {a + 0.1} \right)$ |
    | $a - 0.01$ | $f\left( {a - 0.01} \right)$ |  | $a + 0.01$ | $f\left( {a + 0.01} \right)$ |
    | $a - 0.001$ | $f\left( {a - 0.001} \right)$ |  | $a + 0.001$ | $f\left( {a + 0.001} \right)$ |
    | $a - 0.0001$ | $f\left( {a - 0.0001} \right)$ |  | $a + 0.0001$ | $f\left( {a + 0.0001} \right)$ |
    | Use additional values as necessary. |  |  | Use additional values as necessary. |  |

    Table 2.1 Table of Functional Values for $\underset{x\rightarrow a}{\text{lim}}f(x)$
2.  Next, let’s look at the values in each of the $f(x)$ columns and determine whether the values seem to be approaching a single value as we move down each column. In our columns, we look at the sequence $f\left( {a - 0.1} \right),f\left( {a - 0.01} \right),f\left( {a - 0.001} \right).,f\left( {a - 0.0001} \right),$ and so on, and $f\left( {a + 0.1} \right),f\left( {a + 0.01} \right),f\left( {a + 0.001} \right),f\left( {a + 0.0001} \right),$ and so on. (*Note*: Although we have chosen the *x*-values $a \pm 0.1,a \pm 0.01,a \pm 0.001,a \pm 0.0001,$ and so forth, and these values will probably work nearly every time, on very rare occasions we may need to modify our choices.)
3.  If both columns approach a common *y*-value *L*, we state $\underset{x\rightarrow a}{\text{lim}}f(x) = L.$ We can use the following strategy to confirm the result obtained from the table or as an alternative method for estimating a limit.
4.  Using a graphing calculator or computer software that allows us to graph functions, we can plot the function $f(x),$ making sure the functional values of $f(x)$ for *x*-values near *a* are in our window. We can use the trace feature to move along the graph of the function and watch the *y*-value readout as the *x*-values approach *a*. If the *y*-values approach *L* as our *x*-values approach *a* from both directions, then $\underset{x\rightarrow a}{\text{lim}}f(x) = L.$ We may need to zoom in on our graph and repeat this process several times.

We apply this Problem-Solving Strategy to compute a limit in Example 2.4.

### Example 2.4

#### Evaluating a Limit Using a Table of Functional Values 1

Evaluate $\underset{x\rightarrow 0}{\text{lim}}\frac{\text{sin} x}{x}$ using a table of functional values.

#### Solution

We have calculated the values of $f(x) = {{(\text{sin} x)}\text{/}x}$ for the values of *x* listed in Table 2.2.

| *x* | $\frac{\text{sin} x}{x}$ |  | *x* | $\frac{\text{sin} x}{x}$ |
|----|----|----|----|----|
| −0.1 | 0.998334166468 |  | 0.1 | 0.998334166468 |
| −0.01 | 0.999983333417 |  | 0.01 | 0.999983333417 |
| −0.001 | 0.999999833333 |  | 0.001 | 0.999999833333 |
| −0.0001 | 0.999999998333 |  | 0.0001 | 0.999999998333 |

Table 2.2 Table of Functional Values for $\underset{x\rightarrow 0}{\text{lim}}\frac{\text{sin} x}{x}$

*Note*: The values in this table were obtained using a calculator and using all the places given in the calculator output.

As we read down each $\frac{\left( {\text{sin} x} \right)}{x}$ column, we see that the values in each column appear to be approaching one. Thus, it is fairly reasonable to conclude that $\underset{x\rightarrow 0}{\text{lim}}\frac{\text{sin} x}{x} = 1.$ A calculator or computer-generated graph of $f(x) = \frac{\left( {\text{sin} x} \right)}{x}$ would be similar to that shown in Figure 2.13, and it confirms our estimate.

*Figure 2.13 The graph of f ( x ) = ( sin x ) / x f ( x ) = ( sin x ) / x confirms the estimate from Table 2.2 .*

### Example 2.5

#### Evaluating a Limit Using a Table of Functional Values 2

Evaluate $\underset{x\rightarrow 4}{\text{lim}}\frac{\sqrt{x} - 2}{x - 4}$ using a table of functional values.

#### Solution

As before, we use a table—in this case, Table 2.3—to list the values of the function for the given values of *x*.

| *x* | $\frac{\sqrt{x} - 2}{x - 4}$ |  | *x* | $\frac{\sqrt{x} - 2}{x - 4}$ |
|----|----|----|----|----|
| 3.9 | 0.251582341869 |  | 4.1 | 0.248456731317 |
| 3.99 | 0.25015644562 |  | 4.01 | 0.24984394501 |
| 3.999 | 0.250015627 |  | 4.001 | 0.249984377 |
| 3.9999 | 0.250001563 |  | 4.0001 | 0.249998438 |
| 3.99999 | 0.25000016 |  | 4.00001 | 0.24999984 |

Table 2.3 Table of Functional Values for $\underset{x\rightarrow 4}{\text{lim}}\frac{\sqrt{x} - 2}{x - 4}$

After inspecting this table, we see that the functional values less than 4 appear to be decreasing toward 0.25 whereas the functional values greater than 4 appear to be increasing toward 0.25. We conclude that $\underset{x\rightarrow 4}{\text{lim}}\frac{\sqrt{x} - 2}{x - 4} = 0.25.$ We confirm this estimate using the graph of $f(x) = \frac{\sqrt{x} - 2}{x - 4}$ shown in Figure 2.14.

*Figure 2.14 The graph of f ( x ) = x − 2 x − 4 f ( x ) = x − 2 x − 4 confirms the estimate from Table 2.3 .*

### Checkpoint 2.4

Estimate $\underset{x\rightarrow 1}{\text{lim}}\frac{\frac{1}{x} - 1}{x - 1}$ using a table of functional values. Use a graph to confirm your estimate.

At this point, we see from Example 2.4 and Example 2.5 that it may be just as easy, if not easier, to estimate a limit of a function by inspecting its graph as it is to estimate the limit by using a table of functional values. In Example 2.6, we evaluate a limit exclusively by looking at a graph rather than by using a table of functional values.

### Example 2.6

#### Evaluating a Limit Using a Graph

For $g(x)$ shown in Figure 2.15, evaluate $\underset{x\rightarrow-1}{\text{lim}}g(x).$

*Figure 2.15 The graph of g ( x ) g ( x ) includes one value not on a smooth curve.*

#### Solution

Despite the fact that $g(-1) = 4,$ as the *x*-values approach −1 from either side, the $g(x)$ values approach 3. Therefore, $\underset{x\rightarrow-1}{\text{lim}}g(x) = 3.$ Note that we can determine this limit without even knowing the algebraic expression of the function.

Based on Example 2.6, we make the following observation: It is possible for the limit of a function to exist at a point, and for the function to be defined at this point, but the limit of the function and the value of the function at the point may be different.

### Checkpoint 2.5

Use the graph of $h(x)$ in Figure 2.16 to evaluate $\underset{x\rightarrow 2}{\text{lim}}h(x),$ if possible.

*Figure 2.16*

Looking at a table of functional values or looking at the graph of a function provides us with useful insight into the value of the limit of a function at a given point. However, these techniques rely too much on guesswork. We eventually need to develop alternative methods of evaluating limits. These new methods are more algebraic in nature and we explore them in the next section; however, at this point we introduce two special limits that are foundational to the techniques to come.

### Theorem 2.1

#### Two Important Limits

Let *a* be a real number and *c* be a constant.

1.  
    ``` math
    \underset{x\rightarrow a}{\text{lim}}x = a
    ```
    (2.4)
2.  
    ``` math
    \underset{x\rightarrow a}{\text{lim}}c = c
    ```
    (2.5)

We can make the following observations about these two limits.

1.  For the first limit, observe that as *x* approaches *a*, so does $f(x),$ because $f(x) = x.$ Consequently, $\underset{x\rightarrow a}{\text{lim}}x = a.$
2.  For the second limit, consider Table 2.4.

| *x*            | $f(x) = c$ |     | *x*            | $f(x) = c$ |
|----------------|--------------|-----|----------------|--------------|
| $a - 0.1$    | *c*          |     | $a + 0.1$    | *c*          |
| $a - 0.01$   | *c*          |     | $a + 0.01$   | *c*          |
| $a - 0.001$  | *c*          |     | $a + 0.001$  | *c*          |
| $a - 0.0001$ | *c*          |     | $a + 0.0001$ | *c*          |

Table 2.4 Table of Functional Values for $\underset{x\rightarrow a}{\text{lim}}c = c$

Observe that for all values of *x* (regardless of whether they are approaching *a*), the values $f(x)$ remain constant at *c*. We have no choice but to conclude $\underset{x\rightarrow a}{\text{lim}}c = c.$

### The Existence of a Limit

As we consider the limit in the next example, keep in mind that for the limit of a function to exist at a point, the functional values must approach a single real-number value at that point. If the functional values do not approach a single value, then the limit does not exist.

### Example 2.7

#### Evaluating a Limit That Fails to Exist

Evaluate $\underset{x\rightarrow 0}{\text{lim}}\text{sin}\left( {1\text{/}\textit{x}} \right)$ using a table of values.

#### Solution

Table 2.5 lists values for the function $\text{sin}({1\text{/}{x)}}$ for the given values of *x*.

| *x* | $\text{sin}\left( \frac{1}{x} \right)$ |  | *x* | $\text{sin}\left( \frac{1}{x} \right)$ |
|----|----|----|----|----|
| −0.1 | 0.544021110889 |  | 0.1 | −0.544021110889 |
| −0.01 | 0.50636564111 |  | 0.01 | −0.50636564111 |
| −0.001 | −0.8268795405312 |  | 0.001 | 0.826879540532 |
| −0.0001 | 0.305614388888 |  | 0.0001 | −0.305614388888 |
| −0.00001 | −0.035748797987 |  | 0.00001 | 0.035748797987 |
| −0.000001 | 0.349993504187 |  | 0.000001 | −0.349993504187 |

Table 2.5 Table of Functional Values for $\underset{x\rightarrow 0}{\text{lim}}\text{sin}\left( \frac{1}{x} \right)$

After examining the table of functional values, we can see that the *y*-values do not seem to approach any one single value. It appears the limit does not exist. Before drawing this conclusion, let’s take a more systematic approach. Take the following sequence of *x*-values approaching 0:

$$
\frac{2}{\pi},\frac{2}{3\pi},\frac{2}{5\pi},\frac{2}{7\pi},\frac{2}{9\pi},\frac{2}{11\pi}\text{,….}
$$

The corresponding *y*-values are

$$
1,-1,1,-1,1,-1\text{,….}
$$

At this point we can indeed conclude that $\underset{x\rightarrow 0}{\text{lim}}\text{sin}\left( {1\text{/}\textit{x}} \right)$ does not exist. (Mathematicians frequently abbreviate “does not exist” as DNE. Thus, we would write $\underset{x\rightarrow 0}{\text{lim}}\text{sin}\left( {1\text{/}\textit{x}} \right)$ DNE.) The graph of $f(x) = \text{sin}\left( {1\text{/}x} \right)$ is shown in Figure 2.17 and it gives a clearer picture of the behavior of $\text{sin}({1\text{/}{x)}}$ as *x* approaches 0. You can see that $\text{sin}({1\text{/}{\textit{x})}}$ oscillates ever more wildly between −1 and 1 as *x* approaches 0.

*Figure 2.17 The graph of f ( x ) = sin ( 1 / x ) f ( x ) = sin ( 1 / x ) oscillates rapidly between −1 and 1 as x approaches 0.*

### Checkpoint 2.6

Use a table of functional values to evaluate $\underset{x\rightarrow 2}{\text{lim}}\frac{\left| {x^{2} - 4} \right|}{x - 2},$ if possible.

### One-Sided Limits

Sometimes indicating that the limit of a function fails to exist at a point does not provide us with enough information about the behavior of the function at that particular point. To see this, we now revisit the function $g(x) = {\left| {x - 2} \right|\text{/}\left( {x - 2} \right)}$ introduced at the beginning of the section (see Figure 2.12(b)). As we pick values of *x* close to 2, $g(x)$ does not approach a single value, so the limit as *x* approaches 2 does not exist—that is, $\underset{x\rightarrow 2}{\text{lim}}g(x)$ DNE. However, this statement alone does not give us a complete picture of the behavior of the function around the *x*-value 2. To provide a more accurate description, we introduce the idea of a one-sided limit. For all values to the left of 2 (or *the negative side of* 2), $g(x) = -1.$ Thus, as *x* approaches 2 from the left, $g(x)$ approaches −1. Mathematically, we say that the limit as *x* approaches 2 from the left is −1. Symbolically, we express this idea as

$$
\underset{x\rightarrow 2^{-}}{\text{lim}}g(x) = -1.
$$

Similarly, as *x* approaches 2 from the right (or *from the positive side*), $g(x)$ approaches 1. Symbolically, we express this idea as

$$
\underset{x\rightarrow 2^{+}}{\text{lim}}g(x) = 1.
$$

We can now present an informal definition of one-sided limits.

### Definition

We define two types of **one-sided limits**.

*Limit from the left:* Let $f(x)$ be a function defined at all values in an open interval of the form (*c*, *a*), and let *L* be a real number. If the values of the function $f(x)$ approach the real number *L* as the values of *x* (where $x < \textit{a}\text{)}$ approach the number *a*, then we say that *L* is the limit of $f(x)$ as *x* approaches a from the left. Symbolically, we express this idea as

$$
\underset{x\rightarrow a^{-}}{\text{lim}}f(x) = L.
$$

(2.6)

*Limit from the right:* Let $f(x)$ be a function defined at all values in an open interval of the form $\left( {a,c} \right),$ and let *L* be a real number. If the values of the function $f(x)$ approach the real number L as the values of *x* (where $x > \textit{a}\text{)}$ approach the number *a*, then we say that *L* is the limit of $f(x)$ as *x* approaches *a* from the right. Symbolically, we express this idea as

$$
\underset{x\rightarrow a^{+}}{\text{lim}}f(x) = L.
$$

(2.7)

### Example 2.8

#### Evaluating One-Sided Limits

For the function \$f(x) = \begin{cases}
{x + 1} & {\text{if}\ x < 2} \\
{x^{2} - 4} & {\text{if}\ x \geq 2}
\end{cases},\$ evaluate each of the following limits.

1.  $\underset{x\rightarrow 2^{-}}{\text{lim}}f(x)$
2.  $\underset{x\rightarrow 2^{+}}{\text{lim}}f(x)$

#### Solution

We can use tables of functional values again Table 2.6. Observe that for values of *x* less than 2, we use $f(x) = x + 1$ and for values of *x* greater than 2, we use $f(x) = x^{2} - 4.$

| *x*     | $f(x) = x + 1$ |     | *x*     | $f(x) = x^{2}-4$ |
|---------|------------------|-----|---------|--------------------|
| 1.9     | 2.9              |     | 2.1     | 0.41               |
| 1.99    | 2.99             |     | 2.01    | 0.0401             |
| 1.999   | 2.999            |     | 2.001   | 0.004001           |
| 1.9999  | 2.9999           |     | 2.0001  | 0.00040001         |
| 1.99999 | 2.99999          |     | 2.00001 | 0.0000400001       |

Table 2.6 Table of Functional Values for \$f(x) = \left\{ \begin{array}{l}
{x + 1\ \text{if}\ x < 2} \\
{x^{2} - 4\ \text{if}\ x \geq 2}
\end{array} \right.\$

Based on this table, we can conclude that a. $\underset{x\rightarrow 2^{-}}{\text{lim}}f(x) = 3$ and b. $\underset{x\rightarrow 2^{+}}{\text{lim}}f(x) = 0.$ Therefore, the (two-sided) limit of $f(x)$ does not exist at $x = 2.$ Figure 2.18 shows a graph of $f(x)$ and reinforces our conclusion about these limits.

*Figure 2.18 The graph of f ( x ) = { x + 1 if x \< 2 x 2 − 4 if x ≥ 2 f ( x ) = { x + 1 if x \< 2 x 2 − 4 if x ≥ 2 has a break at x = 2 . x = 2 .*

### Checkpoint 2.7

Use a table of functional values to estimate the following limits, if possible.

1.  $\underset{x\rightarrow 2^{-}}{\text{lim}}\frac{\left| {x^{2} - 4} \right|}{x - 2}$
2.  $\underset{x\rightarrow 2^{+}}{\text{lim}}\frac{\left| {x^{2} - 4} \right|}{x - 2}$

Let us now consider the relationship between the limit of a function at a point and the limits from the right and left at that point. It seems clear that if the limit from the right and the limit from the left have a common value, then that common value is the limit of the function at that point. Similarly, if the limit from the left and the limit from the right take on different values, the limit of the function does not exist. These conclusions are summarized in Relating One-Sided and Two-Sided Limits.

### Theorem 2.2

#### Relating One-Sided and Two-Sided Limits

Let $f(x)$ be a function defined at all values in an open interval containing *a*, with the possible exception of *a* itself, and let *L* be a real number. Then,

$$
\underset{x\rightarrow a}{\text{lim}}f(x) = L\ \text{if and only if}\ \underset{x\rightarrow a^{-}}{\text{lim}}f(x) = L\ \text{and}\ \underset{x\rightarrow a^{+}}{\text{lim}}f(x) = L.
$$

### Infinite Limits

Evaluating the limit of a function at a point or evaluating the limit of a function from the right and left at a point helps us to characterize the behavior of a function around a given value. As we shall see, we can also describe the behavior of functions that do not have finite limits.

We now turn our attention to $h(x) = {1\text{/}{{(x - 2)}^{2},}}$ the third and final function introduced at the beginning of this section (see Figure 2.12(c)). From its graph we see that as the values of *x* approach 2, the values of $h(x) = {1\text{/}{(x - 2)}^{2}}$ become larger and larger and, in fact, become infinite. Mathematically, we say that the limit of $h(x)$ as *x* approaches 2 is positive infinity. Symbolically, we express this idea as

$$
\underset{x\rightarrow 2}{\text{lim}}h(x) = \text{+}\infty.
$$

More generally, we define infinite limits as follows:

### Definition

We define three types of **infinite limits**.

*Infinite limits from the left:* Let $f(x)$ be a function defined at all values in an open interval of the form $\left( {b,a} \right).$

1.  If the values of $f(x)$ increase without bound as the values of *x* (where $x < \textit{a}\text{)}$ approach the number *a*, then we say that the limit as *x* approaches *a* from the left is positive infinity and we write  
    ``` math
    \underset{x\rightarrow a^{-}}{\text{lim}}f(x) = \text{+}\infty.
    ```
    (2.8)
2.  If the values of $f(x)$ decrease without bound as the values of *x* (where $x < \textit{a}\text{)}$ approach the number *a*, then we say that the limit as *x* approaches *a* from the left is negative infinity and we write  
    ``` math
    \underset{x\rightarrow a^{-}}{\text{lim}}f(x) = \text{−}\infty.
    ```
    (2.9)

*Infinite limits from the right*: Let $f(x)$ be a function defined at all values in an open interval of the form $\left( {a,c} \right).$

1.  If the values of $f(x)$ increase without bound as the values of *x* (where $x > \textit{a}\text{)}$ approach the number *a*, then we say that the limit as *x* approaches *a* from the right is positive infinity and we write  
    ``` math
    \underset{x\rightarrow a^{+}}{\text{lim}}f(x) = \text{+}\infty.
    ```
    (2.10)
2.  If the values of $f(x)$ decrease without bound as the values of *x* (where $x > \textit{a}\text{)}$ approach the number *a*, then we say that the limit as *x* approaches *a* from the right is negative infinity and we write  
    ``` math
    \underset{x\rightarrow a^{+}}{\text{lim}}f(x) = \text{−}\infty.
    ```
    (2.11)

*Two-sided infinite limit:* Let $f(x)$ be defined for all $x \neq a$ in an open interval containing *a*.

1.  If the values of $f(x)$ increase without bound as the values of *x* (where $x \neq \textit{a}\text{)}$ approach the number *a*, then we say that the limit as *x* approaches *a* is positive infinity and we write  
    ``` math
    \underset{x\rightarrow a}{\text{lim}}f(x) = \text{+}\infty.
    ```
    (2.12)
2.  If the values of $f(x)$ decrease without bound as the values of *x* (where $x \neq \textit{a}\text{)}$ approach the number *a*, then we say that the limit as *x* approaches *a* is negative infinity and we write  
    ``` math
    \underset{x\rightarrow a}{\text{lim}}f(x) = \text{−}\infty.
    ```
    (2.13)

It is important to understand that when we write statements such as $\underset{x\rightarrow a}{\text{lim}}f(x) = \text{+}\infty$ or $\underset{x\rightarrow a}{\text{lim}}f(x) = \text{-}\infty$ we are describing the behavior of the function, as we have just defined it. We are not asserting that a limit exists. For the limit of a function $f(x)$ to exist at *a*, it must approach a real number *L* as *x* approaches *a*. That said, if, for example, $\underset{x\rightarrow a}{\text{lim}}f(x) = \text{+}\infty,$ we always write $\underset{x\rightarrow a}{\text{lim}}f(x) = \text{+}\infty$ rather than $\underset{x\rightarrow a}{\text{lim}}f(x)$ DNE.

### Example 2.9

#### Recognizing an Infinite Limit

Evaluate each of the following limits, if possible. Use a table of functional values and graph $f(x) = {1\text{/}x}$ to confirm your conclusion.

1.  $\underset{x\rightarrow 0^{-}}{\text{lim}}\frac{1}{x}$
2.  $\underset{x\rightarrow 0^{+}}{\text{lim}}\frac{1}{x}$
3.  $\underset{x\rightarrow 0}{\text{lim}}\frac{1}{x}$

#### Solution

Begin by constructing a table of functional values.

| *x*       | $\frac{1}{x}$ |     | *x*      | $\frac{1}{x}$ |
|-----------|-----------------|-----|----------|-----------------|
| −0.1      | −10             |     | 0.1      | 10              |
| −0.01     | −100            |     | 0.01     | 100             |
| −0.001    | −1000           |     | 0.001    | 1000            |
| −0.0001   | −10,000         |     | 0.0001   | 10,000          |
| −0.00001  | −100,000        |     | 0.00001  | 100,000         |
| −0.000001 | −1,000,000      |     | 0.000001 | 1,000,000       |

Table 2.7 Table of Functional Values for $f(x) = \frac{1}{x}$

1.  The values of $1\text{/}x$ decrease without bound as *x* approaches 0 from the left. We conclude that  
    ``` math
    \underset{x\rightarrow 0^{-}}{\text{lim}}\frac{1}{x} = \text{−}\infty.
    ```
2.  The values of $1\text{/}x$ increase without bound as *x* approaches 0 from the right. We conclude that  
    ``` math
    \underset{x\rightarrow 0^{+}}{\text{lim}}\frac{1}{x} = \text{+}\infty.
    ```
3.  Since $\underset{x\rightarrow 0^{-}}{\text{lim}}\frac{1}{x} = \text{-}\infty$ and $\underset{x\rightarrow 0^{+}}{\text{lim}}\frac{1}{x} = \text{+}\infty$ have different values, we conclude that  
    ``` math
    \underset{x\rightarrow 0}{\text{lim}}\frac{1}{x}\ \text{DNE.}
    ```

The graph of $f(x) = {1\text{/}x}$ in Figure 2.19 confirms these conclusions.

*Figure 2.19 The graph of f ( x ) = 1 / x f ( x ) = 1 / x confirms that the limit as x approaches 0 does not exist.*

### Checkpoint 2.8

Evaluate each of the following limits, if possible. Use a table of functional values and graph $f(x) = {1\text{/}x^{2}}$ to confirm your conclusion.

1.  $\underset{x\rightarrow 0^{-}}{\text{lim}}\frac{1}{x^{2}}$
2.  $\underset{x\rightarrow 0^{+}}{\text{lim}}\frac{1}{x^{2}}$
3.  $\underset{x\rightarrow 0}{\text{lim}}\frac{1}{x^{2}}$

It is useful to point out that functions of the form $f(x) = {1\text{/}\left( {x - a} \right)^{n}},$ where *n* is a positive integer, have infinite limits as *x* approaches *a* from either the left or right (Figure 2.20). These limits are summarized in Infinite Limits from Positive Integers.

*Figure 2.20 The function f ( x ) = 1 / ( x − a ) n f ( x ) = 1 / ( x − a ) n has infinite limits at a .*

### Theorem 2.3

#### Infinite Limits from Positive Integers

If *n* is a positive even integer, then

$$
\underset{x\rightarrow a}{\text{lim}}\frac{1}{\left( {x - a} \right)^{n}} = \text{+}\infty.
$$

If *n* is a positive odd integer, then

$$
\underset{x\rightarrow a^{+}}{\text{lim}}\frac{1}{\left( {x - a} \right)^{n}} = \text{+}\infty
$$

and

$$
\underset{x\rightarrow a^{-}}{\text{lim}}\frac{1}{\left( {x - a} \right)^{n}} = \text{−}\infty.
$$

We should also point out that in the graphs of $f(x) = {1\text{/}{{(x - a)}^{n},}}$ points on the graph having *x*-coordinates very near to *a* are very close to the vertical line $x = a.$ That is, as *x* approaches *a*, the points on the graph of $f(x)$ are closer to the line $x = a.$ The line $x = a$ is called a vertical asymptote of the graph. We formally define a vertical asymptote as follows:

### Definition

Let $f(x)$ be a function. If any of the following conditions hold, then the line $x = a$ is a **vertical asymptote** of $f(x).$

$$
\begin{array}{clc}
{\underset{x\rightarrow a^{-}}{\text{lim}}f(x)} & = & {\text{+}\infty\ \text{or}\ \text{−∞}} \\
{\underset{x\rightarrow a^{+}}{\text{lim}}f(x)} & = & {\text{+}\infty\ \text{or}\ \text{−∞}} \\
 & {\ \text{or}} & \\
{\underset{x\rightarrow a}{\text{lim}}f(x)} & = & {\text{+}\infty\ \text{or}\ \text{−∞}}
\end{array}
$$

### Example 2.10

#### Finding a Vertical Asymptote

Evaluate each of the following limits using Infinite Limits from Positive Integers. Identify any vertical asymptotes of the function $f(x) = {1\text{/}\left( {x + 3} \right)^{4}}.$

1.  $\underset{x\rightarrow-3^{-}}{\text{lim}}\frac{1}{\left( {x + 3} \right)^{4}}$
2.  $\underset{x\rightarrow-3^{+}}{\text{lim}}\frac{1}{\left( {x + 3} \right)^{4}}$
3.  $\underset{x\rightarrow-3}{\text{lim}}\frac{1}{\left( {x + 3} \right)^{4}}$

#### Solution

We can use Infinite Limits from Positive Integers directly.

1.  $\underset{x\rightarrow-3^{-}}{\text{lim}}\frac{1}{\left( {x + 3} \right)^{4}} = \text{+}\infty$
2.  $\underset{x\rightarrow-3^{+}}{\text{lim}}\frac{1}{\left( {x + 3} \right)^{4}} = \text{+}\infty$
3.  $\underset{x\rightarrow-3}{\text{lim}}\frac{1}{\left( {x + 3} \right)^{4}} = \text{+}\infty$

The function $f(x) = {1\text{/}\left( {x + 3} \right)^{4}}$ has a vertical asymptote of $x = -3.$

### Checkpoint 2.9

Evaluate each of the following limits. Identify any vertical asymptotes of the function $f(x) = \frac{1}{\left( {x - 2} \right)^{3}}.$

1.  $\underset{x\rightarrow 2^{-}}{\text{lim}}\frac{1}{\left( {x - 2} \right)^{3}}$
2.  $\underset{x\rightarrow 2^{+}}{\text{lim}}\frac{1}{\left( {x - 2} \right)^{3}}$
3.  $\underset{x\rightarrow 2}{\text{lim}}\frac{1}{\left( {x - 2} \right)^{3}}$

In the next example we put our knowledge of various types of limits to use to analyze the behavior of a function at several different points.

### Example 2.11

#### Behavior of a Function at Different Points

Use the graph of $f(x)$ in Figure 2.21 to determine each of the following values:

1.  $\underset{x\rightarrow-4^{-}}{\text{lim}}f(x);\underset{x\rightarrow-4^{+}}{\text{lim}}f(x);\underset{x\rightarrow-4}{\text{lim}}f(x);f(-4)$
2.  $\underset{x\rightarrow-2^{-}}{\text{lim}}f(x);\underset{x\rightarrow-2^{+}}{\text{lim}}f(x);\underset{x\rightarrow-2}{\text{lim}}f(x);f(-2)$
3.  $\underset{x\rightarrow 1^{-}}{\text{lim}}f(x);\underset{x\rightarrow 1^{+}}{\text{lim}}f(x);\underset{x\rightarrow 1}{\text{lim}}f(x);f(1)$
4.  $\underset{x\rightarrow 3^{-}}{\text{lim}}f(x);\underset{x\rightarrow 3^{+}}{\text{lim}}f(x);\underset{x\rightarrow 3}{\text{lim}}f(x);f(3)$

*Figure 2.21 The graph shows f ( x ) . f ( x ) .*

#### Solution

Using Infinite Limits from Positive Integers and the graph for reference, we arrive at the following values:

1.  $\underset{x\rightarrow-4^{-}}{\text{lim}}f(x) = 0;\underset{x\rightarrow-4^{+}}{\text{lim}}f(x) = 0;\underset{x\rightarrow-4}{\text{lim}}f(x) = 0;f(-4) = 0$
2.  $\underset{x\rightarrow-2^{-}}{\text{lim}}f(x) = 3.;\underset{x\rightarrow-2^{+}}{\text{lim}}f(x) = 3;\underset{x\rightarrow-2}{\text{lim}}f(x) = 3;f(-2)$ is undefined
3.  $\underset{x\rightarrow 1^{-}}{\text{lim}}f(x) = 6;\underset{x\rightarrow 1^{+}}{\text{lim}}f(x) = 3;\underset{x\rightarrow 1}{\text{lim}}f(x)$ DNE; $f(1) = 6$
4.  $\underset{x\rightarrow 3^{-}}{\text{lim}}f(x) = \text{-}\infty;\underset{x\rightarrow 3^{+}}{\text{lim}}f(x) = \text{-}\infty;\underset{x\rightarrow 3}{\text{lim}}f(x) = \text{-}\infty;f(3)$ is undefined

### Checkpoint 2.10

Evaluate $\underset{x\rightarrow 1}{\text{lim}}f(x)$ for $f(x)$ shown here:

### Example 2.12

#### Chapter Opener: Einstein’s Equation

*Figure 2.22 (credit: NASA)*

In the chapter opener we mentioned briefly how Albert Einstein showed that a limit exists to how fast any object can travel. Given Einstein’s equation for the mass of a moving object, what is the value of this bound?

#### Solution

Our starting point is Einstein’s equation for the mass of a moving object,

$$
m = \frac{m_{0}}{\sqrt{1 - \frac{v^{2}}{c^{2}}}},
$$

where $m_{0}$ is the object’s mass at rest, *v* is its speed, and *c* is the speed of light. To see how the mass changes at high speeds, we can graph the ratio of masses $m\text{/}m_{0}$ as a function of the ratio of speeds, $v\text{/}c$ (Figure 2.23).

*Figure 2.23 This graph shows the ratio of masses as a function of the ratio of speeds in Einstein’s equation for the mass of a moving object.*

We can see that as the ratio of speeds approaches 1—that is, as the speed of the object approaches the speed of light—the ratio of masses increases without bound. In other words, the function has a vertical asymptote at ${v\text{/}c} = 1.$ We can try a few values of this ratio to test this idea.

| $\frac{v}{c}$ | $\sqrt{1 - \frac{v^{2}}{c^{2}}}$ | $\frac{m}{m_{0}}$ |
|-----------------|------------------------------------|---------------------|
| 0.99            | 0.1411                             | 7.089               |
| 0.999           | 0.0447                             | 22.37               |
| 0.9999          | 0.0141                             | 70.71               |

Table 2.8 Ratio of Masses and Speeds for a Moving Object

Thus, according to Table 2.8, if an object with mass 100 kg is traveling at 0.9999*c*, its mass becomes 7071 kg. Since no object can have an infinite mass, we conclude that no object can travel at or more than the speed of light.

### Section 2.2 Exercises

For the following exercises, consider the function $f(x) = \frac{x^{2} - 1}{\left| {x - 1} \right|}.$

30\.

**\[T\]** Complete the following table for the function. Round your solutions to four decimal places.

| *x*    | $f(x)$ |     | *x*    | $f(x)$ |
|--------|----------|-----|--------|----------|
| 0.9    | a\.      |     | 1.1    | e\.      |
| 0.99   | b\.      |     | 1.01   | f\.      |
| 0.999  | c\.      |     | 1.001  | g\.      |
| 0.9999 | d\.      |     | 1.0001 | h\.      |

31\.

What do your results in the preceding exercise indicate about the two-sided limit $\underset{x\rightarrow 1}{\text{lim}}f(x)?$ Explain your response.

For the following exercises, consider the function $f(x) = \left( {1 + x} \right)^{1\text{/}x}.$

32\.

**\[T\]** Make a table showing the values of *f* for $x = -0.01,-0.001,-0.0001,-0.00001$ and for $x = 0.01,0.001,0.0001,0.00001.$ Round your solutions to five decimal places.

| *x*      | $f(x)$ |     | *x*     | $f(x)$ |
|----------|----------|-----|---------|----------|
| −0.01    | a\.      |     | 0.01    | e\.      |
| −0.001   | b\.      |     | 0.001   | f\.      |
| −0.0001  | c\.      |     | 0.0001  | g\.      |
| −0.00001 | d\.      |     | 0.00001 | h\.      |

33\.

What does the table of values in the preceding exercise indicate about the function $f(x) = \left( {1 + x} \right)^{1\text{/}x}?$

34\.

To which mathematical constant does the limit in the preceding exercise appear to be getting closer?

In the following exercises, use the given values to set up a table to evaluate the limits. Round your solutions to eight decimal places.

35\.

**\[T\]** $\underset{x\rightarrow 0}{\text{lim}}\frac{\text{sin} 2x}{x};\ \pm 0.1,\pm 0.01,\pm 0.001,\pm.0001$

| *x* | $\frac{\text{sin} 2x}{x}$ |  | *x* | $\frac{\text{sin} 2x}{x}$ |
|----|----|----|----|----|
| −0.1 | a\. |  | 0.1 | e\. |
| −0.01 | b\. |  | 0.01 | f\. |
| −0.001 | c\. |  | 0.001 | g\. |
| −0.0001 | d\. |  | 0.0001 | h\. |

36\.

**\[T\]** $\underset{x\rightarrow 0}{\text{lim}}\frac{\text{sin} 3x}{x}$ ±0.1, ±0.01, ±0.001, ±0.0001

| *X* | $\frac{\text{sin} 3x}{x}$ |  | *x* | $\frac{\text{sin} 3x}{x}$ |
|----|----|----|----|----|
| −0.1 | a\. |  | 0.1 | e\. |
| −0.01 | b\. |  | 0.01 | f\. |
| −0.001 | c\. |  | 0.001 | g\. |
| −0.0001 | d\. |  | 0.0001 | h\. |

37\.

Use the preceding two exercises to conjecture (guess) the value of the following limit: $\underset{x\rightarrow 0}{\text{lim}}\frac{\text{sin} ax}{x}$ for *a*, a positive real value.

**\[T\]** In the following exercises, set up a table of values to find the indicated limit. Round to eight digits.

38\.

$\underset{x\rightarrow 2}{\text{lim}}\frac{x^{2} - 4}{x^{2} + x - 6}$

| *x* | $\frac{x^{2} - 4}{x^{2} + x - 6}$ |  | *x* | $\frac{x^{2} - 4}{x^{2} + x - 6}$ |
|----|----|----|----|----|
| 1.9 | a\. |  | 2.1 | e\. |
| 1.99 | b\. |  | 2.01 | f\. |
| 1.999 | c\. |  | 2.001 | g\. |
| 1.9999 | d\. |  | 2.0001 | h\. |

39\.

$\underset{x\rightarrow 1}{\text{lim}}\left( {1 - 2x} \right)$

| *x*    | $1 - 2x$ |     | *x*    | $1 - 2x$ |
|--------|------------|-----|--------|------------|
| 0.9    | a\.        |     | 1.1    | e\.        |
| 0.99   | b\.        |     | 1.01   | f\.        |
| 0.999  | c\.        |     | 1.001  | g\.        |
| 0.9999 | d\.        |     | 1.0001 | h\.        |

40\.

$\underset{x\rightarrow 0}{\text{lim}}\frac{5}{1 - e^{1\text{/}x}}$

| *x* | $\frac{5}{1 - e^{1\text{/}x}}$ |  | *x* | $\frac{5}{1 - e^{1\text{/}x}}$ |
|----|----|----|----|----|
| −0.1 | a\. |  | 0.1 | e\. |
| −0.01 | b\. |  | 0.01 | f\. |
| −0.001 | c\. |  | 0.001 | g\. |
| −0.0001 | d\. |  | 0.0001 | h\. |

41\.

$\underset{z\rightarrow 0}{\text{lim}}\frac{z - 1}{z^{2}\left( {z + 3} \right)}$

| *z* | $\frac{z - 1}{z^{2}\left( {z + 3} \right)}$ |  | *z* | $\frac{z - 1}{z^{2}\left( {z + 3} \right)}$ |
|----|----|----|----|----|
| −0.1 | a\. |  | 0.1 | e\. |
| −0.01 | b\. |  | 0.01 | f\. |
| −0.001 | c\. |  | 0.001 | g\. |
| −0.0001 | d\. |  | 0.0001 | h\. |

42\.

$\underset{t\rightarrow 0^{+}}{\text{lim}}\frac{\text{cos} t}{t}$

| *t*    | $\frac{\text{cos} t}{t}$ |
|--------|----------------------------------------|
| 0.1    | a\.                                    |
| 0.01   | b\.                                    |
| 0.001  | c\.                                    |
| 0.0001 | d\.                                    |

43\.

$\underset{x\rightarrow 2}{\text{lim}}\frac{1 - \frac{2}{x}}{x^{2} - 4}$

| *x* | $\frac{1 - \frac{2}{x}}{x^{2} - 4}$ |  | *x* | $\frac{1 - \frac{2}{x}}{x^{2} - 4}$ |
|----|----|----|----|----|
| 1.9 | a\. |  | 2.1 | e\. |
| 1.99 | b\. |  | 2.01 | f\. |
| 1.999 | c\. |  | 2.001 | g\. |
| 1.9999 | d\. |  | 2.0001 | h\. |

**\[T\]** In the following exercises, set up a table of values and round to eight significant digits. Based on the table of values, make a guess about what the limit is. Then, use a calculator to graph the function and determine the limit. Was the conjecture correct? If not, why does the method of tables fail?

44\.

$\underset{\theta\rightarrow 0}{\text{lim}}\text{sin}\left( \frac{\pi}{\theta} \right)$

| *θ* | $\text{sin}\left( \frac{\pi}{\theta} \right)$ |  | *θ* | $\text{sin}\left( \frac{\pi}{\theta} \right)$ |
|----|----|----|----|----|
| −0.1 | a\. |  | 0.1 | e\. |
| −0.01 | b\. |  | 0.01 | f\. |
| −0.001 | c\. |  | 0.001 | g\. |
| −0.0001 | d\. |  | 0.0001 | h\. |

45\.

$\underset{\alpha\rightarrow 0^{+}}{\text{lim}}\frac{1}{\alpha}\text{cos}\left( \frac{\pi}{\alpha} \right)$

| $\alpha$ | $\frac{1}{\alpha}\text{cos}\left( \frac{\pi}{\alpha} \right)$ |
|----|----|
| 0.1 | a\. |
| 0.01 | b\. |
| 0.001 | c\. |
| 0.0001 | d\. |

In the following exercises, consider the graph of the function $y = f(x)$ shown here. Which of the statements about $y = f(x)$ are true and which are false? Explain why a statement is false.

46\.

$\underset{x\rightarrow 10}{\text{lim}}f(x) = 0$

47\.

$\underset{x\rightarrow-2^{+}}{\text{lim}}f(x) = 3$

48\.

$\underset{x\rightarrow-8}{\text{lim}}f(x) = f(-8)$

49\.

$\underset{x\rightarrow 6}{\text{lim}}f(x) = 5$

In the following exercises, use the following graph of the function $y = f(x)$ to find the values, if possible. Estimate when necessary.

50\.

$\underset{x\rightarrow 1^{-}}{\text{lim}}f(x)$

51\.

$\underset{x\rightarrow 1^{+}}{\text{lim}}f(x)$

52\.

$\underset{x\rightarrow 1}{\text{lim}}f(x)$

53\.

$\underset{x\rightarrow 2}{\text{lim}}f(x)$

54\.

$f(1)$

In the following exercises, use the graph of the function $y = f(x)$ shown here to find the values, if possible. Estimate when necessary.

55\.

$\underset{x\rightarrow 0^{-}}{\text{lim}}f(x)$

56\.

$\underset{x\rightarrow 0^{+}}{\text{lim}}f(x)$

57\.

$\underset{x\rightarrow 0}{\text{lim}}f(x)$

58\.

$\underset{x\rightarrow 2}{\text{lim}}f(x)$

In the following exercises, use the graph of the function $y = f(x)$ shown here to find the values, if possible. Estimate when necessary.

59\.

$\underset{x\rightarrow-2^{-}}{\text{lim}}f(x)$

60\.

$\underset{x\rightarrow-2^{+}}{\text{lim}}f(x)$

61\.

$\underset{x\rightarrow-2}{\text{lim}}f(x)$

62\.

$\underset{x\rightarrow 2^{-}}{\text{lim}}f(x)$

63\.

$\underset{x\rightarrow 2^{+}}{\text{lim}}f(x)$

64\.

$\underset{x\rightarrow 2}{\text{lim}}f(x)$

In the following exercises, use the graph of the function $y = g(x)$ shown here to find the values, if possible. Estimate when necessary.

65\.

$\underset{x\rightarrow 0^{-}}{\text{lim}}g(x)$

66\.

$\underset{x\rightarrow 0^{+}}{\text{lim}}g(x)$

67\.

$\underset{x\rightarrow 0}{\text{lim}}g(x)$

In the following exercises, use the graph of the function $y = h(x)$ shown here to find the values, if possible. Estimate when necessary.

68\.

$\underset{x\rightarrow 0^{-}}{\text{lim}}h(x)$

69\.

$\underset{x\rightarrow 0^{+}}{\text{lim}}h(x)$

70\.

$\underset{x\rightarrow 0}{\text{lim}}h(x)$

In the following exercises, use the graph of the function $y = f(x)$ shown here to find the values, if possible. Estimate when necessary.

71\.

$\underset{x\rightarrow 0^{-}}{\text{lim}}f(x)$

72\.

$\underset{x\rightarrow 0^{+}}{\text{lim}}f(x)$

73\.

$\underset{x\rightarrow 0}{\text{lim}}f(x)$

74\.

$\underset{x\rightarrow 1}{\text{lim}}f(x)$

75\.

$\underset{x\rightarrow 2}{\text{lim}}f(x)$

In the following exercises, sketch the graph of a function with the given properties.

76\.

$\underset{x\rightarrow 2}{\text{lim}}f(x) = 1,\underset{x\rightarrow 4^{-}}{\text{lim}}f(x) = 3,\underset{x\rightarrow 4^{+}}{\text{lim}}f(x) = 6,f(4)$ is not defined.

77\.

$As~x\rightarrow - \infty~,~f(x)\rightarrow 0,\underset{x\rightarrow-1^{-}}{\text{lim}}f(x) = \text{-}\infty,$ $\underset{x\rightarrow-1^{+}}{\text{lim}}f(x) = \infty,\underset{x\rightarrow 0}{\text{lim}}f(x) = f(0),f(0) = 1,~As~x\rightarrow\infty,~f(x)\rightarrow\text{-}\infty$

78\.

$As~x\rightarrow - \infty,~f(x)\rightarrow 2,\underset{x\rightarrow 3^{-}}{\text{lim}}f(x) = \text{-}\infty,$ $\underset{x\rightarrow 3^{+}}{\text{lim}}f(x) = \infty,~As~x\rightarrow\infty,~f(x)\rightarrow 2,f(0) = \frac{-1}{3}$

79\.

$As~x\rightarrow - \infty,~f(x)\rightarrow 2,\underset{x\rightarrow-2}{\text{lim}}f(x) = \text{-}\infty,$ $As~x\rightarrow\infty,~f(x)\rightarrow 2,f(0) = 0$

80\.

$As~x\rightarrow - \infty,~f(x)\rightarrow 0,\underset{x\rightarrow-1^{-}}{\text{lim}}f(x) = \infty,\underset{x\rightarrow-1^{+}}{\text{lim}}f(x) = \text{-}\infty,$ $f(0) = -1,\underset{x\rightarrow 1^{-}}{\text{lim}}f(x) = \text{-}\infty,\underset{x\rightarrow 1^{+}}{\text{lim}}f(x) = \infty,~As~x\rightarrow\infty,~f(x)\rightarrow 0$

81\.

Shock waves arise in many physical applications, ranging from supernovas to detonation waves. A graph of the density of a shock wave with respect to distance, *x*, is shown here. We are mainly interested in the location of the front of the shock, labeled $x_{\text{SF}}$ in the diagram.

1.  Evaluate $\underset{x\rightarrow x_{SF}{}^{+}}{\text{lim}}\rho(x).$
2.  Evaluate $\underset{x\rightarrow x_{SF}{}^{-}}{\text{lim}}\rho(x).$
3.  Evaluate $\underset{x\rightarrow x_{SF}}{\text{lim}}\rho(x).$ Explain the physical meanings behind your answers.

82\.

A track coach uses a camera with a fast shutter to estimate the position of a runner with respect to time. A table of the values of position of the athlete versus time is given here, where *x* is the position in meters of the runner and *t* is time in seconds. What is $\underset{t\rightarrow 2}{\text{lim}}x(t)?$ What does it mean physically?

| *t* **(sec)** | *x* **(m)** |
|---------------|-------------|
| 1.75          | 4.5         |
| 1.95          | 6.1         |
| 1.99          | 6.42        |
| 2.01          | 6.58        |
| 2.05          | 6.9         |
| 2.25          | 8.5         |
