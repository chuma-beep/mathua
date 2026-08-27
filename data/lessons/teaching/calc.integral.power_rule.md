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

### Section 5.4 Exercises

Use basic integration formulas to compute the following antiderivatives or definite integrals.

207\.

$\left. \int{\left( {\sqrt{x} - \frac{1}{\sqrt{x}}} \right)dx} \right.$

208\.

$\left. \int{\left( {e^{2x} - \frac{1}{2}e^{x\text{/}2}} \right)dx} \right.$

209\.

$\left. \int\frac{dx}{2x} \right.$

210\.

$\left. \int{\frac{x - 1}{x^{2}}dx} \right.$

211\.

$\int_{0}^{\pi}{\left( {\text{sin} x - \text{cos} x} \right)dx}$

212\.

$\int_{0}^{\pi\text{/}2}{\left( {x - \text{sin} x} \right)dx}$

213\.

Write an integral that expresses the increase in the perimeter $P(s)$ of a square when its side length *s* increases from 2 units to 4 units and evaluate the integral.

214\.

Write an integral that quantifies the change in the area $A(s) = s^{2}$ of a square when the side length doubles from *S* units to 2*S* units and evaluate the integral.

215\.

A regular *N*-gon (an *N*-sided polygon with sides that have equal length *s*, such as a pentagon or hexagon) has perimeter *Ns*. Write an integral that expresses the increase in perimeter of a regular *N*-gon when the length of each side increases from 1 unit to 2 units and evaluate the integral.

216\.

The area of a regular pentagon with side length $a > 0$ is *pa*<sup>2</sup> with $p = \frac{1}{4}\sqrt{5\left( {5 + 2\sqrt{5}} \right)}.$ The Pentagon in Washington, DC, has inner sides of length 360 ft and outer sides of length 920 ft. Write an integral to express the area of the roof of the Pentagon according to these dimensions and evaluate this area.

217\.

A dodecahedron is a Platonic solid with a surface that consists of 12 pentagons, each of equal area. By how much does the surface area of a dodecahedron increase as the side length of each pentagon doubles from 1 unit to 2 units?

218\.

An icosahedron is a Platonic solid with a surface that consists of 20 equilateral triangles. By how much does the surface area of an icosahedron increase as the side length of each triangle doubles from *a* unit to 2*a* units?

219\.

Write an integral that quantifies the change in the area of the surface of a cube when its side length doubles from *s* unit to 2*s* units and evaluate the integral.

220\.

Write an integral that quantifies the increase in the volume of a cube when the side length doubles from *s* unit to 2*s* units and evaluate the integral.

221\.

Write an integral that quantifies the increase in the surface area of a sphere as its radius doubles from *R* unit to 2*R* units and evaluate the integral.

222\.

Write an integral that quantifies the increase in the volume of a sphere as its radius doubles from *R* unit to 2*R* units and evaluate the integral.

223\.

Suppose that a particle moves along a straight line with velocity $v(t) = 4 - 2t,$ where $0 \leq t \leq 2$ (in meters per second). Find the displacement at time *t* and the total distance traveled up to $t = 2.$

224\.

Suppose that a particle moves along a straight line with velocity defined by $v(t) = t^{2} - 3t - 18,$ where $0 \leq t \leq 6$ (in meters per second). Find the displacement at time *t* and the total distance traveled up to $t = 6.$

225\.

Suppose that a particle moves along a straight line with velocity defined by $v(t) = \left| {2t - 6} \right|,$ where $0 \leq t \leq 6$ (in meters per second). Find the displacement at time *t* and the total distance traveled up to $t = 6.$

226\.

Suppose that a particle moves along a straight line with acceleration defined by $a(t) = t - 3,$ where $0 \leq t \leq 6$ (in meters per second). Find the velocity and displacement at time *t* and the total distance traveled up to $t = 6$ if $v(0) = 3$ and $d(0) = 0.$

227\.

A ball is thrown upward from a height of 1.5 m at an initial speed of 40 m/sec. Acceleration resulting from gravity is −9.8 m/sec<sup>2</sup>. Neglecting air resistance, solve for the velocity $v(t)$ and the height $h(t)$ of the ball *t* seconds after it is thrown and before it returns to the ground.

228\.

