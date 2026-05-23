> Content sourced from [Applied Calculus](https://www.opentextbookstore.com/appcalc/) by Calaway, Hoffman & Lippman — CC BY 3.0

# Limits and Continuity


Chapter 2    The Derivative
Applied Calculus

This chapter is (c) 2013.  It was remixed by David Lippman from Shana Calaway's remix of Contemporary Calculus
by Dale Hoffman.  It is licensed under the Creative Commons Attribution license.
Section 1: Limits and Continuity

In the last section, we saw that as the interval over which we calculated got smaller, the secant
slopes approached the tangent slope.  The limit gives us better language with which to discuss the
idea of “approaches.”


The limit of a function describes the behavior of the function when the variable is
near, but does not equal , a specified number  (Fig. 1).  If the values of  f(x)  get
closer and closer , as close as we want, to one number  L  as we take values of  x
very close to (but not equal to) a number  c, then





We say  "the limit of f(x), as x approaches  c, is L "   and we write

)
(
lim
x
f
c
x→
= L.
(The symbol  " → "  means  "approaches"  or  "gets very close to.")


f(c)  is a single number that describes the behavior (value) of  f(x)  AT  the point  x = c.

lim
x →c  f(x)  is a single number that describes the behavior of  f(x)  NEAR, BUT NOT AT,  the
point  x = c.

If we have a graph of the function near  x = c , then it is usually easy to determine
)
(
lim
x
f
c
x→
.

Example 1
Use the graph of  y = f(x)  in Fig. 2 to determine the following limits:


(a)
lim
x →1  f(x)
(b)
lim
x →2  f(x)


(c)
lim
x →3  f(x)
(d)
lim
x →4  f(x)


(a) lim
x →1  f(x)  = 2 .
When  x  is very close to  1,  the values of  f(x) are
very close to  y = 2.  In this example, it happens that
f(1) = 2, but that is irrelevant for the limit.  The only thing that matters is what happens for  x
close to 1 but  x ≠ 1.

(b) f(2) is undefined, but we only care about the behavior of  f(x)  for  x  close to 2 and  not
equal to 2.  When  x is close to 2, the values of  f(x)  are close to  3.  If we restrict  x  close
enough to 2, the values of  y  will be as close to  3  as we want,  so  lim
x →2  f(x) = 3.

Chapter 2    The Derivative
Applied Calculus


(c) When  x  is close to  3  (or as  x approaches the value  3), the values of  f(x) are close to  1
(or approach the value 1),  so  lim
x →3 f(x) = 1.  For this limit it is completely irrelevant that  f(3) =
2,  We only care about what happens to  f(x)  for  x  close to and not equal to 3.

(d) This one is harder and we need to be careful.  When  x  is close to  4 and slightly less than
4  (x  is just to the left of  4 on the  x–axis), then the values of  f(x)  are close to  2.  But if x is
close to  4  and slightly larger than  4  then the values of  f(x)  are close to 3.   If we only know
that  x  is very close to  4, then we cannot say whether  y = f(x)  will be close to  2  or close to  3
–– it depends on whether  x  is on the right or the left side of  4.  In this situation,  the  f(x)
values are not close to a single number so we say lim
x →4  f(x)   does not exist.  It is irrelevant that
f(4) = 1.  The limit, as  x  approaches 4, would still be undefined if  f(4)  was 3  or  2 or
anything else.



We can also explore limits using tables and using algebra.

Example 2
Find  lim
x→1 2x 2 −x −1
x −1
.

Solution:  You might try to evaluate
)
(
−
−
−
=
x
x
x
x
f
at  x = 1, but f(x) is not defined at  x = 1.
It is tempting, but wrong, to conclude that this function does not have a limit as x approaches 1.

Using Tables:  Trying some "test" values for  x  which get closer and closer to 1 from both the
left  and the right, we get


The function  f is not defined at x = 1, but when x is close to 1, the values of f(x) are getting
very close to 3.  We can get f(x) as close to 3 as we want by taking x very close to 1 so
lim
x→1 2x 2 −x −1
x −1
= 3.


x
f(x)

x
f(x)
0.9
2.82

1.1
3.2
0.9998
2.9996

1.003
3.006
0.999994
2.999988

1.0001
3.0002
0.9999999
2.9999998

1.000007
3.000014










Chapter 2    The Derivative
Applied Calculus

Using algebra:  We could have found the same result by noting that
)1
(
)1
)(
(
)
(
−
−
+
=
−
−
−
=
x
x
x
x
x
x
x
f
=   2x+1  as long as x ≠ 1.  (If x≠1, then  x–1 ≠ 0  so it is valid
to divide the numerator and denominator by the factor  x–1.)  The "x→1"  part of the limit
means that  x  is close to 1 but not equal to 1, so our division step is valid and
lim
x→1 2x 2 −x −1
x −1
=  lim
x →1  2x + 1  = 3 , the correct answer.

Using a graph:  We can graph  y =
)
(
−
−
−
=
x
x
x
x
f
for
x  close to  1, and notice that whenever  x  is close to  1,
the values of  y = f(x) are close to  3.  f is not defined at  x
= 1, so the graph has a hole above  x = 1,  but we only
care about what  f(x) is doing for  x  close to but not
equal to 1.




One Sided Limits
Sometimes, what happens to us at a place depends on the direction we use to approach that place.
If we approach Niagara Falls from the upstream side, then we will be 182 feet higher and have
different worries than if we approach from the downstream side.  Similarly, the values of a
function near a point may depend on the direction we use to approach that point.




Definition of Left and Right Limits:


The left limit as  x  approaches  c  of  f(x)  is L  if the values of  f(x) get as close to  L as
we want when  x  is very close to and left of  c,  x < c: lim
x→c− f(x)  = L .


The right limit,  written with  x → c+ , requires that  x  lie to the right of c,  x > c:

L
x
f
c
x
=
+
→
)
(
lim

Chapter 2    The Derivative
Applied Calculus

Example 3
Evaluate the one sided limits of the function f(x) graphed here at
x = 0 and x = 1.

As x approach 0 from the left, the value of the function is getting
closer to 1, so
)
(
lim
=
−
→
x
f
x
.

As x approaches 0 from the right, the value of the function is getting closer to 2, so
)
(
lim
=
+
→
x
f
x


Notice that since the limit from the left and limit from the right are different, the general limit,
)
(
lim
x
f
x→
, does not exit.

At x approaches 1 from either direction, the value of the function is approaching 1, so
)
(
lim
)
(
lim
)
(
lim
=
=
=
→
+
→
−
→
x
f
x
f
x
f
x
x
x
.


Continuity
A function that is “friendly” and doesn’t have any breaks or jumps in it is called continuous.
More formally,


Definition of Continuity at a Point



A function  f  is continuous at  x = a   if and only if   lim
x →a f(x) = f(a) .


The graph to the right illustrates some of the different ways a
function can behave at and near a point, and  the table  contains
some numerical information about the function and its behavior.
Based on the information in the table, we can conclude that f is
continuous at  1  since  lim
x →1 f(x) = 2 = f(1).
We can also conclude from the information in the table that
f  is not continuous at  2 or 3 or 4,  because

lim
x →2 f(x) ≠ f(2) , lim
x →3 f(x) ≠ f(3) ,  and  lim
x →4 f(x) ≠ f(4).

The behaviors at x = 2 and x = 4 exhibit a hole in the graph,
sometimes called a removable discontinuity, since the graph
could be made continuous by changing the value of a single
point.  The behavior at x = 3 is called a jump discontinuity,
since the graph jumps between two values.
lim  f(x)
x→a
a
f(a)
does not exist
undefined

Chapter 2    The Derivative
Applied Calculus

So which functions are continuous?  It turns out pretty much every function you’ve studied is
continuous where it is defined:  polynomial, radical, rational, exponential, and logarithmic
functions are all continuous where they are defined.  Moreover, any combination of continuous
functions is also continuous.

This is helpful, because the definition of continuity says that for a continuous function,
)
(
)
(

lim
a
f
x
f
a
x
=
→
.  That means for a continuous function, we can find the limit by direct substitution
(evaluating the function) if the function is continuous at a.

Example 4
Evaluate using continuity, if possible:
a)
x
x
x

lim
−
→


b)

