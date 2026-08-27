> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 6.3 Volumes of Revolution: Cylindrical Shells

### Learning Objectives

- 6.3.1 Calculate the volume of a solid of revolution by using the method of cylindrical shells.
- 6.3.2 Compare the different methods for calculating a volume of revolution.

In this section, we examine the method of cylindrical shells, the final method for finding the volume of a solid of revolution. We can use this method on the same kinds of solids as the disk method or the washer method; however, with the disk and washer methods, we integrate along the coordinate axis parallel to the axis of revolution. With the method of cylindrical shells, we integrate along the coordinate axis *perpendicular* to the axis of revolution. The ability to choose which variable of integration we want to use can be a significant advantage with more complicated functions. Also, the specific geometry of the solid sometimes makes the method of using cylindrical shells more appealing than using the washer method. In the last part of this section, we review all the methods for finding volume that we have studied and lay out some guidelines to help you determine which method to use in a given situation.

### The Method of Cylindrical Shells

Again, we are working with a solid of revolution. As before, we define a region $R,$ bounded above by the graph of a function $y = f(x),$ below by the $x\text{-axis,}$ and on the left and right by the lines $x = a$ and $x = b,$ respectively, as shown in Figure 6.25(a). We then revolve this region around the *y*-axis, as shown in Figure 6.25(b). Note that this is different from what we have done before. Previously, regions defined in terms of functions of $x$ were revolved around the $x\text{-axis}$ or a line parallel to it.

*Figure 6.25 (a) A region bounded by the graph of a function of x . x . (b) The solid of revolution formed when the region is revolved around the y -axis . y -axis .*

As we have done many times before, partition the interval $\left\lbrack {a,b} \right\rbrack$ using a regular partition, $P = \left\{ {x_{0},x_{1}\text{,…},x_{n}} \right\}$ and, for $i = 1,2\text{,…},n,$ choose a point $x_{i}^{*} \in \left\lbrack {x_{i - 1},x_{i}} \right\rbrack.$ Then, construct a rectangle over the interval $\left\lbrack {x_{i - 1},x_{i}} \right\rbrack$ of height $f(x_{i}^{*})$ and width $\text{Δ}x.$ A representative rectangle is shown in Figure 6.26(a). When that rectangle is revolved around the *y*-axis, instead of a disk or a washer, we get a cylindrical shell, as shown in the following figure.

*Figure 6.26 (a) A representative rectangle. (b) When this rectangle is revolved around the y -axis , y -axis , the result is a cylindrical shell. (c) When we put all the shells together, we get an approximation of the original solid.*

To calculate the volume of this shell, consider Figure 6.27.

*Figure 6.27 Calculating the volume of the shell.*

The shell is a cylinder, so its volume is the cross-sectional area multiplied by the height of the cylinder. The cross-sections are annuli (ring-shaped regions—essentially, circles with a hole in the center), with outer radius $x_{i}$ and inner radius $x_{i - 1}.$ Thus, the cross-sectional area is $\pi x_{i}^{2} - \pi x_{i - 1}^{2}.$ The height of the cylinder is $f(x_{i}^{*}).$ Then the volume of the shell is

$$
\begin{array}{cl}
V_{\text{shell}} & {= f(x_{i}^{*})(\pi x_{i}^{2} - \pi x_{i - 1}^{2})} \\
 & {= \pi f(x_{i}^{*})\left( {x_{i}^{2} - x_{i - 1}^{2}} \right)} \\
 & {= \pi f(x_{i}^{*})\left( {x_{i} + x_{i - 1}} \right)\left( {x_{i} - x_{i - 1}} \right)} \\
 & {= 2\pi f(x_{i}^{*})\left( \frac{x_{i} + x_{i - 1}}{2} \right)\left( {x_{i} - x_{i - 1}} \right).}
\end{array}
$$

Note that $x_{i} - x_{i - 1} = \text{Δ}x,$ so we have

$$
V_{\text{shell}} = 2\pi f(x_{i}^{*})\left( \frac{x_{i} + x_{i - 1}}{2} \right)\text{Δ}x.
$$

Furthermore, $\frac{x_{i} + x_{i - 1}}{2}$ is both the midpoint of the interval $\left\lbrack {x_{i - 1},x_{i}} \right\rbrack$ and the average radius of the shell, and we can approximate this by $x_{i}^{*}.$ We then have

