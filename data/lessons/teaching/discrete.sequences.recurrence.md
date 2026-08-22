> Content sourced from [Discrete Mathematics: An Open Introduction, 3e](https://discrete.openmathbooks.org/dmoi3/sec_recurrence.html) by Oscar Levin — CC BY-SA 4.0
##  Section 2.4 Solving Recurrence Relations

######  Investigate!

Consider the recurrence relation

\begin{equation*} a_n = 5a_{n-1} - 6a_{n-2}\text{.} \end{equation*} 

  1. What sequence do you get if the initial conditions are \\(a_0 = 1\text{,}\\) \\(a_1 = 2\text{?}\\) Give a closed formula for this sequence.

🔗

🔗

  2. What sequence do you get if the initial conditions are \\(a_0 = 1\text{,}\\) \\(a_1 = 3\text{?}\\) Give a closed formula.

🔗

🔗

  3. What if \\(a_0 = 2\\) and \\(a_1 = 5\text{?}\\) Find a closed formula.

🔗

🔗

🔗

🔗

We have seen that it is often easier to find recursive definitions than closed formulas. Lucky for us, there are a few techniques for converting recursive definitions to closed formulas. Doing so is called solving a recurrence relation. Recall that the recurrence relation is a recursive definition without the initial conditions. For example, the recurrence relation for the Fibonacci sequence is \\(F_n = F_{n-1} + F_{n-2}\text{.}\\) (This, together with the initial conditions \\(F_0 = 0\\) and \\(F_1 = 1\\) give the entire recursive _definition_ for the sequence.)

🔗

###  Example 2.4.1.

Find a recurrence relation and initial conditions for \\(1, 5, 17, 53, 161, 485\ldots\text{.}\\)

🔗

🔗

We are going to try to _solve_ these recurrence relations. By this we mean something very similar to solving differential equations: we want to find a function of \\(n\\) (a closed formula) which satisfies the recurrence relation, as well as the initial condition. 2 

Recurrence relations are sometimes called difference equations since they can describe the difference between terms and this highlights the relation to differential equations further.

Just like for differential equations, finding a solution might be tricky, but checking that the solution is correct is easy. 

🔗

###  Example 2.4.2.

Check that \\(a_n = 2^n + 1\\) is a solution to the recurrence relation \\(a_n = 2a_{n-1} - 1\\) with \\(a_1 = 3\text{.}\\)

🔗

🔗

Sometimes we can be clever and solve a recurrence relation by inspection. We generate the sequence using the recurrence relation and keep track of what we are doing so that we can see how to jump to finding just the \\(a_n\\) term. Here are two examples of how you might do that.

🔗

Telescoping refers to the phenomenon when many terms in a large sum cancel out—so the sum “telescopes.” For example:

\begin{equation*} (2 - 1) + (3 - 2) + (4 - 3) + \cdots + (100 - 99) + (101 - 100) = -1 + 101 \end{equation*} 

because every third term looks like: \\(2 + -2 = 0\text{,}\\) and then \\(3 + -3 = 0\\) and so on.

🔗

We can use this behavior to solve recurrence relations. Here is an example.

🔗

###  Example 2.4.3.

Solve the recurrence relation \\(a_n = a_{n-1} + n\\) with initial term \\(a_0 = 4\text{.}\\)

🔗

🔗

The above example shows a way to solve recurrence relations of the form \\(a_n = a_{n-1} + f(n)\\) where \\(\sum_{k = 1}^n f(k)\\) has a known closed formula. If you rewrite the recurrence relation as \\(a_n - a_{n-1} = f(n)\text{,}\\) and then add up all the different equations with \\(n\\) ranging between 1 and \\(n\text{,}\\) the left-hand side will always give you \\(a_n - a_0\text{.}\\) The right-hand side will be \\(\sum_{k = 1}^n f(k)\text{,}\\) which is why we need to know the closed formula for that sum.

🔗

However, telescoping will not help us with a recursion such as \\(a_n = 3a_{n-1} + 2\\) since the left-hand side will not telescope. You will have \\(-3a_{n-1}\\)’s but only one \\(a_{n-1}\text{.}\\) However, we can still be clever if we use iteration.

🔗

We have already seen an example of iteration when we found the closed formula for arithmetic and geometric sequences. The idea is, we _iterate_ the process of finding the next term, starting with the known initial condition, up until we have \\(a_n\text{.}\\) Then we simplify. In the arithmetic sequence example, we simplified by multiplying \\(d\\) by the number of times we add it to \\(a\\) when we get to \\(a_n\text{,}\\) to get from \\(a_n = a + d + d + d + \cdots + d\\) to \\(a_n = a + dn\text{.}\\)

🔗

To see how this works, let’s go through the same example we used for telescoping, but this time use iteration.

🔗

###  Example 2.4.4.

Use iteration to solve the recurrence relation \\(a_n = a_{n-1} + n\\) with \\(a_0 = 4\text{.}\\)

🔗

🔗

Of course in this case we still needed to know formula for the sum of \\(1,\ldots,n\text{.}\\) Let’s try iteration with a sequence for which telescoping doesn’t work.

🔗

###  Example 2.4.5.

Solve the recurrence relation \\(a_n = 3a_{n-1} + 2\\) subject to \\(a_0 = 1\text{.}\\)

🔗

🔗

Iteration can be messy, but when the recurrence relation only refers to one previous term (and maybe some function of \\(n\\)) it can work well. However, trying to iterate a recurrence relation such as \\(a_n = 2 a_{n-1} + 3 a_{n-2}\\) will be way too complicated. We would need to keep track of two sets of previous terms, each of which were expressed by two previous terms, and so on. The length of the formula would grow exponentially (double each time, in fact). Luckily there happens to be a method for solving recurrence relations which works very well on relations like this.

🔗

###  Subsection The Characteristic Root Technique

Suppose we want to solve a recurrence relation expressed as a combination of the two previous terms, such as \\(a_n = a_{n-1} + 6a_{n-2}\text{.}\\) In other words, we want to find a function of \\(n\\) which satisfies \\(a_n - a_{n-1} - 6a_{n-2} = 0\text{.}\\) Now iteration is too complicated, but think just for a second what would happen if we _did_ iterate. In each step, we would, among other things, multiply a previous iteration by 6. So our closed formula would include \\(6\\) multiplied some number of times. Thus it is reasonable to guess the solution will contain parts that look geometric. Perhaps the solution will take the form \\(r^n\\) for some constant \\(r\text{.}\\)

🔗

The nice thing is, we know how to check whether a formula is actually a solution to a recurrence relation: plug it in. What happens if we plug in \\(r^n\\) into the recursion above? We get

\begin{equation*} r^n - r^{n-1} - 6r^{n-2} = 0\text{.} \end{equation*} 

🔗

Now solve for \\(r\text{:}\\)

\begin{equation*} r^{n-2}(r^2 - r - 6) = 0\text{,} \end{equation*} 

so by factoring, \\(r = -2\\) or \\(r = 3\\) (or \\(r = 0\text{,}\\) although this does not help us). This tells us that \\(a_n = (-2)^n\\) is a solution to the recurrence relation, as is \\(a_n = 3^n\text{.}\\) Which one is correct? They both are, unless we specify initial conditions. Notice we could also have \\(a_n = (-2)^n + 3^n\text{.}\\) Or \\(a_n = 7(-2)^n + 4\cdot 3^n\text{.}\\) In fact, for any \\(a\\) and \\(b\text{,}\\) \\(a_n = a(-2)^n + b 3^n\\) is a solution (try plugging this into the recurrence relation). To find the values of \\(a\\) and \\(b\text{,}\\) use the initial conditions.

🔗

This points us in the direction of a more general technique for solving recurrence relations. Notice we will always be able to factor out the \\(r^{n-2}\\) as we did above. So we really only care about the other part. We call this other part the characteristic equation for the recurrence relation. We are interested in finding the roots of the characteristic equation, which are called (surprise) the characteristic roots.

🔗

#### Characteristic Roots.

Given a recurrence relation \\(a_n + \alpha a_{n-1} + \beta a_{n-2} = 0\text{,}\\) the characteristic polynomial is 

\begin{equation*} x^2 + \alpha x + \beta \end{equation*} 

giving the characteristic equation: 

\begin{equation*} x^2 + \alpha x + \beta = 0\text{.} \end{equation*} 

🔗

If \\(r_1\\) and \\(r_2\\) are two distinct roots of the characteristic polynomial (i.e., solutions to the characteristic equation), then the solution to the recurrence relation is

\begin{equation*} a_n = ar_1^n + br_2^n\text{,} \end{equation*} 

where \\(a\\) and \\(b\\) are constants determined by the initial conditions.

🔗

🔗

####  Example 2.4.6.

Solve the recurrence relation \\(a_n = 7a_{n-1} - 10 a_{n-2}\\) with \\(a_0 = 2\\) and \\(a_1 = 3\text{.}\\)

🔗

🔗

Perhaps the most famous recurrence relation is \\(F_n = F_{n-1} + F_{n-2}\text{,}\\) which together with the initial conditions \\(F_0 = 0\\) and \\(F_1= 1\\) defines the Fibonacci sequence. But notice that this is precisely the type of recurrence relation on which we can use the characteristic root technique. When you do, the only thing that changes is that the characteristic equation does not factor, so you need to use the quadratic formula to find the characteristic roots. In fact, doing so gives the third most famous irrational number, \\(\varphi\text{,}\\) the golden ratio. 

🔗

Before leaving the characteristic root technique, we should think about what might happen when you solve the characteristic equation. We have an example above in which the characteristic polynomial has two distinct roots. These roots can be integers, or perhaps irrational numbers (requiring the quadratic formula to find them). In these cases, we know what the solution to the recurrence relation looks like.

🔗

However, it is possible for the characteristic polynomial to have only one root. This can happen if the characteristic polynomial factors as \\((x - r)^2\text{.}\\) It is still the case that \\(r^n\\) would be a solution to the recurrence relation, but we won’t be able to find solutions for all initial conditions using the general form \\(a_n = ar_1^n + br_2^n\text{,}\\) since we can’t distinguish between \\(r_1^n\\) and \\(r_2^n\text{.}\\) We are in luck though:

🔗

#### Characteristic Root Technique for Repeated Roots.

Suppose the recurrence relation \\(a_n = \alpha a_{n-1} + \beta a_{n-2}\\) has a characteristic polynomial with only one root \\(r\text{.}\\) Then the solution to the recurrence relation is

\begin{equation*} a_n = ar^n + bnr^n \end{equation*} 

where \\(a\\) and \\(b\\) are constants determined by the initial conditions.

🔗

🔗

Notice the extra \\(n\\) in \\(bnr^n\text{.}\\) This allows us to solve for the constants \\(a\\) and \\(b\\) from the initial conditions.

🔗

####  Example 2.4.7.

Solve the recurrence relation \\(a_n = 6a_{n-1} - 9a_{n-2}\\) with initial conditions \\(a_0 = 1\\) and \\(a_1 = 4\text{.}\\)

🔗

🔗

Although we will not consider examples more complicated than these, this characteristic root technique can be applied to much more complicated recurrence relations. For example, \\(a_n = 2a_{n-1} + a_{n-2} - 3a_{n-3}\\) has characteristic polynomial \\(x^3 - 2 x^2 - x + 3\text{.}\\) Assuming you see how to factor such a degree 3 (or more) polynomial you can easily find the characteristic roots and as such solve the recurrence relation (the solution would look like \\(a_n = ar_1^n + br_2^n + cr_3^n\\) if there were 3 distinct roots). It is also possible that the characteristics roots are complex numbers. 

🔗

However, the characteristic root technique is only useful for solving recurrence relations in a particular form: \\(a_n\\) is given as a linear combination of some number of previous terms. These recurrence relations are called linear homogeneous recurrence relations with constant coefficients. The “homogeneous” refers to the fact that there is no additional term in the recurrence relation other than a multiple of \\(a_j\\) terms. For example, \\(a_n = 2a_{n-1} + 1\\) is _non-homogeneous_ because of the additional constant 1. There are general methods of solving such things, but we will not consider them here, other than through the use of telescoping or iteration described above.

🔗

🔗

🔗
