> Content sourced from [Applied Calculus](https://www.opentextbookstore.com/appcalc/) by Calaway, Hoffman & Lippman — CC BY 3.0

# Derivative Rules


Chapter 2    The Derivative
Applied Calculus

This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 3: Power and Sum Rules for Derivatives

In the next few sections, we’ll get the derivative rules that will let us find formulas for derivatives
when our function comes to us as a formula.  These are very algebraic section, and you should get
lots of practice.  As we learn new rules, we will look at some basic applications.

Building Blocks
These are the simplest rules – rules for the basic functions.  We won’t prove these rules; we’ll just
use them.  But first, let’s look at a few so that we can see they make sense.

Example 1
Find the derivative of
( )
b
mx
x
f
y
+
=
=


This is a linear function, so its graph is its own tangent line!  The slope of the tangent line, the
derivative, is the slope of the line:
( )
m
x
f
=
'


Rule:  The derivative of a linear function is its slope

Example 2
Find the derivative of ( )
.
=
x
f


Think about this one graphically, too.  The graph of f(x) is a horizontal line.  So its slope is zero.
( )
'
=
x
f


Rule: The derivative of a constant is zero

Example 3
Find the derivative of ( )
2x
x
f
=


This question is challenging using limits, as you saw in the previous section.  We will show you
the long way to do it, then give you a shorthand rule to bypass all this.
Recall the formal definition of the derivative:
(
)
( )
h
x
f
h
x
f
x
f
h
−
+
=
→0
lim
)
('
.
Using our function ( )
2x
x
f
=
, (
) (
)
h
xh
x
h
x
h
x
f
+
+
=
+
=
+
.  Then

(
)
( )
(
)
(
)
x
h
x
h
h
x
h
h
h
xh
h
x
h
xh
x
h
x
f
h
x
f
x
f
h
h
h
h
h
lim
lim
lim
lim
lim
)
('
=
+
=
+
=
+
=
−
+
+
=
−
+
=
→
→
→
→
→


From all that, we find that
( )
f
x
x

=

Chapter 2    The Derivative
Applied Calculus

Luckily, there is a handy rule we use to skip using the limit:

Power Rule: The derivative ( )
nx
x
f
=
is
( )
−
=

n
nx
x
f


Example 4
Find the derivative of ( )
4x
x
g
=
.

Using the power rule, we know that if ( )
3x
x
f
=
, then
( )
3x
x
f
=

.  Notice that g is 4 times the
function f.  Think about what this change means to the graph of g – it’s now 4 times as tall as
the graph of f.  If we find the slope of a secant line, it will be
x
f
x
f
x
g


=


=


; each slope will
be 4 times the slope of the secant line on the f graph.  This property will hold for the slopes of
tangent lines, too:
(
)
( )
x
x
x
dx
d
x
dx
d
=

=
=


Rule:  Constants come along for the ride;
(
)
'
kf
kf
dx
d
=


Here are all the basic rules in one place.

Chapter 2    The Derivative
Applied Calculus


Derivative Rules:  Building Blocks


In what follows, f and g are differentiable functions of x.




(a) Constant Multiple Rule:
(
)
'
kf
kf
dx
d
=




(b) Sum (or Difference) Rule:
(
)
'
' g
f
g
f
dx
d
+
=
+
(or
(
)
'
' g
f
g
f
dx
d
−
=
−
)



(c)  Power Rule:
(
)
−
=
n
n
nx
x
dx
d




Special cases:
( )
=
k
dx
d
(because
kx
k =
)




( )
=
x
dx
d
(because
1x
x =
)



(d)  Exponential Functions:
( )
x
x
e
e
dx
d
=





(
)
x
x
a
a
a
dx
d

= ln





(e) Natural Logarithm:
(
)
x
x
dx
d
ln
=





The sum, difference, and constant multiple rule combined with the power rule allow us to easily
find the derivative of any polynomial.

Example 5
Find the derivative of ( )
1003
8.1
+
−
+
=
x
x
x
x
p


(
)
(
)
(
)
(
)
(
)
(
)
( )
( )
(
)
(
)
(
)
( )
8.1
8.1
1003
8.1
1003
8.1
1003
8.1
−
+
=
+
−
+
=
+
−
+
=
+
−
+
=
+
−
+
x
x
x
x
dx
d
x
dx
d
x
dx
d
x
dx
d
dx
d
x
dx
d
x
dx
d
x
dx
d
x
x
x
dx
d

Chapter 2    The Derivative
Applied Calculus

You don’t have to show every single step.  Do be careful when you’re first working with the rules,
but pretty soon you’ll be able to just write down the derivative directly:

Example 6
Find
(
)
+
−
x
x
dx
d


Writing out the rules, we'd write
(
)
)1(
)
(
−
=
+
−
=
+
−
x
x
x
x
dx
d


Once you're familiar with the rules, you can, in your head, multiply the 2 times the 17 and the
33 times 1, and just write
(
)
−
=
+
−
x
x
x
dx
d


The power rule works even if the power is negative or a fraction.  In order to apply it, first translate
all roots and basic rational expressions into exponents:

Example 7
Find the derivative of
te
t
t
y
4 +
−
=


First step – translate into exponents:
t
t
e
t
t
e
t
t
y
/
+
−
=
+
−
=
−


Now you can take the derivative:
(
)
(
)
( )
.
/
/
/
t
t
t
t
e
t
t
e
t
t
e
t
t
dt
d
e
t
t
dt
d
+
+
=
+
−
−






=
+
−
=






+
−
−
−
−
−
−


If there is a reason to, you can rewrite the answer with radicals and positive exponents:
t
t
e
t
t
e
t
t
/
+
+
=
+
+
−
−


Be careful when finding the derivatives with negative exponents.


We can immediately apply these rules to solve the problem we started the chapter with - finding a
tangent line.

Chapter 2    The Derivative
Applied Calculus

Example 8
Find the equation of the line tangent to
)
(
t
t
g
−
=
when t = 2.

The slope of the tangent line is the value of the derivative.  We can compute
t
t
g
)
(
−
=

.  To
find the slope of the tangent line when t = 3, evaluate the derivative at that point.
)
(
)
(
−
=
−
=

g
.  The slope of the tangent line is -4.

To find the equation of the tangent line, we also need a point on the tangent line.  Since the
tangent line touches the original function at t = 2, we can find the point by evaluating the
original function:
(2)
g
=
−
=
.  The tangent line must pass through the point (2, 6).

Using the point-slope equation of a line, the tangent line will
have equation
)
(
−
−
=
−
t
y
.
Simplifying to slope-intercept form, the equation is
4 +
−
=
t
y
.

Graphing, we can verify this line is indeed tangent to the
curve.





We can also use these rules to help us find the derivatives we need to interpret the behavior of a
function.

Example 9
In a memory experiment, a researcher asks the subject to memorize as many words from a list
as possible in 10 seconds.  Recall is tested, then the subject is given 10 more seconds to study,
and so on.  Suppose the number of words remembered after t seconds of studying could be
modeled by
2/5
( )
W t
t
=
.  Find and interpret
(20)
W
.

3/5
3/5
( )
4 5
W t
t
t
−
−

=

=
, so
(
)
3/5
(20)
0.2652
W
−

=



Since W is measured in words, and t is in seconds, W' has units words per second.
(20)
0.2652
W

means that after 20 seconds of studying, the subject is learning about 0.27
more words for each additional second of studying.

Chapter 2    The Derivative
Applied Calculus

Business and Economics
Next we will delve more deeply into some business applications.  To do that, we first need to
review some terminology.

Suppose you are producing and selling some item.  The profit you make is the amount of money
you take in minus what you have to pay to produce the items.  Both of these quantities depend on
how many you make and sell.  (So we have functions here.)  Here is a list of definitions for some
of the terminology, together with their meaning in algebraic terms and in graphical terms.

Your cost is the money you have to spend to produce your items.

The Fixed Cost (FC) is the amount of money you have to spend regardless of how many items
you produce.  FC can include things like rent, purchase costs of machinery, and salaries for office
staff.  You have to pay the fixed costs even if you don’t produce anything.

The Total Variable Cost (TVC) for q items is the amount of money you spend to actually
produce them.  TVC includes things like the materials you use, the electricity to run the machinery,
gasoline for your delivery vans, maybe the wages of your production workers.  These costs will
vary according to how many items you produce.

The Total Cost (TC, or sometimes just C) for q items is the total cost of producing them.  It’s the
sum of the fixed cost and the total variable cost for producing q items.

Why is it OK that are there two definitions for Marginal Cost (and Marginal Revenue, and
Marginal Profit)?

We have been using slopes of secant lines over tiny intervals to approximate derivatives.  In this
example, we’ll turn that around – we’ll use the derivative to approximate the slope of the secant
line.
Notice that the “cost of the next item” definition is actually the slope of a secant line, over an
interval of 1 unit:
( )
(
)
(
)
−
+
=
−
+
=
q
C
q
C
q
MC

So this is approximately the same as the derivative of the cost function at q:
( )
( )
q
C
q
MC
'
=

In practice, these two numbers are so close that there’s no practical reason to make a distinction.
For our purposes, the marginal cost is the derivative is the cost of the next item.
The Marginal Cost (MC) at q items is the cost of producing the next item.  Really, it’s

MC(q) = TC(q + 1) – TC(q).
In many cases, though, it’s easier to approximate this difference using calculus (see Example
below).  And some sources define the marginal cost directly as the derivative,

MC(q) = TC'(q).
In this course, we will use both of these definitions as if they were interchangeable.

The units on marginal cost is cost per item.

Chapter 2    The Derivative
Applied Calculus

Example 10
The table shows the total cost (TC) of producing q items.
a) What is the fixed cost?
b) When 200 items are made, what is the total variable cost?
The average variable cost?
c) When 200 items are made, estimate the marginal cost.

