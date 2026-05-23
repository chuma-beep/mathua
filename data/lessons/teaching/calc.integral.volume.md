> Content sourced from [Applied Calculus](https://www.opentextbookstore.com/appcalc/) by Calaway, Hoffman & Lippman — CC BY 3.0

# Integrals


Chapter 3    The Integral
Applied Calculus


Example 4
Find ∫
dx
e x
.

This is likely one you remember --
x
e is its own derivative, so it is also its own antiderivative.
The integral sign tells me that I need to include the entire family of functions, so I need that + C
on the end:
C
e
dx
e
x
x
+
=
∫


Antiderivatives Graphically or Numerically
Another way to think about the Fundamental Theorem of Calculus is to solve the expression for
F(b):

The Fundamental Theorem of Calculus (restated)



( )
( )
( )
a
F
b
F
dx
x
F
b
a
−
=
∫
'

The definite integral of a derivative from a to b gives the net change in the original function.



( )
( )
( )
∫
+
=
b
a
dx
x
F
a
F
b
F
'

The amount we end up is the amount we start with plus the net change in the function.


This lets us get values for the antiderivative – as long as we have a starting point, and we know
something about the area.

Example 5
Suppose F(t) has the derivative f(t) shown below, and suppose that we know F(0) = 5.  Find
values for F(1), F(2), F(3), and F(4).


Using the second way to think about the Fundamental Theorem of Calculus,
( )
( )
( )
∫
+
=
b
a
dx
x
F
a
F
b
F
'
-- we can see that
( )
( )
( )
∫
+
=
dx
x
f
F
F
.  We know the value of F(0), and we can easily find
( )
∫
dx
x
f
from the
graph – it’s just the area of a triangle.
So
( )
( )
( )
5.5
5.
=
+
=
+
=
∫
dx
x
f
F
F

Chapter 3    The Integral
Applied Calculus


( )
( )
( )
=
+
=
+
=
∫
dx
x
f
F
F

Note that we can start from any place we know the value of – now that we know F(2), we can
use that:
( )
( )
( )
5.5
5.
=
−
=
+
=
∫
dx
x
f
F
F

( )
( )
( )
5.4
5.5
=
−
=
+
=
∫
dx
x
f
F
F



Example 6
F ‘(t) = f(t) is shown below.  Where does F(t) have maximum and minimum values on the
interval [0, 4]?


Since ( )
( )
( )
∫
+
=
b
a
dt
t
f
a
F
b
F
, we know that F is increasing as long as the area accumulating
under F ’ = f is positive (until t = 3), and then decreases when the curve dips below the x-axis so
that negative area starts accumulating.  The area between t = 3 and t = 4 is much smaller than
the positive area that accumulates between 0 and 3, so we know that F(4) must be larger than
F(0).  The maximum value is when t = 3; the minimum value is when t = 0.

Note that this is a different way to look at a problem we already knew how to solve – in Chapter
2, we would have found critical points of F, where f = 0 – there’s only one, when t = 3.  f = F’
goes from positive to negative there, so F has a local max at that point.  It’s the only critical
point, so it must be a global max.  Then we would look at the values of F at the endpoints to find
which was the global min.

We can also attempt to sketch a function based on the graph of the derivative.

Example 7
The graph to the right shows f'(x) - the rate of change of f(x).
Use it sketch a graph of f(x) that satisfies f(0) = 0

Recall from the last chapter the relationships between the
function graph and the derivative graph:


In the graph shown, we can see the derivative is positive on the interval (0, 1) and (3, ∞), so the
graph of f should be increasing on those intervals.  Likewise, f should be decreasing on the
interval (1,3).
f(x)
increasing Decreasing Concave up
Concave down
f '(x)
+
-
Increasing
decreasing
f ''(x)


+
-

Chapter 3    The Integral
Applied Calculus


In the graph, f' is decreasing on the interval (0, 2), so f should be concave down on that interval.
Likewise, f should be concave up on the interval (2, ∞).

The derivative itself is not enough information to know where the function f starts, since there
are a family of antiderivatives, but in this case we are given a specific point to start at.

To start the sketch, we might note first the shapes we need


then sketch the basic shapes.


Now we can attempt to sketch the graph, starting at the point (0, 0).  Notice we are very roughly
sketching this, as we don't have much information to work with.  We can tell, though, from the
graph that the area from x = 0 to x = 1 is about the same as the area from x = 1 to x = 3, so we
would expect the net area from x = 0 to x = 3 to be close to 0.


It turns out this graph isn't horribly bad.  Smoothing it out would give a graph closer to the
actual antiderivative graph, shown below.



increasing
conc down
decreasing
conc down
decreasing
conc up
increasing
conc up
increasing
conc up

Chapter 3    The Integral
Applied Calculus


Derivative of the Integral
There is another important connection between the integral and derivative.


The Fundamental Theorem of Calculus (part 2):


If
∫
=
x
a
dt
t
f
x
A
)
(
)
(
, then
)
(
)
(
)
(
x
f
dt
t
f
dx
d
x
A
x
a
=
=
′
∫



The derivative of the accumulation function is the original function.



Example 8
Let
∫
=
x
dt
t
f
x
F
)
(
)
(
, where f is graphed below.  Estimate
)3
(
F′
.


The function F measures the area from t = 0 to some t = x.  To estimate
)3
(
F′
, we want to
estimate how much the area is increasing when t = 3.  Since the value of the function f  is 0 at t
= 3, the area will not be increasing or decreasing, so we can estimate
)3
(
=
′
F


Directly using the fundamental theorem of calculus part 2,

)
(
)
(
)
(
x
f
dt
t
f
dx
d
x
F
x
=
=
′
∫
, so
)3
(
)3
(
=
=
′
f
F

Chapter 3    The Integral
Applied Calculus



## 3.2 Exercises

In problems 1 – 5, verify that  F(x)  is an antiderivative of the integrand  f(x)  and use Part 2 of the
Fundamental Theorem to evaluate the definite integrals.

1. ⌡⌠
2x dx  ,  F(x) =  x2 + 5
2.

⌡⌠
3x2 dx ,  F(x) =  x3 + 2
3. ⌡⌠
x2 dx  ,  F(x) =
3 x3
4. ⌡⌠
(x2 + 4x – 3 )  dx ,  F(x) =
3 x3 + 2x2 – 3x
5. ⌡⌠
x dx  ,  F(x) = ln( x  )

6.  Given A(x) = ⌡⌠
x
2t  dt,  find A'(x)
7.  Given A(x) = ⌡⌠
x
( 3 – t 2)  dt, find A'(x)

8.  Let A(x) = ⌡⌠
x
f(t)  dt  for the function graphed here.  Evaluate
A'(1), A'(2), A'(3).



For problems 9-10, the graph provided shows g'(x).  Use it sketch a graph of g(x) that satisfies
g(0) = 0.

9.

10.

Chapter 3    The Integral
Applied Calculus
This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 3: Antiderivatives of Formulas

Now we can put the ideas of areas and antiderivatives together to get a way of evaluating definite
integrals that is exact and often easy.   To evaluate a definite integral  ⌡⌠
a
b
f(t) dt  , we can find any
antiderivative  F of  f and evaluate F(b) – F(a).  The problem of finding the exact value of a
definite integral reduces to finding some (any) antiderivative  F  of the integrand and then
evaluating  F(b) – F(a).  Even finding one antiderivative can be difficult, and  we will stick to
functions that have easy antiderivatives.
Building Blocks

Antidifferentiation is going backwards through the derivative process.  So the easiest
antiderivative rules are simply backwards versions of the easiest derivative rules.  Recall from
Chapter 2:

Derivative Rules:  Building Blocks


In what follows, f and g are differentiable functions of x and k and n are constants.



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
⋅
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



Thinking about these basic rules was how we came up with the antiderivatives of 2x and
x
e
before.

Chapter 3    The Integral
Applied Calculus


The corresponding rules for antiderivatives are next – each of the antiderivative rules is simply
rewriting the derivative rule.  All of these antiderivatives can be verified by differentiating.

There is one surprise – the antiderivative of 1/x is actually not simply ln(x), it’s ln|x|.  This is a
good thing – the antiderivative has a domain that matches the domain of 1/x, which is bigger
than the domain of ln(x), so we don’t have to worry about whether our x’s are positive or
negative.  But you must be careful to include those absolute values – otherwise, you could end
up with domain problems.


Antiderivative Rules:  Building Blocks


In what follows, f and g are differentiable functions of x and k, n, and C are constants.




(a) Constant Multiple Rule:
( )
∫
∫
=
dx
x
f
k
dx
x
kf
)
(




(b) Sum (or Difference) Rule:
( )
( )
( )
∫
∫
∫
±
=
±
dx
x
f
dx
x
f
dx
x
g
x
f
)
(




(c)  Power Rule:
C
n
x
dx
x
n
n
+
+
=
+
∫
, provided that n = −1



Special case:
∫
+
=
C
kx
dx
k
(because
kx
k =
)








(d)  Exponential Functions:
∫
+
=
C
e
dx
e
x
x





C
a
a
dx
a
x
x
+
=
∫
ln





(e) Natural Logarithm:
∫
∫
+
=
=
−
C
x
dx
x
dx
x
ln

Chapter 3    The Integral
Applied Calculus


Example 1
Find the antiderivative of
x
x
x
+
−


(
)
C
x
x
x
dx
x
x
x
dx
x
x
x
+
−
+
−
=
+
−
=






+
−
−
−
∫
∫
/
/
/


That’s a little hard to look at, so you might want to simplify a little:
.
/
C
x
x
x
dx
x
x
x
+
−
−
=






+
−
−
∫



Example 2
Find
dx
x
ex
∫






−
+


C
x
x
e
dx
x
e
x
x
+
−
+
=






−
+
∫
ln




Example 3
Find F(x) so that
( )
xe
x
F
=
'
and ( )
0 =
F
.

This time we are looking for a particular antiderivative; we need to find exactly the right
constant.  Let’s start by finding the antiderivative:
∫
+
=
C
e
dx
e
x
x

So we know that ( )
+
=
xe
x
F
some constant; we just need to find which one.  For that, we’ll use
the other piece of information (the initial condition):
( )
( )
=
=
+
=
+
=
+
=
C
C
C
e
F
C
e
x
F
x

The particular constant we need is 9; ( )
.9
+
=
xe
x
F

Chapter 3    The Integral
Applied Calculus


The reason we are looking at antiderivatives right now is so we can evaluate definite integrals
exactly.  Recall the Fundamental Theorem of Calculus:



( )
( )
( )
a
F
b
F
dx
x
F
b
a
−
=
∫
'


If we can find an antiderivative for the integrand, we can use that to evaluate the definite integral.
The evaluation  F(b) – F(a)  is represented by the symbol  ( )]
b
a
x
F
or ( )|
b
a
x
F
.