A ball is thrown upward from a height of 3 m at an initial speed of 60 m/sec. Acceleration resulting from gravity is −9.8 m/sec<sup>2</sup>. Neglecting air resistance, solve for the velocity $v(t)$ and the height $h(t)$ of the ball *t* seconds after it is thrown and before it returns to the ground.

229\.

The area $A(t)$ of a circular shape is growing at a constant rate. If the area increases from 4*π* units to 9*π* units between times $t = 2$ and $t = 3,$ find the net change in the radius during that time.

230\.

A spherical balloon is being inflated at a constant rate. If the volume of the balloon changes from 36*π* in.<sup>3</sup> to 288*π* in.<sup>3</sup> between time $t = 30$ and $t = 60$ seconds, find the net change in the radius of the balloon during that time.

231\.

Water flows into a conical tank with cross-sectional area *πx*<sup>2</sup> at height *x* and volume $\frac{\pi x^{3}}{3}$ up to height *x*. If water flows into the tank at a rate of 1 m<sup>3</sup>/min, find the height of water in the tank after 5 min. Find the change in height between 5 min and 10 min.

232\.

A horizontal cylindrical tank has cross-sectional area $A(x) = 4\left( {6x - x^{2}} \right)m^{2}$ at height *x* meters above the bottom when $x \leq 3.$

1.  The volume *V* between heights *a* and *b* is ${\int_{a}^{b}{A(x)dx}}.$ Find the volume at heights between 2 m and 3 m.
2.  Suppose that oil is being pumped into the tank at a rate of 50 L/min. Using the chain rule, $\frac{dx}{dt} = \frac{dx}{dV}\ \frac{dV}{dt},$ at how many meters per minute is the height of oil in the tank changing, expressed in terms of *x*, when the height is at *x* meters?
3.  How long does it take to fill the tank to 3 m starting from a fill level of 2 m?

233\.

The following table lists the electrical power in gigawatts—the rate at which energy is consumed—used in a certain city for different hours of the day, in a typical 24-hour period, with hour 1 corresponding to midnight to 1 a.m.

| Hour | Power | Hour | Power |
|------|-------|------|-------|
| 1    | 28    | 13   | 48    |
| 2    | 25    | 14   | 49    |
| 3    | 24    | 15   | 49    |
| 4    | 23    | 16   | 50    |
| 5    | 24    | 17   | 50    |
| 6    | 27    | 18   | 50    |
| 7    | 29    | 19   | 46    |
| 8    | 32    | 20   | 43    |
| 9    | 34    | 21   | 42    |
| 10   | 39    | 22   | 40    |
| 11   | 42    | 23   | 37    |
| 12   | 46    | 24   | 34    |

Find the total amount of energy in gigawatt-hours (gW-h) consumed by the city in a typical 24-hour period.

234\.

The average residential electrical power use (in hundreds of watts) per hour is given in the following table.

| Hour | Power | Hour | Power |
|------|-------|------|-------|
| 1    | 8     | 13   | 12    |
| 2    | 6     | 14   | 13    |
| 3    | 5     | 15   | 14    |
| 4    | 4     | 16   | 15    |
| 5    | 5     | 17   | 17    |
| 6    | 6     | 18   | 19    |
| 7    | 7     | 19   | 18    |
| 8    | 8     | 20   | 17    |
| 9    | 9     | 21   | 16    |
| 10   | 10    | 22   | 16    |
| 11   | 10    | 23   | 13    |
| 12   | 11    | 24   | 11    |

1.  Compute the average total energy used in a day in kilowatt-hours (kWh).
2.  If a ton of coal generates 1842 kWh, how long does it take for an average residence to burn a ton of coal?
3.  Explain why the data might fit a plot of the form $p(t) = 11.5 - 7.5\text{sin}\left( \frac{\pi t}{12} \right).$

235\.

The data in the following table are used to estimate the average power output produced by Peter Sagan for each of the last 18 sec of Stage 1 of the 2012 Tour de France.

| Second | Watts | Second | Watts |
|--------|-------|--------|-------|
| 1      | 600   | 10     | 1200  |
| 2      | 500   | 11     | 1170  |
| 3      | 575   | 12     | 1125  |
| 4      | 1050  | 13     | 1100  |
| 5      | 925   | 14     | 1075  |
| 6      | 950   | 15     | 1000  |
| 7      | 1050  | 16     | 950   |
| 8      | 950   | 17     | 900   |
| 9      | 1100  | 18     | 780   |

