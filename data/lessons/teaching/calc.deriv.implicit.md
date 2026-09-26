> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 3.8 Implicit Differentiation

### Learning Objectives

- 3.8.1 Find the derivative of a complicated function by using implicit differentiation.
- 3.8.2 Use implicit differentiation to determine the equation of a tangent line.

We have already studied how to find equations of tangent lines to functions and the rate of change of a function at a specific point. In all these cases we had the explicit equation for the function and differentiated these functions explicitly. Suppose instead that we want to determine the equation of a tangent line to an arbitrary curve or the rate of change of an arbitrary curve at a point. In this section, we solve these problems by finding the derivatives of functions that define $y$ implicitly in terms of $x.$

### Implicit Differentiation

In most discussions of math, if the dependent variable $y$ is a function of the independent variable $x,$ we express *y* in terms of $x.$ If this is the case, we say that $y$ is an explicit function of $x.$ For example, when we write the equation $y = x^{2} + 1,$ we are defining *y* explicitly in terms of $x.$ On the other hand, if the relationship between the function $y$ and the variable $x$ is expressed by an equation where $y$ is not expressed entirely in terms of $x,$ we say that the equation defines *y* implicitly in terms of $x.$ For example, the equation $y - x^{2} = 1$ defines the function $y = x^{2} + 1$ implicitly.

Implicit differentiation allows us to find slopes of tangents to curves that are clearly not functions (they fail the vertical line test). We are using the idea that portions of $y$ are functions that satisfy the given equation, but that $y$ is not actually a function of $x.$

In general, an equation defines a function implicitly if the function satisfies that equation. An equation may define many different functions implicitly. For example, the functions

