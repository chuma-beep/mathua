> Content sourced from [Applied Calculus](https://www.opentextbookstore.com/appcalc/) by Calaway, Hoffman & Lippman — CC BY 3.0

# Implicit Differentiation and Related Rates


Chapter 2    The Derivative
Applied Calculus

Example 2
Find the slope of the tangent line to the circle  x2 + y2 = 25  at the point  (3,4)  using implicit
differentiation.

We differentiate each side of the equation  x2 + y2 = 25 and then
solve for  y'
(
)
)
(
dx
d
y
x
dx
d
=
+

=

+
y
y
x

Solving for  y', we have
y
x
y
x
y
−
=
−
=

, and, at the point  (3,4),
y' = – 3/4.


In the previous example, it would have been easy to explicitly solve for  y,  and then we could
differentiate  y  to  get  y '.  Because we could explicitly solve for  y ,  we had a choice of methods
for calculating   y '.  Sometimes, however, we can not explicitly solve for  y , and the only way of
determining  y '  is implicit differentiation.


Related Rates
If several variables or quantities are related to each other and some of the variables are changing at
a known rate, then we can use derivatives to determine how rapidly the other variables must be
changing.


Example 3
Suppose the border of a town is roughly circular, and the radius of that circle has been
increasing at a rate of 0.1 miles each year.  Find how fast the area of the town has been
increasing when the radius is 5 miles.

We could get an approximate answer by calculating the area of the
circle when the radius is 5 miles  ( A = πr2 = π(5 miles)2 ≈ 78.6
miles2 )  and  1 year later when the radius is 0.1 feet larger than
before  (  A = πr2 = π(5.1 miles)2 ≈ 81.7 miles2 ) and then finding
∆Area/∆time = (81.7 mi2 – 78.6 mi2)/(1 year) = 3.1 mi2/yr.  This
approximate answer represents the average change in area during
the 1 year period when the radius increased from 5 miles to 5.1
miles, and would correspond to the secant slope on the area graph.

To find the exact answer, though, we need derivatives.  In this case both radius and area are
functions of time:


r(t) = radius at time t
A(t) = area at time t