Example 4
Evaluate  ∫
dx
x
in two ways:

(i)
By sketching the graph of  y = x  and geometrically finding the area.

(ii)
By finding an antiderivative of  F(x)  of  the integrand  and  evaluating  F(3)–F(1).




(i)  The graph of  y = x is shown to the right, and the shaded region
corresponding to the integral has area  4.





(ii)  One antiderivative of  x  is
)
(
x
x
F
=
, and

( )
( )
.4
]
=
−
=




−




=
=
∫
x
dx
x


Note that this answer agrees with the answer we got geometrically.

If we had used another antiderivative of  x,  say
)
(
2 +
=
x
x
F
,
then
( )
( )
.4
]
=
−
−
+
=




+
−




+
=






+
=
∫
x
dx
x

Whatever constant you choose, it gets subtracted away during the evaluation; we might as well
always choose the easiest one, where the constant = 0.


Example 5
Find the area between the graph of  y = 3x2  and the horizontal axis for  x  between  1  and  2.

This is
( ) ( )
.7
]
=
−
=
=
∫
x
dx
x

Chapter 3    The Integral
Applied Calculus



Example 6
A robot has been programmed so that when it starts to move, its velocity after  t  seconds will be
3t feet/second.

(a)  How far will the robot travel during its first 4 seconds of movement?