$$
V_{\text{shell}} \approx 2\pi f(x_{i}^{*})x_{i}^{*}\text{Δ}x.
$$

Another way to think of this is to think of making a vertical cut in the shell and then opening it up to form a flat plate (Figure 6.28).

*Figure 6.28 (a) Make a vertical cut in a representative shell. (b) Open the shell up to form a flat plate.*

In reality, the outer radius of the shell is greater than the inner radius, and hence the back edge of the plate would be slightly longer than the front edge of the plate. However, we can approximate the flattened shell by a flat plate of height $f(x_{i}^{*}),$ width $2\pi x_{i}^{*},$ and thickness $\text{Δ}x$ (Figure 6.28). The volume of the shell, then, is approximately the volume of the flat plate. Multiplying the height, width, and depth of the plate, we get

$$
V_{\text{shell}} \approx f(x_{i}^{*})\left( {2\pi x_{i}^{*}} \right)\text{Δ}x,
$$

which is the same formula we had before.

To calculate the volume of the entire solid, we then add the volumes of all the shells and obtain

$$
V \approx \sum\limits_{i = 1}^{n}\left( {2\pi x_{i}^{*}f(x_{i}^{*})\text{Δ}x} \right).
$$

Here we have another Riemann sum, this time for the function $2\pi xf(x).$ Taking the limit as $n\rightarrow\infty$ gives us

$$
V = \underset{n\rightarrow\infty}{\text{lim}}\sum\limits_{i = 1}^{n}\left( {2\pi x_{i}^{*}f(x_{i}^{*})\text{Δ}x} \right) = {\int_{a}^{b}\left( {2\pi xf(x)} \right)}dx.
$$

This leads to the following rule for the method of cylindrical shells.

### Rule: The Method of Cylindrical Shells

Let $f(x)$ be continuous and nonnegative. Define $R$ as the region bounded above by the graph of $f(x),$ below by the $x\text{-axis},$ on the left by the line $x = a,$ and on the right by the line $x = b.$ Then the volume of the solid of revolution formed by revolving $R$ around the *y*-axis is given by

$$
V = {\int_{a}^{b}\left( {2\pi xf(x)} \right)}dx.
$$

(6.6)

Now let’s consider an example.

### Example 6.12

#### The Method of Cylindrical Shells 1

Define $R$ as the region bounded above by the graph of $f(x) = {1\text{/}x}$ and below by the $x\text{-axis}$ over the interval $\left\lbrack {1,3} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the $y\text{-axis}.$

#### Solution

First we must graph the region $R$ and the associated solid of revolution, as shown in the following figure.

*Figure 6.29 (a) The region R R under the graph of f ( x ) = 1 / x f ( x ) = 1 / x over the interval \[ 1 , 3 \] . \[ 1 , 3 \] . (b) The solid of revolution generated by revolving R R about the y -axis . y -axis .*

Then the volume of the solid is given by

$$
\begin{array}{cl}
V & {= {\int_{a}^{b}\left( {2\pi xf(x)} \right)}dx} \\
 & {= {\int_{1}^{3}{\left( {2\pi x\left( \frac{1}{x} \right)} \right)dx}}} \\
 & {= {\int_{1}^{3}2}\pi\ dx = \left. {2\pi x} \right|_{1}^{3} = 4\pi\ \text{units}^{3}\text{.}}
\end{array}
$$

### Checkpoint 6.12

Define *R* as the region bounded above by the graph of $f(x) = x^{2}$ and below by the *x*-axis over the interval $\left\lbrack {1,2} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the $y\text{-axis}.$

### Example 6.13

#### The Method of Cylindrical Shells 2

Define *R* as the region bounded above by the graph of $f(x) = 2x - x^{2}$ and below by the $x\text{-axis}$ over the interval $\left\lbrack {0,2} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the $y\text{-axis}.$

#### Solution

First graph the region $R$ and the associated solid of revolution, as shown in the following figure.

*Figure 6.30 (a) The region R R under the graph of f ( x ) = 2 x − x 2 f ( x ) = 2 x − x 2 over the interval \[ 0 , 2 \] . \[ 0 , 2 \] . (b) The volume of revolution obtained by revolving R R about the y -axis . y -axis .*