a) The fixed cost is $20,000, the cost even when no items are made.

b) When 200 items are made, the total cost is $45,000.  Subtracting the fixed cost, the total
variable cost is $45,000 - $20,000 = $25,000.

The average variable cost is the total variable cost divided by the number of items, so we would
divide the $25,000 total variable cost by the 200 items made.  $25,000 ÷ 200= $125.  On
average, each item had a variable cost of $125.

c) We need to estimate the value of the derivative, or the slope of the tangent line at q = 200.
Finding the secant line from q=100 to q=200 gives a slope of
,
,
=
−
−
.  Finding
the secant line from q=200 to q=300 gives a slope of
,
,
=
−
−
.   We could estimate
the tangent slope by averaging these secant slopes, giving us an estimate of $90/item.

This tells us that after 200 items have been made, it will cost about $90 to make one more item.


Example 11
The cost to produce x items is
x  hundred dollars.

(a)  What is the cost for producing 100 items?  101 items?   What is cost of the 101st item?

(b)  For   C(x) = x   , calculate  C'(x)  and evaluate  C'  at  x = 100.  How does C '(100)
compare with the last answer in part (a)?

(a) Put  C(x) = x   = x1/2 hundred dollars,  the cost for  x  items.  Then C (100) = $1000 and
C(101) = $1004.99, so it costs  $4.99  for that 101st item.  Using this definition, the marginal
cost is $4.99.