lim
+
−
→x
x
x


b)

lim
−
−
→x
x
x


a) The given function is polynomial, and is defined for all values of x, so we can find the limit
by direct substitution:
)
(

lim
=
−
=
−
→
x
x
x


b) The given function is rational.  It is not defined at x = -3, but we are taking the limit as x
approaches 2, and the function is defined at that point, so we can use direct substitution:

lim
−
=
+
−
=
+
−
→x
x
x


c) This function is not defined at x = 2, and so is not continuous at x = 2.  We cannot use direct
substitution.



## 2.1 Exercises

1. Use the graph to determine the following limits.

(a)
lim
x→1 f (x)
(b)
lim
x→2 f (x)

(c)
lim
x→3 f (x)
(d)
lim
x→4 f (x)


2. Use the graph to determine the following limits.

(a)
lim
x→1 f (x)
(b)
lim
x→2 f (x)


(c)
lim
x→3 f (x)
(d)
lim
x→4 f (x)

Chapter 2    The Derivative
Applied Calculus


5. Evaluate     (a) lim
x→1 x 2 + 3x + 3
x −2

(b) lim
x→2 x 2 + 3x + 3
x −2


6. Evaluate  (a) lim
x→0
x + 7
x 2 + 9x +14
(b)   lim
x→3
x + 7
x 2 + 9x +14

(c) lim
x→−4
x + 7
x 2 + 9x +14
(d) lim
x→−7
x + 7
x 2 + 9x +14



7. At which points is the function shown
discontinuous?





8. At which points is the function shown  discontinuous?




9. Find at least one point at which each function is not continuous and  state which of the  3
conditions in the definition of continuity is

violated at that point.
(a) x + 5
x – 3
(b) x2 + x – 6
x – 2

(c) x
x

(d)
π
x2 –6x + 9
(e) ln( x2 )