Table 5.6 Average Power Output *Source*: sportsexercisengineering.com

Estimate the net energy used in kilojoules (kJ), noting that 1W = 1 J/s, and the average power output by Sagan during this time interval.

236\.

The data in the following table are used to estimate the average power output produced by Peter Sagan for each 15-min interval of Stage 1 of the 2012 Tour de France.

| Minutes | Watts | Minutes | Watts |
|---------|-------|---------|-------|
| 15      | 200   | 165     | 170   |
| 30      | 180   | 180     | 220   |
| 45      | 190   | 195     | 140   |
| 60      | 230   | 210     | 225   |
| 75      | 240   | 225     | 170   |
| 90      | 210   | 240     | 210   |
| 105     | 210   | 255     | 200   |
| 120     | 220   | 270     | 220   |
| 135     | 210   | 285     | 250   |
| 150     | 150   | 300     | 400   |

Table 5.7 Average Power Output *Source*: sportsexercisengineering.com

Estimate the net energy used in kilojoules, noting that 1W = 1 J/s.

237\.

The distribution of incomes as of 2012 in the United States in \$5000 increments is given in the following table. The *k*th row denotes the percentage of households with incomes between $\$ 5000xk$ and $5000xk + 4999.$ The row $k = 40$ contains all households with income between \$200,000 and \$250,000.

|     |     |     |     |     |      |     |     |
|-----|-----|-----|-----|-----|------|-----|-----|
| 0   | 3.5 | 11  | 3.5 | 21  | 1.5  | 31  | 0.6 |
| 1   | 4.1 | 12  | 3.7 | 22  | 1.4  | 32  | 0.5 |
| 2   | 5.9 | 13  | 3.2 | 23  | 1.3  | 33  | 0.5 |
| 3   | 5.7 | 14  | 3.0 | 24  | 1.3  | 34  | 0.4 |
| 4   | 5.9 | 15  | 2.8 | 25  | 1.1  | 35  | 0.3 |
| 5   | 5.4 | 16  | 2.5 | 26  | 1.0  | 36  | 0.3 |
| 6   | 5.5 | 17  | 2.2 | 27  | 0.75 | 37  | 0.3 |
| 7   | 5.1 | 18  | 2.2 | 28  | 0.8  | 38  | 0.2 |
| 8   | 4.8 | 19  | 1.8 | 29  | 1.0  | 39  | 1.8 |
| 9   | 4.1 | 20  | 2.1 | 30  | 0.6  | 40  | 2.3 |
| 10  | 4.3 |     |     |     |      |     |     |

Table 5.8 Income Distributions *Source*: http://www.census.gov/prod/2013pubs/p60-245.pdf

1.  Estimate the percentage of U.S. households in 2012 with incomes less than \$55,000.
2.  What percentage of households had incomes exceeding \$85,000?
3.  Plot the data and try to fit its shape to that of a graph of the form $a\left( {x + c} \right)e^{\text{-}b{({x + e})}}$ for suitable $a,b,c.$

238\.

Newton’s law of gravity states that the gravitational force exerted by an object of mass *M* and one of mass *m* with centers that are separated by a distance *r* is $F = G\frac{mM}{r^{2}},$ with *G* an empirical constant $G = 6.67x10^{-11}\ m^{3}\text{/}\left( {kg \cdot s^{2}} \right).$ The work done by a variable force over an interval $\left\lbrack {a,b} \right\rbrack$ is defined as $W = {\int_{a}^{b}{F(x)dx}}.$ If Earth has mass $5.97219\  \times \ 10^{24}$ and radius 6371 km, compute the amount of work to elevate a polar weather satellite of mass 1400 kg to its orbiting altitude of 850 km above Earth.

239\.

For a given motor vehicle, the maximum achievable deceleration from braking is approximately 7 m/sec<sup>2</sup> on dry concrete. On wet asphalt, it is approximately 2.5 m/sec<sup>2</sup>. Given that 1 mph corresponds to 0.447 m/sec, find the total distance that a car travels in meters on dry concrete after the brakes are applied until it comes to a complete stop if the initial velocity is 67 mph (30 m/sec) or if the initial braking velocity is 56 mph (25 m/sec). Find the corresponding distances if the surface is slippery wet asphalt.