(b)
1/2
( )
C x
x
x
−

=
=
so
(100)
2 100
C
=
=
hundred dollars  =  $5.00.

Note how close these answers are!  This shows (again) why it’s OK that we use both definitions
for marginal cost.

Items, q
Total Cost, TC
$20,000
$35,000
$45,000
$53,000

Chapter 2    The Derivative
Applied Calculus

Demand is the functional relationship between the price p and the quantity q that can be sold (that
is demanded).  Depending on your situation, you might think of p as a function of q, or of q as a
function of p.

Your revenue is the amount of money you actually take in from selling your products.  Revenue is
price × quantity.

The Total Revenue (TR, or just R) for q items is the total amount of money you take in for selling
q items.









The Profit (P) for q items is TR(q) – TC(q), the difference between total revenue and total costs

The average profit for q items is P/q.  The marginal profit at q items is P(q + 1) – P(q), or
( )
q
P


Graphical Interpretations of the Basic Business Math Terms

Illustration/Example:
Here are the graphs of TR and TC for producing and selling a certain item.  The horizontal axis is
the number of items, in thousands.  The vertical axis is the number of dollars, also in thousands.



First, notice how to find the fixed cost and variable cost from the graph here.  FC is the y-
intercept of the TC graph.  (FC = TC(0).)  The graph of TVC would have the same shape as the
graph of TC, shifted down.  (TVC = TC – FC.)
The Marginal Revenue (MR) at q items is the cost of producing the next item,

MR(q) = TR(q + 1) – TR(q).
Just as with marginal cost, we will use both this definition and the derivative definition

MR(q) = TR’(q).

Your profit is what’s left over from total revenue after costs have been subtracted.

Chapter 2    The Derivative
Applied Calculus

MC(q) = TC(q + 1) – TC(q), but that’s impossible to read on this graph.  How could you
distinguish between TC(4022) and TC(4023)?  On this graph, that interval is too small to see, and
our best guess at the secant line is actually the tangent line to the TC curve at that point.  (This is
the reason we want to have the derivative definition handy.)

MC(q) is the slope of the tangent line to the TC curve at (q, TC(q)).
MR(q) is the slope of the tangent line to the TR curve at (q, TR(q)).

Profit is the distance between the TR and TC curve.  If you experiment with your clear plastic
ruler, you’ll see that the biggest profit occurs exactly when the tangent lines to the TR and TC
curves are parallel.  This is the rule “profit is maximized when MR = MC.” which we'll explore
later in the chapter.

Example 12
The demand, D, for a product at a price of p dollars is given by
( )
0.2
D p
p
=
−
.  Find the
marginal revenue when the price is $10.

First we need to form a revenue equation.  Since Revenue = Price×Quantity, and the demand
equation shows the quantity of product that can be sold, we have
(
)
( )
( )
0.2
0.2
R p
D p
p
p
p
p
p
=

=
−
=
−


Now we can find marginal revenue by finding the derivative
( )
(
)
( )
200 1

## 0.2 3
0.6
R p
p
p

=
−
=
−


At a price of $10,
(
)
(10)

## 0.6 10
R
=
−
=
.

Notice the units for R' are dollars of Revenue
dollar of price
, so
(10)
R
=
means that when the price is
$10, the revenue will increase by $140 for each dollar the price was increased.



## 2.3 Exercises

1. Fill in the values in the table for
( )
(
)
x
f
dx
d 3
,
( )
( )
(
)
x
g
x
f
dx
d
+
, and
( )
( )
(
)
x
f
x
g
dx
d
−
.

x
f(x) f '(x) g(x) g '(x)
( )
(
)
x
f
dx
d 3

( )
( )
(
)
x
g
x
f
dx
d
+

( )
( )
(
)
x
f
x
g
dx
d
−











–2
–4

–1

Chapter 2    The Derivative
Applied Calculus

2. Find  (a)  D( x12 )      (b)  d
dx ( 7 x  )
(c)  D( 1
x3  )
(d)  d xe
dx

