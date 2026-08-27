> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 3.5 Derivatives of Trigonometric Functions

### Learning Objectives

- 3.5.1 Find the derivatives of the sine and cosine function.
- 3.5.2 Find the derivatives of the standard trigonometric functions.
- 3.5.3 Calculate the higher-order derivatives of the sine and cosine.

One of the most important types of motion in physics is simple harmonic motion, which is associated with such systems as an object with mass oscillating on a spring. Simple harmonic motion can be described by using either sine or cosine functions. In this section we expand our knowledge of derivative formulas to include derivatives of these and other trigonometric functions. We begin with the derivatives of the sine and cosine functions and then use them to obtain formulas for the derivatives of the remaining four trigonometric functions. Being able to calculate the derivatives of the sine and cosine functions will enable us to find the velocity and acceleration of simple harmonic motion.

### Derivatives of the Sine and Cosine Functions

We begin our exploration of the derivative for the sine function by using the formula to make a reasonable guess at its derivative. Recall that for a function $f(x),$

$$
f'(x) = \underset{h\rightarrow 0}{\text{lim}}\frac{f\left( {x + h} \right) - f(x)}{h}.
$$

Consequently, for values of $h$ very close to 0, $f'(x) \approx \frac{f\left( {x + h} \right) - f(x)}{h}.$ We see that by using $h = 0.01,$

$$
\frac{d}{dx}(\text{sin} x) \approx \frac{\text{sin}\left( {x + 0.01} \right) - \text{sin} x}{0.01}
$$

By setting $D(x) = \frac{\text{sin}\left( {x + 0.01} \right) - \text{sin} x}{0.01}$ and using a graphing utility, we can get a graph of an approximation to the derivative of $\text{sin} x$ (Figure 3.25).

*Figure 3.25 The graph of the function D ( x ) D ( x ) looks a lot like a cosine curve.*

Upon inspection, the graph of $D(x)$ appears to be very close to the graph of the cosine function. Indeed, we will show that

$$
\frac{d}{dx}\left( {\text{sin} x} \right) = \text{cos} x.
$$

If we were to follow the same steps to approximate the derivative of the cosine function, we would find that

$$
\frac{d}{dx}(\text{cos} x) = \text{−}\text{sin}{x.}
$$

### Theorem 3.8

#### The Derivatives of sin *x* and cos *x*

The derivative of the sine function is the cosine and the derivative of the cosine function is the negative sine.

$$
\frac{d}{dx}\left( {\text{sin} x} \right) = \text{cos} x
$$

(3.11)
$$
\frac{d}{dx}\left( {\text{cos} x} \right) = \text{−}\text{sin} x
$$

(3.12)

#### Proof

Because the proofs for $\frac{d}{dx}\left( {\text{sin} x} \right) = \text{cos} x$ and $\frac{d}{dx}\left( {\text{cos} x} \right) = \text{-}\text{sin} x$ use similar techniques, we provide only the proof for $\frac{d}{dx}\left( {\text{sin} x} \right) = \text{cos} x.$ Before beginning, recall two important trigonometric limits we learned in Introduction to Limits:

$$
\underset{h\rightarrow 0}{\text{lim}}\frac{\text{sin} h}{h} = 1\ \text{and}\ \underset{h\rightarrow 0}{\text{lim}}\frac{\text{cos} h - 1}{h} = 0.
$$

The graphs of $y = \frac{\left( {\text{sin} h} \right)}{h}$ and $y = \frac{\left( {\text{cos} h - 1} \right)}{h}$ are shown in Figure 3.26.

*Figure 3.26 These graphs show two important limits needed to establish the derivative formulas for the sine and cosine functions.*

We also recall the following trigonometric identity for the sine of the sum of two angles:

$$
\text{sin}\left( {x + h} \right) = \text{sin} x\text{cos} h + \text{cos} x\text{sin} h.
$$