(b)  How far will the robot travel during its next  4  seconds of movement?

(a) The distance during the first 4 seconds will be the area under
the graph of velocity, from t = 0  to  t = 4.




That area is the definite integral ∫
dt
t
.  An antiderivative of
3t   is
3t  , so
]
=
−
=
=
∫
t
dt
t
feet.

(b)
]
=
−
=
−
=
=
∫
t
dt
t
feet.


Example 7
Suppose that  t  minutes after putting 1000 bacteria on a Petri plate the rate of growth of the
population is  6t  bacteria per minute.
(a)  How many new bacteria are added to the population during the first 7 minutes?
(b)  What is the total population after  7  minutes?

(a)  The number of new bacteria is the area under the rate of
growth graph, and one antiderivative of  6t  is  3t2 .



So new bacteria = ⌡⌠
6t dt    =  3t2 |
= 3(7)2 – 3(0)2  = 147
(b) The new population = (old population) + (new bacteria)
= 1000 + 147 = 1147 bacteria.



Example 8
A company determines their marginal cost for production, in dollars per item, is
)
(
+
=
x
x
MC
when producing x thousand items.  Find the cost of increasing production
from 4 thousand items to 5 thousand items.

Remember that marginal cost is the rate of change of cost, and so the fundamental theorem tells
us that
)
(
)
(
)
(
)
(
a
C
b
C
dx
x
C
dx
x
MC
b
a
b
a
−
=
′
= ∫
∫
.  In other words, the integral of marginal cost will

Chapter 3    The Integral
Applied Calculus


give us a net change in cost.  To find the cost of increasing production from 4 thousand items to
5 thousand items, we need to integrate ∫
)
(
dx
x
MC
.

We can write the marginal cost as
)
(
/
+
=
−
x
x
MC
.  We can then use the basic rules to find
an antiderivative:
x
x
x
x
x
C
/
)
(
/
+
=
+
=
.  Using this,
Net change in cost =
(
)]
(
) (
)
.3
≈
⋅
+
−
⋅
+
=
+
=