3. Find  (a)  D( x9  )
(b)  d x2/3
dx
(c)  D( 1
x4  )    (d)  D( xπ  )

In problems  4 – 8,  (a) calculate  f '(1)  and  (b) determine when  f '(x) = 0.

4. f(x) =  x2  – 5x + 13

5. f(x) = 5x2  – 40x + 73
6. f(x) =  x3  + 9x2  + 6

7. f(x) =  x3 + 3x2 + 3x  – 1
8. f(x) =  x3  + 2x2  + 2x  – 1


9. Where do f(x) = x2 – 10x + 3  and  g(x) = x3 – 12x have horizontal tangent lines ?

10. It takes  T(x) = x2  hours to weave  x  small rugs.  What is the marginal production time to
weave a rug?  (Be sure to include the units with your answer.)

11. It costs  C(x) = x   dollars to produce  x  golf balls.  What is the marginal production cost to
make a golf ball?  What is the marginal production cost when  x = 25?  when x= 100?  (Include
units.)

12. An arrow shot straight up from ground level with an initial velocity of  128 feet per second will
be at height  h(x) = –16x2 + 128x  feet at  x  seconds.

(a) Determine the velocity of the arrow when  x = 0, 1 and 2
seconds.

(b) What is the velocity of the arrow, v(x), at any time  x?

(c) At what time  x  will the velocity of the arrow be  0?

(d) What is the greatest height the arrow reaches?

(e) How long will the arrow be aloft?

(f) Use the answer for the velocity in part (b) to determine the


acceleration, a(x) = v '(x), at any time  x.

13. If an arrow is shot straight up from ground level on the moon with an initial velocity of 128
feet per second, its height will be  h(x) = –2.65x2 + 128x  feet at  x  seconds.  Do parts (a) – (e)
of problem 40 using this new equation for  h.

14. f(x) = x3  +  A x2  + B x  + C  with constants  A, B  and  C.  Can you find conditions on the

constants  A, B  and  C which will guarantee that the graph of  y = f(x)  has two distinct
"vertices"? (Here a "vertex" means a place where the curve changes from increasing to
decreasing or from decreasing to increasing.)

Chapter 2    The Derivative
Applied Calculus

This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 4: Product and Quotient Rules

The basic rules will let us tackle simple functions.  But what happens if we need the derivative of a
combination of these functions?

Example 1
Find the derivative of ( ) (
)(
)
h x
x
x
=
−
+


This function is not a simple sum or difference of polynomials.  It’s a product of polynomials.
We can simply multiply it out to find its derivative:
( ) (
)(
)
( )
'
11 36
h x
x
x
x
x
x
h
x
x
x
=
−
+
=
−
+
−
=
−
+



Now suppose we wanted to find the derivative of
( ) (
)(
)3
.7
5.1
+
+
−
−
−
+
=
x
x
x
x
x
x
x
f


This function is not a simple sum or difference of polynomials.  It’s a product of polynomials.  We
could simply multiply it out to find its derivative as before – who wants to volunteer?  Nobody?

We’ll need a rule for finding the derivative of a product so we don’t have to multiply everything
out.

It would be great if we can just take the derivatives of the factors and multiply them, but
unfortunately that won’t give the right answer.  to see that, consider finding derivative of
( ) (
)(
)3
+
−
=
x
x
x
g
.  We already worked out the derivative.  It’s
( )
'
x
x
x
g
+
−
=
.
What if we try differentiating the factors and multiplying them?  We’d get (
)( )
x
x
=
, which
is totally different from the correct answer.

The rules for finding derivatives of products and quotients are a little complicated, but they save us
the much more complicated algebra we might face if we were to try to multiply things out.  They
also let us deal with products where the factors are not polynomials.  We can use these rules,
together with the basic rules, to find derivatives of many complicated looking functions.

Chapter 2    The Derivative
Applied Calculus


Derivative Rules:  Product and Quotient Rules


In what follows, f and g are differentiable functions of x.




(a) Product Rule:




(
)
'
'
fg
g
f
fg
dx
d
+
=






The derivative of the first factor times the second left alone, plus the first left alone times
the derivative of the second.







The product rule can extend to a product of several functions; the pattern continues – take
the derivative of each factor in turn, multiplied by all the other factors left alone, and add
them up.



(b) Quotient Rule:



'
'
g
fg
g
f
g
f
dx
d
−
=












The numerator of the result resembles the product rule, but there is a minus instead of a
plus; the minus sign goes with the g’.  The denominator is simply the square of the original
denominator – no derivatives there.



Example 2
Find the derivative of ( ) (
)(
)
h x
x
x
=
−
+


This is the same function we found the derivative of in Example 1, but let's use the product rule
and check to see if we get the same answer.  For this first example, we will provide a lot more
detail and steps than one usually actually shows when working a problem like this.

Notice we can think of h(x) as the product of two functions
( )
f x
x
=
−
and ( )
g x
x
=
+ .
Finding the derivative of each of these,
( )
f
x
x

=

( )
g x

=

Using the product rule,
( )
(
)( ) ( )(
) (
)(
) (
)( )
11 1
h x
f
g
f
g
x
x
x



=
+
=
+
+
−