Then the volume of the solid is given by

$$
\begin{array}{cl}
V & {= {\int_{a}^{b}\left( {2\pi xf(x)} \right)}dx} \\
 & {= {\int_{0}^{2}\left( {2\pi x\left( {2x - x^{2}} \right)} \right)}dx = 2\pi{\int_{0}^{2}\left( {2x^{2} - x^{3}} \right)}dx} \\
 & {= \left. {2\pi\left\lbrack {\frac{2x^{3}}{3} - \frac{x^{4}}{4}} \right\rbrack}\  \right|_{0}^{2} = \frac{8\pi}{3}\ \text{units}^{3}\text{.}}
\end{array}
$$

### Checkpoint 6.13

Define $R$ as the region bounded above by the graph of $f(x) = 3x - x^{2}$ and below by the $x\text{-axis}$ over the interval $\left\lbrack {0,2} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the $y\text{-axis}.$

As with the disk method and the washer method, we can use the method of cylindrical shells with solids of revolution, revolved around the $x\text{-axis},$ when we want to integrate with respect to $y.$ The analogous rule for this type of solid is given here.

### Rule: The Method of Cylindrical Shells for Solids of Revolution around the *x*-axis

Let $g(y)$ be continuous and nonnegative. Define $Q$ as the region bounded on the right by the graph of $g(y),$ on the left by the $y\text{-axis},$ below by the line $y = c,$ and above by the line $y = d.$ Then, the volume of the solid of revolution formed by revolving $Q$ around the $x\text{-axis}$ is given by

$$
V = {\int_{c}^{d}\left( {2\pi yg(y)} \right)}dy.
$$

### Example 6.14

#### The Method of Cylindrical Shells for a Solid Revolved around the *x*-axis

Define $Q$ as the region bounded on the right by the graph of $g(y) = 2\sqrt{y}$ and on the left by the $y\text{-axis}$ for $y \in \left\lbrack {0,4} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $Q$ around the *x*-axis.

#### Solution

First, we need to graph the region $Q$ and the associated solid of revolution, as shown in the following figure.

*Figure 6.31 (a) The region Q Q to the left of the function g ( y ) g ( y ) over the interval \[ 0 , 4 \] . \[ 0 , 4 \] . (b) The solid of revolution generated by revolving Q Q around the x -axis . x -axis .*

Label the shaded region $Q.$ Then the volume of the solid is given by

$$
\begin{array}{cl}
V & {= {\int_{c}^{d}\left( {2\pi yg(y)} \right)}dy} \\
 & {= {\int_{0}^{4}\left( {2\pi y\left( {2\sqrt{y}} \right)} \right)}dy = 4\pi{\int_{0}^{4}y^{3\text{/}2}}dy} \\
 & {= {\left. {4\pi\left\lbrack \frac{2y^{5\text{/}2}}{5} \right.} \right\rbrack\left. \  \right|}_{0}^{4} = \frac{256\pi}{5}\ \text{units}^{3}\text{.}}
\end{array}
$$

### Checkpoint 6.14

Define $Q$ as the region bounded on the right by the graph of $g(y) = {3\text{/}y}$ and on the left by the $y\text{-axis}$ for $y \in \left\lbrack {1,3} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $Q$ around the $x\text{-axis}.$

For the next example, we look at a solid of revolution for which the graph of a function is revolved around a line other than one of the two coordinate axes. To set this up, we need to revisit the development of the method of cylindrical shells. Recall that we found the volume of one of the shells to be given by

$$
\begin{array}{cl}
V_{\text{shell}} & {= f(x_{i}^{*})(\pi x_{i}^{2} - \pi x_{i - 1}^{2})} \\
 & {= \pi f(x_{i}^{*})\left( {x_{i}^{2} - x_{i - 1}^{2}} \right)} \\
 & {= \pi f(x_{i}^{*})\left( {x_{i} + x_{i - 1}} \right)\left( {x_{i} - x_{i - 1}} \right)} \\
 & {= 2\pi f(x_{i}^{*})\left( \frac{x_{i} + x_{i - 1}}{2} \right)\left( {x_{i} - x_{i - 1}} \right).}