+
∫
x
x
dx
x

It will cost 3.889 thousand dollars to increase production from 4 thousand items to 5 thousand
items.

Chapter 3    The Integral
Applied Calculus



## 3.3 Exercises

For problems  1-10, find the indicated antiderivative.

1. (
)
∫
+
−
dx
x
x

2. (
)dx
x
x
∫
−
−
.1
5.2

3. ∫
dy
3.

4. ∫
dw
π

5.∫
dP
e P

6.
dx
x
e
x
x
∫






−
+

7. ∫
dx
x

8.∫
dx
x 2

9. (
)(
)
∫
+
−
dx
x
x

10. ∫
−
dt
t
t
t


For problems 11-18, find an antiderivative of the integrand and use the Fundamental Theorem to
evaluate the definite integral.
11. ⌡⌠
3x2 dx
12.

⌡⌠
–1
x2 dx
13. ⌡⌠
(x2 + 4x – 3 )  dx
14.

⌡⌠
e
x dx

15. ∫
dx
x

16. ∫
dx
x

17. ∫
1 dx
x

18. ∫
1000
1 dx
x




For problems 19 - 21 find the area shown in the figure.

19.
20.
21.

Chapter 3    The Integral
Applied Calculus
This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 4: Substitution

We don’t have many integration rules.  For quite a few of the problems we see, the rules won’t
directly apply; we’ll have to do some algebraic manipulation first.  In practice, it is much harder
to write down the antiderivative of a function than it is to find a derivative.  (In fact, it’s really
easy to write a function that doesn’t have any antiderivative you can find with algebra.)

The Substitution Method is one way of algebraically manipulating an integrand so that the rules
apply.  This is a way to unwind the Chain Rule for derivatives.  When you find the derivative of
a function using the Chain Rule, you end up with a product of something like the original
function TIMES a derivative.  We can reverse this to write an integral:

(
)
(
)
dx
x
g
x
g
f
x
g
f
dx
d
)
(
)
(
)
(
′
′
=
,      so    (
)
(
)
∫
′
′
=
dx
x
g
x
g
f
x
g
f
)
(
)
(
)
(


With substitution, we will substitute
)
(x
g
u =
.   This means
)
(x
g
dx
du
′
=
, so
dx
x
g
du
)
(′
=
.
Making this substitutions,
(
)
∫
′
′
dx
x
g
x
g
f
)
(
)
(
becomes
( )
∫
′
du
u
f
, which will probably be easier
to ingegrate.

Try Substitution when you see a product in your integral, especially if you recognize one factor
as the derivative of some part of the other factor.

The Substitution Method for Antiderivatives:

The goal is to turn
(
)
∫
dx
x
g
f
)
(
into
( )
∫
du
u
f
, where f(u) is much less messy than f(g(x)).

1. Let u be some part of the integrand.  A good first choice is “one step inside the messiest
bit.”

2.  Compute
dx
dx
du
du =


3. Translate all your x’s into u’s everywhere in the integral, including the dx.  When you’re
done, you should have a new integral that is entirely in u.  If you have any x’s left, then
that’s an indication that the substitution didn’t work or isn't complete; you may need to
go back to step 1 and try a different choice for u.

4. Integrate the new u-integral, if possible.  If you still can’t integrate it, go back to step 1
and try a different choice for u.

5. Finally, substitute back x’s for u’s everywhere in your answer.

Chapter 3    The Integral
Applied Calculus



Example 1
Evaluate ∫
−
dx
x
x
.

This integrand is more complicated than anything in our list of basic integral formulas, so we’ll
have to try something else.  The only tool we have is substitution, so let’s try that!

1. Let u be some part of the integrand.  A good first choice is “one step inside the messiest bit.”
In this case, the square root in the denominator is the messiest part, so let’s let u be one step
inside:
Let
x
u
−
=

2.  Compute
dx
dx
du
du =

xdx
du
−
=

There is x dx in the integrand, so that’s a good sign; that will be −½du.

3. Translate all your x’s into u’s everywhere in the integral, including the dx.

(
)
∫
∫
∫
∫
∫
−
−
=
−
=





−
=
−
=
−
du
u
du
u
du
u
xdx
x
dx
x
x
/


Alternatively, we could have solved
xdx
du
−
=
for dx and substituted that and simplified:
x
du
dx
−
=
, so

∫
∫
∫
∫
∫
−
−
=
−
=





−
=






−
−
=
−
du
u
du
u
du
u
x
du
x
x
dx
x
x
/


4. Integrate the new u-integral, if possible.
C
u
C
u
du
u
+
−
=
+
−
=
−∫
−
/
/
/
/


5. Finally, undo our
x
u
−
=
substitution, putting back x’s for u’s everywhere in your answer.
C
x
C
u
+
−
−
=
+
−
/
.  So we have found

C
x
dx
x
x
+
−
−
=
−
∫
.

How would we check this?  By differentiating:
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
/
/
/
x
x
x
x
x
x
C
x
dx
d
C
x
dx
d
−
=
−
=
−
−
−
=
+
−
−
=
+
−
−
−
−
.

Chapter 3    The Integral
Applied Calculus



Example 2
Evaluate (
)
∫
+
x
x
e
dx
e


This integral is not in our list of building blocks.  But notice that the derivative of
+
x
e
(that
we see in the denominator) is just
x
e (which I see in the numerator), so substitution will be a
good choice for this.

Let
+
=
xe
u
.  Then
dx
e
du
x
=
, and this integral becomes ∫
∫
−
=
du
u
u
du
.
Luckily, that is on our list of building block formulas:  ∫
+
−
=
+
−
=
−
−
.
C
u
C
u
du
u

Finally, translating back:
(
)
(
)
C
e
e
dx
e
x
x
x
+
+
−
=
+
∫



Example 3
Evaluate  a) ∫
+
dx
x
x

b) ∫
+
dx
x
x


a) This is not a basic integral, but the composition is less obvious.  Here, we can treat the
denominator as the inside of the 1/x function.

Let
3 +
= x
u
.  Then
dx
x
du
=
.  Solving for dx,
3x
du
dx =
.   Substituting,
∫
∫
∫
∫
=
=
=
+
du
u
du
u
x
du
u
x
dx
x
x
.

Using our basic formulas,
∫
=
u
du
u
ln


Undoing the substitution,
C
x
dx
x
x
+
+
=
+
∫
ln


b)  It is tempting to start this problem the same way we did the last, but if we try it will not
work, since the numerator of this fraction is not the derivative of the denominator.  Instead, we
need to try a different approach.  For this problem, we can use some basic algebra.

Chapter 3    The Integral
Applied Calculus


(
)
∫
∫
∫
−
+
=






+
=
+
dx
x
x
dx
x
x
x
dx
x
x


We can integrate this using our basic rules, without needing substitution.
(
)
C
x
x
x
x
dx
x
x
+
−
=
−
+
=
+
−
−
∫


Substitution and Definite Integrals

When you use substitution to help evaluate a definite integrals, you have a choice for how to
handle the limits of integration.  You can do either of these, whichever seems better to you.  The
important thing to remember is – the original limits of integration were values of the original
variable (say, x), not values of the new variable (say, u).
(a)
You can solve the antiderivative as a side problem, translating back to x’s, and then use
the antiderivative with the original limits of integration.  Or
(b)
You can substitute for the limits of integration at the same time as you’re substituting for
everything inside the integral, and then skip the “translate back into x” step.  If the
original integral had endpoints  x =a  and
x =b,  and we make the substitution  u = g(x ) and  du = g'(x )dx,  then the new integral
will have endpoints
u= g(a) and  u=g(b)  and