To check if this is equivalent to the answer we found in Example 1 we could simplify:
( ) (
)(
) (
)( )
11 1
11 16
h x
x
x
x
x
x
x
x
x

=
+
+
−
=
+
+
−
=
+
−


From this, we can see the answers are equivalent.

Chapter 2    The Derivative
Applied Calculus


Example 3
Find the derivative of ( )
t
e
t
F
t ln
=


This is a product, so we need to use the product rule.  I like to put down empty parentheses to
remind myself of the pattern; that way I don’t forget anything.

( ) ( )( )
( )( )
+
=
t
F '

Then I fill in the parentheses – the first set gets the derivative of
te , the second gets
t
ln left
alone, the third gets
te left alone, and the fourth gets the derivative of
t
ln .
( ) ( )(
) ( )
t
e
t
e
t
e
t
e
t
F
t
t
t
t
+
=






+
=
ln
ln
'




Notice that this was one we couldn’t have done by “multiplying out” before taking the derivative.


Example 4
Find the derivative of
x
x
y
x
+
+
=


This is a quotient, so we need to use the quotient rule.  Again, you find it helpful to put down
the empty parentheses as a template:
( )( ) ( )( )
( )
'
−
=
y


Then fill in all the pieces:
(
)(
) (
)(
)
(
)
ln 4
3 16
'
3 16
x
x
x
x
x
x
y
x
+

+
−
+
=
+


Now for goodness’ sakes don’t try to simplify that!  Remember that “simple” depends on what you
will do next; in this case, we were asked to find the derivative, and we’ve done that.  Please STOP,
unless there is a reason to simplify further.


Example 5
Suppose a large tank contains 8 kg of a chemical dissolved in 50 liters of water.  If a tap is
opened and water is added to the tank at a rate of 5 liters per minute, at what rate is the
concentration of chemical in the tank changing after 4 minutes?

Chapter 2    The Derivative
Applied Calculus

First we need to set up a model for the concentration of chemical.  The concentration would be
measured as kg of chemical per liter of water, kg
L .  The number of kg of chemical stays
constant at 8 kg, but the quantity of water in the tank is increasing by 5 L/min.  The total
volume of water in the tank after t minutes is 50 + 5t, so the concentration after t minutes is
( )
c t
t
=
+


To find the rate at which the concentration is changing, we need the derivative:
( ) (
) ( )
(
)
(
)
( )
d
d
t
t
dt
dt
c t
t

+
−
+

=
+


( ) (
) ( )( )
(
)
(
)
t
t
t

+
−
−
=
=
+
+


At t = 4,
(
)
(4)
0.00816
5 4
c
−

=
−
+ 


Note that the units here are kg per liter, per minute, or
/
min
kg L .  In other words, this tells us that
after 4 minutes, the concentration of chemical is decreasing by 0.00816 kg/L each minute.


Returning to our discussion of business and economics topics, in addition to total cost and
marginal cost, we often also want to talk about average cost or average revenue.

The Average Cost (AC) for q items is the total cost divided by q, or
( )
( )
TC q
AC q
q
=
.  You can
also talk about the average fixed cost, FC/q, or the average variable cost, TVC/q.

The Average Revenue (AR) for q items is the total revenue divided by q, or TR/q.
We already know that we can find average rates of change by finding slopes of secant lines. AC,
AR, MC, and MR are all rates of change, and we can find them with slopes, too.

AC(q) is the slope of a diagonal line, from (0, 0) to (q, TC(q)).
AR(q) is the slope of the line from (0, 0) to (q, TR(q)).

Chapter 2    The Derivative
Applied Calculus



Just as we found marginal Total Cost, we can also find marginal Average Cost.

Example 6
The cost, in thousands of dollars, for producing x thousand cellphone cases is given by
( )
0.004
C x
x
x
=
+
−
.  Find
a) The Fixed costs
b) The Average Cost when 5 thousand, 10 thousand, or 20 thousand cases are produced
c) The Marginal Average Cost when 5 thousand cases are produced

a) The fixed costs are the costs when no items are produced:
(0)
C
=
thousand dollars

b) The average cost function is total cost divided by number of items, so
( )
0.004
( )
C x
x
x
AC x
x
x
+
−
=
=


Note the units are thousands of dollars per thousands of items, which simplifies to just dollars
per item.
At a production of 5 thousand items:
0.004(5)
(5)
5.38
AC
+ −
=
=
dollars per item
At a production of 10 thousand items:
22 10
0.004(10)
(10)
3.16
AC
+
−
=
=
dollars per item
At a production of 20 thousand items:
0.004(20)
(20)
2.02
AC
+
−
=
=
dollars per item
Notice that while the total cost increases with production, the average cost per item decreases,
because the initial fixed costs are being distributed across more items.

slope = AR
slope = AC

Chapter 2    The Derivative
Applied Calculus

c) For the marginal average cost, we need to find the derivative of the average cost function.
We can either calculate this using the quotient rule, or we could use algebra to simplify the
equation first:
0.004
0.004
( )
1 0.004
1 0.004
x
x
x
x
AC x
x
x
x
x
x
x
x
x
−
+
−
=
=
+
−
=
+ −
=
+ −
.

