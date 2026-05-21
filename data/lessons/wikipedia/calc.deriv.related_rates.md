> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Related_rates) — CC BY-SA 4.0

# Related rates

In differential calculus, **related rates** problems involve finding a rate at which a quantity changes by relating that quantity to other quantities whose rates of change are known. The rate of change is usually with respect to time. Because science and engineering often relate quantities to each other, the methods of related rates have broad applications in these fields. Differentiation with respect to time or one of the other variables requires application of the chain rule; since most problems involve several variables with respect to the input variable.

Fundamentally, if a function \(F\) is defined such that \(F = f(x)\), then the derivative of the function \(F\) can be taken with respect to another variable. We assume \(x\) is a function of \(t\), i.e. \(x=g(t)\). Then \(F=f(g(t))\), so

\(F'(t) =f'(g(t)) \cdot g'(t)\)

Written in Leibniz notation, this is:

\(\frac{dF}{dt} = \frac{df}{dx} \cdot \frac{dx}{dt}.\)

Thus, if it is known how \(x\) changes with respect to \(t\), then we can determine how \(F\) changes with respect to \(t\) and vice versa. We can extend this application of the chain rule with the sum, difference, product and quotient rules of calculus, etc.

For example, if \(F(x)= G(y)+ H(z)\) then
\(\frac{dF}{dx}\cdot \frac{dx}{dt} =\frac{dG}{dy} \cdot \frac{dy}{dt}+ \frac{dH}{dz} \cdot \frac{dz}{dt}.\)

## Procedure
The most common way to approach related rates problems is the following:
#Identify the known variables, including rates of change and the rate of change that is to be found. (Drawing a picture or representation of the problem can help to keep everything in order)
#Construct an equation relating the quantities whose rates of change are known to the quantity whose rate of change is to be found.
#Differentiate both sides of the equation with respect to time (or other rate of change). Often, the chain rule is employed at this step.
#Substitute the known rates of change and the known quantities into the equation.
#Solve for the wanted rate of change.

Errors in this procedure are often caused by plugging in the known values for the variables *before* (rather than after) finding the derivative with respect to time. Doing so will yield an incorrect result; since if those values are substituted for the variables before differentiation, those variables will become constants; and when the equation is differentiated, zeroes appear in places of all variables for which the values were plugged in.

## Example
A 10-meter ladder is leaning against the wall of a building, and the base of the ladder is sliding away from the building at a rate of 3 meters per second. How fast is the top of the ladder sliding down the wall when the base of the ladder is 6 meters from the wall?

The distance between the base of the ladder and the wall, *x*, and the height of the ladder on the wall, *y*, represent the sides of a right triangle with the ladder as the hypotenuse, *h*. The objective is to find *dy*/*dt*, the rate of change of *y* with respect to time, *t*, when *h*, *x* and *dx*/*dt*, the rate of change of *x*, are known.

Step 1:
*\(x=6\)
*\(h=10\)
*\(\frac{dx}{dt}=3\)
*\(\frac{dh}{dt}=0\)
*\(\frac{dy}{dt}=\text{?}\)

Step 2:
From the Pythagorean theorem, the equation
\(x^2+y^2=h^2,\)
describes the relationship between *x*, *y* and *h*, for a right triangle. Differentiating both sides of this equation with respect to time, *t*, yields
\(\frac{d}{dt}\left(x^2+y^2\right)=\frac{d}{dt}\left(h^2\right)\)

Step 3:
When solved for the wanted rate of change, *dy*/*dt*, gives us
\(\frac{d}{dt}\left(x^2\right)+\frac{d}{dt}\left(y^2\right)=\frac{d}{dt}\left(h^2\right)\)
\((2x)\frac{dx}{dt}+(2y)\frac{dy}{dt}=(2h)\frac{dh}{dt}\)
\(x\frac{dx}{dt}+y\frac{dy}{dt}=h\frac{dh}{dt}\)
\(\frac{dy}{dt}=\frac{h\frac{dh}{dt}-x\frac{dx}{dt{y}.\)

Step 4 & 5:
Using the variables from step 1 gives us:
\(\frac{dy}{dt}=\frac{h\frac{dh}{dt}-x\frac{dx}{dt{y}.\)
\(\frac{dy}{dt}=\frac{10\times0-6\times3}{y}=-\frac{18}{y}.\)
Solving for y using the Pythagorean Theorem gives:
\(x^2+y^2=h^2\)
\(6^2+y^2=10^2\)
\(y=8\)

Plugging in 8 for the equation:
\(-\frac{18}{y}=-\frac{18}{8}=-\frac{9}{4}\)

It is generally assumed that negative values represent the downward direction. In doing such, the top of the ladder is sliding down the wall at a rate of meters per second.

## Physics examples
Because one physical quantity often depends on another, which, in turn depends on others, such as time, related-rates methods have broad applications in Physics. This section presents an example of related rates kinematics and electromagnetic induction.

### Relative kinematics of two vehicles

For example, one can consider the kinematics problem where one vehicle is heading West toward an intersection at 80 miles per hour while another is heading North away from the intersection at 60 miles per hour. One can ask whether the vehicles are getting closer or further apart and at what rate at the moment when the North bound vehicle is 3 miles North of the intersection and the West bound vehicle is 4 miles East of the intersection.

**Big idea:** use chain rule to compute rate of change of distance between two vehicles.

**Plan:**
#Choose coordinate system
#Identify variables
#Draw picture
#Big idea: use chain rule to compute rate of change of distance between two vehicles
#Express *c* in terms of *x* and *y* via Pythagorean theorem
#Express *dc*/*dt* using chain rule in terms of *dx*/*d*t and *dy*/*dt*
#Substitute in *x*, *y*, *dx*/*dt*, *dy*/*dt*
#Simplify.

**Choose coordinate system:**
Let the *y*-axis point North and the *x*-axis point East.

**Identify variables:**
Define *y*(*t*) to be the distance of the vehicle heading North from the origin and *x*(*t*) to be the distance of the vehicle heading West from the origin.

**Express** *c* **in terms of** *x* and *y* **via the Pythagorean theorem:**

\(c = \left(x^2 + y^2\right)^{1/2}\)

**Express** *dc*/*dt* **using chain rule in terms of** *dx*/*dt* **and** *dy/dt:*

<dl>
  <dt>\(\frac{dc}{dt} = \frac{d}{dt}\left(x^2 + y^2\right)^{1/2}\)</dt>
  <dd>| Apply derivative operator to entire function</dd>
  <dt>\(= \frac{1}{2}\left(x^2 + y^2\right)^{-1/2}\frac{d}{dt}\left(x^2 + y^2\right)\)</dt>
  <dd>Square root is outside function; Sum of squares is inside function</dd>
  <dt>\(=\frac{1}{2}\left(x^2 + y^2\right)^{-1/2}\left[\frac{d}{dt}(x^2) + \frac{d}{dt}(y^2) \right]\)</dt>
  <dd>Distribute differentiation operator</dd>
  <dt>\(= \frac{1}{2}\left(x^2 + y^2\right)^{-1/2}\left[ 2x\frac{dx}{dt} + 2y\frac{dy}{dt}\right]\)</dt>
  <dd>Apply chain rule to *x*(*t*) and *y*(*t*)}</dd>
  <dt>\(= \frac{x\frac{dx}{dt} + y\frac{dy}{dt{\sqrt{x^2 + y^2}\)</dt>
  <dd>Simplify.</dd>
</dl>

**Substitute in** *x* = 4 mi, *y* = 3 mi, *dx*/*dt* = −80 mi/hr, *dy*/*dt* = 60 mi/hr **and simplify**

\(\begin{align}
\frac{dc}{dt} & = \frac{4 \text{ mi} \cdot (-80 \text{ mi}/\text{hr}) + 3 \text{ mi} \cdot (60) \text{mi}/\text{hr{\sqrt{(4 \text{ mi})^2 + (3 \text{ mi})^2\\
& = \frac{-320 \text{ mi}^2/\text{hr} + 180 \text{ mi}^2/\text{hr{5\text{ mi\\
&= \frac{-140 \text{ mi}^2/\text{hr{5\text{ mi\\
& = -28 \text{ mi}/\text{hr}
\end{align}\)

Consequently, the two vehicles are getting closer together at a rate of 28 mi/hr.

### Electromagnetic induction of conducting loop spinning in magnetic field
The magnetic flux through a loop of area *A* whose normal is at an angle *θ* to a magnetic field of strength *B* is

\(\Phi_B = B A \cos(\theta),\)

Faraday's law of electromagnetic induction states that the induced electromotive force \(\mathcal{E}\) is the negative rate of change of magnetic flux \(\Phi_B\) through a conducting loop.

\(\mathcal{E} = -\frac{d\Phi_B}{dt},\)

If the loop area *A* and magnetic field *B* are held constant, but the loop is rotated so that the angle *θ* is a known function of time, the rate of change of *θ* can be related to the rate of change of \(\Phi_B\) (and therefore the electromotive force) by taking the time derivative of the flux relation

\(\mathcal{E} = -\frac{d\Phi_B}{dt} = B A \sin\theta \frac{d\theta}{dt}\)

If for example, the loop is rotating at a constant angular velocity *ω*, so that *θ* = *ωt*, then

\(\mathcal{E}= \omega B A \sin\omega t\)
