> Content sourced from [Linear Algebra](https://hefferon.net/linearalgebra/) by Jim Hefferon — CC BY-SA 2.5

# Solving Linear Systems


Chapter One. Linear Systems
must equal the number present afterward. Applying that in turn to the elements
C, H, N, and O gives this system.
7x = 7z
8x + 1y = 5z + 2w
1y = 3z
3y = 6z + 1w
Both examples come down to solving a system of equations. In each system,
the equations involve only the ﬁrst power of each variable. This chapter shows
how to solve any such system of equations.
I.1
Gauss’s Method

## 1.1 Deﬁnition A linear combination of x1, . . . , xn has the form
a1x1 + a2x2 + a3x3 + · · · + anxn
where the numbers a1, . . . , an ∈R are the combination’s coeﬃcients. A linear
equation in the variables x1, . . . , xn has the form a1x1 + a2x2 + a3x3 + · · · +
anxn = d where d ∈R is the constant.
An n-tuple (s1, s2, . . . , sn) ∈Rn is a solution of, or satisﬁes, that equation
if substituting the numbers s1, . . . , sn for the variables gives a true statement:
a1s1 + a2s2 + · · · + ansn = d. A system of linear equations
a1,1x1 + a1,2x2 + · · · + a1,nxn =
d1
a2,1x1 + a2,2x2 + · · · + a2,nxn =
d2
...
am,1x1 + am,2x2 + · · · + am,nxn = dm
has the solution (s1, s2, . . . , sn) if that n-tuple is a solution of all of the equations.

## 1.2 Example The combination 3x1 + 2x2 of x1 and x2 is linear. The combination
3x2
1 + 2x2 is not a linear function of x1 and x2, nor is 3x1 + 2 sin(x2).
We usually take x1, . . . , xn to be unequal to each other because in a
sum with repeats we can rearrange to make the elements unique, as with
2x + 3y + 4x = 6x + 3y. We sometimes include terms with a zero coeﬃcient, as
in x −2y + 0z, and at other times omit them, depending on what is convenient.

Section I. Solving Linear Systems

## 1.3 Example The ordered pair (−1, 5) is a solution of this system.
3x1 + 2x2 = 7
−x1 + x2 = 6
In contrast, (5, −1) is not a solution.
Finding the set of all solutions is solving the system. We don’t need guesswork
or good luck; there is an algorithm that always works. This algorithm is Gauss’s
Method (or Gaussian elimination or linear elimination).

## 1.4 Example To solve this system
3x3 = 9
x1 + 5x2 −2x3 = 2
3x1 + 2x2
= 3
we transform it, step by step, until it is in a form that we can easily solve.
The ﬁrst transformation rewrites the system by interchanging the ﬁrst and
third row.
swap row 1 with row 3
−→
3x1 + 2x2
= 3
x1 + 5x2 −2x3 = 2
3x3 = 9
The second transformation rescales the ﬁrst row by a factor of 3.
multiply row 1 by 3
−→
x1 + 6x2
= 9
x1 + 5x2 −2x3 = 2
3x3 = 9
The third transformation is the only nontrivial one in this example. We mentally
multiply both sides of the ﬁrst row by −1, mentally add that to the second row,
and write the result in as the new second row.
add −1 times row 1 to row 2
−→
x1 + 6x2
=
−x2 −2x3 = −7
3x3 =
These steps have brought the system to a form where we can easily ﬁnd the
value of each variable. The bottom equation shows that x3 = 3. Substituting 3
for x3 in the middle equation shows that x2 = 1. Substituting those two into
the top equation gives that x1 = 3. Thus the system has a unique solution; the
solution set is {(3, 1, 3)}.
We will use Gauss’s Method throughout the book. It is fast and easy. We
will now show that it is also safe: Gauss’s Method never loses solutions nor does
it ever pick up extraneous solutions, so that a tuple is a solution to the system
before we apply the method if and only if it is a solution after.

Chapter One. Linear Systems

## 1.5 Theorem (Gauss’s Method) If a linear system is changed to another by one of
these operations
(1) an equation is swapped with another
(2) an equation has both sides multiplied by a nonzero constant
(3) an equation is replaced by the sum of itself and a multiple of another
then the two systems have the same set of solutions.
Each of the three operations has a restriction. Multiplying a row by 0 is not
allowed because obviously that can change the solution set. Similarly, adding a
multiple of a row to itself is not allowed because adding −1 times the row to
itself has the eﬀect of multiplying the row by 0. And we disallow swapping a
row with itself, to make some results in the fourth chapter easier. Besides, it’s
pointless.
Proof We will cover the equation swap operation here. The other two cases
are similar and are Exercise 33.
Consider a linear system.
a1,1x1 + a1,2x2 + · · · + a1,nxn = d1
...
ai,1x1 + ai,2x2 + · · · + ai,nxn = di
...
aj,1x1 + aj,2x2 + · · · + aj,nxn = dj
...
am,1x1 + am,2x2 + · · · + am,nxn = dm
The tuple (s1, . . . , sn) satisﬁes this system if and only if substituting the values
for the variables, the s’s for the x’s, gives a conjunction of true statements:
a1,1s1+a1,2s2+· · ·+a1,nsn = d1 and . . . ai,1s1+ai,2s2+· · ·+ai,nsn = di and
. . . aj,1s1 +aj,2s2 +· · ·+aj,nsn = dj and . . . am,1s1 +am,2s2 +· · ·+am,nsn =
dm.
In a list of statements joined with ‘and’ we can rearrange the order of the
statements. Thus this requirement is met if and only if a1,1s1 + a1,2s2 + · · · +
a1,nsn = d1 and . . . aj,1s1 +aj,2s2 +· · ·+aj,nsn = dj and . . . ai,1s1 +ai,2s2 +
· · ·+ai,nsn = di and . . . am,1s1 +am,2s2 +· · ·+am,nsn = dm. This is exactly
the requirement that (s1, . . . , sn) solves the system after the row swap.
QED

Section I. Solving Linear Systems

## 1.6 Deﬁnition The three operations from Theorem 1.5 are the elementary re-
duction operations, or row operations, or Gaussian operations. They are
swapping, multiplying by a scalar (or rescaling), and row combination.
When writing out the calculations, we will abbreviate ‘row i’ by ‘ρi’ (this is
the Greek letter rho, pronounced aloud as “row”). For instance, we will denote
a row combination operation by kρi + ρj, with the row that changes written
second. To save writing we will often combine addition steps when they use the
same ρi, as in the next example.

## 1.7 Example Gauss’s Method systematically applies the row operations to solve
a system. Here is a typical case.
x + y
= 0
2x −y + 3z = 3
x −2y −z = 3
We begin by using the ﬁrst row to eliminate the 2x in the second row and the x
in the third. To get rid of the 2x we mentally multiply the entire ﬁrst row by
−2, add that to the second row, and write the result in as the new second row.
To eliminate the x in the third row we multiply the ﬁrst row by −1, add that to
the third row, and write the result in as the new third row.
−2ρ1+ρ2
−→
−ρ1+ρ3
x +
y
= 0
−3y + 3z = 3
−3y −z = 3
We ﬁnish by transforming the second system into a third, where the bottom
equation involves only one unknown. We do that by using the second row to
eliminate the y term from the third row.
−ρ2+ρ3
−→
x +
y
= 0
−3y +
3z = 3
−4z = 0
Now ﬁnding the system’s solution is easy. The third row gives z = 0. Substitute
that back into the second row to get y = −1. Then substitute back into the ﬁrst
row to get x = 1.

## 1.8 Example For the Physics problem from the start of this chapter, Gauss’s
Method gives this.
40h + 15c = 100
−50h + 25c = 50
5/4ρ1+ρ2
−→
40h +
15c = 100
(175/4)c = 175
So c = 4, and back-substitution gives that h = 1. (We will solve the Chemistry
problem later.)

Chapter One. Linear Systems

## 1.9 Example The reduction
x + y + z = 9
2x + 4y −3z = 1
3x + 6y −5z = 0
−2ρ1+ρ2
−→
−3ρ1+ρ3
x + y + z =
2y −5z = −17
3y −8z = −27
−(3/2)ρ2+ρ3
−→
x + y +
z =
2y −
5z =
−17
−(1/2)z = −(3/2)
shows that z = 3, y = −1, and x = 7.
As illustrated above, the point of Gauss’s Method is to use the elementary
reduction operations to set up back-substitution.

## 1.10 Deﬁnition In each row of a system, the ﬁrst variable with a nonzero coeﬃcient
is the row’s leading variable. A system is in echelon form if each leading
variable is to the right of the leading variable in the row above it, except for the
leading variable in the ﬁrst row, and any rows with all-zero coeﬃcients are at
the bottom.

## 1.11 Example The prior three examples only used the operation of row combina-
tion. This linear system requires the swap operation to get it into echelon form
because after the ﬁrst combination
x −y
= 0
2x −2y + z + 2w = 4
y
+ w = 0
2z + w = 5
−2ρ1+ρ2
−→
x −y
= 0
z + 2w = 4
y
+ w = 0
2z + w = 5
the second equation has no leading y. We exchange it for a lower-down row that
has a leading y.
ρ2↔ρ3
−→
x −y
= 0
y
+ w = 0
z + 2w = 4
2z + w = 5
(Had there been more than one suitable row below the second then we could
have used any one.) With that, Gauss’s Method proceeds as before.
−2ρ3+ρ4
−→
x −y
=
y
+
w =
z +
2w =
−3w = −3
Back-substitution gives w = 1, z = 2 , y = −1, and x = −1.

Section I. Solving Linear Systems
Strictly speaking, to solve linear systems we don’t need the row rescaling
operation. We have introduced it here because it is convenient and because we
will use it later in this chapter as part of a variation of Gauss’s Method, the
Gauss-Jordan Method.
All of the systems so far have the same number of equations as unknowns.
All of them have a solution and for all of them there is only one solution. We
ﬁnish this subsection by seeing other things that can happen.

## 1.12 Example This system has more equations than variables.
x + 3y =
2x + y = −3
2x + 2y = −2
Gauss’s Method helps us understand this system also, since this
−2ρ1+ρ2
−→
−2ρ1+ρ3
x +
3y =
−5y = −5
−4y = −4
shows that one of the equations is redundant. Echelon form
−(4/5)ρ2+ρ3
−→
x +
3y =
−5y = −5
0 =
gives that y = 1 and x = −2. The ‘0 = 0’ reﬂects the redundancy.
Gauss’s Method is also useful on systems with more variables than equations.
The next subsection has many examples.
Another way that linear systems can diﬀer from the examples shown above
is that some linear systems do not have a unique solution. This can happen in
two ways. The ﬁrst is that a system can fail to have any solution at all.

## 1.13 Example Contrast the system in the last example with this one.
x + 3y =
2x + y = −3
2x + 2y =
−2ρ1+ρ2
−→
−2ρ1+ρ3
x +
3y =
−5y = −5
−4y = −2
Here the system is inconsistent: no pair of numbers (s1, s2) satisﬁes all three
equations simultaneously. Echelon form makes the inconsistency obvious.
−(4/5)ρ2+ρ3
−→
x +
3y =
−5y = −5
0 =
The solution set is empty.

Chapter One. Linear Systems

## 1.14 Example The prior system has more equations than unknowns but that
is not what causes the inconsistency — Example 1.12 has more equations than
unknowns and yet is consistent. Nor is having more equations than unknowns
necessary for inconsistency, as we see with this inconsistent system that has the
same number of equations as unknowns.
x + 2y = 8
2x + 4y = 8
−2ρ1+ρ2
−→
x + 2y =
0 = −8
Instead, inconsistency has to do with the interaction of the left and right sides;
in the ﬁrst system above the left side’s second equation is twice the ﬁrst but the
right side’s second constant is not twice the ﬁrst. Later we will have more to
say about dependencies between a system’s parts.
The other way that a linear system can fail to have a unique solution, besides
having no solutions, is to have many solutions.

## 1.15 Example In this system
x + y = 4
2x + 2y = 8
any pair of numbers satisfying the ﬁrst equation also satisﬁes the second. The
solution set {(x, y) | x + y = 4} is inﬁnite; some example member pairs are (0, 4),
(−1, 5), and (2.5, 1.5).
The result of applying Gauss’s Method here contrasts with the prior example
because we do not get a contradictory equation.
−2ρ1+ρ2
−→
x + y = 4
0 = 0
Don’t be fooled by that example: a 0 = 0 equation is not the signal that a
system has many solutions.

## 1.16 Example The absence of a 0 = 0 equation does not keep a system from having
many diﬀerent solutions. This system is in echelon form, has no 0 = 0, but
has inﬁnitely many solutions, including (0, 1, −1), (0, 1/2, −1/2), (0, 0, 0), and
(0, −π, π) (any triple whose ﬁrst component is 0 and whose second component
is the negative of the third is a solution).
x + y + z = 0
y + z = 0
Nor does the presence of 0 = 0 mean that the system must have many
solutions. Example 1.12 shows that. So does this system, which does not have

Section I. Solving Linear Systems
any solutions at all despite that in echelon form it has a 0 = 0 row.
2x
−2z = 6
y + z = 1
2x + y −z = 7
3y + 3z = 0
−ρ1+ρ3
−→
2x
−2z = 6
y + z = 1
y + z = 1
3y + 3z = 0
−ρ2+ρ3
−→
−3ρ2+ρ4
2x
−2z =
y + z =
0 =
0 = −3
In summary, Gauss’s Method uses the row operations to set a system up for
back substitution. If any step shows a contradictory equation then we can stop
with the conclusion that the system has no solutions. If we reach echelon form
without a contradictory equation, and each variable is a leading variable in its
row, then the system has a unique solution and we ﬁnd it by back substitution.
Finally, if we reach echelon form without a contradictory equation, and there is
not a unique solution — that is, at least one variable is not a leading variable —
then the system has many solutions.
The next subsection explores the third case. We will see that such a system
must have inﬁnitely many solutions and we will describe the solution set.
Note. In the exercises here, and in the rest of the book, you must justify all
of your answers. For instance, if a question asks whether a system has a
solution then you must justify a yes response by producing the solution and
must justify a no response by showing that no solution exists.
Exercises
✓1.17 Use Gauss’s Method to ﬁnd the unique solution for each system.
(a) 2x + 3y = 13
x −y = −1
(b)
x
−z = 0
3x + y
= 1
−x + y + z = 4

## 1.18 Each system is in echelon form. For each, say whether the system has a unique
solution, no solution, or inﬁnitely many solutions.
(a) −3x +
2y = 0
−2y = 0
(b) x + y
= 4
y −z = 0
(c) x + y
= 4
y −z = 0
0 = 0
(d) x + y = 4
0 = 4
(e) 3x + 6y +
z = −0.5
−z =
2.5
(f) x −3y = 2
0 = 0
(g) 2x + 2y = 4
y = 1
0 = 4
(h) 2x + y = 0
(i) x −y = −1
0 =
0 =
(j) x + y −3z = −1
y −z =
z =
0 =

Chapter One. Linear Systems
✓1.19 Use Gauss’s Method to solve each system or conclude ‘many solutions’ or ‘no
solutions’.
(a) 2x + 2y = 5
x −4y = 0
(b) −x + y = 1
x + y = 2
(c) x −3y + z = 1
x + y + 2z = 14
(d)
−x −y = 1
−3x −3y = 2
(e)
4y + z = 20
2x −2y + z = 0
x
+ z = 5
x + y −z = 10
(f) 2x
+ z + w =
y
−w = −1
3x
−z −w =
4x + y + 2z + w =

## 1.20 Solve each system or conclude ‘many solutions’ or ‘no solutions’. Use Gauss’s
Method.
(a) x + y + z = 5
x −y
= 0
y + 2z = 7
(b) 3x
+ z =
x −y + 3z =
x + 2y −5z = −1
(c)
x + 3y + z = 0
−x −y
= 2
−x + y + 2z = 8
✓1.21 We can solve linear systems by methods other than Gauss’s. One often taught
in high school is to solve one of the equations for a variable, then substitute the
resulting expression into other equations. Then we repeat that step until there
is an equation with only one variable. From that we get the ﬁrst number in the
solution and then we get the rest with back-substitution. This method takes longer
than Gauss’s Method, since it involves more arithmetic operations, and is also
more likely to lead to errors. To illustrate how it can lead to wrong conclusions,
we will use the system
x + 3y =
2x + y = −3
2x + 2y =
from Example 1.13.
(a) Solve the ﬁrst equation for x and substitute that expression into the second
equation. Find the resulting y.
(b) Again solve the ﬁrst equation for x, but this time substitute that expression
into the third equation. Find this y.
What extra step must a user of this method take to avoid erroneously concluding a
system has a solution?
✓1.22 For which values of k are there no solutions, many solutions, or a unique
solution to this system?
x −y = 1
3x −3y = k

## 1.23 This system is not linear in that it says sin α instead of α
2 sin α −
cos β + 3 tan γ = 3
4 sin α + 2 cos β −2 tan γ = 10
6 sin α −3 cos β +
tan γ = 9
and yet we can apply Gauss’s Method. Do so. Does the system have a solution?
✓1.24 What conditions must the constants, the b’s, satisfy so that each of these
systems has a solution? Hint. Apply Gauss’s Method and see what happens to the
right side.

Section I. Solving Linear Systems
(a)
x −3y = b1
3x + y = b2
x + 7y = b3
2x + 4y = b4
(b)
x1 + 2x2 + 3x3 = b1
2x1 + 5x2 + 3x3 = b2
x1
+ 8x3 = b3

## 1.25 True or false: a system with more unknowns than equations has at least one
solution. (As always, to say ‘true’ you must prove it, while to say ‘false’ you must
produce a counterexample.)

## 1.26 Must any Chemistry problem like the one that starts this subsection — a balance
the reaction problem — have inﬁnitely many solutions?
✓1.27 Find the coeﬃcients a, b, and c so that the graph of f(x) = ax2 + bx + c passes
through the points (1, 2), (−1, 6), and (2, 3).

## 1.28 After Theorem 1.5 we note that multiplying a row by 0 is not allowed because
that could change a solution set. Give an example of a system with solution set S0
where after multiplying a row by 0 the new system has a solution set S1 and S0 is
a proper subset of S1, that is, S0 ̸= S1. Give an example where S0 = S1.

## 1.29 Gauss’s Method works by combining the equations in a system to make new
equations.
(a) Can we derive the equation 3x −2y = 5 by a sequence of Gaussian reduction
steps from the equations in this system?
x + y = 1
4x −y = 6
(b) Can we derive the equation 5x−3y = 2 with a sequence of Gaussian reduction
steps from the equations in this system?
2x + 2y = 5
3x + y = 4
(c) Can we derive 6x −9y + 5z = −2 by a sequence of Gaussian reduction steps
from the equations in the system?
2x + y −z = 4
6x −3y + z = 5

## 1.30 Prove that, where a, b, c, d, e are real numbers with a ̸= 0, if this linear equation
ax + by = c
has the same solution set as this one
ax + dy = e
then they are the same equation. What if a = 0?

## 1.31 Show that if ad −bc ̸= 0 then
ax + by = j
cx + dy = k
has a unique solution.
✓1.32 In the system
ax + by = c
dx + ey = f
each of the equations describes a line in the xy-plane. By geometrical reasoning,
show that there are three possibilities: there is a unique solution, there is no
solution, and there are inﬁnitely many solutions.

Chapter One. Linear Systems

## 1.33 Finish the proof of Theorem 1.5.

## 1.34 Is there a two-unknowns linear system whose solution set is all of R2?
✓1.35 Are any of the operations used in Gauss’s Method redundant? That is, can we
make any of the operations from a combination of the others?

## 1.36 Prove that each operation of Gauss’s Method is reversible. That is, show that
if two systems are related by a row operation S1 →S2 then there is a row operation
to go back S2 →S1.
? 1.37 [Anton] A box holding pennies, nickels and dimes contains thirteen coins with
a total value of 83 cents. How many coins of each type are in the box? (These are
US coins; a penny is 1 cent, a nickel is 5 cents, and a dime is 10 cents.)
? 1.38 [Con. Prob. 1955] Four positive integers are given. Select any three of the
integers, ﬁnd their arithmetic average, and add this result to the fourth integer.
Thus the numbers 29, 23, 21, and 17 are obtained. One of the original integers
is:
(a) 19
(b) 21
(c) 23
(d) 29
(e) 17
? 1.39 [Am. Math. Mon., Jan. 1935] Laugh at this: AHAHA + TEHE = TEHAW. It
resulted from substituting a code letter for each digit of a simple example in
addition, and it is required to identify the letters and prove the solution unique.
? 1.40 [Wohascum no. 2] The Wohascum County Board of Commissioners, which has
20 members, recently had to elect a President. There were three candidates (A, B,
and C); on each ballot the three candidates were to be listed in order of preference,
with no abstentions. It was found that 11 members, a majority, preferred A over
B (thus the other 9 preferred B over A). Similarly, it was found that 12 members
preferred C over A. Given these results, it was suggested that B should withdraw,
to enable a runoﬀelection between A and C. However, B protested, and it was
then found that 14 members preferred B over C! The Board has not yet recovered
from the resulting confusion. Given that every possible order of A, B, C appeared
on at least one ballot, how many members voted for B as their ﬁrst choice?
? 1.41 [Am. Math. Mon., Jan. 1963] “This system of n linear equations with n un-
knowns,” said the Great Mathematician, “has a curious property.”
“Good heavens!” said the Poor Nut, “What is it?”
“Note,” said the Great Mathematician, “that the constants are in arithmetic
progression.”
“It’s all so clear when you explain it!” said the Poor Nut. “Do you mean like
6x + 9y = 12 and 15x + 18y = 21?”
“Quite so,” said the Great Mathematician, pulling out his bassoon. “Indeed,
the system has a unique solution. Can you ﬁnd it?”
“Good heavens!” cried the Poor Nut, “I am baﬄed.”
Are you?

Section I. Solving Linear Systems
I.2
Describing the Solution Set
A linear system with a unique solution has a solution set with one element. A
linear system with no solution has a solution set that is empty. In these cases
the solution set is easy to describe. Solution sets are a challenge to describe only
when they contain many elements.

## 2.1 Example This system has many solutions because in echelon form
2x
+ z = 3
x −y −z = 1
3x −y
= 4
−(1/2)ρ1+ρ2
−→
−(3/2)ρ1+ρ3
2x
+
z =
−y −(3/2)z = −1/2
−y −(3/2)z = −1/2
−ρ2+ρ3
−→
2x
+
z =
−y −(3/2)z = −1/2
0 =
not all of the variables are leading variables. Theorem 1.5 shows that an (x, y, z)
satisﬁes the ﬁrst system if and only if it satisﬁes the third. So we can describe
the solution set {(x, y, z) | 2x + z = 3 and x −y −z = 1 and 3x −y = 4} in this
way.
{(x, y, z) | 2x + z = 3 and −y −3z/2 = −1/2}
(∗)
This description is better because it has two equations instead of three but it is
not optimal because it still has some hard to understand interactions among the
variables.
To improve it, use the variable that does not lead any equation, z, to describe
the variables that do lead, x and y. The second equation gives y = (1/2)−(3/2)z
and the ﬁrst equation gives x = (3/2)−(1/2)z. Thus we can describe the solution
set as this set of triples.
{((3/2) −(1/2)z, (1/2) −(3/2)z, z) | z ∈R}
(∗∗)
Compared with (∗), the advantage of (∗∗) is that z can be any real number.
This makes the job of deciding which tuples are in the solution set much easier.
For instance, taking z = 2 shows that (1/2, −5/2, 2) is a solution.

## 2.2 Deﬁnition In an echelon form linear system the variables that are not leading
are free.

## 2.3 Example Reduction of a linear system can end with more than one variable

Chapter One. Linear Systems
free. Gauss’s Method on this system
x +
y + z −w =
y −z + w = −1
3x
+ 6z −6w =
−y + z −w =
−3ρ1+ρ3
−→
x +
y + z −w =
y −z + w = −1
−3y + 3z −3w =
−y + z −w =
3ρ2+ρ3
−→
ρ2+ρ4
x + y + z −w =
y −z + w = −1
0 =
0 =
leaves x and y leading and both z and w free. To get the description that we
prefer, we work from the bottom. We ﬁrst express the leading variable y in terms
of z and w, as y = −1 + z −w. Moving up to the top equation, substituting for
y gives x + (−1 + z −w) + z −w = 1 and solving for x leaves x = 2 −2z + 2w.
The solution set
{(2 −2z + 2w, −1 + z −w, z, w) | z, w ∈R}
(∗∗)
has the leading variables expressed in terms of the variables that are free.

## 2.4 Example The list of leading variables may skip over some columns. After
this reduction
2x −2y
= 0
z + 3w = 2
3x −3y
= 0
x −y + 2z + 6w = 4
−(3/2)ρ1+ρ3
−→
−(1/2)ρ1+ρ4
2x −2y
= 0
z + 3w = 2
0 = 0
2z + 6w = 4
−2ρ2+ρ4
−→
2x −2y
= 0
z + 3w = 2
0 = 0
0 = 0
x and z are the leading variables, not x and y. The free variables are y and w
and so we can describe the solution set as {(y, y, 2 −3w, w) | y, w ∈R}. For
instance, (1, 1, 2, 0) satisﬁes the system — take y = 1 and w = 0. The four-tuple
(1, 0, 5, 4) is not a solution since its ﬁrst coordinate does not equal its second.
A variable that we use to describe a family of solutions is a parameter. We
say that the solution set in the prior example is parametrized with y and w.
The terms ‘parameter’ and ‘free variable’ do not mean the same thing. In the
prior example y and w are free because in the echelon form system they do not
lead. They are parameters because we used them to describe the set of solutions.
Had we instead rewritten the second equation as w = 2/3 −(1/3)z then the free
variables would still be y and w but the parameters would be y and z.

Section I. Solving Linear Systems
In the rest of this book we will solve linear systems by bringing them to
echelon form and then parametrizing with the free variables.

## 2.5 Example This is another system with inﬁnitely many solutions.
x + 2y
= 1
2x
+ z
= 2
3x + 2y + z −w = 4
−2ρ1+ρ2
−→
−3ρ1+ρ3
x +
2y
= 1
−4y + z
= 0
−4y + z −w = 1
−ρ2+ρ3
−→
x +
2y
= 1
−4y + z
= 0
−w = 1
The leading variables are x, y, and w. The variable z is free. Notice that,
although there are inﬁnitely many solutions, the value of w doesn’t vary but
is constant w = −1. To parametrize, write w in terms of z with w = −1 + 0z.
Then y = (1/4)z. Substitute for y in the ﬁrst equation to get x = 1 −(1/2)z.
The solution set is {(1 −(1/2)z, (1/4)z, z, −1) | z ∈R}.
Parametrizing solution sets shows that systems with free variables have
inﬁnitely many solutions. For instance, above z takes on all of inﬁnitely many
real number values, each associated with a diﬀerent solution.
We ﬁnish this subsection by developing a streamlined notation for linear
systems and their solution sets.

## 2.6 Deﬁnition An m×n matrix is a rectangular array of numbers with m rows
and n columns. Each number in the matrix is an entry.
We usually denote a matrix with an upper case roman letter. For instance,
A =

2.2
−7
!
has 2 rows and 3 columns and so is a 2×3 matrix. Read that aloud as “two-by-
three”; the number of rows is always stated ﬁrst. (The matrix has parentheses
around it so that when two matrices are adjacent we can tell where one ends and
the other begins.) We name matrix entries with the corresponding lower-case
letter, so that the entry in the second row and ﬁrst column of the above array
is a2,1 = 3. Note that the order of the subscripts matters: a1,2 ̸= a2,1 since
a1,2 = 2.2. We denote the set of all m×n matrices by Mm×n.
We do Gauss’s Method using matrices in essentially the same way that we
did it for systems of equations: a matrix row’s leading entry is its ﬁrst nonzero
entry (if it has one) and we perform row operations to arrive at matrix echelon
form, where the leading entry in lower rows are to the right of those in the rows

Chapter One. Linear Systems
above. We like matrix notation because it lightens the clerical load, the copying
of variables and the writing of +’s and =’s.

## 2.7 Example We can abbreviate this linear system
x + 2y
= 4
y −z = 0
x
+ 2z = 4
with this matrix.



−1



The vertical bar reminds a reader of the diﬀerence between the coeﬃcients on
the system’s left hand side and the constants on the right. With a bar, this is
an augmented matrix.



−1



−ρ1+ρ3
−→



−1
−2



2ρ2+ρ3
−→



−1



The second row stands for y −z = 0 and the ﬁrst row stands for x + 2y = 4 so
the solution set is {(4 −2z, z, z) | z ∈R}.
Matrix notation also clariﬁes the descriptions of solution sets. Example 2.3’s
{(2 −2z + 2w, −1 + z −w, z, w) | z, w ∈R} is hard to read. We will rewrite it
to group all of the constants together, all of the coeﬃcients of z together, and
all of the coeﬃcients of w together. We write them vertically, in one-column
matrices.
{





−1




+





−2




· z +





−1




· w | z, w ∈R}
For instance, the top line says that x = 2 −2z + 2w and the second line says
that y = −1 + z −w. (Our next section gives a geometric interpretation that
will help us picture the solution sets.)

## 2.8 Deﬁnition A column vector, often just called a vector, is a matrix with a
single column. A matrix with a single row is a row vector. The entries of
a vector are sometimes called components. A column or row vector whose
components are all zeros is a zero vector.
Vectors are an exception to the convention of representing matrices with
capital roman letters. We use lower-case roman or greek letters overlined with an

Section I. Solving Linear Systems
arrow: ⃗a, ⃗b, . . . or ⃗α, ⃗β, . . . (boldface is also common: a or α). For instance,
this is a column vector with a third component of 7.
⃗v =






A zero vector is denoted ⃗0. There are many diﬀerent zero vectors — the one-tall
zero vector, the two-tall zero vector, etc. — but nonetheless we will often say
“the” zero vector, expecting that the size will be clear from the context.

## 2.9 Deﬁnition The linear equation a1x1 +a2x2 + · · · +anxn = d with unknowns
x1, . . . , xn is satisﬁed by
⃗s =



s1
...
sn



if a1s1 + a2s2 + · · · + ansn = d. A vector satisﬁes a linear system if it satisﬁes
each equation in the system.
The style of description of solution sets that we use involves adding the
vectors, and also multiplying them by real numbers. Before we give the examples
showing the style we ﬁrst need to deﬁne these operations.

## 2.10 Deﬁnition The vector sum of ⃗u and ⃗v is the vector of the sums.
⃗u +⃗v =



u1
...
un


+



v1
...
vn


=



u1 + v1
...
un + vn



Note that for the addition to be deﬁned the vectors must have the same
number of entries. This entry-by-entry addition works for any pair of matrices,
not just vectors, provided that they have the same number of rows and columns.

## 2.11 Deﬁnition The scalar multiplication of the real number r and the vector ⃗v
is the vector of the multiples.
r ·⃗v = r ·



v1
...
vn


=



rv1
...
rvn



As with the addition operation, the entry-by-entry scalar multiplication
operation extends beyond vectors to apply to any matrix.

Chapter One. Linear Systems
We write scalar multiplication either as r · ⃗v or ⃗v · r, and sometimes even
omit the ‘·’ symbol: r⃗v. (Do not refer to scalar multiplication as ‘scalar product’
because that name is for a diﬀerent operation.)

## 2.12 Example





+



−1


=



2 + 3
3 −1
1 + 4


=






7 ·





−1
−3




=





−7
−21





Observe that the deﬁnitions of addition and scalar multiplication agree where
they overlap; for instance, ⃗v +⃗v = 2⃗v.
With these deﬁnitions, we are set to use matrix and vector notation to both
solve systems and express the solution.

## 2.13 Example This system
2x + y
−w
= 4
y
+ w + u = 4
x
−z + 2w
= 0
reduces in this way.



−1
−1



−(1/2)ρ1+ρ3
−→



−1
−1/2
−1
5/2
−2



(1/2)ρ2+ρ3
−→



−1
−1
1/2



The solution set is {(w + (1/2)u, 4 −w −u, 3w + (1/2)u, w, u) | w, u ∈R}. We
write that in vector form.
{







x
y
z
w
u







=














+







−1







w +







1/2
−1
1/2







u | w, u ∈R}
Note how well vector notation sets oﬀthe coeﬃcients of each parameter. For
instance, the third row of the vector form shows plainly that if u is ﬁxed then z
increases three times as fast as w. Another thing shown plainly is that setting

Section I. Solving Linear Systems
both w and u to zero gives that







x
y
z
w
u







=














is a particular solution of the linear system.

## 2.14 Example In the same way, the system
x −y + z = 1
3x
+ z = 3
5x −2y + 3z = 5
reduces



−1
−2



−3ρ1+ρ2
−→
−5ρ1+ρ3



−1
−2
−2



−ρ2+ρ3
−→



−1
−2



to give a one-parameter solution set.
{





+



−1/3
2/3


z | z ∈R}
As in the prior example, the vector not associated with the parameter






is a particular solution of the system.
Before the exercises, we will consider what we have accomplished and what
we will do in the remainder of the chapter. So far we have done the mechanics
of Gauss’s Method. We have not stopped to consider any of the questions that
arise, except for proving Theorem 1.5 — which justiﬁes the method by showing
that it gives the right answers.
For example, can we always describe solution sets as above, with a particular
solution vector added to an unrestricted linear combination of some other vectors?

Chapter One. Linear Systems
We’ve noted that the solution sets described in this way have inﬁnitely many
members so answering this question would tell us about the size of solution sets.
The following subsection shows that the answer is “yes.” This chapter’s second
section then uses that answer to describe the geometry of solution sets.
Other questions arise from the observation that we can do Gauss’s Method
in more than one way (for instance, when swapping rows we may have a choice
of rows to swap with). Theorem 1.5 says that we must get the same solution set
no matter how we proceed but if we do Gauss’s Method in two ways must we
get the same number of free variables in each echelon form system? Must those
be the same variables, that is, is it impossible to solve a problem one way to get
y and w free and solve it another way to get y and z free? The third section
of this chapter answers “yes,” that from any starting linear system, all derived
echelon form versions have the same free variables.
Thus, by the end of the chapter we will not only have a solid grounding in
the practice of Gauss’s Method but we will also have a solid grounding in the
theory. We will know exactly what can and cannot happen in a reduction.
Exercises
✓2.15 Find the indicated entry of the matrix, if it is deﬁned.
A =
1
−1

(a) a2,1
(b) a1,2
(c) a2,2
(d) a3,1
✓2.16 Give the size of each matrix.
(a)
1

(b)


−1
−1


(c)
 5

✓2.17 Do the indicated vector operation, if it is deﬁned.
(a)



+




(b) 5
 4
−1

(c)



−




(d) 7
2

+ 9
3

(e)
1

+




(f) 6



−4



+ 2




✓2.18 Solve each system using matrix notation.
Express the solution using vec-
tors.
(a) 3x + 6y = 18
x + 2y = 6
(b) x + y =
x −y = −1
(c)
x1
+ x3 = 4
x1 −x2 + 2x3 = 5
4x1 −x2 + 5x3 = 17
(d) 2a + b −c = 2
2a
+ c = 3
a −b
= 0
(e)
x + 2y −z
= 3
2x + y
+ w = 4
x −y + z + w = 1
(f)
x
+ z + w = 4
2x + y
−w = 2
3x + y + z
= 7

## 2.19 Solve each system using matrix notation. Give each solution set in vector
notation.

Section I. Solving Linear Systems
(a) 2x + y −z = 1
4x −y
= 3
(b) x
−z
= 1
y + 2z −w = 3
x + 2y + 3z −w = 7
(c)
x −
y + z
= 0
y
+ w = 0
3x −2y + 3z + w = 0
−y
−w = 0
(d)
a + 2b + 3c + d −e = 1
3a −b + c + d + e = 3

## 2.20 Solve each system using matrix notation.
Express the solution set using
vectors.
(a)
3x + 2y + z = 1
x −y + z = 2
5x + 5y + z = 0
(b)
x + y −2z =
x −y
= −3
3x −y −2z = −6
2y −2z =
(c) 2x −y −z + w =
x + y + z
= −1
(d)
x + y −2z =
x −y
= −3
3x −y −2z =
✓2.21 The vector is in the set. What value of the parameters produces that vec-
tor?
(a)
 5
−5

, {
 1
−1

k | k ∈R}
(b)


−1

, {


−2

i +



j | i, j ∈R}
(c)


−4

, {



m +



n | m, n ∈R}

## 2.22 Decide if the vector is in the set.
(a)
 3
−1

, {
−6

k | k ∈R}
(b)
5

, {
 5
−4

j | j ∈R}
(c)


−1

, {


−7

+


−1

r | r ∈R}
(d)



, {



j +


−3
−1

k | j, k ∈R}

## 2.23 [Cleary] A farmer with 1200 acres is considering planting three diﬀerent crops,
corn, soybeans, and oats. The farmer wants to use all 1200 acres. Seed corn costs
$20 per acre, while soybean and oat seed cost $50 and $12 per acre respectively.
The farmer has $40 000 available to buy seed and intends to spend it all.
(a) Use the information above to formulate two linear equations with three
unknowns and solve it.
(b) Solutions to the system are choices that the farmer can make. Write down
two reasonable solutions.
(c) Suppose that in the fall when the crops mature, the farmer can bring in

Chapter One. Linear Systems
revenue of $100 per acre for corn, $300 per acre for soybeans and $80 per acre
for oats. Which of your two solutions in the prior part would have resulted in a
larger revenue?

## 2.24 Parametrize the solution set of this one-equation system.
x1 + x2 + · · · + xn = 0
✓2.25
(a) Apply Gauss’s Method to the left-hand side to solve
x + 2y
−w = a
2x
+ z
= b
x + y
+ 2w = c
for x, y, z, and w, in terms of the constants a, b, and c.
(b) Use your answer from the prior part to solve this.
x + 2y
−w =
2x
+ z
=
x + y
+ 2w = −2

## 2.26 Why is the comma needed in the notation ‘ai,j’ for matrix entries?
✓2.27 Give the 4×4 matrix whose i, j-th entry is
(a) i + j;
(b) −1 to the i + j power.

## 2.28 For any matrix A, the transpose of A, written AT, is the matrix whose columns
are the rows of A. Find the transpose of each of these.
(a)
1

(b)
2
−3

(c)
 5

(d)




✓2.29
(a) Describe all functions f(x) = ax2 +bx+c such that f(1) = 2 and f(−1) = 6.
(b) Describe all functions f(x) = ax2 + bx + c such that f(1) = 2.

## 2.30 Show that any set of ﬁve points from the plane R2 lie on a common conic section,
that is, they all satisfy some equation of the form ax2 + by2 + cxy + dx + ey + f = 0
where some of a, . . . , f are nonzero.

## 2.31 Make up a four equations/four unknowns system having
(a) a one-parameter solution set;
(b) a two-parameter solution set;
(c) a three-parameter solution set.
? 2.32 [Shepelev] This puzzle is from a Russian web-site http://www.arbuz.uz/ and
there are many solutions to it, but mine uses linear algebra and is very naive.
There’s a planet inhabited by arbuzoids (watermeloners, to translate from Russian).
Those creatures are found in three colors: red, green and blue. There are 13 red
arbuzoids, 15 blue ones, and 17 green. When two diﬀerently colored arbuzoids
meet, they both change to the third color.
The question is, can it ever happen that all of them assume the same color?
? 2.33 [USSR Olympiad no. 174]
(a) Solve the system of equations.
ax +
y = a2
x + ay =
For what values of a does the system fail to have solutions, and for what values
of a are there inﬁnitely many solutions?