⌡⌠
x=a
x=b
(original integrand) dx     becomes
u=g(b)
(new integrand) du
⌡⌠
u=g(a)
.

Method (a) seems more straightforward for most students.  But it can involve some messy
algebra.  Method (b) is often neater and usually involves fewer steps.


Example 4
Evaluate ⌡⌠
(3x –1)4 dx

We’ll need substitution to find an antiderivative, so we’ll need to handle the limits of
integration carefully.  Let's solve this example both ways.

(a) Doing the antiderivative as a side problem:

Chapter 3    The Integral
Applied Calculus


Step One – find the antiderivative, using substitution:

(
)
∫
−
dx
x


Let
3 −
= x
u
.  Then
dx
du
=
and (
)
C
u
du
u
dx
x
+
=






=
−
∫
∫


Translating back to x: (
)
(
)
C
x
dx
x
+
−
=
−
∫

Step Two – evaluate the definite integral:

(
)
(
)
( )
(
)
( )
(
)
=
−
−
=
−
−
−
=



−
=
−
∫
x
dx
x
.

(b)  Substituting for the limits of integration:
(
)
∫
−
dx
x

Let
3 −
= x
u
.  Then
dx
du
=
, and (substituting for the limits of integration) when x = 0,
u = -1, when x = 1, u = 2.
(
)
( )
(
)
]
=
−
−
=
−
−
=
=






=
−
=
−
=
=
−
=
=
=
∫
∫
u
u
u
u
x
x
u
du
u
dx
x
.


Example 5
Evaluate (
)
∫
ln
dx
x
x


I can see the derivative of
x
ln
in the integrand, so I can tell that substitution is a good choice.
Let
x
u
ln
=
.  Then
dx
x
du
=
.  When
=
x
,
ln
=
u
.  When
=
x
,
ln
=
u
.  So the new
definite integral is
(
)
(
)
(
)
(
)
.
.
ln
ln
ln
ln
ln
ln
ln
≅
−
=


=
=
=
=
=
=
=
=
∫
∫
u
u
u
u
x
x
u
du
u
dx
x
x

Chapter 3    The Integral
Applied Calculus



## 3.4 Exercises

For problems  1-8, find the indicated antiderivative.

1. (
)
∫
+
dx
x

2. ∫
dx
e
x

3. (
)
∫
dt
t
0003
.1

4. ∫
dx
x
e
x
/

5. ∫
+
dw
w

6.
dx
x
x
∫
−1

7. ∫
x
x
dx
ln

8. ∫
+
−
−
dx
x
x
x



For problems 9-12, find an antiderivative of the integrand and use the Fundamental Theorem to
evaluate the definite integral.

9. ∫
−
+
dx
x
x

10. ∫
2 dx
e x

11. ⌡⌠
(x – 2)3 dx
12.

⌡⌠
x
1 – x2  dx

Chapter 3    The Integral
Applied Calculus
This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 5: Additional Integration Techniques
Integration By Parts
Integration by parts is an integration method which enables us to find antiderivatives of some
new functions such as
)
ln(x  as well as antiderivatives of products of functions such as
)
ln(
x
x
and
x
xe .

If the function we're trying to integrate can be written as a product of two functions, u, and dv,
then integration by parts lets us trade out a complicated integral for hopefully simpler one.




INTEGRATION BY PARTS FORMULA


∫
∫
−
=
vdu
uv
udv



For definite integrals,

] ∫
∫
−
=
b
a
b
a
b
a
vdu
uv
udv



Example 1
Integrate ∫
dx
xex


To use the By Parts method, we break apart the product into two parts:
x
u =
and
dx
e
dv
x
=


We now calculate du, the derivative of u, and v, the integral of dv.
dx
dx
x
dx
d
du
=






=
and
x
x
e
dx
e
v
=
= ∫
.

Using the By Parts formula,
∫
∫
∫
−
=
−
=
dx
e
xe
vdu
uv
dx
xe
x
x
x


Notice the remaining integral is simpler that the original, and one we can easily evaluate.
C
e
xe
dx
e
xe
dx
xe
x
x
x
x
x
+
−
=
−
=
∫
∫


We could have chosen either x or
xe  as our u in the last example, but had we chosen
xe , the
second integral would have become messier, rather than simpler.