240\.

John is a 25-year old man who weighs 160 lb. He burns $500 - 50t$ calories/hr while riding his bike for *t* hours. If an oatmeal cookie has 55 cal and John eats cookies at a rate of 4*t* calories during the *t*th hour, how many net calories has he lost after 3 hours riding his bike?

241\.

Sandra is a 25-year old woman who weighs 120 lb. She burns $300 - 50t$ cal/hr while walking on her treadmill. Her caloric intake from drinking Gatorade is 100*t* calories/hour during the *t*th hour. What is her net decrease in calories after walking for 3 hours?

242\.

A motor vehicle has a maximum efficiency of 33 mpg at a cruising speed of 40 mph. The efficiency drops at a rate of 0.1 mpg/mph between 40 mph and 50 mph, and at a rate of 0.4 mpg/mph between 50 mph and 80 mph. What is the efficiency in miles per gallon if the car is cruising at 50 mph? What is the efficiency in miles per gallon if the car is cruising at 80 mph? If gasoline costs \$3.50/gal, what is the cost of fuel to drive 50 mi at 40 mph, at 50 mph, and at 80 mph?

243\.

Although some engines are more efficient at given a horsepower than others, on average, fuel efficiency decreases with horsepower at a rate of $1\text{/}25$ mpg/horsepower. If a typical 50-horsepower engine has an average fuel efficiency of 32 mpg, what is the average fuel efficiency of an engine with the following horsepower: 150, 300, 450?

244\.

**\[T\]** The following table lists the 2013 schedule of federal income tax versus taxable income.

| Taxable Income Range | The Tax Is …         | … Of the Amount Over |
|----------------------|----------------------|----------------------|
| \$0-\$8925           | 10%                  | \$0                  |
| \$8925-\$36,250      | \$892.50 + 15%       | \$8925               |
| \$36,250-\$87,850    | \$4,991.25 + 25%     | \$36,250             |
| \$87,850-\$183,250   | \$17,891.25 + 28%    | \$87,850             |
| \$183,250-\$398,350  | \$44,603.25 + 33%    | \$183,250            |
| \$398,350-\$400,000  | \$115,586.25 + 35%   | \$398,350            |
| \> \$400,000         | \$116,163.75 + 39.6% | \$400,000            |

Table 5.9 Federal Income Tax Versus Taxable Income *Source*: http://www.irs.gov/pub/irs-prior/i1040tt--2013.pdf.

Suppose that Steve just received a \$10,000 raise. How much of this raise is left after federal taxes if Steve’s salary before receiving the raise was \$40,000? If it was \$90,000? If it was \$385,000?

245\.

**\[T\]** The following table provides hypothetical data regarding the level of service for a certain highway.

| Highway Speed Range (mph) | Vehicles per Hour per Lane | Density Range (vehicles/mi) |
|----|----|----|
| \> 60 | \< 600 | \< 10 |
| 60–57 | 600–1000 | 10–20 |
| 57–54 | 1000–1500 | 20–30 |
| 54–46 | 1500–1900 | 30–45 |
| 46–30 | 1900**–**2100 | 45–70 |
| \<30 | Unstable | 70–200 |

Table 5.10

1.  Plot vehicles per hour per lane on the *x*-axis and highway speed on the *y*-axis.
2.  Compute the average decrease in speed (in miles per hour) per unit increase in congestion (vehicles per hour per lane) as the latter increases from 600 to 1000, from 1000 to 1500, and from 1500 to 2100. Does the decrease in miles per hour depend linearly on the increase in vehicles per hour per lane?
3.  Plot minutes per mile (60 times the reciprocal of miles per hour) as a function of vehicles per hour per lane. Is this function linear?

For the next two exercises use the data in the following table, which displays bald eagle populations from 1963 to 2000 in the continental United States.

| Year | Population of Breeding Pairs of Bald Eagles |
|------|---------------------------------------------|
| 1963 | 487                                         |
| 1974 | 791                                         |
| 1981 | 1188                                        |
| 1986 | 1875                                        |
| 1992 | 3749                                        |
| 1996 | 5094                                        |
| 2000 | 6471                                        |

