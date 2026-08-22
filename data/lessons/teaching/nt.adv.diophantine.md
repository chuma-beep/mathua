> Content sourced from [Discrete Mathematics: An Open Introduction, 3e](https://discrete.openmathbooks.org/dmoi3/sec_addtops-numbth.html) by Oscar Levin — CC BY-SA 4.0
##  Section 5.2 Introduction to Number Theory

We have used the natural numbers to solve problems. This was the right set of numbers to work with in discrete mathematics because we always dealt with a whole number of things. The natural numbers have been a tool. Let’s take a moment now to inspect that tool. What mathematical discoveries can we make _about_ the natural numbers themselves?

🔗

This is the main question of number theory: a huge, ancient, complex, and above all, beautiful branch of mathematics. Historically, number theory was known as the Queen of Mathematics and was very much a branch of _pure_ mathematics, studied for its own sake instead of as a means to understanding real world applications. This has changed in recent years however, as applications of number theory have been unearthed. Probably the most well known example of this is RSA cryptography, one of the methods used in encrypt data on the internet. It is number theory that makes this possible.

🔗

What sorts of questions belong to the realm of number theory? Here is a motivating example. Recall in our study of induction, we asked:

🔗

> Which amounts of postage can be made exactly using just 5-cent and 8-cent stamps?
> 
> 🔗
> 
> 🔗

We were able to prove that _any_ amount greater than 27 cents could be made. You might wonder what would happen if we changed the denomination of the stamps. What if we instead had 4- and 9-cent stamps? Would there be some amount after which all amounts would be possible? Well, again, we could replace two 4-cent stamps with a 9-cent stamp, or three 9-cent stamps with seven 4-cent stamps. In each case we can create one more cent of postage. Using this as the inductive case would allow us to prove that any amount of postage greater than 23 cents can be made.

🔗

What if we had 2-cent and 4-cent stamps. Here it looks less promising. If we take some number of 2-cent stamps and some number of 4-cent stamps, what can we say about the total? Could it ever be odd? Doesn’t look like it.

🔗

_Why_ does 5 and 8 work, 4 and 9 work, but 2 and 4 not work? What is it about these numbers? If I gave you a pair of numbers, could you tell me right away if they would work or not? We will answer these questions, and more, after first investigating some simpler properties of numbers themselves.

🔗

###  Subsection Divisibility

It is easy to add and multiply natural numbers. If we extend our focus to all integers, then subtraction is also easy (we need the negative numbers so we can subtract any number from any other number, even larger from smaller). Division is the first operation that presents a challenge. If we wanted to extend our set of numbers so any division would be possible (maybe excluding division by 0) we would need to look at the rational numbers (the set of all numbers which can be written as fractions). This would be going too far, so we will refuse this option.

🔗

In fact, it is a good thing that not every number can be divided by other numbers. This helps us understand the structure of the natural numbers and opens the door to many interesting questions and applications.

🔗

If given numbers \\(a\\) and \\(b\text{,}\\) it is possible that \\(a \div b\\) gives a whole number. In this case, we say that \\(b\\) _divides_ \\(a\text{,}\\) in symbols, we write \\(b \mid a\text{.}\\) If this holds, then \\(b\\) is a divisor or factor of \\(a\text{,}\\) and \\(a\\) is a multiple of \\(b\text{.}\\) In other words, if \\(b \mid a\text{,}\\) then \\(a = bk\\) for some integer \\(k\\) (this is saying \\(a\\) is some multiple of \\(b\\)).

🔗

#### The Divisibility Relation.

Given integers \\(m\\) and \\(n\text{,}\\) we say “\\(m\\) divides \\(n\\)” and write

\begin{equation*} m \mid n \end{equation*} 

provided \\(n \div m\\) is an integer. Thus the following assertions mean the same thing:

  1. \\(\displaystyle m \mid n\\)

🔗

  2. \\(n = mk\\) for some integer \\(k\\)

🔗

  3. \\(m\\) is a factor (or divisor) of \\(n\\)

🔗

  4. \\(n\\) is a multiple of \\(m\text{.}\\)

🔗

🔗

🔗

Notice that \\(m \mid n\\) is a statement. It is either true or false. On the other hand, \\(n \div m\\) or \\(n/m\\) is some number. If we want to claim that \\(n/m\\) is not an integer, so \\(m\\) does not divide \\(n\text{,}\\) then we can write \\(m \nmid n\text{.}\\)

🔗

####  Example 5.2.1.

Decide whether each of the statements below are true or false.

  1. \\(\displaystyle 4 \mid 20\\)

🔗

  2. \\(\displaystyle 20 \mid 4\\)

🔗

  3. \\(\displaystyle 0 \mid 5\\)

🔗

  4. \\(\displaystyle 5 \mid 0\\)

🔗

  5. \\(\displaystyle 7 \mid 7\\)

🔗

  6. \\(\displaystyle 1 \mid 37\\)

🔗

  7. \\(\displaystyle -3 \mid 12\\)

🔗

  8. \\(\displaystyle 8 \mid 12\\)

🔗

  9. \\(\displaystyle 1642 \mid 136299\\)

🔗

🔗

🔗

This last example raises a question: how might one decide whether \\(m \mid n\text{?}\\) Of course, if you had a trusted calculator, you could ask it for the value of \\(n \div m\text{.}\\) If it spits out anything other than an integer, you know \\(m \nmid n\text{.}\\) This seems a little like cheating though: we don’t have division, so should we really use division to check divisibility?

🔗

While we don’t really know how to divide, we do know how to multiply. We might try multiplying \\(m\\) by larger and larger numbers until we get close to \\(n\text{.}\\) How close? Well, we want to be sure that if we multiply \\(m\\) by the next larger integer, we go over \\(n\text{.}\\)

🔗

For example, let’s try this to decide whether \\(1642 \mid 136299\text{.}\\) Start finding multiples of 1642:

\begin{equation*} 1642 \cdot 2 = 3284 \qquad 1642 \cdot 3 = 4926 \qquad 1642\cdot 4 = 6568 \qquad \cdots\text{.} \end{equation*} 

🔗

All of these are well less than 136299. I suppose we can jump ahead a bit:

\begin{equation*} 1642 \cdot 50 = 82100 \qquad 1642 \cdot 80 = 131360 \qquad 1642 \cdot 85 = 139570\text{.} \end{equation*} 

🔗

Ah, so we need to look somewhere between 80 and 85. Try 83:

\begin{equation*} 1642 \cdot 83 = 136286\text{.} \end{equation*} 

🔗

Is this the best we can do? How far are we from our desired 136299? If we subtract, we get \\(136299 - 136286 = 13\text{.}\\) So we know we cannot go up to 84, that will be too much. In other words, we have found that

\begin{equation*} 136299 = 83 \cdot 1642 + 13\text{.} \end{equation*} 

🔗

Since \\(13 \lt 1642\text{,}\\) we can now safely say that \\(1642 \nmid 136299\text{.}\\)

🔗

It turns out that the process we went through above can be repeated for any pair of numbers. We can always write the number \\(a\\) as some multiple of the number \\(b\\) plus some remainder. We know this because we know about division with remainder from elementary school. This is just a way of saying it using multiplication. Due to the procedural nature that can be used to find the remainder, this fact is usually called the division algorithm:

🔗

#### The Division Algorithm.

Given any two integers \\(a\\) and \\(b\text{,}\\) we can always find an integer \\(q\\) such that

\begin{equation*} a = qb + r \end{equation*} 

where \\(r\\) is an integer satisfying \\(0 \le r \lt |b|\\)

🔗

🔗

The idea is that we can always take a large enough multiple of \\(b\\) so that the remainder \\(r\\) is as small as possible. We do allow the possibility of \\(r = 0\text{,}\\) in which case we have \\(b \mid a\text{.}\\)

🔗

🔗

###  Subsection Remainder Classes

The division algorithm tells us that there are only \\(b\\) possible remainders when dividing by \\(b\text{.}\\) If we fix this divisor, we can group integers by the remainder. Each group is called a remainder class modulo \\(b\\) (or sometimes residue class).

🔗

####  Example 5.2.2.

Describe the remainder classes modulo \\(5\text{.}\\)

🔗

🔗

Note that in the example above, _every_ integer is in exactly one remainder class. The technical way to say this is that the remainder classes modulo \\(b\\) form a partition of the integers. 1 

It is possible to develop a mathematical theory of partitions, prove statements about all partitions in general and then apply those observations to our case here.

The most important fact about partitions, is that it is possible to define an equivalence relation from a partition: this is a relationship between pairs of numbers which acts in all the important ways like the “equals” relationship. 2 

Again, there is a mathematical theory of equivalence relations which applies in many more instances than the one we look at here.

🔗

All fun technical language aside, the idea is really simple. If two numbers belong to the same remainder class, then in some way, they are the same. That is, they are the same _up to division by \\(b\\)_. In the case where \\(b = 5\\) above, the numbers \\(8\\) and \\(23\text{,}\\) while not the same number, are the same when it comes to dividing by 5, because both have remainder \\(3\text{.}\\)

🔗

It matters what the divisor is: \\(8\\) and \\(23\\) are the same up to division by \\(5\text{,}\\) but not up to division by \\(7\text{,}\\) since \\(8\\) has remainder of 1 when divided by 7 while 23 has a remainder of 2.

🔗

With all this in mind, let’s introduce some notation. We want to say that \\(8\\) and 23 are basically the same, even though they are not equal. It would be wrong to say \\(8 = 23\text{.}\\) Instead, we write \\(8 \equiv 23\text{.}\\) But this is not always true. It works if we are thinking division by 5, so we need to denote that somehow. What we will actually write is this:

\begin{equation*} 8 \equiv 23 \pmod{5} \end{equation*} 

which is read, “8 is congruent to 23 modulo 5” (or just “mod 5”). Of course then we could observe that

\begin{equation*} 8 \not\equiv 23 \pmod{7}\text{.} \end{equation*} 

🔗

#### Congruence Modulo \\(n\\).

We say \\(a\\) is congruent to \\(b\\) modulo \\(n\\), and write,

\begin{equation*} a \equiv b \pmod{n} \end{equation*} 

provided \\(a\\) and \\(b\\) have the same remainder when divided by \\(n\text{.}\\) In other words, provided \\(a\\) and \\(b\\) belong to the same remainder class modulo \\(n\text{.}\\)

🔗

🔗

Many books define congruence modulo \\(n\\) slightly differently. They say that \\(a \equiv b \pmod{n}\\) if and only if \\(n \mid a-b\text{.}\\) In other words, two numbers are congruent modulo \\(n\text{,}\\) if their difference is a multiple of \\(n\text{.}\\) So which definition is correct? Turns out, it doesn’t matter: they are equivalent.

🔗

To see why, consider two numbers \\(a\\) and \\(b\\) which are congruent modulo \\(n\text{.}\\) Then \\(a\\) and \\(b\\) have the same remainder when divided by \\(n\text{.}\\) We have

\begin{equation*} a = q_1 n + r \qquad\qquad b = q_2 n + r\text{.} \end{equation*} 

🔗

Here the two \\(r\\)’s really are the same. Consider what we get when we take the difference of \\(a\\) and \\(b\text{:}\\)

\begin{equation*} a-b = q_1n + r - (q_2n + r) = q_1n - q_2 n = (q_1-q_2)n\text{.} \end{equation*} 

🔗

So \\(a-b\\) is a multiple of \\(n\text{,}\\) or equivalently, \\(n \mid a-b\text{.}\\)

🔗

On the other hand, if we assume first that \\(n \mid a-b\text{,}\\) so \\(a-b = kn\text{,}\\) then consider what happens if we divide each term by \\(n\text{.}\\) Dividing \\(a\\) by \\(n\\) will leave some remainder, as will dividing \\(b\\) by \\(n\text{.}\\) However, dividing \\(kn\\) by \\(n\\) will leave 0 remainder. So the remainders on the left-hand side must cancel out. That is, the remainders must be the same.

🔗

Thus we have:

🔗

#### Congruence and Divisibility.

For any integers \\(a\text{,}\\) \\(b\text{,}\\) and \\(n\text{,}\\) we have

\begin{equation*} a \equiv b \pmod{n} \qquad \text{ if and only if } \qquad n \mid (a-b)\text{.} \end{equation*} 

🔗

🔗

It will also be useful to switch back and forth between congruences and regular equations. The above fact helps with this. We know that \\(a \equiv b \pmod{n}\\) if and only if \\(n \mid a-b\text{,}\\) if and only if \\(a-b = kn\\) for some integer \\(k\text{.}\\) Rearranging that equation, we get \\(a = b + kn\text{.}\\) In other words, if \\(a\\) and \\(b\\) are congruent modulo \\(n\text{,}\\) then \\(a\\) is \\(b\\) more than some multiple of \\(n\text{.}\\) This conforms with our earlier observation that all the numbers in a particular remainder class are the same amount larger than the multiples of \\(n\text{.}\\)

🔗

#### Congruence and Equality.

For any integers \\(a\text{,}\\) \\(b\text{,}\\) and \\(n\text{,}\\) we have

\begin{equation*} a \equiv b \pmod{n} \qquad \text{ if and only if } \qquad a = b + kn \mbox{ for some integer } k \text{.} \end{equation*} 

🔗

🔗

🔗

###  Subsection Properties of Congruence

We said earlier that congruence modulo \\(n\\) behaves, in many important ways, the same way equality does. Specifically, we could prove that congruence modulo \\(n\\) is an equivalence relation, which would require checking the following three facts:

🔗

#### Congruence Modulo \\(n\\) is an Equivalence Relation.

Given any integers \\(a\text{,}\\) \\(b\text{,}\\) and \\(c\text{,}\\) and any positive integer \\(n\text{,}\\) the following hold:

🔗

  1. \\(a \equiv a \pmod{n}\text{.}\\)

🔗

  2. If \\(a \equiv b \pmod{n}\\) then \\(b \equiv a \pmod{n}\text{.}\\)

🔗

🔗

  3. If \\(a \equiv b \pmod{n}\\) and \\(b \equiv c \pmod{n}\text{,}\\) then \\(a \equiv c \pmod{n}\text{.}\\)

🔗

🔗

🔗

In other words, congruence modulo \\(n\\) is reflexive, symmetric, and transitive, so is an equivalence relation.

🔗

🔗

You should take a minute to convince yourself that each of the properties above actually hold of congruence. Try explaining each using both the remainder and divisibility definitions.

🔗

Next, consider how congruence behaves when doing basic arithmetic. We already know that if you subtract two congruent numbers, the result will be congruent to 0 (be a multiple of \\(n\\)). What if we add something congruent to 1 to something congruent to 2? Will we get something congruent to 3?

🔗

#### Congruence and Arithmetic.

Suppose \\(a \equiv b \pmod{n}\\) and \\(c \equiv d \pmod{n}\text{.}\\) Then the following hold:

🔗

  1. \\(a+c \equiv b+d \pmod{n}\text{.}\\)

🔗

  2. \\(a-c \equiv b-d \pmod{n}\text{.}\\)

🔗

  3. \\(ac \equiv bd \pmod{n}\text{.}\\)

🔗

🔗

🔗

The above facts might be written a little strangely, but the idea is simple. If we have a true congruence, and we add the same thing to both sides, the result is still a true congruence. This sounds like we are saying:

🔗

> If \\(a \equiv b \pmod{n}\\) then \\(a+c \equiv b+c \pmod{n}\text{.}\\)
> 
> 🔗
> 
> 🔗

Of course this is true as well, it is the special case where \\(c = d\text{.}\\) But what we have works in more generality. Think of congruence as being “basically equal.” If we have two numbers which are basically equal, and we add basically the same thing to both sides, the result will be basically equal.

🔗

This seems reasonable. Is it really true? Let’s prove the first fact:

🔗

#### Proof.

Suppose \\(a \equiv b \pmod{n}\\) and \\(c \equiv d \pmod{n}\text{.}\\) That means \\(a = b + kn\\) and \\(c = d + jn\\) for integers \\(k\\) and \\(j\text{.}\\) Add these equations:

\begin{equation*} a+c = b+d + kn + jn\text{.} \end{equation*} 

🔗

But \\(kn + jn = (k+j)n\text{,}\\) which is just a multiple of \\(n\text{.}\\) So \\(a+c = b+d + (j+k)n\text{,}\\) or in other words, \\(a+c \equiv b+d \pmod{n}\\)

🔗

🔗

The other two facts can be proved in a similar way. 

🔗

One of the important consequences of these facts about congruences, is that we can basically replace any number in a congruence with any other number it is congruent to. Here are some examples to see how (and why) that works:

🔗

####  Example 5.2.3.

Find the remainder of \\(3491\\) divided by \\(9\text{.}\\)

🔗

🔗

The above example should convince you that the well known divisibility test for 9 is true: the sum of the digits of a number is divisible by 9 if and only if the original number is divisible by 9. In fact, we now know something more: any number is congruent to the sum of its digits, modulo 9. 3 

This works for 3 as well, but definitely not for any modulus in general.

🔗

Let’s try another:

🔗

####  Example 5.2.4.

Find the remainder when \\(3^{123}\\) is divided by 7.

🔗

🔗

In the above example, we are using the fact that if \\(a \equiv b \pmod n\text{,}\\) then \\(a^p \equiv b^p \pmod n\text{.}\\) This is just applying property 3 a bunch of times.

🔗

So far we have seen how to add, subtract and multiply with congruences. What about division? There is a reason we have waited to discuss it. It turns out that we cannot simply divide. In other words, even if \\(ad \equiv bd \pmod n\text{,}\\) we do not know that \\(a \equiv b \pmod n\text{.}\\) Consider, for example:

\begin{equation*} 18 \equiv 42 \pmod 8\text{.} \end{equation*} 

🔗

This is true. Now \\(18\\) and \\(42\\) are both divisible by 6. However,

\begin{equation*} 3 \not\equiv 7 \pmod 8\text{.} \end{equation*} 

🔗

While this doesn’t work, note that \\(3 \equiv 7 \pmod 4\text{.}\\) We cannot divide \\(8\\) by 6, but we can divide 8 by the greatest common factor of \\(8\\) and \\(6\text{.}\\) Will this always happen?

🔗

Suppose \\(ad \equiv bd \pmod n\text{.}\\) In other words, we have \\(ad = bd + kn\\) for some integer \\(k\text{.}\\) Of course \\(ad\\) is divisible by \\(d\text{,}\\) as is \\(bd\text{.}\\) So \\(kn\\) must also be divisible by \\(d\text{.}\\) Now if \\(n\\) and \\(d\\) have no common factors (other than 1), then we must have \\(d \mid k\text{.}\\) But in general, if we try to divide \\(kn\\) by \\(d\text{,}\\) we don’t know that we will get an integer multiple of \\(n\text{.}\\) Some of the \\(n\\) might get divided as well. To be safe, let’s divide as much of \\(n\\) as we can. Take the largest factor of both \\(d\\) and \\(n\text{,}\\) and cancel that out from \\(n\text{.}\\) The rest of the factors of \\(d\\) will come from \\(k\text{,}\\) no problem.

🔗

We will call the largest factor of both \\(d\\) and \\(n\\) the \\(\gcd(d,n)\text{,}\\) for _greatest common divisor_. In our example above, \\(\gcd(6,8) = 2\\) since the greatest divisor common to 6 and 8 is 2.

🔗

#### Congruence and Division.

Suppose \\(ad \equiv bd \pmod n\text{.}\\) Then \\(a \equiv b \pmod{\frac{n}{\gcd(d,n)}}\text{.}\\)

🔗

If \\(d\\) and \\(n\\) have no common factors then \\(\gcd(d,n) = 1\text{,}\\) so \\(a \equiv b \pmod n\text{.}\\)

🔗

🔗

####  Example 5.2.5.

Simplify the following congruences using division: (a) \\(24 \equiv 39 \pmod 5\\) and (b) \\(24 \equiv 39 \pmod{15}\text{.}\\)

🔗

🔗

🔗

###  Subsection Solving Congruences

Now that we have some algebraic rules to govern congruence relations, we can attempt to solve for an unknown in a congruence. For example, is there a value of \\(x\\) that satisfies,

\begin{equation*} 3x + 2 \equiv 4 \pmod{5}\text{,} \end{equation*} 

and if so, what is it?

🔗

In this example, since the modulus is small, we could simply try every possible value for \\(x\text{.}\\) There are really only 5 to consider, since any integer that satisfied the congruence could be replaced with any other integer it was congruent to modulo 5. Here, when \\(x = 4\\) we get \\(3x + 2 = 14\\) which is indeed congruent to 4 modulo 5. This means that \\(x = 9\\) and \\(x = 14\\) and \\(x = 19\\) and so on will each also be a solution because as we saw above, replacing any number in a congruence with a congruent number does not change the truth of the congruence.

🔗

So in this example, simply compute \\(3x + 2\\) for values of \\(x \in \\{0,1,2,3,4\\}\text{.}\\) This gives 2, 5, 8, 11, and 14 respectively, for which only 14 is congruent to 4.

🔗

Let’s also see how you could solve this using our rules for the algebra of congruences. Such an approach would be much simpler than the trial and error tactic if the modulus was larger. First, we know we can subtract 2 from both sides:

\begin{equation*} 3x \equiv 2 \pmod{5}\text{.} \end{equation*} 

🔗

Then to divide both sides by 3, we first add 0 to both sides. Of course, on the right-hand side, we want that 0 to be a 10 (yes, \\(10\\) really is 0 since they are congruent modulo 5). This gives,

\begin{equation*} 3x \equiv 12 \pmod{5}\text{.} \end{equation*} 

🔗

Now divide both sides by 3. Since \\(\gcd(3,5) = 1\text{,}\\) we do not need to change the modulus:

\begin{equation*} x \equiv 4 \pmod{5}\text{.} \end{equation*} 

🔗

Notice that this in fact gives the _general solution_ : not only can \\(x = 4\text{,}\\) but \\(x\\) can be any number which is congruent to 4. We can leave it like this, or write “\\(x = 4 + 5k\\) for any integer \\(k\text{.}\\)”

🔗

####  Example 5.2.6.

Solve the following congruences for \\(x\text{.}\\)

🔗

  1. \\(7x \equiv 12 \pmod{13}\text{.}\\)

🔗

  2. \\(84x - 38 \equiv 79 \pmod{15}\text{.}\\)

🔗

  3. \\(20x \equiv 23 \pmod{14}\text{.}\\)

🔗

🔗

🔗

The last congruence above illustrates the way in which congruences might not have solutions. We could have seen this immediately in fact. Look at the original congruence:

\begin{equation*} 20x \equiv 23 \pmod{14}\text{.} \end{equation*} 

🔗

If we write this as an equation, we get

\begin{equation*} 20x = 23 + 14k\text{,} \end{equation*} 

or equivalently \\(20x - 14k = 23\text{.}\\) We can easily see there will be no solution to this equation in integers. The left-hand side will always be even, but the right-hand side is odd. A similar problem would occur if the right-hand side was divisible by _any_ number the left-hand side was not.

🔗

So in general, given the congruence

\begin{equation*} ax \equiv b \pmod{n}\text{,} \end{equation*} 

if \\(a\\) and \\(n\\) are divisible by a number which \\(b\\) is not divisible by, then there will be no solutions. In fact, we really only need to check one divisor of \\(a\\) and \\(n\text{:}\\) the greatest common divisor. Thus, a more compact way to say this is:

🔗

#### Congruences with no solutions.

If \\(\gcd(a,n) \nmid b\text{,}\\) then \\(ax \equiv b \pmod{n}\\) has no solutions.

🔗

🔗

🔗

###  Subsection Solving Linear Diophantine Equations

Discrete math deals with whole numbers of things. So when we want to solve equations, we usually are looking for _integer_ solutions. Equations which are intended to only have integer solutions were first studied by in the third century by the Greek mathematician Diophantus of Alexandria, and as such are called _Diophantine equations_. Probably the most famous example of a Diophantine equation is \\(a^2 + b^2 = c^2\text{.}\\) The integer solutions to this equation are called Pythagorean triples. In general, solving Diophantine equations is hard (in fact, there is provably no general algorithm for deciding whether a Diophantine equation has a solution, a result known as Matiyasevich’s Theorem). We will restrict our focus to _linear_ Diophantine equations, which are considerably easier to work with.

🔗

#### Diophantine Equations.

An equation in two or more variables is called a Diophantine equation if only integers solutions are of interest. A linear Diophantine equation takes the form \\(a_1x_1 + a_2x_2 + \cdots + a_nx_n = b\\) for constants \\(a_1,\ldots, a_n, b\text{.}\\)

🔗

A solution to a Diophantine equation is a solution to the equation consisting only of integers.

🔗

🔗

We have the tools we need to solve linear Diophantine equations. We will consider, as a main example, the equation

\begin{equation*} 51x + 87y = 123\text{.} \end{equation*} 

🔗

The general strategy will be to convert the equation to a congruence, then solve that congruence. 4 

This is certainly not the only way to proceed. A more common technique would be to apply the Euclidean algorithm. Our way can be a little faster, and is presented here primarily for variety.

Let’s work this particular example to see how this might go.

🔗

First, check if perhaps there are no solutions because a divisor of \\(51\\) and \\(87\\) is not a divisor of \\(123\text{.}\\) Really, we just need to check whether \\(\gcd(51, 87) \mid 123\text{.}\\) This greatest common divisor is 3, and yes \\(3 \mid 123\text{.}\\) At this point, we might as well factor out this greatest common divisor. So instead, we will solve:

\begin{equation*} 17x + 29y = 41\text{.} \end{equation*} 

🔗

Now observe that if there are going to be solutions, then for those values of \\(x\\) and \\(y\text{,}\\) the two sides of the equation must have the same remainder as each other, no matter what we divide by. In particular, if we divide both sides by 17, we must get the same remainder. Thus we can safely write

\begin{equation*} 17x + 29y \equiv 41 \pmod{17}\text{.} \end{equation*} 

🔗

We choose 17 because \\(17x\\) will have remainder 0. This will allow us to reduce the congruence to just one variable. We could have also moved to a congruence modulo 29, although there is usually a good reason to select the smaller choice, as this will allow us to reduce the other coefficient. In our case, we reduce the congruence as follows:

\begin{equation*} \begin{aligned}17x + 29y \amp \equiv 41 \pmod{17} \\\ 0x + 12y \amp \equiv 7 \pmod{17} \\\ 12 y \amp \equiv 24 \pmod{17} \\\ y \amp \equiv 2 \pmod{17}. \end{aligned} \end{equation*} 

🔗

Now at this point we know \\(y = 2 + 17k\\) will work for any integer \\(k\text{.}\\) If we haven’t made a mistake, we should be able to plug this back into our original Diophantine equation to find \\(x\text{:}\\)

\begin{equation*} \begin{aligned}17x + 29(2 + 17k) \amp = 41\\\ 17x \amp = -17 - 29\cdot 17k\\\ x \amp = -1-29k. \end{aligned} \end{equation*} 

🔗

We have now found all solutions to the Diophantine equation. For each \\(k\text{,}\\) \\(x = -1-29k\\) and \\(y = 2 + 17k\\) will satisfy the equation. We could check this for a few cases. If \\(k = 0\text{,}\\) the solution is \\((-1,2)\text{,}\\) and yes, \\(-17 + 2\cdot 29 = 41\text{.}\\) If \\(k = 3\text{,}\\) the solution is \\((-88, 53)\text{.}\\) If \\(k = -2\text{,}\\) we get \\((57, -32)\text{.}\\)

🔗

To summarize this process, to solve \\(ax + by = c\text{,}\\) we,

🔗

  1. Divide both sides of the equation by \\(\gcd(a,b)\\) (if this does not leave the right-hand side as an integer, there are no solutions). Let’s assume that \\(ax + by = c\\) has already been reduced in this way.

🔗

🔗

  2. Pick the smaller of \\(a\\) and \\(b\\) (here, assume it is \\(b\\)), and convert to a congruence modulo \\(b\text{:}\\)

\begin{equation*} ax + by \equiv c \pmod{b}\text{.} \end{equation*} 

This will reduce to a congruence with one variable, \\(x\text{:}\\)

\begin{equation*} ax \equiv c \pmod{b}\text{.} \end{equation*} 

🔗

🔗

  3. Solve the congruence as we did in the previous section. Write your solution as an equation, such as,

\begin{equation*} x = n + kb\text{.} \end{equation*} 

🔗

🔗

  4. Plug this into the original Diophantine equation, and solve for \\(y\text{.}\\)

🔗

🔗

  5. If we want to know solutions in a particular range (for example, \\(0 \le x, y \le 20\\)), pick different values of \\(k\\) until you have all required solutions.

🔗

🔗

🔗

Here is another example:

🔗

####  Example 5.2.7.

How can you make $6.37 using just 5-cent and 8-cent stamps? What is the smallest and largest number of stamps you could use?

🔗

🔗

Using this method, as long as you can solve linear congruences in one variable, you can solve linear Diophantine equations of two variables. There are times though that solving the linear congruence is a lot of work. For example, suppose you need to solve,

\begin{equation*} 13x \equiv 6 \pmod{51}\text{.} \end{equation*} 

🔗

You _could_ keep adding 51 to the right side until you get a multiple of 13: You would get 57, 108, 159, 210, 261, 312, and 312 is the first of these that is divisible by 13. This works, but is really too much work. Instead we could convert _back_ to a Diophantine equation:

\begin{equation*} 13x = 6 + 51k\text{.} \end{equation*} 

🔗

Now solve _this_ like we have in this section. Write it as a congruence modulo 13:

\begin{equation*} \begin{aligned}0 \amp \equiv 6 + 51k \pmod{13}\\\ -12k \amp \equiv 6 \pmod{13}\\\ 2k \amp \equiv -1 \pmod{13}\\\ 2k \amp \equiv 12 \pmod{13}\\\ k \amp \equiv 6 \pmod{13}. \end{aligned} \end{equation*} 

so \\(k = 6 + 13j\text{.}\\) Now go back and figure out \\(x\text{:}\\)

\begin{equation*} \begin{aligned}13x \amp = 6 + 51(6+13j)\\\ x \amp = 24 + 51j. \end{aligned} \end{equation*} 

🔗

Of course you could do this switching back and forth between congruences and Diophantine equations as many times as you like. If you _only_ used this technique, you would essentially replicate the Euclidean algorithm, a more standard way to solve Diophantine equations.

🔗

🔗

🔗