Now that we have gathered all the necessary equations and identities, we proceed with the proof.

$$
\begin{array}{clccl}
{\frac{d}{dx}\text{sin} x} & {= \underset{h\rightarrow 0}{\text{lim}}\frac{\text{sin}\left( {x + h} \right) - \text{sin} x}{h}} & & & {\text{Apply the definition}\ \text{of the derivative.}} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\frac{\text{sin} x\text{cos} h + \text{cos} x\text{sin} h - \text{sin} x}{h}} & & & \text{Use trig identity for the sine of the sum of two angles.} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\left( {\frac{\text{sin} x\text{cos} h - \text{sin} x}{h} + \frac{\text{cos} x\text{sin} h}{h}} \right)} & & & \text{Regroup.} \\
 & {= \underset{h\rightarrow 0}{\text{lim}}\left( {\text{sin} x\left( \frac{\text{cos} h - 1}{h} \right) + \text{cos} x\left( \frac{\text{sin} h}{h} \right)} \right)} & & & {\text{Factor out}\ \text{sin} x\ \text{and}\ \text{cos} x.} \\
 & {= \text{sin} x{\cdot 0} + \text{cos} x{\cdot 1}} & & & \text{Apply trig limit formulas.} \\
 & {= \text{cos} x} & & & \text{Simplify.}
\end{array}
$$

□

Figure 3.27 shows the relationship between the graph of $f(x) = \text{sin} x$ and its derivative $f'(x) = \text{cos} x.$ Notice that at the points where $f(x) = \text{sin} x$ has a horizontal tangent, its derivative $f'(x) = \text{cos} x$ takes on the value zero. We also see that where $f(x) = \text{sin} x$ is increasing, $f'(x) = \text{cos} x > 0$ and where $f(x) = \text{sin} x$ is decreasing, $f'(x) = \text{cos} x < 0.$

*Figure 3.27 Where f ( x ) f ( x ) has a maximum or a minimum, f ′ ( x ) = 0 f ′ ( x ) = 0 that is, f ′ ( x ) = 0 f ′ ( x ) = 0 where f ( x ) f ( x ) has a horizontal tangent. These points are noted with dots on the graphs.*

### Example 3.39

#### Differentiating a Function Containing sin *x*

Find the derivative of $f(x) = 5x^{3}\text{sin} x.$

#### Solution

Using the product rule, we have

