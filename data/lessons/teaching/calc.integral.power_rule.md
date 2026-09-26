> Content sourced from [OpenStax Calculus Volume 1](https://openstax.org/books/calculus-volume-1/pages/1-introduction) by Gilbert Strang & Edwin "Jed" Herman — CC BY-NC-SA 4.0

## 5.4 Integration Formulas and the Net Change Theorem

### Learning Objectives

- 5.4.1 Apply the basic integration formulas.
- 5.4.2 Explain the significance of the net change theorem.
- 5.4.3 Use the net change theorem to solve applied problems.
- 5.4.4 Apply the integrals of odd and even functions.

In this section, we use some basic integration formulas studied previously to solve some key applied problems. It is important to note that these formulas are presented in terms of *indefinite* integrals. Although definite and indefinite integrals are closely related, there are some key differences to keep in mind. A definite integral is either a number (when the limits of integration are constants) or a single function (when one or both of the limits of integration are variables). An indefinite integral represents a family of functions, all of which differ by a constant. As you become more familiar with integration, you will get a feel for when to use definite integrals and when to use indefinite integrals. You will naturally select the correct approach for a given problem without thinking too much about it. However, until these concepts are cemented in your mind, think carefully about whether you need a definite integral or an indefinite integral and make sure you are using the proper notation based on your choice.

### Basic Integration Formulas

Recall the integration formulas given in the table in Antiderivatives and the rule on properties of definite integrals. Let’s look at a few examples of how to apply these rules.

### Example 5.23

#### Integrating a Function Using the Power Rule

Use the power rule to integrate the function ${\int_{1}^{4}{\sqrt{t}\left( {1 + t} \right)dt}}.$

#### Solution

The first step is to rewrite the function and simplify it so we can apply the power rule:

$$
\begin{array}{cl}
{\int_{1}^{4}{\sqrt{t}(1 + t)dt}} & {= {\int_{1}^{4}{t^{1\text{/}2}(1 + t)dt}}} \\
 & \\
 & {= {\int_{1}^{4}{\left( {t^{1\text{/}2} + t^{3\text{/}2}} \right)dt}}.}
\end{array}
$$

Now apply the power rule:

$$
\begin{array}{cl}
{\int_{1}^{4}{\left( {t^{1\text{/}2} + t^{3\text{/}2}} \right)dt}} & {= \left. \left( {\frac{2}{3}t^{3\text{/}2} + \frac{2}{5}t^{5\text{/}2}} \right) \right|_{1}^{4}} \\
 & {= \left\lbrack {\frac{2}{3}{(4)}^{3\text{/}2} + \frac{2}{5}{(4)}^{5\text{/}2}} \right\rbrack - \left\lbrack {\frac{2}{3}{(1)}^{3\text{/}2} + \frac{2}{5}{(1)}^{5\text{/}2}} \right\rbrack} \\
 & {= \frac{256}{15}.}
\end{array}
$$

### Checkpoint 5.21

Find the definite integral of $f(x) = x^{2} - 3x$ over the interval $\left\lbrack {1,3} \right\rbrack.$

### The Net Change Theorem

The net change theorem considers the integral of a *rate of change*. It says that when a quantity changes, the new value equals the initial value plus the integral of the rate of change of that quantity. The formula can be expressed in two ways. The second is more familiar; it is simply the definite integral.

### Theorem 5.6

#### Net Change Theorem

The new value of a changing quantity equals the initial value plus the integral of the rate of change:

$$
\begin{matrix}
 \\
 \\
{F(b) = F(a) + {\int_{a}^{b}{F'(x)dx}}} \\
\text{or} \\
{{\int_{a}^{b}{F'(x)dx = F(b) - F(a)}}.}
\end{matrix}
$$

(5.18)

Subtracting $F(a)$ from both sides of the first equation yields the second equation. Since they are equivalent formulas, which one we use depends on the application.

The significance of the net change theorem lies in the results. Net change can be applied to area, distance, and volume, to name only a few applications. Net change accounts for negative quantities automatically without having to write more than one integral. To illustrate, let’s apply the net change theorem to a velocity function in which the result is displacement.

We looked at a simple example of this in The Definite Integral. Suppose a car is moving due north (the positive direction) at 40 mph between 2 p.m. and 4 p.m., then the car moves south at 30 mph between 4 p.m. and 5 p.m. We can graph this motion as shown in Figure 5.32.

*Figure 5.32 The graph shows speed versus time for the given motion of a car.*

Just as we did before, we can use definite integrals to calculate the net displacement as well as the total distance traveled. The net displacement is given by

$$
\begin{array}{cl}
{{\int_{2}^{5}v}(t)dt} & {= {\int_{2}^{4}4}0dt + \int_{4}^{5}-30dt} \\
 & {= 80 - 30} \\
 & {= 50.}
\end{array}
$$

Thus, at 5 p.m. the car is 50 mi north of its starting position. The total distance traveled is given by

$$
\begin{array}{cl}
 & \\
 & \\
{\int_{2}^{5}\left| {v(t)} \right|dt} & {= {\int_{2}^{4}4}0dt + \int_{4}^{5}30dt} \\
 & {= 80 + 30} \\
 & {= 110.}
\end{array}
$$

Therefore, between 2 p.m. and 5 p.m., the car traveled a total of 110 mi.

To summarize, net displacement may include both positive and negative values. In other words, the velocity function accounts for both forward distance and backward distance. To find net displacement, integrate the velocity function over the interval. Total distance traveled, on the other hand, is always positive. To find the total distance traveled by an object, regardless of direction, we need to integrate the absolute value of the velocity function.

### Example 5.24

#### Finding Net Displacement

Given a velocity function $v(t) = 3t - 5$ (in meters per second) for a particle in motion from time $t = 0$ to time $t = 3,$ find the net displacement of the particle.

#### Solution

Applying the net change theorem, we have

$$
\begin{array}{ll}
{\int_{0}^{3}{\left( {3t - 5} \right)dt}} & {= \frac{3t^{2}}{2} - 5t|_{0}^{3}} \\
 & \\
 & {= \left\lbrack {\frac{3(3)^{2}}{2} - 5(3)} \right\rbrack - 0} \\
 & {= \frac{27}{2} - 15} \\
 & {= \frac{27}{2} - \frac{30}{2}} \\
 & {= - \frac{3}{2}.}
\end{array}
$$

The net displacement is $- \frac{3}{2}$ m (Figure 5.33).

*Figure 5.33 The graph shows velocity versus time for a particle moving with a linear velocity function.*

### Example 5.25

#### Finding the Total Distance Traveled

Use Example 5.24 to find the total distance traveled by a particle according to the velocity function $v(t) = 3t - 5$ m/sec over a time interval $\left\lbrack {0,3} \right\rbrack.$

#### Solution

The total distance traveled includes both the positive and the negative values. Therefore, we must integrate the absolute value of the velocity function to find the total distance traveled.

To continue with the example, use two integrals to find the total distance. First, find the *t*-intercept of the function, since that is where the division of the interval occurs. Set the equation equal to zero and solve for *t*. Thus,

$$
\begin{array}{cll}
{3t - 5} & = & 0 \\
{3t} & = & 5 \\
t & = & {\frac{5}{3}.}
\end{array}
$$

The two subintervals are $\left\lbrack {0,\frac{5}{3}} \right\rbrack$ and $\left\lbrack {\frac{5}{3},3} \right\rbrack.$ To find the total distance traveled, integrate the absolute value of the function. Since the function is negative over the interval $\left\lbrack {0,\frac{5}{3}} \right\rbrack,$ we have $\left| {v(t)} \right| = \text{-}v(t)$ over that interval. Over $\left\lbrack {\frac{5}{3},3} \right\rbrack,$ the function is positive, so $\left| {v(t)} \right| = v(t).$ Thus, we have

$$
\begin{array}{cl}
 & \\
 & \\
{\int_{0}^{3}{\left| {v(t)} \right|dt}} & {= {\int_{0}^{5\text{/}3}{\text{−}v(t)dt + {\int_{5\text{/}3}^{3}v}(t)dt}}} \\
 & \\
 & {= {\int_{0}^{5\text{/}3}5} - 3tdt + {\int_{5\text{/}3}^{3}3}t - 5dt} \\
 & {= \left. \left( {5t - \frac{3t^{2}}{2}} \right) \right|_{0}^{5\text{/}3} + \left. \left( {\frac{3t^{2}}{2} - 5t} \right) \right|_{5\text{/}3}^{3}} \\
 & {= \left\lbrack {5\left( \frac{5}{3} \right) - \frac{3\left( {5\text{/}3} \right)^{2}}{2}} \right\rbrack - 0 + \left\lbrack {\frac{27}{2} - 15} \right\rbrack - \left\lbrack {\frac{3\left( {5\text{/}3} \right)^{2}}{2} - \frac{25}{3}} \right\rbrack} \\
 & {= \frac{25}{3} - \frac{25}{6} + \frac{27}{2} - 15 - \frac{25}{6} + \frac{25}{3}} \\
 & {= \frac{41}{6}.}
\end{array}
$$

So, the total distance traveled is $\frac{41}{6}$ m.

### Checkpoint 5.22

Find the net displacement and total distance traveled in meters given the velocity function $f(t) = \frac{1}{2}e^{t} - 2$ over the interval $\left\lbrack {0,2} \right\rbrack.$

### Applying the Net Change Theorem

The net change theorem can be applied to the flow and consumption of fluids, as shown in Example 5.26.

### Example 5.26

#### How Many Gallons of Gasoline Are Consumed?

If the motor on a motorboat is started at $t = 0$ and the boat consumes gasoline at the rate of $5 - 0.1t^{3}$ gal/hr, how much gasoline is used in the first 2 hours?

#### Solution

Express the problem as a definite integral, integrate, and evaluate using the Fundamental Theorem of Calculus. The limits of integration are the endpoints of the interval $\left\lbrack {0,2} \right\rbrack.$ We have

$$
\int_{0}^{2}\left( 5 - 0.1t^{3} \right)dt = \left( 5t–0.1\frac{t^{4}}{4} \right)\left. \middle| {}_{2} \right._{0} = \left\lbrack 5(2)–0.1\frac{(2)^{4}}{4} \right\rbrack –0 = 10–0.4 = 9.6
$$

Thus, the motorboat uses 9.6 gal of gas in 2 hours.

### Example 5.27

#### Chapter Opener: Iceboats

*Figure 5.34 (credit: modification of work by Carter Brown, Flickr)*

As we saw at the beginning of the chapter, top iceboat racers (Figure 5.1) can attain speeds of up to five times the wind speed. Andrew is an intermediate iceboater, though, so he attains speeds equal to only twice the wind speed. Suppose Andrew takes his iceboat out one morning when a light 5-mph breeze has been blowing all morning. As Andrew gets his iceboat set up, though, the wind begins to pick up. During his first half hour of iceboating, the wind speed increases according to the function $v(t) = 20t + 5.$ For the second half hour of Andrew’s outing, the wind remains steady at 15 mph. In other words, the wind speed is given by

$$
v(t) = \left\{ \begin{array}{lll}
{20t + 5} & \text{for} & {0 \leq t \leq \frac{1}{2}} \\
15 & \text{for} & {\frac{1}{2} \leq t \leq 1.}
\end{array} \right.
$$

Recalling that Andrew’s iceboat travels at twice the wind speed, and assuming he moves in a straight line away from his starting point, how far is Andrew from his starting point after 1 hour?

#### Solution

To figure out how far Andrew has traveled, we need to integrate his velocity, which is twice the wind speed. Then

Distance $= {\int_{0}^{1}{2v(t)dt}}.$

Substituting the expressions we were given for $v(t),$ we get

$$
\begin{array}{cl}
{\int_{0}^{1}{2v(t)dt}} & {= {\int_{0}^{1\text{/}2}{2v(t)dt + {\int_{1\text{/}2}^{1}{2v(t)dt}}}}} \\
 & {= {\int_{0}^{1\text{/}2}{2\left( {20t + 5} \right)dt + {\int_{1\text{/}2}^{1}{2(15)dt}}}}} \\
 & {= {\int_{0}^{1\text{/}2}{\left( {40t + 10} \right)dt + {\int_{1\text{/}2}^{1}{30dt}}}}} \\
 & {= \left\lbrack {20t^{2} + 10t} \right\rbrack{|_{0}^{1\text{/}2} + \left\lbrack {30t} \right\rbrack|_{1\text{/}2}^{1}}} \\
 & {= \left( {\frac{20}{4} + 5} \right) - 0 + \left( {30 - 15} \right)} \\
 & {= 25.}
\end{array}
$$

Andrew is 25 mi from his starting point after 1 hour.

### Checkpoint 5.23

Suppose that, instead of remaining steady during the second half hour of Andrew’s outing, the wind starts to die down according to the function $v(t) = -10t + 20.$ In other words, the wind speed is given by

$$
v(t) = \left\{ \begin{array}{lll}
{20t + 5} & \text{for} & {0 \leq t \leq \frac{1}{2}} \\
{- 10t + 15} & \text{for} & {\frac{1}{2} \leq t \leq 1.}
\end{array} \right.
$$

Under these conditions, how far from his starting point is Andrew after 1 hour?

### Integrating Even and Odd Functions

We saw in Functions and Graphs that an even function is a function in which $f\left( {\text{-}x} \right) = f(x)$ for all *x* in the domain—that is, the graph of the curve is unchanged when *x* is replaced with −*x*. The graphs of even functions are symmetric about the *y*-axis. An odd function is one in which $f\left( {\text{-}x} \right) = \text{-}f(x)$ for all *x* in the domain, and the graph of the function is symmetric about the origin.

Integrals of even functions, when the limits of integration are from −*a* to *a*, involve two equal areas, because they are symmetric about the *y*-axis. Integrals of odd functions, when the limits of integration are similarly $\left\lbrack {\text{-}a,a} \right\rbrack,$ evaluate to zero because the areas above and below the *x*-axis are equal.

### Rule: Integrals of Even and Odd Functions

For continuous even functions such that $f\left( {\text{-}x} \right) = f(x),$

$$
{\int_{\text{−}a}^{a}{f(x)dx = 2{\int_{0}^{a}{f(x)dx}}}}.
$$

For continuous odd functions such that $f\left( {\text{-}x} \right) = \text{-}f(x),$

$$
{\int_{\text{−}a}^{a}{f(x)dx = 0}}.
$$

### Example 5.28

#### Integrating an Even Function

Integrate the even function $\int_{-2}^{2}{\left( {3x^{8} - 2} \right)dx}$ and verify that the integration formula for even functions holds.

#### Solution

The symmetry appears in the graphs in Figure 5.35. Graph (a) shows the region below the curve and above the *x*-axis. We have to zoom in to this graph by a huge amount to see the region. Graph (b) shows the region above the curve and below the *x*-axis. The signed area of this region is negative. Both views illustrate the symmetry about the *y*-axis of an even function. We have

$$
\begin{array}{ll}
{\int_{-2}^{2}{\left( {3x^{8} - 2} \right)dx}} & {= \left( {\frac{x^{9}}{3} - 2x} \right)|_{-2}^{2}} \\
 & \\
 & \\
 & {= \left\lbrack {\frac{(2)^{9}}{3} - 2(2)} \right\rbrack - \left\lbrack {\frac{(-2)^{9}}{3} - 2(-2)} \right\rbrack} \\
 & {= \left( {\frac{512}{3} - 4} \right) - \left( {- \frac{512}{3} + 4} \right)} \\
 & {= \frac{1000}{3}.}
\end{array}
$$

To verify the integration formula for even functions, we can calculate the integral from 0 to 2 and double it, then check to make sure we get the same answer.

$$
\begin{array}{ll}
{\int_{0}^{2}{\left( {3x^{8} - 2} \right)dx}} & {= \left( {\frac{x^{9}}{3} - 2x} \right)|_{0}^{2}} \\
 & \\
 & {= \frac{512}{3} - 4} \\
 & {= \frac{500}{3}}
\end{array}
$$

Since $2 \cdot \frac{500}{3} = \frac{1000}{3},$ we have verified the formula for even functions in this particular example.

*Figure 5.35 Graph (a) shows the positive area between the curve and the x -axis, whereas graph (b) shows the negative area between the curve and the x -axis. Both views show the symmetry about the y -axis.*

### Example 5.29

#### Integrating an Odd Function

Evaluate the definite integral of the odd function $-5\text{sin} x$ over the interval $\left\lbrack {\text{-}\pi,\pi} \right\rbrack.$

#### Solution

The graph is shown in Figure 5.36. We can see the symmetry about the origin by the positive area above the *x*-axis over $\left\lbrack {\text{-}\pi,0} \right\rbrack,$ and the negative area below the *x*-axis over $\left\lbrack {0,\pi} \right\rbrack.$ We have

$$
\begin{array}{ll}
{\int_{\text{−}\pi}^{\pi}{-5\text{sin} xdx}} & {= -5\left( {\text{−}\text{cos} x} \right)|_{\text{−}\pi}^{\pi}} \\
 & \\
 & \\
 & {= 5\text{cos} x|_{\text{−}\pi}^{\pi}} \\
 & {= \left\lbrack {5\text{cos}\pi} \right\rbrack - \left\lbrack {5\text{cos}\left( {\text{−}\pi} \right)} \right\rbrack} \\
 & {= -5 - (-5)} \\
 & {= 0.}
\end{array}
$$

*Figure 5.36 The graph shows areas between a curve and the x -axis for an odd function.*

### Checkpoint 5.24

Integrate the function ${\int_{-2}^{2}{x^{4}dx}}.$