Chapter 3    The Integral
Applied Calculus







RULE OF THUMB


When selecting the u for By Parts, select a logarithmic expression if one
is present.  If not, select an algebraic expression (like x or dx).


Example 2
Integrate ∫
2 ln
dx
x
x


Since this contains a logarithmic expression, we'll use it for our u.
x
u
ln
=
and
dx
x
dv
=


We now calculate du and v.
dx
x
du
=
and
x
x
dx
x
v
=
=
= ∫
.

Using the By Parts formula,
]
∫
∫
−
=
ln
ln
dx
x
x
x
x
dx
x
x


We can simplify the expression in the integral on the right:
]
∫
∫
−
=
ln
ln
xdx
x
x
dx
x
x


The remaining integral is a basic one we can now evaluate.
]
]
ln
ln
x
x
x
dx
x
x
−
=
∫


Finally, we can evaluate the expressions
(
) (
)
[
] [
]
( )
.
ln
)
(
)
(
ln
ln
ln
≈
−
=
⋅
−
⋅
−
⋅
−
⋅
=
∫
dx
x
x


Integration Using Tables of Integrals
There are many techniques of integration we will not be studying.  Many of them lead to general
formulas which can be compiled into a Table of Integrals - a type of cheat-sheet for integration.

Chapter 3    The Integral
Applied Calculus


For example, here are two entries you might find in a table of integrals:




TABLE OF INTEGRAL EXAMPLES


∫
+
+
−
=
−
C
a
x
a
x
a
a
x
ln



∫
+
+
+
=
+
C
a
x
x
a
x
ln





Example 3
Integrate
dx
x∫
−9


This integral looks very similar to the form of the first integral in the examples table.  By
employing the rule that allows us to pull out constants, and by rewriting 9 as 32, we can better
see the match.
∫
∫
−
=
−
dx
x
dx
x


Now we simply use the formula from the table, with a = 3.

C
x
x
x
x
dx
x
dx
x
+
+
−
=






+
−
⋅
=
−
=
−
∫
∫
ln
ln


Sometimes we have to combine the table with other techniques we've learned, like substitution.
Example 4
Integrate
dx
x
x
∫
+16


This integral looks somewhat like the second integral in the example table, but the power of x is
incorrect, and there is an x2 in the numerator which does not match.  Trying to utilize this rule,
we can try to rewrite the denominator to look like (something)2.  Luckily,
( )
x
x =

dx
x
x
∫
+16
=
( )
dx
x
x
∫
+16


Now we can use substitution, letting
3x
u =
, so
dx
x
du
=
.

Chapter 3    The Integral
Applied Calculus




Making the subsitution,
( )
du
u
du
u
dx
x
x
∫
∫
∫
+
=
+
=
+


Now we can use the table entry.
C
u
u
du
u
+
+
+
=
+
∫
ln


Undoing the substitution,
C
x
x
dx
x
x
+
+
+
=
+
∫
ln



## 3.5 Exericses

In problems  1–4,  a function  u  or  dv  is given.  Find the piece  u  or  dv  which is not given,
calculate  du  and  v, and apply the Integration by Parts Formula.

1. ⌡⌠   12x.ln(x) dx
u = ln(x)
2.

⌡⌠  x.e–x  dx
u = x

3. ⌡⌠   x4 ln(x) dx
dv = x4 dx
4.

⌡⌠  x.(5x + 1)19 dx
u = x

In problems 5 - 10 evaluate the integrals

5. ⌡⌠
x
e3x    dx

6. ⌡⌠
10x.e3x dx

7.

⌡⌠
ln(2x + 5)  dx


8. ⌡⌠   x3 ln(5x) dx

9.

⌡⌠  x ln(x + 1) dx

10. ⌡⌠
ln(x)
x2    dx

For problems 11 - 14 integrate each function.

11. ⌡⌠
4 – x2
12. ⌡⌠
9 – x2
13. ⌡⌠
4 + x2
14.

⌡⌠
9 + x2

Chapter 3    The Integral
Applied Calculus
This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 6: Area, Volume, and Average Value
Area

We have already used integrals to find the area between the graph of a function and the
horizontal axis.  Integrals can also be used to find the area between two graphs.

If  f(x) ≥ g(x) for all  x  in  [a,b], then we can approximate the area between  f  and  g  by partitioning
the interval  [a,b]  and forming a Riemann sum, as shown in the picture.  The height of each
rectangle is  top – bottom,  f(ci) – g(ci)  so the area of the ith rectangle is
(height).(base) = (f(ci) – g(ci)).∆x .  Adding up this rectangles gives an approximation of the total
area as
( )
( )
(
)
∑
=
∆
−
n
i
i
i
x
c
g
c
f
, a Riemann sum.



The limit of this Riemann sum, as the number of rectangles gets larger and their width gets
smaller, is the definite integral
( )
( )
(
)
∫
−
b
a
dx
x
g
x
f
.


The area between two curves f(x) and g(x), where f(x) ≥ g(x), between x = a and x = b
is


( )
( )
(
)
∫
−
b
a
dx
x
g
x
f


The integrand is “top – bottom.”  Make a graph to be sure which curve is which.


Example 1
Find the area bounded between the graphs of  f(x) = x  and  g(x) = 3  for  1 ≤ x ≤ 4.

Chapter 3    The Integral
Applied Calculus


Always start with a graph so you can see which graph is the top and which is the bottom.  In
this example, the two curves cross, and they change positions; we’ll need to split the area into
two pieces.  Geometrically, we can see that the area is 2 + ½ = 2.5.

Writing the area as a sum of definite integrals, we get:
Area =
(
)
(
)
∫
∫
−
+
−
dx
x
dx
x

These integrals are easy to evaluate using antiderivatives:
(
)
.2
=












−
−






−
=






−
=
−
∫
x
x
dx
x



(
)
.
=












−
−






−
=






−
=
−
∫
x
x
dx
x




The two integrals also tell us that the total area between  f  and  g  is  2.5  square units, which
we already knew.

Note that the single integral
(
)
5.1
=
−
∫
dx
x
is not the area we want in the last example.  The
value of the integral is 1.5, and the value of the area is 2.5.  That’s because for the triangle on
the right, the graph of y = x is above the graph of y = 3, so the integrand 3 – x is negative; in the
definite integral, the area of that triangle comes in with a negative sign.

In this example, it was easy to see exactly where the two curves crossed so we could break the
region into the two pieces to figure separately.  In other examples, you might need to solve an
equation to find where the curves cross.

Example 2
Two objects start from the same location and travel along the
same path with velocities
( )
+
= t
t
v A
and
( )
+
−
=
t
t
t
vB

meters per second.  How far ahead is  A  after 3 seconds?

Since
( )
( )t
v
t
v
B
A
≥
, the "area" between the graphs of
( )t
vA

and
( )t
vB
represents the distance between the objects.

After 3 seconds, the distance apart
( )
( )
(
)
(
) (
)
(
)
(
)
∫
∫
∫
−
=
+
−
−
+
=
−
=
dt
t
t
dt
t
t
t
dt
t
v
t
v
B
A

( )
5.
=
−






−
⋅
=






−
=
t
t
meters.

Chapter 3    The Integral
Applied Calculus


Volume
Just as we can partition an interval and imagine approximating an area with rectangles to find a
formula for the area between curves, we can partition an interval and imagine approximating a
volume with simple shapes to find a formula for the volume of a solid.  While this approach
works for a variety of shapes, our focus will be on shapes formed by revolving a curve around
the horizontal axis.

We start with an area, the region below a function on the interval a ≤ x ≤ b.  We are going to take
that region, and rotate it around the x axis, creating the solid shape shown.


To find the volume of this solid, we can start by partitioning the
interval [0,1] and approximating the area with rectangles.  As before,
the width of each rectangle would be ∆x and the height f(ci).

If we took just one of these rectangles and rotated it about the
horizontal axis, it would form a cylindrical shape.  The radius of that
cylinder would be f(ci), so the volume would be
(
)
x
c
f
h
r
V
i
∆
=
=
)
(
π
π


The volume of the whole solid could be approximated by rotating each
of the rectangles about the x axis.  Adding up the volume of each of the
little cylindrical discs gives an approximation of the total volume as
( )
(
)
∑
=
∆
n
i
i
x
c
f
π
, a Riemann sum.

The limit of this sum as the width of the rectanges becomes small is the
definite integral
( )
(
) dx
x
f
b
a∫
π
.


The volume of the solid obtained by rotating about the x-axis the area bounded by the
curve f(x), the x-axis, x = a, and x = b is


( )
(
)
∫
b
a
dx
x
f
π





a
b
Rotate about
the axis
a
b
a
b