$$
\begin{array}{cl}
{f'(x)} & {= \frac{d}{dx}\left( {5x^{3}} \right) \cdot \text{sin} x + \frac{d}{dx}\left( {\text{sin} x} \right) \cdot 5x^{3}} \\
 & {= 15x^{2} \cdot \text{sin} x + \text{cos} x \cdot 5x^{3}.}
\end{array}
$$

After simplifying, we obtain

$$
f'(x) = 15x^{2}\text{sin} x + 5x^{3}\text{cos} x.
$$

### Checkpoint 3.25

Find the derivative of $f(x) = \text{sin} x\text{cos} x.$

### Example 3.40

#### Finding the Derivative of a Function Containing cos *x*

Find the derivative of $g(x) = \frac{\text{cos} x}{4x^{2}}.$

#### Solution

By applying the quotient rule, we have

$$
g'(x) = \frac{(\text{−}\text{sin} x)4x^{2} - 8x(\text{cos} x)}{\left( {4x^{2}} \right)^{2}}.
$$

Simplifying, we obtain

$$
\begin{array}{cl}
{g'(x)} & {= \frac{-4x^{2}\text{sin} x - 8x\text{cos} x}{16x^{4}}} \\
 & {= \frac{\text{−}x\text{sin} x - 2\text{cos} x}{4x^{3}}.}
\end{array}
$$

### Checkpoint 3.26

Find the derivative of $f(x) = \frac{x}{\text{cos} x}.$

### Example 3.41

#### An Application to Velocity

A particle moves along a coordinate axis in such a way that its position at time $t$ is given by $s(t) = 2\text{sin} t - t$ for $0 \leq t \leq 2\pi.$ At what times is the particle at rest?

#### Solution

To determine when the particle is at rest, set $s'(t) = v(t) = 0.$ Begin by finding $s'(t).$ We obtain

$$
s'(t) = 2\text{cos} t - 1,
$$

so we must solve

$$
2\text{cos} t - 1 = 0\ \text{for}\ 0 \leq t \leq 2\pi.
$$

The solutions to this equation are $t = \frac{\pi}{3}$ and $t = \frac{5\pi}{3}.$ Thus the particle is at rest at times $t = \frac{\pi}{3}$ and $t = \frac{5\pi}{3}.$

### Checkpoint 3.27

A particle moves along a coordinate axis. Its position at time $t$ is given by $s(t) = \sqrt{3}t + 2\text{cos} t$ for $0 \leq t \leq 2\pi.$ At what times is the particle at rest?

### Derivatives of Other Trigonometric Functions

Since the remaining four trigonometric functions may be expressed as quotients involving sine, cosine, or both, we can use the quotient rule to find formulas for their derivatives.

### Example 3.42

#### The Derivative of the Tangent Function

Find the derivative of $f(x) = \text{tan} x.$

#### Solution

Start by expressing $\text{tan} x$ as the quotient of $\text{sin} x$ and $\text{cos} x:$

$$
f(x) = \text{tan} x = \frac{\text{sin} x}{\text{cos} x}.
$$

Now apply the quotient rule to obtain

$$
f'(x) = \frac{\text{cos} x\text{cos} x - (\text{−}\text{sin} x)\text{sin} x}{\left( {\text{cos} x} \right)^{2}}.
$$

Simplifying, we obtain

$$
f'(x) = \frac{\text{cos}^{2}x + {\ \text{sin}}^{2}x}{\text{cos}^{2}x}.
$$

Recognizing that $\text{cos}^{2}x + \text{sin}^{2}x = 1,$ by the Pythagorean theorem, we now have

$$
f'(x) = \frac{1}{\text{cos}^{2}x}.
$$

Finally, use the identity $\text{sec} x = \frac{1}{\text{cos} x}$ to obtain

$$
f'(x) = \text{sec}^{2}x.
$$

### Checkpoint 3.28

Find the derivative of $f(x) = \text{cot} x.$

The derivatives of the remaining trigonometric functions may be obtained by using similar techniques. We provide these formulas in the following theorem.

### Theorem 3.9

#### Derivatives of $\text{tan} x,\text{cot} x,\text{sec} x,$ and $\text{csc} x$

The derivatives of the remaining trigonometric functions are as follows:

$$
\frac{d}{dx}\left( {\text{tan} x} \right) = \text{sec}^{2}x
$$

(3.13)
$$
\frac{d}{dx}(\text{cot} x) = \text{−}\text{csc}^{2}x
$$

(3.14)
$$
\frac{d}{dx}(\text{sec} x) = \text{sec} x\text{tan} x
$$

(3.15)
$$
\frac{d}{dx}(\text{csc} x) = \text{−}\text{csc} x\text{cot}{x.}
$$

(3.16)

### Example 3.43

#### Finding the Equation of a Tangent Line

Find the equation of a line tangent to the graph of $f(x) = \text{cot} x$ at $x = \frac{\text{π}}{4}.$

#### Solution

To find an equation of the tangent line, we need a point and a slope at that point. To find the point, compute

$$
f\left( \frac{\pi}{4} \right) = \text{cot}\ \frac{\pi}{4} = 1.
$$

Thus the tangent line passes through the point $\left( {\frac{\pi}{4},1} \right).$ Next, find the slope by finding the derivative of $f(x) = \text{cot} x$ and evaluating it at $\frac{\pi}{4}\text{:}$

$$
f'(x) = \text{−}\text{csc}^{2}x\ \text{and}\ f'\left( \frac{\pi}{4} \right) = \text{−}\text{csc}^{2}\left( \frac{\pi}{4} \right) = -2.
$$

Using the point-slope equation of the line, we obtain

$$
y - 1 = -2\left( {x - \frac{\pi}{4}} \right)
$$

or equivalently,

$$
y = -2x + 1 + \frac{\pi}{2}.
$$

### Example 3.44

#### Finding the Derivative of Trigonometric Functions

Find the derivative of $f(x) = \text{csc} x + x\text{tan} x.$

#### Solution

To find this derivative, we must use both the sum rule and the product rule. Using the sum rule, we find

$$
f'(x) = \frac{d}{dx}\left( {\text{csc} x} \right) + \frac{d}{dx}(x\text{tan} x).
$$

In the first term, $\frac{d}{dx}\left( {\text{csc} x} \right) = \text{-}\text{csc} x\text{cot} x,$ and by applying the product rule to the second term we obtain

$$
\frac{d}{dx}(x\text{tan} x) = (1)(\text{tan} x) + (\text{sec}^{2}x)(x).
$$

Therefore, we have

$$
f'(x) = \text{−}\text{csc} x\text{cot} x + \text{tan} x + x\text{sec}^{2}x.
$$

### Checkpoint 3.29

Find the derivative of $f(x) = 2\text{tan} x - 3\text{cot} x.$

### Checkpoint 3.30

Find the slope of the line tangent to the graph of $f(x) = \text{tan} x$ at $x = \frac{\pi}{6}.$

### Higher-Order Derivatives

The higher-order derivatives of $\text{sin} x$ and $\text{cos} x$ follow a repeating pattern. By following the pattern, we can find any higher-order derivative of $\text{sin} x$ and $\text{cos} x.$

### Example 3.45

#### Finding Higher-Order Derivatives of $y = \text{sin} x$

Find the first four derivatives of $y = \text{sin} x.$

#### Solution

Each step in the chain is straightforward:

$$
\begin{array}{rll}
y & = & {\text{sin} x} \\
\frac{dy}{dx} & = & {\text{cos} x} \\
\frac{d^{2}y}{dx^{2}} & = & {\text{−}\text{sin} x} \\
\frac{d^{3}y}{dx^{3}} & = & {\text{−}\text{cos} x} \\
\frac{d^{4}y}{dx^{4}} & = & {\text{sin} x.}
\end{array}
$$

#### Analysis

Once we recognize the pattern of derivatives, we can find any higher-order derivative by determining the step in the pattern to which it corresponds. For example, every fourth derivative of sin *x* equals sin *x*, so

$$
\begin{array}{l}
{\frac{d^{4}}{dx^{4}}\left( {\text{sin} x} \right) = \frac{d^{8}}{dx^{8}}\left( {\text{sin} x} \right) = \frac{d^{12}}{dx^{12}}\left( {\text{sin} x} \right) = \text{…} = \frac{d^{4n}}{dx^{4n}}\left( {\text{sin} x} \right) = \text{sin} x} \\
{\frac{d^{5}}{dx^{5}}\left( {\text{sin} x} \right) = \frac{d^{9}}{dx^{9}}\left( {\text{sin} x} \right) = \frac{d^{13}}{dx^{13}}\left( {\text{sin} x} \right) = \text{…} = \frac{d^{4n + 1}}{dx^{4n + 1}}\left( {\text{sin} x} \right) = \text{cos} x.}
\end{array}
$$

### Checkpoint 3.31

For $y = \text{cos} x,$ find $\frac{d^{4}y}{dx^{4}}.$

### Example 3.46

#### Using the Pattern for Higher-Order Derivatives of $y = \text{sin} x$

Find $\frac{d^{74}}{dx^{74}}\left( {\text{sin} x} \right).$

#### Solution

We can see right away that for the 74th derivative of $\text{sin} x,74 = 4(18) + 2,$ so

$$
\frac{d^{74}}{dx^{74}}\left( {\text{sin} x} \right) = \frac{d^{72 + 2}}{dx^{72 + 2}}\left( {\text{sin} x} \right) = \frac{d^{2}}{dx^{2}}\left( {\text{sin} x} \right) = \text{−}\text{sin} x.
$$

### Checkpoint 3.32

For $y = \text{sin} x,$ find $\frac{d^{59}}{dx^{59}}\left( {\text{sin} x} \right).$

### Example 3.47

#### An Application to Acceleration

A particle moves along a coordinate axis in such a way that its position at time $t$ is given by $s(t) = 2 - \text{sin} t.$ Find $v\left( {\pi\text{/}4} \right)$ and $a\left( {\pi\text{/}4} \right).$ Compare these values and decide whether the particle is speeding up or slowing down.

#### Solution

First find $v(t) = s'(t)\text{:}$

$$
v(t) = s'(t) = \text{−}\text{cos} t.
$$

Thus,

$$
v\left( \frac{\pi}{4} \right) = - \frac{1}{\sqrt{2}}.
$$

Next, find $a(t) = v'(t).$ Thus, $a(t) = v'(t) = \text{sin} t$ and we have

$$
a\left( \frac{\pi}{4} \right) = \frac{1}{\sqrt{2}}.
$$

Since $v\left( \frac{\pi}{4} \right) = - \frac{1}{\sqrt{2}} < 0$ and $a\left( \frac{\pi}{4} \right) = \frac{1}{\sqrt{2}} > 0,$ we see that velocity and acceleration are acting in opposite directions; that is, the object is being accelerated in the direction opposite to the direction in which it is travelling. Consequently, the particle is slowing down.

### Checkpoint 3.33

A block attached to a spring is moving vertically. Its position at time $t$ is given by $s(t) = 2\text{sin} t.$ Find $v\left( \frac{5\pi}{6} \right)$ and $a\left( \frac{5\pi}{6} \right).$ Compare these values and decide whether the block is speeding up or slowing down.

### Section 3.5 Exercises

For the following exercises, find $\frac{dy}{dx}$ for the given functions.

175\.

$y = x^{2} - \text{sec} x + 1$

176\.

$y = 3\text{csc} x + \frac{5}{x}$

177\.

$y = x^{2}\text{cot} x$

178\.

$y = x - x^{3}\text{sin} x$

179\.

$y = \frac{\text{sec} x}{x}$

180\.

$y = \text{sin} x\text{tan} x$

181\.

$y = \left( {x + \text{cos} x} \right)\left( {1 - \text{sin} x} \right)$

182\.

$y = \frac{\text{tan} x}{1 - \text{sec} x}$

183\.

$y = \frac{1 - \text{cot} x}{1 + \text{cot} x}$

184\.

$y = \text{cos} x\left( {1 + \text{csc} x} \right)$

For the following exercises, find an equation of the tangent line to each of the given functions at the indicated values of $x.$ Then use a calculator to graph both the function and the tangent line to ensure the equation for the tangent line is correct.

185\.

**\[T\]** $f(x) = \text{-}\text{sin} x,x = 0$

186\.

**\[T\]** $f(x) = \text{csc} x,x = \frac{\pi}{2}$

187\.

**\[T\]** $f(x) = 1 + \text{cos} x,x = \frac{3\pi}{2}$

188\.

**\[T\]** $f(x) = \text{sec} x,x = \frac{\pi}{4}$

189\.

**\[T\]** $f(x) = x^{2} - \text{tan} x,\ x = 0$

190\.

**\[T\]** $f(x) = 5\text{cot} x,\ x = \frac{\pi}{4}$

For the following exercises, find $\frac{d^{2}y}{dx^{2}}$ for the given functions.

191\.

$y = x\text{sin} x - \text{cos} x$

192\.

$y = \text{sin} x\text{cos} x$

193\.

$y = x - \frac{1}{2}\text{sin} x$

194\.

$y = \frac{1}{x} + \text{tan} x$

195\.

$y = 2\text{csc} x$

196\.

$y = \text{sec}^{2}x$

197\.

Find all $x$ values on the graph of $f(x) = -3\text{sin} x\text{cos} x$ where the tangent line is horizontal.

198\.

Find all $x$ values on the graph of $f(x) = x - 2\text{cos} x$ for $0 < x < 2\pi$ where the tangent line has slope 2.

199\.

Let $f(x) = \text{cot} x.$ Determine the points on the graph of $f$ for $0 < x < 2\pi$ where the tangent line(s) is (are) parallel to the line $y = -2x.$

200\.

**\[T\]** A mass on a spring bounces up and down in simple harmonic motion, modeled by the function $s(t) = -6\text{cos} t$ where $s$ is measured in inches and $t$ is measured in seconds. Find the rate at which the spring is oscillating at $t = 5$ s.

201\.

Let the position of a swinging pendulum in simple harmonic motion be given by $s(t) = a\text{cos} t + b\text{sin} t$ where $a$ and $b$ are constants, $t$ measures time in seconds, and $s$ measures position in centimeters. If the position is 0 cm and the velocity is 3 cm/s when $t = 0$, find the values of $a$ and $b$.

202\.

After a diver jumps off a diving board, the edge of the board oscillates with position given by $s(t) = -5\text{cos} t$ cm at $t$ seconds after the jump.

1.  Sketch one period of the position function for $t \geq 0.$
2.  Find the velocity function.
3.  Sketch one period of the velocity function for $t \geq 0.$
4.  Determine the times when the velocity is 0 over one period.
5.  Find the acceleration function.
6.  Sketch one period of the acceleration function for $t \geq 0.$

203\.

The number of hamburgers sold at a fast-food restaurant in Pasadena, California, is given by $y = 10 + 5\text{sin} x$ where $y$ is the number of hamburgers sold and $x$ represents the number of hours after the restaurant opened at 11 a.m. until 11 p.m., when the store closes. Find $y'$ and determine the intervals where the number of burgers being sold is increasing.

204\.

**\[T\]** The amount of rainfall per month in Phoenix, Arizona, can be approximated by $y(t) = 0.5 + 0.3\text{cos} t,$ where $t$ is months since January. Find $y'$ and use a calculator to determine the intervals where the amount of rain falling is decreasing.

For the following exercises, use the quotient rule to derive the given equations.

205\.

$\frac{d}{dx}(\text{cot} x) = \text{-}\text{csc}^{2}x$

206\.

$\frac{d}{dx}(\text{sec} x) = \text{sec} x\text{tan} x$

207\.

$\frac{d}{dx}(\text{csc} x) = \text{-}\text{csc} x\text{cot} x$

208\.

Use the definition of derivative and the identity

$\text{cos}\left( {x + h} \right) = \text{cos} x\text{cos} h - \text{sin} x\text{sin} h$ to prove that $\frac{d\left( {\text{cos} x} \right)}{dx} = \text{-}\text{sin} x.$

For the following exercises, find the requested higher-order derivative for the given functions.

209\.

$\frac{d^{3}y}{dx^{3}}$ of $y = 3\text{cos} x$

210\.

$\frac{d^{2}y}{dx^{2}}$ of $y = 3\text{sin} x + x^{2}\text{cos} x$

211\.

$\frac{d^{4}y}{dx^{4}}$ of $y = 5\text{cos} x$

212\.

$\frac{d^{2}y}{dx^{2}}$ of $y = \text{sec} x + \text{cot} x$

213\.

$\frac{d^{3}y}{dx^{3}}$ of $y = x^{10} - \text{sec} x$