\end{array}
$$

This was based on a shell with an outer radius of $x_{i}$ and an inner radius of $x_{i - 1}.$ If, however, we rotate the region around a line other than the $y\text{-axis},$ we have a different outer and inner radius. Suppose, for example, that we rotate the region around the line $x = \text{-}k,$ where $k$ is some positive constant. Then, the outer radius of the shell is $x_{i} + k$ and the inner radius of the shell is $x_{i - 1} + k.$ Substituting these terms into the expression for volume, we see that when a plane region is rotated around the line $x = \text{-}k,$ the volume of a shell is given by

$$
\begin{array}{cl}
V_{\text{shell}} & {= 2\pi f(x_{i}^{*})\left( \frac{\left( {x_{i} + k} \right) + \left( {x_{i - 1} + k} \right)}{2} \right)\left( {\left( {x_{i} + k} \right) - \left( {x_{i - 1} + k} \right)} \right)} \\
 & {= 2\pi f(x_{i}^{*})\left( {\left( \frac{x_{i} + x_{i - 1}}{2} \right) + k} \right)\text{Δ}x.}
\end{array}
$$

As before, we notice that $\frac{x_{i} + x_{i - 1}}{2}$ is the midpoint of the interval $\left\lbrack {x_{i - 1},x_{i}} \right\rbrack$ and can be approximated by $x_{i}^{*}.$ Then, the approximate volume of the shell is

$$
V_{\text{shell}} \approx 2\pi\left( {x_{i}^{*} + k} \right)f(x_{i}^{*})\text{Δ}x.
$$

The remainder of the development proceeds as before, and we see that

$$
V = {\int_{a}^{b}\left( {2\pi\left( {x + k} \right)f(x)} \right)}dx.
$$

We could also rotate the region around other horizontal or vertical lines, such as a vertical line in the right half plane. In each case, the volume formula must be adjusted accordingly. Specifically, the $x\text{-term}$ in the integral must be replaced with an expression representing the radius of a shell. To see how this works, consider the following example.

### Example 6.15

#### A Region of Revolution Revolved around a Line

Define $R$ as the region bounded above by the graph of $f(x) = x$ and below by the $x\text{-axis}$ over the interval $\left\lbrack {1,2} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the line $x = -1.$

#### Solution

First, graph the region $R$ and the associated solid of revolution, as shown in the following figure.

*Figure 6.32 (a) The region R R between the graph of f ( x ) f ( x ) and the x -axis x -axis over the interval \[ 1 , 2 \] . \[ 1 , 2 \] . (b) The solid of revolution generated by revolving R R around the line x = −1 . x = −1 .*

Note that the radius of a shell is given by $x + 1.$ Then the volume of the solid is given by

$$
\begin{array}{cl}
V & {= {\int_{1}^{2}\left( {2\pi\left( {x + 1} \right)f(x)} \right)}dx} \\
 & {= {\int_{1}^{2}\left( {2\pi\left( {x + 1} \right)x} \right)}dx = 2\pi{\int_{1}^{2}\left( {x^{2} + x} \right)}dx} \\
 & {= \left. {2\pi\left\lbrack {\frac{x^{3}}{3} + \frac{x^{2}}{2}} \right\rbrack}\  \right|_{1}^{2} = \frac{23\pi}{3}\ \text{units}^{3}\text{.}}
\end{array}
$$

### Checkpoint 6.15

Define $R$ as the region bounded above by the graph of $f(x) = x^{2}$ and below by the $x\text{-axis}$ over the interval $\left\lbrack {0,1} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the line $x = -2.$

For our final example in this section, let’s look at the volume of a solid of revolution for which the region of revolution is bounded by the graphs of two functions.

### Example 6.16

#### A Region of Revolution Bounded by the Graphs of Two Functions

Define $R$ as the region bounded above by the graph of the function $f(x) = \sqrt{x}$ and below by the graph of the function $g(x) = {1\text{/}x}$ over the interval $\left\lbrack {1,4} \right\rbrack.$ Find the volume of the solid of revolution generated by revolving $R$ around the $y\text{-axis}.$

#### Solution

First, graph the region $R$ and the associated solid of revolution, as shown in the following figure.