${y = \sqrt{25 - x^{2}}}{,~y = - \sqrt{25 - x^{2}},}$ and \$y = \left\{ \begin{matrix}
{\sqrt{25 - x^{2}}\ \text{if} - 5 < x < 0} \\
{\text{−}\sqrt{25 - x^{2}}\ \text{if}\ 0 < x < 25}
\end{matrix} \right.,$ which are illustrated in Figure 3.30, are just three of the many functions defined implicitly by the equation $x^{2} + y^{2} = 25.\$

*Figure 3.30 The equation x 2 + y 2 = 25 x 2 + y 2 = 25 defines many functions implicitly.*

If we want to find the slope of the line tangent to the graph of $x^{2} + y^{2} = 25$ at the point $\left( {3,4} \right),$ we could evaluate the derivative of the function $y = \sqrt{25 - x^{2}}$ at $x = 3.$ On the other hand, if we want the slope of the tangent line at the point $\left( {3,-4} \right),$ we could use the derivative of $y = \text{-}\sqrt{25 - x^{2}}.$ However, it is not always easy to solve for a function defined implicitly by an equation. Fortunately, the technique of implicit differentiation allows us to find the derivative of an implicitly defined function without ever solving for the function explicitly. The process of finding $\frac{dy}{dx}$ using implicit differentiation is described in the following problem-solving strategy.

### Problem-Solving Strategy

#### Implicit Differentiation

To perform implicit differentiation on an equation that defines a function $y$ implicitly in terms of a variable $x,$ use the following steps:

1.  Take the derivative of both sides of the equation. Keep in mind that *y* is a function of *x*. Consequently, whereas $\frac{d}{dx}\left( {\text{sin} x} \right) = \text{cos} x,\frac{d}{dx}(\text{sin} y) = \text{cos} y\frac{dy}{dx}$ because we must use the chain rule to differentiate $\text{sin} y$ with respect to $x.$
2.  Rewrite the equation so that all terms containing $\frac{dy}{dx}$ are on the left and all terms that do not contain $\frac{dy}{dx}$ are on the right.
3.  Factor out $\frac{dy}{dx}$ on the left.
4.  Solve for $\frac{dy}{dx}$ by dividing both sides of the equation by an appropriate algebraic expression.

### Example 3.68

#### Using Implicit Differentiation

Assuming that $y$ is defined implicitly by the equation $x^{2} + y^{2} = 25,$ find $\frac{dy}{dx}.$

#### Solution

Follow the steps in the problem-solving strategy.

$$
\begin{array}{rllccl}
{\frac{d}{dx}\left( {x^{2} + y^{2}} \right)} & = & {\frac{d}{dx}(25)} & & & \text{Step 1. Differentiate both sides of the equation.} \\
{\frac{d}{dx}\left( x^{2} \right) + \frac{d}{dx}\left( y^{2} \right)} & = & 0 & & & \begin{array}{l}
\text{Step 1.1. Use the sum rule on the left.} \\
{\text{On the right}\ \frac{d}{dx}(25) = 0.}
\end{array} \\
{2x + 2y\frac{dy}{dx}} & = & 0 & & & \begin{array}{l}
{\text{Step 1.2. Take the derivatives, so}\ \frac{d}{dx}\left( x^{2} \right) = 2x} \\
{\text{and}\ \frac{d}{dx}\left( y^{2} \right) = 2y\frac{dy}{dx}.}
\end{array} \\
{2y\frac{dy}{dx}} & = & {-2x} & & & \begin{array}{l}
{\text{Step 2. Keep the terms with}\ \frac{dy}{dx}\ \text{on the left.}} \\
\text{Move the remaining terms to the right.}
\end{array} \\
\frac{dy}{dx} & = & {- \frac{x}{y}} & & & \begin{array}{l}
\text{Step 4. Divide both sides of the equation by} \\
{2y.\ \text{(Step 3 does not apply in this case.)}}
\end{array}
\end{array}
$$

#### Analysis

Note that the resulting expression for $\frac{dy}{dx}$ is in terms of both the independent variable $x$ and the dependent variable $y.$ Although in some cases it may be possible to express $\frac{dy}{dx}$ in terms of $x$ only, it is generally not possible to do so.

### Example 3.69

#### Using Implicit Differentiation and the Product Rule

Assuming that $y$ is defined implicitly by the equation $x^{3}\text{sin} y + y = 4x + 3,$ find $\frac{dy}{dx}.$

#### Solution

$$
\begin{array}{rllccl}
{\frac{d}{dx}\left( {x^{3}\text{sin} y + y} \right)} & = & {\frac{d}{dx}\left( {4x + 3} \right)} & & & \text{Step 1: Differentiate both sides of the equation.} \\
{\frac{d}{dx}\left( {x^{3}\text{sin} y} \right) + \frac{d}{dx}(y)} & = & 4 & & & \begin{array}{l}
\text{Step 1.1: Apply the sum rule on the left.} \\
{\text{On the right,}\ \frac{d}{dx}\left( {4x + 3} \right) = 4.}
\end{array} \\
{\left( {\frac{d}{dx}\left( x^{3} \right) \cdot \text{sin} y + \frac{d}{dx}\left( {\text{sin} y} \right) \cdot x^{3}} \right) + \frac{dy}{dx}} & = & 4 & & & \begin{array}{l}
\text{Step 1.2: Use the product rule to find} \\
{\frac{d}{dx}\left( {x^{3}\text{sin} y} \right).\ \text{Observe that}\ \frac{d}{dx}(y) = \frac{dy}{dx}.}
\end{array} \\
{3x^{2}\text{sin} y + \left( {\text{cos} y\frac{dy}{dx}} \right) \cdot x^{3} + \frac{dy}{dx}} & = & 4 & & & \begin{array}{l}
{\text{Step 1.3: We know}\ \frac{d}{dx}\left( x^{3} \right) = 3x^{2}.\ \text{Use the}} \\
{\text{chain rule to obtain}\ \frac{d}{dx}\left( {\text{sin} y} \right) = \text{cos} y\frac{dy}{dx}.}
\end{array} \\
{\text{x}^{3}\text{cos} y\frac{dy}{dx} + \frac{dy}{dx}} & = & {4 - 3x^{2}\text{sin} y} & & & \begin{array}{l}
{\text{Step 2: Keep all terms containing}\ \frac{dy}{dx}\ \text{on the}} \\
\text{left. Move all other terms to the right.}
\end{array} \\
{\frac{dy}{dx}\left( {\text{x}^{3}\text{cos} y + 1} \right)} & = & {4 - 3x^{2}\text{sin} y} & & & {\text{Step 3: Factor out}\ \frac{dy}{dx}\ \text{on the left.}} \\
\frac{dy}{dx} & = & \frac{4 - 3x^{2}\text{sin} y}{x^{3}\text{cos} y + 1} & & & \begin{array}{l}
{\text{Step 4: Solve for}\ \frac{dy}{dx}\ \text{by dividing both sides of}} \\
{\text{the equation by}\ \text{x}^{3}\text{cos} y + 1.}
\end{array}
\end{array}
$$

### Example 3.70

#### Using Implicit Differentiation to Find a Second Derivative

Find $\frac{d^{2}y}{dx^{2}}$ if $x^{2} + y^{2} = 25.$

#### Solution

In Example 3.68, we showed that $\frac{dy}{dx} = - \frac{x}{y}.$ We can take the derivative of both sides of this equation to find $\frac{d^{2}y}{dx^{2}}.$

$$
\begin{array}{clccl}
\frac{d^{2}y}{dx^{2}} & {= \frac{d}{dx}\left( {- \frac{x}{y}} \right)} & & & {\text{Differentiate both sides of}\ \frac{dy}{dx} = - \frac{x}{y}.} \\
 & {= - \frac{\left( {1 \cdot y - x\frac{dy}{dx}} \right)}{y^{2}}} & & & {\text{Use the quotient rule to find}\ \frac{d}{dy}\left( {- \frac{x}{y}} \right).} \\
 & {= \frac{\text{−}y + x\frac{dy}{dx}}{y^{2}}} & & & \text{Simplify.} \\
 & {= \frac{\text{−}y + x\left( {- \frac{x}{y}} \right)}{y^{2}}} & & & {\text{Substitute}\ \frac{dy}{dx} = - \frac{x}{y}.} \\
 & {= \frac{\text{−}y^{2} - x^{2}}{y^{3}}} & & & \text{Simplify.}
\end{array}
$$

At this point we have found an expression for $\frac{d^{2}y}{dx^{2}}.$ If we choose, we can simplify the expression further by recalling that $x^{2} + y^{2} = 25$ and making this substitution in the numerator to obtain $\frac{d^{2}y}{dx^{2}} = - \frac{25}{y^{3}}.$

### Checkpoint 3.48

Find $\frac{dy}{dx}$ for $y$ defined implicitly by the equation $4x^{5} + \text{tan} y = y^{2} + 5x.$

### Finding Tangent Lines Implicitly

Now that we have seen the technique of implicit differentiation, we can apply it to the problem of finding equations of tangent lines to curves described by equations.

### Example 3.71

#### Finding a Tangent Line to a Circle

Find an equation of the line tangent to the curve $x^{2} + y^{2} = 25$ at the point $\left( {3,-4} \right).$

#### Solution

Although we could find this equation without using implicit differentiation, using that method makes it much easier. In Example 3.68, we found $\frac{dy}{dx} = - \frac{x}{y}.$

The slope of the tangent line is found by substituting $\left( {3,-4} \right)$ into this expression. Consequently, the slope of the tangent line is
$$\frac{dy}{dx}\left| \begin{array}{l}
\\
{}_{({3,-4})}
\end{array} \right. = - \frac{3}{-4} = \frac{3}{4}.$$

Using the point $\left( {3,-4} \right)$ and the slope $\frac{3}{4}$ in the point-slope equation of the line, we obtain the equation $y = \frac{3}{4}x - \frac{25}{4}$ (Figure 3.31).

*Figure 3.31 The line y = 3 4 x − 25 4 y = 3 4 x − 25 4 is tangent to x 2 + y 2 = 25 x 2 + y 2 = 25 at the point (3, −4).*

### Example 3.72

#### Finding the Equation of the Tangent Line to a Curve

Find an equation of the line tangent to the graph of $y^{3} + x^{3} - 3xy = 0$ at the point $\left( {\frac{3}{2},\frac{3}{2}} \right)$ (Figure 3.32). This curve is known as the folium (or leaf) of Descartes.

*Figure 3.32 Finding the tangent line to the folium of Descartes at ( 3 2 , 3 2 ) . ( 3 2 , 3 2 ) .*

#### Solution

Begin by finding $\frac{dy}{dx}.$

$$
\begin{array}{rll}
{\frac{d}{dx}\left( {y^{3} + x^{3} - 3xy} \right)} & = & {\frac{d}{dx}(0)} \\
{3y^{2}\frac{dy}{dx} + 3x^{2} - \left( {3y + \frac{dy}{dx}3x} \right)} & = & 0 \\
\frac{dy}{dx} & = & {\frac{3y - 3x^{2}}{3y^{2} - 3x}.}
\end{array}
$$

Next, substitute $\left( {\frac{3}{2},\frac{3}{2}} \right)$ into $\frac{dy}{dx} = \frac{3y - 3x^{2}}{3y^{2} - 3x}$ to find the slope of the tangent line:

$$
\frac{dy}{dx}\left| \begin{array}{l}
 \\
{}_{({\frac{3}{2},\frac{3}{2}})}
\end{array} \right. = -1.
$$

Finally, substitute into the point-slope equation of the line to obtain

$$
y = \text{−}x + 3.
$$

### Example 3.73

#### Applying Implicit Differentiation

In a simple video game, a rocket travels in an elliptical orbit whose path is described by the equation $4x^{2} + 25y^{2} = 100.$ The rocket can fire missiles along lines tangent to its path. The object of the game is to destroy an incoming asteroid traveling along the positive *x*-axis toward $\left( {0,0} \right).$ If the rocket fires a missile when it is located at $\left( {3,\frac{8}{5}} \right),$ where will it intersect the *x*-axis?

#### Solution

To solve this problem, we must determine where the line tangent to the graph of

$4x^{2} + 25y^{2} = 100$ at $\left( {3,\frac{8}{5}} \right)$ intersects the *x*-axis. Begin by finding $\frac{dy}{dx}$ implicitly.

Differentiating, we have

$$
8x + 50y\frac{dy}{dx} = 0.
$$

Solving for $\frac{dy}{dx},$ we have

$$
\frac{dy}{dx} = - \frac{4x}{25y}.
$$

The slope of the tangent line is $\frac{dy}{dx}\left| {}_{({3,\frac{8}{5}})} \right. = - \frac{3}{10}.$ The equation of the tangent line is $y = - \frac{3}{10}x + \frac{5}{2}.$ To determine where the line intersects the *x*-axis, solve $0 = - \frac{3}{10}x + \frac{5}{2}.$ The solution is $x = \frac{25}{3}.$ The missile intersects the *x*-axis at the point $\left( {\frac{25}{3},0} \right).$

### Checkpoint 3.49

Find an equation of the line tangent to the hyperbola $x^{2} - y^{2} = 16$ at the point $\left( {5,3} \right).$