Taking the derivative,
( )
0.004
0.004
AC x
x
x
−

= −
−
= −
−


When 5 thousand items are produced,
(5)
0.004
0.884
AC
= −
−
= −


Since the units on AC are dollars per item, and the units on x are thousands of items, the units
on AC' dollars per item per thousands of items.  This tells us that when 5 thousand items are
produced, the average cost per item is decreasing by $0.884 for each additional thousand items
produced.


## 2.4 Exercises

1. Use the values in the table to fill in the rest of the table.

x
f(x) f '(x)
g(x) g '(x)
( )
( )
(
)
x
g
x
f
dx
d


( )
( )





x
g
x
f
dx
d

( )
( )





x
f
x
g
dx
d











–2
–4

–1

2.   Use the information in the graph  to plot the
values of the functions  f + g, f.g and  f/g  and
their derivatives  at  x = 1, 2  and  3 .

3. Use the information in the graph  to plot the
values of the functions  2f, f – g  and  g/f  and
their derivatives at  x = 1, 2  and  3 .

4. Calculate
(
)(
)
(
)
+
−
x
x
dx
d
by  (a)  using the product rule  and  (b)  expanding the product
and then differentiating.  Verify that both methods give the same result.

Chapter 2    The Derivative
Applied Calculus

5. If the product of f and  g  is a constant  (  f(x) ∙ g(x) = k  for all  x), then how are
( )
(
)
( )
x
f
x
f
dx
d


and
( )
(
)
( )
x
g
x
g
dx
d
related?

6. If the quotient of f  and  g  is a constant  (
( )
( )
k
x
g
x
f
=
for all x), then how are  g . f '  and  f . g '
related?

In problems  7 – 8,  (a) calculate  f '(1)  and  (b) determine when  f '(x) = 0.
7. f(x) =
7x
x2  + 4
8.
f(x) =
x
x −


9. Determine
(
)(
)
−
+
x
x
dx
d
and   d
dt (  3t – 2
5t + 1  )    .
10. Find  (a)
(
)
xe
x
dx
d
and  (b)
( )
x
e
dx
d
.
11. Find (a)
(
)
t
te
dt
d
, (b) ( )
xe
d


12. A manufacturer has determined that an employee with d days of production experience will
be able to

produce approximately   P(d) = 3 + 15( 1 – e–0.2d )  items per day.  Graph  P(d).

(a) Approximately how many items will a beginning employee be able to produce each day?

(b) How many items will an experienced employee be able to produce each day?
(c) What is the marginal production rate of an employee with 5 days of experience?  (What
are the units of your answer, and what does this answer mean?)

Chapter 2    The Derivative
Applied Calculus

This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 5: Chain Rule

There is one more type of complicated function that we will want to know how to differentiate:
composition.  The Chain Rule will let us find the derivative of a composition.  (This is the last
derivative rule we will learn!)

Example 1
Find the derivative of
(
)
x
x
y
+
=
.

This is not a simple polynomial, so we can’t use the basic building block rules yet.  It is a
product, so we could write it as
(
)
(
)(
)x
x
x
x
x
x
y
+
+
=
+
=
and use the product
rule.  Or we could multiply it out and simply differentiate the resulting polynomial.  I’ll do it the
second way:
(
)
x
x
x
y
x
x
x
x
x
y
'
+
+
=
+
+
=
+
=



Now suppose we want to find the derivative of
(
)
x
x
y
+
=
.  We could write it as a product with
20 factors and use the product rule, or we could multiply it out.  But I don’t want to do that, do you?

We need an easier way, a rule that will handle a composition like this.  The Chain Rule is a little
complicated, but it saves us the much more complicated algebra of multiplying something like this out.
It will also handle compositions where it wouldn’t be possible to “multiply it out.”

The Chain Rule is the most common place for students to make mistakes.  Part of the reason is that the
notation takes a little getting used to.  And part of the reason is that students often forget to use it when
they should.  When should you use the Chain Rule?  Almost every time you take a derivative.

Chapter 2    The Derivative
Applied Calculus


Derivative Rules:  Chain Rule


In what follows, f and g are differentiable functions with
( )
u
f
y =
and
( )
x
g
u =




(h) Chain Rule (Leibniz notation):



dx
du
du
dy
dx
dy

=



Notice that the du’s seem to cancel.  This is one advantage of the Leibniz notation; it can
remind you of how the chain rule chains together.



(h)  Chain Rule (using prime notation):



( )
( )
( )
( )
(
)
( )
x
g
x
g
f
x
g
u
f
x
f
'
'
'
'
'

=

=




(h) Chain Rule (in words):


The derivative of a composition is the derivative of the outside, with the inside staying the
same, TIMES the derivative of what’s inside.



I recite the version in words each time I take a derivative, especially if the function is complicated.




Example 2
Find the derivative of
(
)
x
x
y
+
=
.

This is the same one we did before by multiplying out.  This time, let’s use the Chain Rule:  The
inside function is what appears inside the parentheses:
x
x
3 +
.  The outside function is the
first thing we find as we come in from the outside – it’s the square function, (inside)2.