*Figure 6.33 (a) The region R R between the graph of f ( x ) f ( x ) and the graph of g ( x ) g ( x ) over the interval \[ 1 , 4 \] . \[ 1 , 4 \] . (b) The solid of revolution generated by revolving R R around the y -axis . y -axis .*

Note that the axis of revolution is the $y\text{-axis},$ so the radius of a shell is given simply by $x.$ We don’t need to make any adjustments to the *x*-term of our integrand. The height of a shell, though, is given by $f(x) - g(x),$ so in this case we need to adjust the $f(x)$ term of the integrand. Then the volume of the solid is given by

$$
\begin{array}{cl}
V & {= {\int_{1}^{4}\left( {2\pi x\left( {f(x) - g(x)} \right)} \right)}dx} \\
 & {= {\int_{1}^{4}{\left( {2\pi x\left( {\sqrt{x} - \frac{1}{x}} \right)} \right)dx}} = 2\pi{\int_{1}^{4}\left( {x^{3\text{/}2} - 1} \right)}dx} \\
 & {= \left. {2\pi\left\lbrack {\frac{2x^{5\text{/}2}}{5} - x} \right\rbrack}\  \right|_{1}^{4} = \frac{94\pi}{5}\ \text{units}^{3}.}
\end{array}
$$

### Checkpoint 6.16

Define $R$ as the region bounded above by the graph of $f(x) = x$ and below by the graph of $g(x) = x^{2}$ over the interval $\left\lbrack {0,1} \right\rbrack.$ Find the volume of the solid of revolution formed by revolving $R$ around the $y\text{-axis}.$

### Which Method Should We Use?

We have studied several methods for finding the volume of a solid of revolution, but how do we know which method to use? It often comes down to a choice of which integral is easiest to evaluate. Figure 6.34 describes the different approaches for solids of revolution around the $x\text{-axis}.$ It’s up to you to develop the analogous table for solids of revolution around the $y\text{-axis}.$

*Figure 6.34*

Let’s take a look at a couple of additional problems and decide on the best approach to take for solving them.

### Example 6.17

#### Selecting the Best Method

For each of the following problems, select the best method to find the volume of a solid of revolution generated by revolving the given region around the $x\text{-axis},$ and set up the integral to find the volume (do not evaluate the integral).

1.  The region bounded by the graphs of $y = x,$ $y = 2 - x,$ and the $x\text{-axis}.$
2.  The region bounded by the graphs of $y = 4x - x^{2}$ and the $x\text{-axis}.$

#### Solution

1.  First, sketch the region and the solid of revolution as shown.  

    *Figure 6.35 (a) The region R R bounded by two lines and the x -axis . x -axis . (b) The solid of revolution generated by revolving R R about the x -axis . x -axis .*

      
    Looking at the region, if we want to integrate with respect to $x,$ we would have to break the integral into two pieces, because we have different functions bounding the region over $\left\lbrack {0,1} \right\rbrack$ and $\left\lbrack {1,2} \right\rbrack.$ In this case, using the disk method, we would have  
    ``` math
    V = {\int_{0}^{1}\left( {\pi x^{2}} \right)}dx + {\int_{1}^{2}\left( {\pi{(2 - x)}^{2}} \right)}dx.
    ```
      
    If we used the shell method instead, we would use functions of $y$ to represent the curves, producing  
    ``` math
    \begin{array}{cl}
    V & {= {\int_{0}^{1}\left( {2\pi y\left\lbrack {\left( {2 - y} \right) - y} \right\rbrack} \right)}dy} \\
     & {= {\int_{0}^{1}\left( {2\pi y\left\lbrack {2 - 2y} \right\rbrack} \right)}dy.}
    \end{array}
    ```
      
    Neither of these integrals is particularly onerous, but since the shell method requires only one integral, and the integrand requires less simplification, we should probably go with the shell method in this case.

2.  First, sketch the region and the solid of revolution as shown.  

    *Figure 6.36 (a) The region R R between the curve and the x -axis . x -axis . (b) The solid of revolution generated by revolving R R about the x -axis . x -axis .*

      
    Looking at the region, it would be problematic to define a horizontal rectangle; the region is bounded on the left and right by the same function. Therefore, we can dismiss the method of shells. The solid has no cavity in the middle, so we can use the method of disks. Then  
    ``` math
    V = {\int_{0}^{4}\pi}\left( {4x - x^{2}} \right)^{2}dx.
    ```