Table 5.11 Population of Breeding Bald Eagle Pairs *Source*: http://www.fws.gov/Midwest/eagle/population/chtofprs.html. 246.

**\[T\]** The graph below plots the quadratic ${p(t) = 6.48t^{2} - 80.3}\ {1t + 585.69}$ against the data in preceding table, normalized so that $t = 0$ corresponds to 1963. Estimate the average number of bald eagles per year present for the 37 years by computing the average value of *p* over $\left\lbrack {0,37} \right\rbrack.$

247\.

**\[T\]** The graph below plots the cubic $p(t) = 0.07t^{3} + 2.42t^{2} - 25.63t + 521.23$ against the data in the preceding table, normalized so that $t = 0$ corresponds to 1963. Estimate the average number of bald eagles per year present for the 37 years by computing the average value of *p* over $\left\lbrack {0,37} \right\rbrack.$

248\.

**\[T\]** Suppose you go on a road trip and record your speed at every half hour, as compiled in the following table. The best quadratic fit to the data is $q(t) = 5x^{2} - 11x + 49\text{,}$ shown in the accompanying graph. Integrate *q* to estimate the total distance driven over the 3 hours.

| Time (hr) | Speed (mph) |
|-----------|-------------|
| 0 (start) | 50          |
| 1         | 40          |
| 2         | 50          |
| 3         | 60          |

As a car accelerates, it does not accelerate at a constant rate; rather, the acceleration is variable. For the following exercises, use the following table, which contains the acceleration measured at every second as a driver merges onto a freeway.

| Time (sec) | Acceleration (mph/sec) |
|------------|------------------------|
| 1          | 11.2                   |
| 2          | 10.6                   |
| 3          | 8.1                    |
| 4          | 5.4                    |
| 5          | 0                      |

As a car accelerates, it does not accelerate at a constant rate; rather, the acceleration is variable. For the next three exercises use the following table, which contains the acceleration measured at every second as a driver merges onto a freeway.

| Time (sec) | Acceleration (mph/sec) |
|------------|------------------------|
| 1          | 11.2                   |
| 2          | 10.6                   |
| 3          | 8.1                    |
| 4          | 5.4                    |
| 5          | 0                      |

Table 5.12 249.

**\[T\]** The accompanying graph plots the best quadratic fit, $a(t) = -0.70t^{2} + 1.44t + 10.44,$ to the data from the preceding table. Compute the average value of $a(t)$ to estimate the average acceleration between $t = 0$ and $t = 5.$

250\.

**\[T\]** Using your acceleration equation from the previous exercise, find the corresponding velocity equation. Assuming the initial velocity is 65 mph, find the velocity at time $t = 0.$

251\.

**\[T\]** Using your velocity equation from the previous exercise, find the corresponding distance equation, assuming your initial distance is 0 mi. How far did you travel while you accelerated your car? (*Hint:* You will need to convert time units.)

252\.

**\[T\]** The number of hamburgers sold at a restaurant throughout the day is given in the following table, with the accompanying graph plotting the best cubic fit to the data, $b(t) = 0.12t^{3} - 2.13t^{2} + 12.13t + 3.91,$ with $t = 0$ corresponding to 9 a.m. and $t = 12$ corresponding to 9 p.m. Compute the average value of $b(t)$ to estimate the average number of hamburgers sold per hour.

| Hours Past Midnight | No. of Burgers Sold |
|---------------------|---------------------|
| 9                   | 3                   |
| 12                  | 28                  |
| 15                  | 20                  |
| 18                  | 30                  |
| 21                  | 45                  |

253\.

**\[T\]** An athlete runs by a motion detector, which records her speed, as displayed in the following table. The best linear fit to this data, $\ell(t) = -0.068t + 5.14\text{,}$ is shown in the accompanying graph. Use the average value of $\ell(t)$ between $t = 0$ and $t = 40$ to estimate the runner’s average speed.

| Minutes | Speed (m/sec) |
|---------|---------------|
| 0       | 5             |
| 10      | 4.8           |
| 20      | 3.6           |
| 30      | 3.0           |
| 40      | 2.5           |