The derivative of this outside function is (2*inside).  Now using the chain rule, the derivative of
our original function is:
(2*inside) TIMES the derivative of what’s inside (which is
2 +
x
):
(
)
(
) (
)
'
+

+
=
+
=
x
x
x
y
x
x
y



If you multiply this out, you get the same answer we got before.  Hurray!  Algebra works!

Chapter 2    The Derivative
Applied Calculus

Example 3
Find the derivative of
(
)
x
x
y
+
=


Now we have a way to handle this one.  It’s the derivative of the outside TIMES the derivative
of what’s inside.

The outside function is (inside)20, which has the derivative 20(inside)19.
(
)
(
)
(
)
'
+

+
=
+
=
x
x
x
y
x
x
y



Example 4
Differentiate
2 +
xe
.

This isn’t a simple exponential function; it’s a composition.  Typical calculator or computer
syntax can help you see what the “inside” function is here.  On a TI calculator, for example,
when you push the
x
e key, it opens up parentheses:
(
^
e
This tells you that the “inside” of the
exponential function is the exponent.  Here, the inside is the exponent
2 +
x
.  Now we can use
the Chain Rule:  We want the derivative of the outside TIMES the derivative of what’s inside.
The outside is the “e to the something” function, so its derivative is the same thing.  The
derivative of what’s inside is 2x.  So
(
) (
) (
)
x
e
e
dx
d
x
x

=
+
+



Example 5
The table gives values for  f , f ' , g  and g '  at a number of points.  Use these values to
determine   ( fg )(x)  and  ( fg ) '(x)   at  x = –1  and  0.



( fg )(–1)  =  f( g(–1) )  =  f( 3 )  =  0
( fg )(0)  =  f( g(0) )  =  f( 1 )  =  1.
( fg ) '(–1)  =  f '( g(–1) ).g '( –1 )  =  f '( 3 ).(0)  =  (2)(0)  = 0   and
( fg ) '( 0 )  =  f '( g( 0 ) ).g '( 0 )  =  f '( 1 ).( 2 )  =  (–1)(2)  = –2 .

x
f(x)
g(x)
f'(x)
g'(x)
( fg )(x)
( gf )(x)
-1


-1


-1


-1


-1

Chapter 2    The Derivative
Applied Calculus

Example 6
If 2400 people now have a disease, and the number of people with the disease appears to double
every 3 years, then the number of people expected to have the disease in t years is
/3
2400 2
t
y =


(a)  How many people are expected to have the disease in 2 years?
(b)  When are  50,000 people expected to have the disease?
(c)  How fast is the number of people with the disease expected to grow now and 2 years from
now?

(a)  In 2 years,  y = 2400 . 22/3 ≈ 3,810 people.

(b) We know  y = 50,000 , and we need to solve  50,000 = 2400 . 2t/3  for  t.  We could start by
isolating the exponential by dividing both sides by 2400,
/3
50000
2400
t
=


Taking the logarithm of both sides,
(
)
/3
50000
ln
ln 2
2400
t

=





Using the exponent property for logs,
( )
50000
ln
ln 2
2400
t

=





Solving for t,
50000
3ln
2400
13.14
ln(2)
t






=

years
We expect 50,000 people to have the disease about 13.14  years from now.

(c) This is asking for  dy/dt  when  t = 0 and 2 years.  Using the chain rule,
(
)
/3
/3
/3
2400 2
2400 2
ln(2)

## 554.5 2
t
t
t
dy
d
dt
dt
=

=






Now, at t = 0, the rate of growth of the disease is approximately  554.5.20 ≈ 554.5 people/year.
In 2 years the rate of growth will be approximately  554.5 . 22/3 ≈ 880  people/year.


Derivatives of Complicated Functions

You’re now ready to take the derivative of some mighty complicated functions.  But how do you
tell what rule applies first?  Come in from the outside – what do you encounter first?  That’s the
first rule you need.  Use the Product, Quotient, and Chain Rules to peel off the layers, one at a
time, until you’re all the way inside.

Chapter 2    The Derivative
Applied Calculus

Example 7
Find
(
)
(
)
ln
+

x
e
dx
d
x


Coming in from the outside, I see that this is a product of two (complicated) functions.  So I’ll
need the Product Rule first.  I’ll fill in the pieces I know, and then I can figure the rest as
separate steps and substitute in at the end:

(
)
(
)
(
)
(
)
(
) (
)
(
)
(
)





+
+
+






=
+

ln
ln
ln
x
dx
d
e
x
e
dx
d
x
e
dx
d
x
x
x


Now as separate steps, I’ll find
(
)
x
x
e
e
dx
d
=
(using the Chain Rule) and

(
)
(
)
ln

+
=
+
x
x
dx
d
(also using the Chain Rule).

Finally, to substitute these in their places:
(
)
(
) (
)
(
)
(
) (
)







+
+
+
=
+

ln
ln
x
e
x
e
x
e
dx
d
x
x
x


(And please don’t try to simplify that!)


Example 8
Differentiate
(
)






−
=
t
e
t
z
t


Don’t panic!  As you come in from the outside, what’s the first thing you encounter?  It’s that
4th power.  That tells you that this is a composition, a (complicated) function raised to the 4th
power.

Step One:  Use the Chain Rule.  The derivative of the outside TIMES the derivative of what’s
inside.

(
)
(
)
(
)





−







−
=






−
=
t
e
t
dt
d
t
e
t
t
e
t
dt
d
dt
dz
t
t
t

Chapter 2    The Derivative
Applied Calculus

Now we’re one step inside, and we can concentrate on just the
(
)





−1
3 3
t
e
t
dt
d
t
part.  Now, as you
come in from the outside, the first thing you encounter is a quotient – this is the quotient of two
(complicated) functions.

Step Two: Use the Quotient Rule.  The derivative of the numerator is straightforward, so we
can just calculate it. The derivative of the denominator is a bit trickier, so we'll leave it for now.
(
)
(
)
(
)
(
) (
)
(
)
(
)
(
)
(
)
−






−
−
−
=






−
t
e
t
e
dt
d
t
t
e
t
t
e
t
dt
d
t
t
t
t


Now we’ve gone one more step inside, and we can concentrate on just the
(
)
(
)
−
t
e
dt
d
t
part.
Now we have a product.

Step Three: Use the Product Rule:

(
)
(
) ( )(
) ( )( )1
t
t
t
e
t
e
t
e
dt
d
+
−
=
−


And now we’re all the way in – no more derivatives to take.

Step Four:  Now it’s just a question of substituting back – be careful now!
(
)
(
) ( )(
) ( )( )1
t
t
t
e
t
e
t
e
dt
d
+
−
=
−
, so

(
)
(
)
(
)
(
) (
)( )(
) ( )( )
(
)
(
)
(
)
−
+
−
−
−
=






−
t
e
e
t
e
t
t
e
t
t
e
t
dt
d
t
t
t
t
t
, so


(
)
(
)
(
)
(
)
(
) (
)( )(
) ( )( )
(
)
(
)
(
)








−
+
−
−
−







−
=






−
=
t
e
e
t
e
t
t
e
t
t
e
t
t
e
t
dt
d
dt
dz
t
t
t
t
t
t
.

Phew!

Chapter 2    The Derivative
Applied Calculus

What if the Derivative Doesn’t Exist?
A function is called differentiable at a point if its derivative exists at that point.

We’ve been acting as if derivatives exist everywhere for every function.  This is true for most of
the functions that you will run into in this class.  But there are some common places where the
derivative doesn’t exist.

Remember that the derivative is the slope of the tangent line to the curve.  That’s what to think
about.

Where can a slope not exist?  If the tangent line is vertical, the derivative will not exist.

Example 9
Show that
/
)
(
x
x
x
f
=
=
is not differentiable at x = 0.

Finding the derivative,
/
/
)
(
x
x
x
f
=
=
−
.  At x = 0, this function is undefined.  From the
graph, we can see that the tangent line to this curve at x = 0 is vertical with undefined slope,
which is why the derivative does not exist at x = 0.



Where can a tangent line not exist?  If there is a sharp corner (cusp) in the graph, the derivative
will not exist at that point because there is no well-defined tangent line (a teetering tangent, if you
will).  If there is a jump in the graph, the tangent line will be different on either side and the
derivative can’t exist.

Example 10
Show that
x
x
f
=
)
(
is not differentiable at x = 0.

On the left side of the graph, the slope of the line is -1.  On the
right side of the graph, the slope is +1.  There is no well-
defined tangent line at the sharp corner at x = 0, so the function
is not differentiable at that point.

Chapter 2    The Derivative
Applied Calculus


## 2.5 Exercises

1. The graph of  y = f(x)  is shown.

(a)  At which integers is  f  continuous?

(b)  At which integers is  f  differentiable?

2. The graph of  y = g(x)  is shown.

(a)  At which integers is  g  continuous?

(b)  At which integers is  g  differentiable?


Problems 3 and 4 refer to the values given in this table:

x
f(x)
g(x)
f '(x)
g '(x)
( fg )(x)
( fg )' (x)






–2
–1

–1

–2
–1

–2
–1

–1

3. Use the table of values to determine  ( fg )(x)  and  ( fg )' (x)  at  x = 1  and  2.

4. Use the table of values to determine  ( fg )(x)  and  ( fg )' (x)  at  x = –2, –1  and  0.

5. Use the graphs  to estimate the values of  g(x),
g '(x), (fg)(x), f '( g(x) ), and ( fg ) '( x )  at
x = 1.

6. Use the graphs  to estimate the values of  g(x),
g '(x), (fg)(x), f '( g(x) ), and ( fg ) '( x )  for
x = 2.

In problems  7 –  12 , find the derivative of each function.

7. f(x) = (2x – 8)5
8. f(x) = (6x – x2)10
9. f(x) = x .(3x + 7)5


10. f(x) = (2x + 3)6.(x – 2)4
11. f(x) =
x2 + 6x – 1
12. f(x)  =   x – 5
(x + 3)4

13. If  f  is a differentiable function,

(a)  how are the graphs of  y = f(x)  and   y = f(x) + k  related?

(b)  how are the derivatives of  f(x)  and  f(x) + k  related?