### Checkpoint 6.17

Select the best method to find the volume of a solid of revolution generated by revolving the given region around the $x\text{-axis},$ and set up the integral to find the volume (do not evaluate the integral): the region bounded by the graphs of $y = 2 - x^{2}$ and $y = x^{2}.$

### Section 6.3 Exercises

For the following exercises, find the volume generated when the region between the two curves is rotated around the given axis. Use both the shell method and the washer method. Use technology to graph the functions and draw a typical slice by hand.

114\.

**\[T\]** Bounded by the curves $y = 3x,x = 0,$ and $y = 3$ rotated around the $y\text{-axis}.$

115\.

**\[T\]** Bounded by the curves $y = 3x,y = 0,\ \text{and}\ x = 3$ rotated around the $y\text{-axis}.$

116\.

**\[T\]** Bounded by the curves $y = 3x,y = 0,\ \text{and}\ y = 3$ rotated around the $x\text{-axis}.$

117\.

**\[T\]** Bounded by the curves $y = 3x,y = 0,\ \text{and}\ x = 3$ rotated around the $x\text{-axis}.$

118\.

**\[T\]** Bounded by the curves $y = 2x^{3},y = 0,\ \text{and}\ x = 2$ rotated around the $y\text{-axis}.$

119\.

**\[T\]** Bounded by the curves $y = 2x^{3},y = 0,\ \text{and}\ x = 2$ rotated around the $x\text{-axis}.$

For the following exercises, use shells to find the volumes of the given solids. Note that the rotated regions lie between the curve and the $x\text{-axis}$ and are rotated around the $y\text{-axis}.$

120\.

$y = 1 - x^{2},x = 0,\ \text{and}\ x = 1$

121\.

$y = 5x^{3},x = 0,\ \text{and}\ x = 1$

122\.

$y = \frac{1}{x},x = 1,\ \text{and}\ x = 100$

123\.

$y = \sqrt{1 - x^{2}},x = 0,\ \text{and}\ x = 1$

124\.

$y = \frac{1}{1 + x^{2}},x = 0,\ \text{and}\ x = 3$

125\.

$y = \text{sin}x^{2},x = 0,\ \text{and}\ x = \sqrt{\pi}$

126\.

$y = \frac{1}{\sqrt{1 - x^{2}}},x = 0,\ \text{and}\ x = \frac{1}{2}$

127\.

$y = \sqrt{x},x = 0,\ \text{and}\ x = 1$

128\.

$y = \left( {1 + x^{2}} \right)^{3},x = 0,\ \text{and}\ x = 1$

129\.

$y = 5x^{3} - 2x^{4},x = 0,\ \text{and}\ x = 2$

For the following exercises, use shells to find the volume generated by rotating the regions between the given curve and $y = 0$ around the $x\text{-axis}.$

130\.

$y = \sqrt{1 - x^{2}},x = 0,\ x = 1$ and the *x*-axis

131\.

$y = x^{2},x = 0,\ x = 2$ and the *x*-axis

132\.

$y = \frac{x^{3}}{2},\ x = 0,\ x = 2,$ and the *x*-axis

133\.

$y = \frac{2}{x^{2}},\ x = 1,\ x = 2,$ and the *x*-axis

134\.

$x = \frac{1}{1 + y^{2}},y = 4$

135\.

$x = \frac{1 + y^{2}}{y},y = 1,\ y = 4,$ and the *y*-axis

136\.

$x = \sqrt{4 - y^{2}}\text{,}x = 0\text{,}y = 0$

137\.

$x = y^{3} - 2y^{2},\ x = 0,\ x = 9$

138\.

$x = \sqrt{y} + 1,\ x = 1,\ x = 3,$ and the *x*-axis

139\.

$x = \sqrt[3]{27y}\text{and}\ x = \frac{3y}{4}$

For the following exercises, find the volume generated when the region between the curves is rotated around the given axis.

140\.

$y = 3 - x,y = 0,x = 0,\ \text{and}\ x = 2$ rotated around the $y\text{-axis}.$

141\.

$y = x^{3},x = 0,\ \text{and}\ y = 8$ rotated around the $y\text{-axis}.$

142\.

$y = x^{2},y = x,$ rotated around the $y\text{-axis}.$

143\.

$y = \sqrt{x},y = 0,\ \text{and}\ x = 1$ rotated around the line $x = 2.$

144\.

$y = \frac{1}{4 - x},x = 1,\ x = 2\ \text{and}\ y = 0$ rotated around the line $x = 4.$

145\.

$y = \sqrt{x}\ \text{and}\ y = x^{2}$ rotated around the $y\text{-axis}.$

146\.

$y = \sqrt{x}\ \text{and}\ y = x^{2}$ rotated around the line $x = 2.$

147\.

$x = y^{3},x = \frac{1}{y},x = 1,\ \text{and}\ x = 2$ rotated around the $x\text{-axis}.$

148\.

$x = y^{2}\ \text{and}\ y = x$ rotated around the line $y = 2.$

149\.

**\[T\]** Left of $x = \text{sin}\left( {\pi y} \right),$ right of $y = x,$ around the $y\text{-axis}.$

For the following exercises, use technology to graph the region. Determine which method you think would be easiest to use to calculate the volume generated when the function is rotated around the specified axis. Then, use your chosen method to find the volume.

150\.

**\[T\]** $y = x^{2}$ and $y = 4x$ rotated around the $y\text{-axis}.$

151\.

**\[T\]** $y = \text{cos}\left( {\pi x} \right),y = \text{sin}\left( {\pi x} \right),x = \frac{1}{4},\ \text{and}\ x = \frac{5}{4}$ rotated around the $y\text{-axis}.$ This exercise requires advanced technique. You may use technology to perform the integration.

152\.

**\[T\]** $y = x^{2} - 2x,x = 2,\ \text{and}\ x = 4$ rotated around the $y\text{-axis}.$

153\.

**\[T\]** $y = x^{2} - 2x,x = 2,\ \text{and}\ x = 4$ rotated around the $x\text{-axis}.$

154\.

**\[T\]** $y = 3x^{3} - 2,y = x,\ \text{and}\ x = 2$ rotated around the $x\text{-axis}.$

155\.

**\[T\]** $y = 3x^{3} - 2,y = x,\ \text{and}\ x = 2$ rotated around the $y\text{-axis}.$

156\.

**\[T\]** $x = \text{sin}\left( {\pi y^{2}} \right)$ and $x = \sqrt{2}y$ rotated around the $x\text{-axis}.$

157\.

**\[T\]** $x = y^{2},x = y^{2} - 2y + 1,\ \text{and}\ x = 2$ rotated around the $y\text{-axis}.$

For the following exercises, use the method of shells to approximate the volumes of some common objects, which are pictured in accompanying figures.

158\.

Use the method of shells to find the volume of a sphere of radius $r.$

159\.

Use the method of shells to find the volume of a cone with radius $r$ and height $h.$

160\.

Use the method of shells to find the volume of an ellipsoid $\left( {x^{2}\text{/}a^{2}} \right) + \left( {y^{2}\text{/}b^{2}} \right) = 1$ rotated around the $x\text{-axis}.$

161\.

Use the method of shells to find the volume of a cylinder with radius $r$ and height $h.$

162\.

Use the method of shells to find the volume of the donut created when the circle $x^{2} + y^{2} = 4$ is rotated around the line $x = 4.$

163\.

Consider the region enclosed by the graphs of $y = f(x),y = 1 + f(x),x = 0,y = 0,$ and $x = a > 0.$ What is the volume of the solid generated when this region is rotated around the $y\text{-axis}?$ Assume that the function is defined over the interval $\lbrack 0,a\rbrack.$

164\.

Consider the function $y = f(x),$ which decreases from $f(0) = b$ to $f(1) = 0.$ Set up the integrals for determining the volume, using both the shell method and the disk method, of the solid generated when this region, with $x = 0$ and $y = 0,$ is rotated around the $y\text{-axis}.$ Prove that both methods approximate the same volume. Which method is easier to apply? (*Hint:* Since $f(x)$ is one-to-one, there exists an inverse $f^{-1}(y).)$
