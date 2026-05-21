> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Euclidean_algorithm) — CC BY-SA 4.0

# Euclidean algorithm for GCD

In mathematics, the **Euclidean algorithm**, or **Euclid's algorithm**, is an efficient method for computing the greatest common divisor (GCD) of two integers, the largest number that divides them both without a remainder. It is named after the ancient Greek mathematician Euclid, who first described it in his *Elements* ().
It is an example of an *algorithm*, and is one of the oldest algorithms in common use. It can be used to reduce fractions to their simplest form, and is a part of many other number-theoretic and cryptographic calculations.

The Euclidean algorithm is based on the principle that the greatest common divisor of two numbers does not change if the larger number is replaced by its difference with the smaller number. For example, 21 is the GCD of 252 and 105 (as 252 and 105 , and the same number 21 is also the GCD of 105 and 252 − 105. Since this replacement reduces the larger of the two numbers, repeating this process gives successively smaller pairs of numbers until the two numbers become equal. When that occurs, that number is the GCD of the original two numbers. By reversing the steps or using the extended Euclidean algorithm, the GCD can be expressed as a linear combination of the two original numbers, that is the sum of the two numbers, each multiplied by an integer (for example, 21 ). The fact that the GCD can always be expressed in this way is known as Bézout's identity.

The version of the Euclidean algorithm described above—which follows Euclid's original presentation—may require many subtraction steps to find the GCD when one of the given numbers is much bigger than the other. A more efficient version of the algorithm shortcuts these steps, instead replacing the larger of the two numbers by its remainder when divided by the smaller of the two (with this version, the algorithm stops when reaching a zero remainder). With this improvement, the algorithm never requires more steps than five times the number of digits (base 10) of the smaller integer. This was proven by Gabriel Lamé in 1844 (Lamé's Theorem), and marks the beginning of computational complexity theory. Additional methods for improving the algorithm's efficiency were developed in the 20th century.

The Euclidean algorithm has many theoretical and practical applications. It is used for reducing fractions to their simplest form and for performing division in modular arithmetic. Computations using this algorithm form part of the cryptographic protocols that are used to secure internet communications, and in methods for breaking these cryptosystems by factoring large composite numbers. The Euclidean algorithm may be used to solve Diophantine equations, such as finding numbers that satisfy multiple congruences according to the Chinese remainder theorem, to construct continued fractions, and to find accurate rational approximations to real numbers. Finally, it can be used as a basic tool for proving theorems in number theory such as Lagrange's four-square theorem and the uniqueness of prime factorizations.

The original algorithm was described only for natural numbers and geometric lengths (real numbers), but the algorithm was generalized in the 19th century to other types of numbers, such as Gaussian integers and polynomials of one variable. This led to modern abstract algebraic notions such as Euclidean domains.

## Background: greatest common divisor

The Euclidean algorithm calculates the greatest common divisor (GCD) of two natural numbers and. The greatest common divisor is the largest natural number that divides both and without leaving a remainder. Synonyms for GCD include *greatest common factor* (GCF), *highest common factor* (HCF), *highest common divisor* (HCD), and *greatest common measure* (GCM). The greatest common divisor is often written as gcd(*a*, *b*) or, more simply, as (*a*, *b*), although the latter notation is ambiguous, also used for concepts such as an ideal in the ring of integers, which is closely related to GCD.

If gcd(*a*, *b*) , then and are said to be coprime (or relatively prime). This property does not imply that or are themselves prime numbers. For example, 6 and 35 factor as 1=6 = 2 × 3 and 1=35 = 5 × 7, so they are not prime, but their prime factors are different, so 6 and 35 are coprime, with no common factors other than 1.

Let *g*. Since and are both multiples of , they can be written *a* and *ng*, and there is no larger number *G* > *g* for which this is true. The natural numbers and must be coprime; since any common factor could be factored out of and to make greater. Thus, any other number that divides both and must also divide. The greatest common divisor of and is the unique (positive) common divisor of and that is divisible by any other common divisor. The greatest common divisor can be visualized as follows. Consider a rectangular area by , and any common divisor that divides both and exactly. The sides of the rectangle can be divided into segments of length , which divides the rectangle into a grid of squares of side length. The GCD is the largest value of for which this is possible. For illustration, a 24×60 rectangular area can be divided into a grid of: 1×1 squares, 2×2 squares, 3×3 squares, 4×4 squares, 6×6 squares or 12×12 squares. Therefore, 12 is the GCD of 24 and 60. A 24×60 rectangular area can be divided into a grid of 12×12 squares, with two squares along one edge (24/12 ) and five squares along the other (60/12 ).

The greatest common divisor of two numbers and is the product of the prime factors shared by the two numbers, where each prime factor can be repeated as many times as it divides both and. For example; since 1386 can be factored into 2 × 3 × 3 × 7 × 11, and 3213 can be factored into 3 × 3 × 3 × 7 × 17, the GCD of 1386 and 3213 equals 63 , the product of their shared prime factors (with 3 repeated since 3 × 3 divides both). If two numbers have no common prime factors, their GCD is 1 (obtained here as an instance of the empty product); in other words, they are coprime. A key advantage of the Euclidean algorithm is that it can find the GCD efficiently without having to compute the prime factors. Factorization of large integers is believed to be a computationally very difficult problem, and the security of many widely used cryptographic protocols is based upon its infeasibility.

Another definition of the GCD is helpful in advanced mathematics, particularly ring theory. The greatest common divisor of two nonzero numbers and is also their smallest positive integral linear combination, that is, the smallest positive number of the form *ua* + *vb* where and are integers. The set of all integral linear combinations of and is actually the same as the set of all multiples of (, where is an integer). In modern mathematical language, the ideal generated by and is the ideal generated by  alone (an ideal generated by a single element is called a principal ideal, and all ideals of the integers are principal ideals). Some properties of the GCD are in fact easier to see with this description, for instance the fact that any common divisor of and also divides the GCD (it divides both terms of *ua* + *vb*). The equivalence of this GCD definition with the other definitions is described below.

The GCD of three or more numbers equals the product of the prime factors common to all the numbers, but it can also be calculated by repeatedly taking the GCDs of pairs of numbers. For example,
gcdgcd(*a*, gcd(*b*, *c*)) = gcd(gcd(*a*, *b*), *c*) = gcd(gcd(*a*, *c*), *b*).

Thus, Euclid's algorithm, which computes the GCD of two integers, suffices to calculate the GCD of arbitrarily many integers.

## Description
### Procedure

The Euclidean algorithm can be thought of as constructing a sequence of non-negative integers that begins with the two given integers \(r_{-2} = a\) and \(r_{-1} = b\) and will eventually terminate with the integer zero: \(\{ r_{-2} = a,\ r_{-1} = b,\ r_0,\ r_1,\ \cdots,\ r_{n-1},\ r_n = 0 \}\) with \(r_{k+1} < r_k.\) The integer \(r_{n-1}\) will then be the GCD and we can state \(\text{gcd}(a,b) = r_{n-1}.\) The algorithm indicates how to construct the intermediate remainders \(r_k\) via division-with-remainder on the preceding pair \((r_{k-2},\ r_{k-1})\) by finding an integer quotient \(q_k\) so that:
\(r_{k-2} = q_k \cdot r_{k-1} + r_k \text{, with } \ r_{k-1} > r_k \geq 0.\)

Because the sequence of non-negative integers \(\{ r_k \}\) is strictly decreasing, it eventually must terminate. In other words; since \(r_k \ge 0\) for every \(k,\) and each \(r_k\) is an integer that is strictly smaller than the preceding \(r_{k-1},\) there eventually cannot be a non-negative integer smaller than zero, and hence the algorithm must terminate. In fact, the algorithm will always terminate at the *n*th step with \(r_n\) equal to zero.

To illustrate, suppose the GCD of 1071 and 462 is requested. The sequence is initially \(\{r_{-2} = 1071,\ r_{-1} = 462 \}\) and in order to find \(r_0,\) we need to find integers \(q_0\) and \(r_0 < r_{-1}\) such that:
\(1071 = q_0 \cdot 462 + r_0.\)

This is the quotient \(q_0 = 2\) since \(1071 = 2 \cdot 462 + 147.\) This determines \(r_0 = 147\) and so the sequence is now \(\{1071,\ 462,\ r_0 = 147 \}.\) The next step is to continue the sequence to find \(r_1\) by finding integers \(q_1\) and \(r_1 < r_0\) such that:
\(462 = q_1 \cdot 147 + r_1.\)

This is the quotient \(q_1 = 3\) since \(462 = 3 \cdot 147 + 21.\) This determines \(r_1 = 21\) and so the sequence is now \(\{1071,\ 462,\ 147,\ r_1 = 21 \}.\) The next step is to continue the sequence to find \(r_2\) by finding integers \(q_2\) and \(r_2 < r_1\) such that:
\(147 = q_2 \cdot 21 + r_2.\)

This is the quotient \(q_2 = 7\) since \(147 = 7 \cdot 21 + 0.\) This determines \(r_2 = 0\) and so the sequence is completed as \(\{1071,\ 462,\ 147,\ 21,\ r_2 = 0 \}\) as no further non-negative integer smaller than \(0\) can be found. The penultimate remainder \(21\) is therefore the requested GCD:
\(\text{gcd}(1071,\ 462) = 21.\)

We can generalize slightly by dropping any ordering requirement on the initial two values \(a\) and \(b.\) If \(a = b,\) the algorithm may continue and trivially find that \(\text{gcd}(a,\ a) = a\) as the sequence of remainders will be \(\{a,\ a,\ 0\}.\) If \(a < b,\) then we can also continue since \(a \equiv 0 \cdot b + a,\) suggesting the next remainder should be \(a\) itself, and the sequence is \(\{a,\ b,\ a,\ \cdots \}.\) Normally, this would be invalid because it breaks the requirement \(r_0 < r_{-1}\) but now we have \(a < b\) by construction, so the requirement is automatically satisfied and the Euclidean algorithm can continue as normal. Therefore, dropping any ordering between the first two integers does not affect the conclusion that the sequence must eventually terminate because the next remainder will always satisfy \(r_0 < b\) and everything continues as above. The only modifications that need to be made are that \(r_{k} < r_{k-1}\) only for \(k \ge 0,\) and that the sub-sequence of non-negative integers \(\{ r_{k-1} \}\) for \(k \ge 0\) is strictly decreasing, therefore excluding \(a = r_{-2}\) from both statements.

### Proof of validity
The existence of a step \(N\) such that \(r_N=0\) follows from the condition \(|r_i|<|r_{i-1}|\), for \(i\geq1\). This ensures that the algorithm terminates.

The last non-zero remainder \(r_{N-1}\) being equal to \(\gcd(a,b)\) follows from the following properties of \(\gcd\)

1. \(\gcd(x,0)=x\), for all \(x\neq0\).

In fact, the common divisors of \(x\) and \(0\) are exactly all the divisors of \(x\), of which \(x\) itself is the largest.

2. \(\gcd(x,y)=\gcd(y,x-zy)\)

In fact, the sets of common divisors of the pairs \((x,y)\) and \((y, x-zy)\) are the same. In particular, the largest of their common divisors are the same.

To see this, assume that \(d\) is a common divisor of \(x\) and \(y\). This means that there are integers \(X,Y\) such that \(x=dX\) and \(y=dY\). It follows that \(x-zy=dX-zdY=d(X-zY)\), where \(X-zY\) is an integer. Therefore, \(d\) is also a divisor of \(x-zy\).

Conversely, if \(d\) is a common divisor of \(x-zy\) and \(y\), then there are integers \(Z\) and \(Y\) such that \(x-zy=dZ\) and \(y=dY\). From this \(x=(x-zy)+zy=dZ+zdY=d(Z+zY)\), where \(Z+zY\) is an integer. Therefore, \(d\) is also a common divisor of \(x\) and \(y\).

Now, the Euclidean algorithm starts with the pair \((r_{-2},r_{-1})=(a,b)\) and at each step, the pair of remainders \((r_{k-2},r_{k-1})\) is replaced by \((r_{k-1},r_{k-2}-q_kr_{k-1})=(r_{k-1},r_{k})\). Property 2 shows that the \(\gcd\) of the pairs remains invariant. In particular, the \(\gcd\) of the first pair \((a,b)\) and the last pair \((r_{N-1},0)\) are the same. Applying Property 1, \(\gcd(a,b)=\gcd(r_{N-1},0)=r_{N-1}\).

### Worked example

For illustration, the Euclidean algorithm can be used to find the greatest common divisor of 1=*a* = 1071 and 1=*b* = 462. To begin, multiples of 462 are subtracted from 1071 until the remainder is less than 462. Two such multiples can be subtracted (1=*q*\(_{0}\) = 2), leaving a remainder of 147:
1=1071 = 2 × 462 + 147.

Then multiples of 147 are subtracted from 462 until the remainder is less than 147. Three multiples can be subtracted (1=*q*\(_{1}\) = 3), leaving a remainder of 21:
1=462 = 3 × 147 + 21.

Then multiples of 21 are subtracted from 147 until the remainder is less than 21. Seven multiples can be subtracted (1=*q*\(_{2}\) = 7), leaving no remainder:
1=147 = 7 × 21 + 0.

Since the last remainder is zero, the algorithm ends with 21 as the greatest common divisor of 1071 and 462. This agrees with the gcd(1071, 462) found by prime factorization above. In tabular form, the steps are:

{| class="wikitable" id="basic_Euclidean_algorithm" style="margin-left:auto; margin-right:auto; text-align:center"
|-
!Step *k*!!Equation!!Quotient and remainder
|-
| 0 || 1=1071 = *q*\(_{0}\) 462 + *r*\(_{0}\) || 1=*q*\(_{0}\) = 2 and 1=*r*\(_{0}\) = 147
|-
| 1 || 1=462 = *q*\(_{1}\) 147 + *r*\(_{1}\) || 1=*q*\(_{1}\) = 3 and 1=*r*\(_{1}\) = 21
|-
| 2 || 1=147 = *q*\(_{2}\) 21 + *r*\(_{2}\) || 1=*q*\(_{2}\) = 7 and 1=*r*\(_{2}\) = 0; algorithm ends
|}

### Visualization
The Euclidean algorithm can be visualized in terms of the tiling analogy given above for the greatest common divisor. Assume that we wish to cover an *a*×*b* rectangle with square tiles exactly, where *a* is the larger of the two numbers. We first attempt to tile the rectangle using *b*×*b* square tiles; however, this leaves an *r*\(_{0}\)×*b* residual rectangle untiled, where *r*\(_{0}\) < *b*. We then attempt to tile the residual rectangle with *r*\(_{0}\)×*r*\(_{0}\) square tiles. This leaves a second residual rectangle *r*\(_{1}\)×*r*\(_{0}\), which we attempt to tile using *r*\(_{1}\)×*r*\(_{1}\) square tiles, and so on. The sequence ends when there is no residual rectangle, i.e., when the square tiles cover the previous residual rectangle exactly. The length of the sides of the smallest square tile is the GCD of the dimensions of the original rectangle. For example, the smallest square tile in the adjacent figure is 21×21 (shown in red), and 21 is the GCD of 1071 and 462, the dimensions of the original rectangle (shown in green).

### Euclidean division

At every step *k*, the Euclidean algorithm computes a quotient *q*\(_{*k*}\) and remainder *r*\(_{*k*}\) from two numbers *r*\(_{*k*−1}\) and *r*\(_{*k*−2}\)
1=*r*\(_{*k*−2}\) = *q*\(_{*k*}\) *r*\(_{*k*−1}\) + *r*\(_{*k*}\),
where the *r*\(_{*k*}\) is non-negative and is strictly less than the absolute value of *r*\(_{*k*−1}\). The theorem which underlies the definition of the Euclidean division ensures that such a quotient and remainder always exist and are unique.

In Euclid's original version of the algorithm, the quotient and remainder are found by repeated subtraction; that is, *r*\(_{*k*−1}\) is subtracted from *r*\(_{*k*−2}\) repeatedly until the remainder *r*\(_{*k*}\) is smaller than *r*\(_{*k*−1}\). After that *r*\(_{*k*}\) and *r*\(_{*k*−1}\) are exchanged and the process is iterated. Euclidean division reduces all the steps between two exchanges into a single step, which is thus more efficient. Moreover, the quotients are not needed, thus one may replace Euclidean division by the modulo operation, which gives only the remainder. Thus the iteration of the Euclidean algorithm becomes simply
1=*r*\(_{*k*}\) = *r*\(_{*k*−2}\) mod *r*\(_{*k*−1}\).

### Implementations
Implementations of the algorithm may be expressed in pseudocode. For example, the division-based version may be programmed as

 **function** gcd(a, b)
 **while** b ≠ 0
 t := b
 b := a **mod** b
 = t
 **return** a

At the beginning of the *k*th iteration, the variable *b* holds the latest remainder *r*\(_{*k*−1}\), whereas the variable *a* holds its predecessor, *r*\(_{*k*−2}\). The step 1=*b* := *a* mod *b* is equivalent to the above recursion formula *r*\(_{*k*}\) ≡ *r*\(_{*k*−2}\) mod *r*\(_{*k*−1}\). The temporary variable *t* holds the value of *r*\(_{*k*−1}\) while the next remainder *r*\(_{*k*}\) is being calculated. At the end of the loop iteration, the variable *b* holds the remainder *r*\(_{*k*}\), whereas the variable *a* holds its predecessor, *r*\(_{*k*−1}\).

(If negative inputs are allowed, or if the **mod** function may return negative values, the last line must be replaced with **return abs**(a).)

In the subtraction-based version, which was Euclid's original version, the remainder calculation (1=b := a **mod** b) is replaced by repeated subtraction. Contrary to the division-based version, which works with arbitrary integers as input, the subtraction-based version supposes that the input consists of positive integers and stops when 1=*a* = *b*:

 **function** gcd(a, b)
 **while** a ≠ b
 **if** a > b
 = a − b
 **else**
 b := b − a
 **return** a

The variables *a* and *b* alternate holding the previous remainders *r*\(_{*k*−1}\) and *r*\(_{*k*−2}\). Assume that *a* is larger than *b* at the beginning of an iteration; then *a* equals *r*\(_{*k*−2}\); since *r*\(_{*k*−2}\) > *r*\(_{*k*−1}\). During the loop iteration, *a* is reduced by multiples of the previous remainder *b* until *a* is smaller than *b*. Then *a* is the next remainder *r*\(_{*k*}\). Then *b* is reduced by multiples of *a* until it is again smaller than *a*, giving the next remainder *r*\(_{*k*+1}\), and so on.

The recursive version is based on the equality of the GCDs of successive remainders and the stopping condition 1=gcd(*r*\(_{*N*−1}\), 0) = *r*\(_{*N*−1}\).

 **function** gcd(a, b)
 **if** b = 0
 **return** a
 **else**
 **return** gcd(b, a **mod** b)

(As above, if negative inputs are allowed, or if the **mod** function may return negative values, the instruction **return** a must be replaced by **return max**(a, −a).)

For illustration, the gcd(1071, 462) is calculated from the equivalent gcdgcd(462, 147). The latter GCD is calculated from the gcdgcd(147, 21), which in turn is calculated from the gcdgcd(21, 0) = 21.

### Method of least absolute remainders
In another version of Euclid's algorithm, the quotient at each step is increased by one if the resulting negative remainder is smaller in magnitude than the typical positive remainder. Previously, the equation
1=*r*\(_{*k*−2}\) = *q*\(_{*k*}\) *r*\(_{*k*−1}\) + *r*\(_{*k*}\)
assumed that 1=|*r*\(_{*k*−1}\)| > *r*\(_{*k*}\) > 0. However, an alternative negative remainder 1=*e*\(_{*k*}\) can be computed:
1=*r*\(_{*k*−2}\) = (*q*\(_{*k*}\) + 1) *r*\(_{*k*−1}\) + *e*\(_{*k*}\)
if 1=*r*\(_{*k*−1}\) > 0 or
1=*r*\(_{*k*−2}\) = (*q*\(_{*k*}\) – 1) *r*\(_{*k*−1}\) + *e*\(_{*k*}\)
if 1=*r*\(_{*k*−1}\) < 0.

If 1=*r*\(_{*k*}\) is replaced by 1=*e*\(_{*k*}\). when 1=|*e*\(_{*k*}\)| < |*r*\(_{*k*}\)|, then one gets a variant of Euclidean algorithm such that
1=|*r*\(_{*k*}\)| ≤ |*r*\(_{*k*−1}\)| / 2
at each step.

Leopold Kronecker has shown that this version requires the fewest steps of any version of Euclid's algorithm. More generally, it has been proven that, for every input numbers *a* and *b*, the number of steps is minimal if and only if *q*\(_{*k*}\) is chosen in order that \(\left |\frac{r_{k+1{r_k}\right |<\frac{1}{\varphi}\sim 0.618,\) where \(\varphi\) is the golden ratio.

## Historical development

The Euclidean algorithm is one of the oldest algorithms in common use. It appears in Euclid's *Elements* (c. 300 BC), specifically in Book 7 (Propositions 1–2) and Book 10 (Propositions 2–3). In Book 7, the algorithm is formulated for integers, whereas in Book 10, it is formulated for lengths of line segments. (In modern usage, one would say it was formulated there for real numbers. But lengths, areas, and volumes, represented as real numbers in modern usage, are not measured in the same units and there is no natural unit of length, area, or volume; the concept of real numbers was unknown at that time.) The latter algorithm is geometrical. The GCD of two lengths *a* and *b* corresponds to the greatest length *g* that measures *a* and *b* evenly; in other words, the lengths *a* and *b* are both integer multiples of the length *g*.

The algorithm was probably not discovered by Euclid, who compiled results from earlier mathematicians in his *Elements*. The mathematician and historian B. L. van der Waerden suggests that Book VII derives from a textbook on number theory written by mathematicians in the school of Pythagoras. The algorithm was probably known by Eudoxus of Cnidus (about 375 BC). The algorithm may even pre-date Eudoxus, judging from the use of the technical term ἀνθυφαίρεσις (*anthyphairesis*, reciprocal subtraction) in works by Euclid and Aristotle. Claude Brezinski, following remarks by Pappus of Alexandria, credits the algorithm to Theaetetus (c. 417 – c. 369 BC).

Centuries later, Euclid's algorithm was discovered independently both in India and in China, primarily to solve Diophantine equations that arose in astronomy and making accurate calendars. In the late 5th century, the Indian mathematician and astronomer Aryabhata described the algorithm as the "pulverizer", perhaps because of its effectiveness in solving Diophantine equations. Although a special case of the Chinese remainder theorem had already been described in the Chinese book *Sunzi Suanjing*, the general solution was published by Qin Jiushao in his 1247 book *Shushu Jiuzhang* (數書九章 *Mathematical Treatise in Nine Sections*). The Euclidean algorithm was first described *numerically* and popularized in Europe in the second edition of Bachet's *Problèmes plaisants et délectables* (*Pleasant and enjoyable problems*, 1624). In Europe, it was likewise used to solve Diophantine equations and in developing continued fractions. The extended Euclidean algorithm was published by the English mathematician Nicholas Saunderson, who attributed it to Roger Cotes as a method for computing continued fractions efficiently.

In the 19th century, the Euclidean algorithm led to the development of new number systems, such as Gaussian integers and Eisenstein integers. In 1815, Carl Gauss used the Euclidean algorithm to demonstrate unique factorization of Gaussian integers, although his work was first published in 1832. Gauss mentioned the algorithm in his *Disquisitiones Arithmeticae* (published 1801), but only as a method for continued fractions. Peter Gustav Lejeune Dirichlet seems to have been the first to describe the Euclidean algorithm as the basis for much of number theory. Lejeune Dirichlet noted that many results of number theory, such as unique factorization, would hold true for any other system of numbers to which the Euclidean algorithm could be applied. Lejeune Dirichlet's lectures on number theory were edited and extended by Richard Dedekind, who used Euclid's algorithm to study algebraic integers, a new general type of number. For example, Dedekind was the first to prove Fermat's two-square theorem using the unique factorization of Gaussian integers. Dedekind also defined the concept of a Euclidean domain, a number system in which a generalized version of the Euclidean algorithm can be defined (as described below). In the closing decades of the 19th century, the Euclidean algorithm gradually became eclipsed by Dedekind's more general theory of ideals.

Other applications of Euclid's algorithm were developed in the 19th century. In 1829, Charles Sturm showed that the algorithm was useful in the Sturm chain method for counting the real roots of polynomials in any given interval.

The Euclidean algorithm was the first integer relation algorithm, which is a method for finding integer relations between commensurate real numbers. Several novel integer relation algorithms have been developed, such as the algorithm of Helaman Ferguson and R.W. Forcade (1979) and the LLL algorithm.

In 1969, Cole and Davie developed a two-player game based on the Euclidean algorithm, called *The Game of Euclid*, which has an optimal strategy. The players begin with two piles of *a* and *b* stones. The players take turns removing *m* multiples of the smaller pile from the larger. Thus, if the two piles consist of *x* and *y* stones, where *x* is larger than *y*, the next player can reduce the larger pile from *x* stones to *x* − *my* stones, as long as the latter is a nonnegative integer. The winner is the first player to reduce one pile to zero stones.

## Mathematical applications
### Bézout's identity

Bézout's identity states that the greatest common divisor *g* of two integers *a* and *b* can be represented as a linear sum of the original two numbers *a* and *b*. In other words, it is always possible to find integers *s* and *t* such that 1=*g* = *sa* + *tb*.

The integers *s* and *t* can be calculated from the quotients *q*\(_{0}\), *q*\(_{1}\), etc. by reversing the order of equations in Euclid's algorithm. Beginning with the next-to-last equation, *g* can be expressed in terms of the quotient *q*\(_{*N*−1}\) and the two preceding remainders, *r*\(_{*N*−2}\) and *r*\(_{*N*−3}\):
1=*g* = *r*\(_{*N*−1}\) = *r*\(_{*N*−3}\) − *q*\(_{*N*−1}\) *r*\(_{*N*−2}\).

Those two remainders can be likewise expressed in terms of their quotients and preceding remainders,
1=*r*\(_{*N*−2}\) = *r*\(_{*N*−4}\) − *q*\(_{*N*−2}\) *r*\(_{*N*−3}\) and
1=*r*\(_{*N*−3}\) = *r*\(_{*N*−5}\) − *q*\(_{*N*−3}\) *r*\(_{*N*−4}\).

Substituting these formulae for *r*\(_{*N*−2}\) and *r*\(_{*N*−3}\) into the first equation yields *g* as a linear sum of the remainders *r*\(_{*N*−4}\) and *r*\(_{*N*−5}\). The process of substituting remainders by formulae involving their predecessors can be continued until the original numbers *a* and *b* are reached:
1=*r*\(_{2}\) = *r*\(_{0}\) − *q*\(_{2}\) *r*\(_{1}\)
1=*r*\(_{1}\) = *b* − *q*\(_{1}\) *r*\(_{0}\)
1=*r*\(_{0}\) = *a* − *q*\(_{0}\) *b*.

After all the remainders *r*\(_{0}\), *r*\(_{1}\), etc. have been substituted, the final equation expresses *g* as a linear sum of *a* and *b*, so that 1=*g* = *sa* + *tb*.

The Euclidean algorithm, and thus Bézout's identity, can be generalized to the context of Euclidean domains.

### Principal ideals and related problems
Bézout's identity provides yet another definition of the greatest common divisor *g* of two numbers *a* and *b*. Consider the set of all numbers *ua* + *vb*, where *u* and *v* are any two integers. Since *a* and *b* are both divisible by *g*, every number in the set is divisible by *g*. In other words, every number of the set is an integer multiple of *g*. This is true for every common divisor of *a* and *b*. However, unlike other common divisors, the greatest common divisor is a member of the set; by Bézout's identity, choosing 1=*u* = *s* and 1=*v* = *t* gives *g*. A smaller common divisor cannot be a member of the set; since every member of the set must be divisible by *g*. Conversely, any multiple *m* of *g* can be obtained by choosing 1=*u* = *ms* and 1=*v* = *mt*, where *s* and *t* are the integers of Bézout's identity. This may be seen by multiplying Bézout's identity by *m*,
1=*mg* = *msa* + *mtb*.

Therefore, the set of all numbers *ua* + *vb* is equivalent to the set of multiples *m* of *g*. In other words, the set of all possible sums of integer multiples of two numbers (*a* and *b*) is equivalent to the set of multiples of gcd(*a*, *b*). The GCD is said to be the generator of the ideal of *a* and *b*. This GCD definition led to the modern abstract algebraic concepts of a principal ideal (an ideal generated by a single element) and a principal ideal domain (a domain in which every ideal is a principal ideal).

Certain problems can be solved using this result. For example, consider two measuring cups of volume *a* and *b*. By adding/subtracting *u* multiples of the first cup and *v* multiples of the second cup, any volume *ua* + *vb* can be measured out. These volumes are all multiples of 1=*g* = gcd(*a*, *b*).

### Extended Euclidean algorithm

The integers *s* and *t* of Bézout's identity can be computed efficiently using the extended Euclidean algorithm. This extension adds two recursive equations to Euclid's algorithm
1=*s*\(_{*k*}\) = *s*\(_{*k*−2}\) − *q*\(_{*k*}\)*s*\(_{*k*−1}\)
1=*t*\(_{*k*}\) = *t*\(_{*k*−2}\) − *q*\(_{*k*}\)*t*\(_{*k*−1}\)
with the starting values
1=*s*\(_{−2}\) = 1, *t*\(_{−2}\) = 0
1=*s*\(_{−1}\) = 0, *t*\(_{−1}\) = 1.

Using this recursion, Bézout's integers *s* and *t* are given by 1=*s* = *s*\(_{*N*}\) and 1=*t* = *t*\(_{*N*}\), where *N* + 1 is the step on which the algorithm terminates with 1=*r*\(_{*N*+1}\) = 0.

The validity of this approach can be shown by induction. Assume that the recursion formula is correct up to step *k* − 1 of the algorithm; in other words, assume that
1=*r*\(_{*j*}\) = *s*\(_{*j*}\) *a* + *t*\(_{*j*}\) *b*
for all *j* less than *k*. The *k*th step of the algorithm gives the equation
1=*r*\(_{*k*}\) = *r*\(_{*k*−2}\) − *q*\(_{*k*}\)*r*\(_{*k*−1}\).

Since the recursion formula has been assumed to be correct for *r*\(_{*k*−2}\) and *r*\(_{*k*−1}\), they may be expressed in terms of the corresponding *s* and *t* variables
1=*r*\(_{*k*}\) = (*s*\(_{*k*−2}\) *a* + *t*\(_{*k*−2}\) *b*) − *q*\(_{*k*}\)(*s*\(_{*k*−1}\) *a* + *t*\(_{*k*−1}\) *b*).

Rearranging this equation yields the recursion formula for step *k*, as required
1=*r*\(_{*k*}\) = *s*\(_{*k*}\) *a* + *t*\(_{*k*}\) *b* = (*s*\(_{*k*−2}\) − *q*\(_{*k*}\)*s*\(_{*k*−1}\)) *a* + (*t*\(_{*k*−2}\) − *q*\(_{*k*}\)*t*\(_{*k*−1}\)) *b*.

### Matrix method
The integers *s* and *t* can also be found using an equivalent matrix method. The sequence of equations of Euclid's algorithm
\(\begin{align}
a & = q_0 b + r_0 \\
b & = q_1 r_0 + r_1 \\
& \,\,\,\vdots \\
r_{N-2} & = q_N r_{N-1} + 0
\end{align}\)
can be written as a product of 2×2 quotient matrices multiplying a two-dimensional remainder vector

\(\begin{pmatrix} a \\ b \end{pmatrix} =
\begin{pmatrix} q_0 & 1 \\ 1 & 0 \end{pmatrix} \begin{pmatrix} b \\ r_0 \end{pmatrix} =
\begin{pmatrix} q_0 & 1 \\ 1 & 0 \end{pmatrix} \begin{pmatrix} q_1 & 1 \\ 1 & 0 \end{pmatrix} \begin{pmatrix} r_0 \\ r_1 \end{pmatrix} =
\cdots =
\prod_{i=0}^N \begin{pmatrix} q_i & 1 \\ 1 & 0 \end{pmatrix} \begin{pmatrix} r_{N-1} \\ 0 \end{pmatrix} \,.\)

Let **M** represent the product of all the quotient matrices

\(\mathbf{M} = \begin{pmatrix} m_{11} & m_{12} \\ m_{21} & m_{22} \end{pmatrix} =
\prod_{i=0}^N \begin{pmatrix} q_i & 1 \\ 1 & 0 \end{pmatrix} =
\begin{pmatrix} q_0 & 1 \\ 1 & 0 \end{pmatrix} \begin{pmatrix} q_1 & 1 \\ 1 & 0 \end{pmatrix} \cdots \begin{pmatrix} q_{N} & 1 \\ 1 & 0 \end{pmatrix} \,.\)

This simplifies the Euclidean algorithm to the form

\(\begin{pmatrix} a \\ b \end{pmatrix} =
\mathbf{M} \begin{pmatrix} r_{N-1} \\ 0 \end{pmatrix} =
\mathbf{M} \begin{pmatrix} g \\ 0 \end{pmatrix} \,.\)

To express *g* as a linear sum of *a* and *b*, both sides of this equation can be multiplied by the inverse of the matrix **M**. The determinant of **M** equals (−1)\(^{*N*+1}\); since it equals the product of the determinants of the quotient matrices, each of which is negative one. Since the determinant of **M** is never zero, the vector of the final remainders can be solved using the inverse of **M**

\(\begin{pmatrix} g \\ 0 \end{pmatrix} =
\mathbf{M}^{-1} \begin{pmatrix} a \\ b \end{pmatrix} =
(-1)^{N+1} \begin{pmatrix} m_{22} & -m_{12} \\ -m_{21} & m_{11} \end{pmatrix} \begin{pmatrix} a \\ b \end{pmatrix} \,.\)

Since the top equation gives
1=*g* = (−1)\(^{*N*+1}\) ( *m*\(_{22}\) *a* − *m*\(_{12}\) *b*),
the two integers of Bézout's identity are 1=*s* = (−1)\(^{*N*+1}\)*m*\(_{22}\) and 1=*t* = (−1)\(^{*N*}\)*m*\(_{12}\). The matrix method is as efficient as the equivalent recursion, with two multiplications and two additions per step of the Euclidean algorithm.

### Euclid's lemma and unique factorization
Bézout's identity is essential to many applications of Euclid's algorithm, such as demonstrating the unique factorization of numbers into prime factors. To illustrate this, suppose that a number *L* can be written as a product of two factors *u* and *v*, that is, 1=*L* = *uv*. If another number *w* also divides *L* but is coprime with *u*, then *w* must divide *v*, by the following argument: If the greatest common divisor of *u* and *w* is 1, then integers *s* and *t* can be found such that
1=1 = *su* + *tw*
by Bézout's identity. Multiplying both sides by *v* gives the relation:
1=*v* = *suv* + *twv* = *sL* + *twv*

Since *w* divides both terms on the right-hand side, it must also divide the left-hand side, *v*. This result is known as Euclid's lemma. Specifically, if a prime number divides *L*, then it must divide at least one factor of *L*. Conversely, if a number *w* is coprime to each of a series of numbers *a*\(_{1}\), *a*\(_{2}\),.. ., *a*\(_{*n*}\), then *w* is also coprime to their product, *a*\(_{1}\) × *a*\(_{2}\) ×.. . × *a*\(_{*n*}\).

Euclid's lemma suffices to prove that every number has a unique factorization into prime numbers. To see this, assume the contrary, that there are two independent factorizations of *L* into *m* and *n* prime factors, respectively
1=*L* = *p*\(_{1}\)*p*\(_{2}\)...*p*\(_{*m*}\) = *q*\(_{1}\)*q*\(_{2}\)...*q*\(_{*n*}\). Since each prime *p* divides *L* by assumption, it must also divide one of the *q* factors; since each *q* is prime as well, it must be that 1=*p* = *q*. Iteratively dividing by the *p* factors shows that each *p* has an equal counterpart *q*; the two prime factorizations are identical except for their order. The unique factorization of numbers into primes has many applications in mathematical proofs, as shown below.

### Linear Diophantine equations

Diophantine equations are equations in which the solutions are restricted to integers; they are named after the 3rd-century Alexandrian mathematician Diophantus. A typical *linear* Diophantine equation seeks integers *x* and *y* such that
1=*ax* + *by* = *c*
where *a*, *b* and *c* are given integers. This can be written as an equation for *x* in modular arithmetic:
1=*ax* ≡ *c* mod *b*.

Let *g* be the greatest common divisor of *a* and *b*. Both terms in *ax* + *by* are divisible by *g*; therefore, *c* must also be divisible by *g*, or the equation has no solutions. By dividing both sides by *c*/*g*, the equation can be reduced to Bezout's identity
1=*sa* + *tb* = *g*,
where *s* and *t* can be found by the extended Euclidean algorithm. This provides one solution to the Diophantine equation, 1=*x*\(_{1}\) = *s* (*c*/*g*) and 1=*y*\(_{1}\) = *t* (*c*/*g*).

In general, a linear Diophantine equation has no solutions, or an infinite number of solutions. To find the latter, consider two solutions, (*x*\(_{1}\), *y*\(_{1}\)) and (*x*\(_{2}\), *y*\(_{2}\)), where
1=*ax*\(_{1}\) + *by*\(_{1}\) = *c* = *ax*\(_{2}\) + *by*\(_{2}\)
or equivalently
1=*a*(*x*\(_{1}\) − *x*\(_{2}\)) = *b*(*y*\(_{2}\) − *y*\(_{1}\)).

Therefore, the smallest difference between two *x* solutions is *b*/*g*, whereas the smallest difference between two *y* solutions is *a*/*g*. Thus, the solutions may be expressed as
1=*x* = *x*\(_{1}\) − *bu*/*g*
1=*y* = *y*\(_{1}\) + *au*/*g*.

By allowing *u* to vary over all possible integers, an infinite family of solutions can be generated from a single solution (*x*\(_{1}\), *y*\(_{1}\)). If the solutions are required to be *positive* integers (*x* > 0, *y* > 0), only a finite number of solutions may be possible. This restriction on the acceptable solutions allows some systems of Diophantine equations with more unknowns than equations to have a finite number of solutions; this is impossible for a system of linear equations when the solutions can be any real number (see Underdetermined system).

### Multiplicative inverses and the RSA algorithm
A finite field is a set of numbers with four generalized operations. The operations are called addition, subtraction, multiplication and division and have their usual properties, such as commutativity, associativity and distributivity. An example of a finite field is the set of 13 numbers using modular arithmetic. In this field, the results of any mathematical operation (addition, subtraction, multiplication, or division) is reduced modulo 13; that is, multiples of 13 are added or subtracted until the result is brought within the range 0–12. For example, the result of 1=5 × 7 = 35 mod 13 = 9. Such finite fields can be defined for any prime *p*; using more sophisticated definitions, they can also be defined for any power *m* of a prime *p*\(^{*m*}\). Finite fields are often called Galois fields, and are abbreviated as GF(*p*) or GF(*p*\(^{*m*}\)).

In such a field with *m* numbers, every nonzero element *a* has a unique modular multiplicative inverse, *a*\(^{−1}\) such that 1=*aa*\(^{−1}\) = *a*\(^{−1}\)*a* ≡ 1 mod *m*. This inverse can be found by solving the congruence equation *ax* ≡ 1 mod *m*, or the equivalent linear Diophantine equation
1=*ax* + *my* = 1.

This equation can be solved by the Euclidean algorithm, as described above. Finding multiplicative inverses is an essential step in the RSA algorithm, which is widely used in electronic commerce; specifically, the equation determines the integer used to decrypt the message. Although the RSA algorithm uses rings rather than fields, the Euclidean algorithm can still be used to find a multiplicative inverse where one exists. The Euclidean algorithm also has other applications in error-correcting codes; for example, it can be used as an alternative to the Berlekamp–Massey algorithm for decoding BCH and Reed–Solomon codes, which are based on Galois fields.

### Chinese remainder theorem
Euclid's algorithm can also be used to solve multiple linear Diophantine equations. Such equations arise in the Chinese remainder theorem, which describes a novel method to represent an integer *x*. Instead of representing an integer by its digits, it may be represented by its remainders *x*\(_{*i*}\) modulo a set of *N* coprime numbers *m*\(_{*i*}\):
\(\begin{align}
x_1 & \equiv x \pmod {m_1} \\
x_2 & \equiv x \pmod {m_2} \\
& \,\,\,\vdots \\
x_N & \equiv x \pmod {m_N} \,.
\end{align}\)

The goal is to determine *x* from its *N* remainders *x*\(_{*i*}\). The solution is to combine the multiple equations into a single linear Diophantine equation with a much larger modulus *M* that is the product of all the individual moduli *m*\(_{*i*}\), and define *M*\(_{*i*}\) as
\(M_i = \frac M {m_i}.\)

Thus, each *M*\(_{*i*}\) is the product of all the moduli *except* *m*\(_{*i*}\). The solution depends on finding *N* new numbers *h*\(_{*i*}\) such that
\(M_i h_i \equiv 1 \pmod {m_i} \,.\)

With these numbers *h*\(_{*i*}\), any integer *x* can be reconstructed from its remainders *x*\(_{*i*}\) by the equation
\(x \equiv (x_1 M_1 h_1 + x_2 M_2 h_2 + \cdots + x_N M_N h_N) \pmod M \,.\)

Since these numbers *h*\(_{*i*}\) are the multiplicative inverses of the *M*\(_{*i*}\), they may be found using Euclid's algorithm as described in the previous subsection.

### Stern–Brocot tree

The Euclidean algorithm can be used to arrange the set of all positive rational numbers into an infinite binary search tree, called the Stern–Brocot tree.
The number 1 (expressed as a fraction 1/1) is placed at the root of the tree, and the location of any other number *a*/*b* can be found by computing gcd(*a*,*b*) using the original form of the Euclidean algorithm, in which each step replaces the larger of the two given numbers by its difference with the smaller number (not its remainder), stopping when two equal numbers are reached. A step of the Euclidean algorithm that replaces the first of the two numbers corresponds to a step in the tree from a node to its right child, and a step that replaces the second of the two numbers corresponds to a step in the tree from a node to its left child. The sequence of steps constructed in this way does not depend on whether *a*/*b* is given in lowest terms, and forms a path from the root to a node containing the number *a*/*b*. This fact can be used to prove that each positive rational number appears exactly once in this tree.

For example, 3/4 can be found by starting at the root, going to the left once, then to the right twice:

\(\begin{align}
 & \gcd(3,4) & \leftarrow \\
= {} & \gcd(3,1) & \rightarrow \\
= {} & \gcd(2,1) & \rightarrow \\
= {} & \gcd(1,1).
\end{align}\)

The Euclidean algorithm has almost the same relationship to another binary tree on the rational numbers called the Calkin–Wilf tree. The difference is that the path is reversed: instead of producing a path from the root of the tree to a target, it produces a path from the target to the root.

### Continued fractions
The Euclidean algorithm has a close relationship with continued fractions. The sequence of equations can be written in the form
\(\begin{align}
\frac a b &= q_0 + \frac{r_0} b \\
\frac b {r_0} &= q_1 + \frac{r_1}{r_0} \\
\frac{r_0}{r_1} &= q_2 + \frac{r_2}{r_1} \\
& \,\,\, \vdots \\
\frac{r_{k-2{r_{k-1&= q_k + \frac{r_k}{r_{k-1\\
& \,\,\, \vdots \\
\frac{r_{N-2{r_{N-1&= q_N\,.
\end{align}\)

The last term on the right-hand side always equals the inverse of the left-hand side of the next equation. Thus, the first two equations may be combined to form
\(\frac a b = q_0 + \cfrac 1 {q_1 + \cfrac{r_1}{r_0\,.\)

The third equation may be used to substitute the denominator term *r*\(_{1}\)/*r*\(_{0}\), yielding
\(\frac a b = q_0 + \cfrac 1 {q_1 + \cfrac 1 {q_2 + \cfrac{r_2}{r_1}\,.\)

The final ratio of remainders *r*\(_{*k*}\)/*r*\(_{*k*−1}\) can always be replaced using the next equation in the series, up to the final equation. The result is a continued fraction
\(\frac a b = q_0 + \cfrac 1 {q_1 + \cfrac 1 {q_2 + \cfrac{1}{\ddots + \cfrac 1 {q_N= [ q_0; q_1, q_2, \ldots , q_N ] \,.\)

In the worked example above, the gcd(1071, 462) was calculated, and the quotients *q*\(_{*k*}\) were 2, 3 and 7, respectively. Therefore, the fraction 1071/462 may be written
\(\frac{1071}{462} = 2 + \cfrac 1 {3 + \cfrac 1 7} = [2; 3, 7]\)
as can be confirmed by calculation.

### Factorization algorithms
Calculating a greatest common divisor is an essential step in several integer factorization algorithms, such as Pollard's rho algorithm, Shor's algorithm, Dixon's factorization method and the Lenstra elliptic curve factorization. The Euclidean algorithm may be used to find this GCD efficiently. Continued fraction factorization uses continued fractions, which are determined using Euclid's algorithm.

## Algorithmic efficiency

The computational efficiency of Euclid's algorithm has been studied thoroughly. This efficiency can be described by the number of division steps the algorithm requires, multiplied by the computational expense of each step. The first known analysis of Euclid's algorithm is due to A. A. L. Reynaud in 1811, who showed that the number of division steps on input (*u*, *v*) is bounded by *v*; later he improved this to *v*/2 + 2. Later, in 1841, P. J. E. Finck showed that the number of division steps is at most 2 log\(_{2}\) *v* + 1, and hence Euclid's algorithm runs in time polynomial in the size of the input. Émile Léger, in 1837, studied the worst case, which is when the inputs are consecutive Fibonacci numbers. Finck's analysis was refined by Gabriel Lamé in 1844, who showed that the number of steps required for completion is never more than five times the number *h* of base-10 digits of the smaller number *b*.

In the uniform cost model (suitable for analyzing the complexity of gcd calculation on numbers that fit into a single machine word), each step of the algorithm takes constant time, and Lamé's analysis implies that the total running time is also *O*(*h*). However, in a model of computation suitable for computation with larger numbers, the computational expense of a single remainder computation in the algorithm can be as large as *O*(*h*\(^{2}\)). In this case the total time for all of the steps of the algorithm can be analyzed using a telescoping series, showing that it is also *O*(*h*\(^{2}\)). Modern algorithmic techniques based on the Schönhage–Strassen algorithm for fast integer multiplication can be used to speed this up, leading to quasilinear algorithms for the GCD.

### Number of steps
The number of steps to calculate the GCD of two natural numbers, *a* and *b*, may be denoted by *T*(*a*, *b*). If *g* is the GCD of *a* and *b*, then *a* = *mg* and *b* = *ng* for two coprime numbers *m* and *n*. Then
1=*T*(*a*, *b*) = *T*(*m*, *n*)
as may be seen by dividing all the steps in the Euclidean algorithm by *g*. By the same argument, the number of steps remains the same if *a* and *b* are multiplied by a common factor *w*: *T*(*a*, *b*) = *T*(*wa*, *wb*). Therefore, the number of steps *T* may vary dramatically between neighboring pairs of numbers, such as T(*a*, *b*) and T(*a*, *b* + 1), depending on the size of the two GCDs.

The recursive nature of the Euclidean algorithm gives another equation
1=*T*(*a*, *b*) = 1 + *T*(*b*, *r*\(_{0}\)) = 2 + *T*(*r*\(_{0}\), *r*\(_{1}\)) = … = *N* + *T*(*r*\(_{*N*−2}\), *r*\(_{*N*−1}\)) = *N* + 1
where *T*(*x*, 0) = 0 by assumption.

#### Worst-case
If the Euclidean algorithm requires *N* steps for a pair of natural numbers *a* > *b* > 0, the smallest values of *a* and *b* for which this is true are the Fibonacci numbers *F*\(_{*N*+2}\) and *F*\(_{*N*+1}\), respectively. More precisely, if the Euclidean algorithm requires *N* steps for the pair *a* > *b*, then one has *a* ≥ *F*\(_{*N*+2}\) and *b* ≥ *F*\(_{*N*+1}\). This can be shown by induction. If *N* = 1, *b* divides *a* with no remainder; the smallest natural numbers for which this is true is *b* = 1 and *a* = 2, which are *F*\(_{2}\) and *F*\(_{3}\), respectively. Now assume that the result holds for all values of *N* up to *M* − 1. The first step of the *M*-step algorithm is *a* = *q*\(_{0}\)*b* + *r*\(_{0}\), and the Euclidean algorithm requires *M* − 1 steps for the pair *b* > *r*\(_{0}\). By induction hypothesis, one has *b* ≥ *F*\(_{*M*+1}\) and *r*\(_{0}\) ≥ *F*\(_{*M*}\). Therefore, *a* = *q*\(_{0}\)*b* + *r*\(_{0}\) ≥ *b* + *r*\(_{0}\) ≥ *F*\(_{*M*+1}\) + *F*\(_{*M*}\) = *F*\(_{*M*+2}\),
which is the desired inequality.
This proof, published by Gabriel Lamé in 1844, represents the beginning of computational complexity theory, and also the first practical application of the Fibonacci numbers.

This result suffices to show that the number of steps in Euclid's algorithm can never be more than five times the number of its digits (base 10). For if the algorithm requires *N* steps, then *b* is greater than or equal to *F*\(_{*N*+1}\) which in turn is greater than or equal to *φ*\(^{*N*−1}\), where *φ* is the golden ratio. Since *b* ≥ *φ*\(^{*N*−1}\), then *N* − 1 ≤ log\(_{*φ*}\)*b*. Since log\(_{10}\)*φ* > 1/5, (*N* − 1)/5 < log\(_{10}\)*φ* log\(_{*φ*}\)*b* = log\(_{10}\)*b*. Thus, *N* ≤ 5 log\(_{10}\)*b*. Thus, the Euclidean algorithm always needs less than *O*(*h*) divisions, where *h* is the number of digits in the smaller number *b*.

#### Average
The average number of steps taken by the Euclidean algorithm has been defined in three different ways. The first definition is the average time *T*(*a*) required to calculate the GCD of a given number *a* and a smaller natural number *b* chosen with equal probability from the integers 0 to *a* − 1

\(T(a) = \frac 1 a \sum_{0 \leq b−(1/6)+*ε*, where *ε* is infinitesimal. The constant *C* in this formula is called Porter's constant and equals
\(C= -\frac 1 2 + \frac{6 \ln 2}{\pi^2}\left(4\gamma -\frac{24}{\pi^2}\zeta'(2) + 3\ln 2 - 2\right) \approx 1.467\)
where *γ* is the Euler–Mascheroni constant and *ζ* is the derivative of the Riemann zeta function. The leading coefficient (12/π\(^{2}\)) ln 2 was determined by two independent methods.

Since the first average can be calculated from the tau average by summing over the divisors *d* of *a*

\(T(a) = \frac 1 a \sum_{d \mid a} \varphi(d) \tau(d)\)
it can be approximated by the formula
\(T(a) \approx C + \frac{12}{\pi^2} \ln 2\, \biggl({\ln a} - \sum_{d \mid a} \frac{\Lambda(d)} d\biggr)\)
where Λ(*d*) is the Mangoldt function.

A third average *Y*(*n*) is defined as the mean number of steps required when both *a* and *b* are chosen randomly (with uniform distribution) from 1 to *n*

\(Y(n) = \frac 1 {n^2} \sum_{a=1}^n \sum_{b=1}^n T(a, b) = \frac 1 n \sum_{a=1}^n T(a).\)

Substituting the approximate formula for *T*(*a*) into this equation yields an estimate for *Y*(*n*)
\(Y(n) \approx \frac{12}{\pi^2} \ln 2 \ln n + 0.06.\)

### Computational expense per step
In each step *k* of the Euclidean algorithm, the quotient *q*\(_{*k*}\) and remainder *r*\(_{*k*}\) are computed for a given pair of integers *r*\(_{*k*−2}\) and *r*\(_{*k*−1}\)
1=*r*\(_{*k*−2}\) = *q*\(_{*k*}\) *r*\(_{*k*−1}\) + *r*\(_{*k*}\).

The computational expense per step is associated chiefly with finding *q*\(_{*k*}\); since the remainder *r*\(_{*k*}\) can be calculated quickly from *r*\(_{*k*−2}\), *r*\(_{*k*−1}\), and *q*\(_{*k*}\)
1=*r*\(_{*k*}\) = *r*\(_{*k*−2}\) − *q*\(_{*k*}\) *r*\(_{*k*−1}\).

The computational expense of dividing *h*-bit numbers scales as *O*(*h*(*ℓ* + 1)), where is the length of the quotient.

For comparison, Euclid's original subtraction-based algorithm can be much slower. A single integer division is equivalent to the quotient *q* number of subtractions. If the ratio of *a* and *b* is very large, the quotient is large and many subtractions will be required. On the other hand, it has been shown that the quotients are very likely to be small integers. The probability of a given quotient *q* is approximately ln |*u*/(*u* − 1)| where 1=*u* = (*q* + 1)\(^{2}\). For illustration, the probability of a quotient of 1, 2, 3, or 4 is roughly 41.5%, 17.0%, 9.3%, and 5.9%, respectively. Since the operation of subtraction is faster than division, particularly for large numbers, the subtraction-based Euclid's algorithm is competitive with the division-based version. This is exploited in the binary version of Euclid's algorithm.

Combining the estimated number of steps with the estimated computational expense per step shows that the Euclid's algorithm grows quadratically (*h*\(^{2}\)) with the average number of digits *h* in the initial two numbers *a* and *b*. Let *h*\(_{0}\), *h*\(_{1}\),.. ., *h*\(_{*N*−1}\) represent the number of digits in the successive remainders *r*\(_{0}\), *r*\(_{1}\),.. ., *r*\(_{*N*−1}\). Since the number of steps *N* grows linearly with *h*, the running time is bounded by

\(O\Big(\sum_{i<N}h_i(h_i-h_{i+1}+2)\Big)\subseteq O\Big(h\sum_{i<N}(h_i-h_{i+1}+2) \Big) \subseteq O(h(h_0+2N))\subseteq O(h^2).\)

### Alternative methods
Euclid's algorithm is widely used in practice, especially for small numbers, due to its simplicity. For comparison, the efficiency of alternatives to Euclid's algorithm may be determined.

One inefficient approach to finding the GCD of two natural numbers *a* and *b* is to calculate all their common divisors; the GCD is then the largest common divisor. The common divisors can be found by dividing both numbers by successive integers from 2 to the smaller number *b*. The number of steps of this approach grows linearly with *b*, or exponentially in the number of digits. Another inefficient approach is to find the prime factors of one or both numbers. As noted above, the GCD equals the product of the prime factors shared by the two numbers *a* and *b*. Present methods for prime factorization are also inefficient; many modern cryptography systems even rely on that inefficiency.

The binary GCD algorithm is an efficient alternative that substitutes division with faster operations by exploiting the binary representation used by computers. However, this alternative also scales like *O*(*h*²). It is generally faster than the Euclidean algorithm on real computers, even though it scales in the same way. Additional efficiency can be gleaned by examining only the leading digits of the two numbers *a* and *b*. The binary algorithm can be extended to other bases (*k*-ary algorithms), with up to fivefold increases in speed. Lehmer's GCD algorithm uses the same general principle as the binary algorithm to speed up GCD computations in arbitrary bases.

A recursive approach for very large integers (with more than 25,000 digits) leads to quasilinear integer GCD algorithms, such as those of Schönhage, and Stehlé and Zimmermann. These algorithms exploit the 2×2 matrix form of the Euclidean algorithm given above. These quasilinear methods generally scale as *O*(*h* log *h*

## Generalizations
Although the Euclidean algorithm is used to find the greatest common divisor of two natural numbers (positive integers), it may be generalized to the real numbers, and to other mathematical objects, such as polynomials, quadratic integers and Hurwitz quaternions. In the latter cases, the Euclidean algorithm is used to demonstrate the crucial property of unique factorization, i.e., that such numbers can be factored uniquely into irreducible elements, the counterparts of prime numbers. Unique factorization is essential to many proofs of number theory.

### Rational and real numbers
Euclid's algorithm can be applied to real numbers, as described by Euclid in Book 10 of his *Elements*. The goal of the algorithm is to identify a real number such that two given real numbers, and , are integer multiples of it: 1=*a* = *mg* and 1=*b* = *ng*, where and are integers. This identification is equivalent to finding an integer relation among the real numbers and ; that is, it determines integers and such that 1=*sa* + *tb* = 0. If such an equation is possible, *a* and *b* are called commensurable lengths, otherwise they are incommensurable lengths.

The real-number Euclidean algorithm differs from its integer counterpart in two respects. First, the remainders *r*\(_{*k*}\) are real numbers, although the quotients *q*\(_{*k*}\) are integers as before. Second, the algorithm is not guaranteed to end in a finite number of steps. If it does, the fraction *a*/*b* is a rational number, i.e., the ratio of two integers
\(\frac{a}{b} = \frac{mg}{ng} = \frac{m}{n},\)
and can be written as a finite continued fraction 1=[*q*\(_{0}\); *q*\(_{1}\), *q*\(_{2}\),.. ., *q*\(_{*N*}\)]. If the algorithm does not stop, the fraction *a*/*b* is an irrational number and can be described by an infinite continued fraction 1=[*q*\(_{0}\); *q*\(_{1}\), *q*\(_{2}\), …]. Examples of infinite continued fractions are the golden ratio 1=*φ* = [1; 1, 1,.. .] and the square root of two, sqrt[1; 2, 2,.. .]. When applied to two arbitrary real numbers, the algorithm is unlikely to stop; since almost all ratios *a*/*b* of two real numbers are irrational.

An infinite continued fraction may be truncated at a step 1=*k* [*q*\(_{0}\); *q*\(_{1}\), *q*\(_{2}\),.. ., *q*\(_{*k*}\)] to yield an approximation to *a*/*b* that improves as is increased. The approximation is described by convergents *m*\(_{*k*}\)/*n*\(_{*k*}\); the numerator and denominators are coprime and obey the recurrence relation
\(\begin{align}
 m_k &= q_k m_{k-1} + m_{k-2} \\
 n_k &= q_k n_{k-1} + n_{k-2},
 \end{align}\)
where 1=*m*\(_{−1}\) = *n*\(_{−2}\) = 1 and 1=*m*\(_{−2}\) = *n*\(_{−1}\) = 0 are the initial values of the recursion. The convergent *m*\(_{*k*}\)/*n*\(_{*k*}\) is the best rational number approximation to *a*/*b* with denominator *n*\(_{*k*}\):
\(\left|\frac{a}{b} - \frac{m_k}{n_k}\right| < \frac{1}{n_k^2}.\)

### Polynomials

Polynomials in a single variable *x* can be added, multiplied and factored into irreducible polynomials, which are the analogs of the prime numbers for integers. The greatest common divisor polynomial *g*(*x*) of two polynomials *a*(*x*) and *b*(*x*) is defined as the product of their shared irreducible polynomials, which can be identified using the Euclidean algorithm. The basic procedure is similar to that for integers. At each step , a quotient polynomial *q*\(_{*k*}\)(*x*) and a remainder polynomial *r*\(_{*k*}\)(*x*) are identified to satisfy the recursive equation
\(r_{k-2}(x) = q_k(x)r_{k-1}(x) + r_k(x),\)
where 1=*r*\(_{−2}\)(*x*) = *a*(*x*) and 1=*r*\(_{−1}\)(*x*) = *b*(*x*). Each quotient polynomial is chosen such that each remainder is either zero or has a degree that is smaller than the degree of its predecessor: deg[*r*\(_{*k*}\)(*x*)] < deg[*r*\(_{*k*−1}\)(*x*)]. Since the degree is a nonnegative integer, and since it decreases with every step, the Euclidean algorithm concludes in a finite number of steps. The last nonzero remainder is the greatest common divisor of the original two polynomials, *a*(*x*) and *b*(*x*).

For example, consider the following two quartic polynomials, which each factor into two quadratic polynomials
\(\begin{align}
 a(x) &= x^4 - 4x^3 + 4x^2 - 3x + 14 = (x^2 - 5x + 7)(x^2 + x + 2) \qquad \text{and}\\
 b(x) &= x^4 + 8x^3 + 12x^2 + 17x + 6 = (x^2 + 7x + 3)(x^2 + x + 2).
 \end{align}\)

Dividing *a*(*x*) by *b*(*x*) yields a remainder 1=*r*\(_{0}\)(*x*) = *x*\(^{3}\) + (2/3)*x*\(^{2}\) + (5/3)*x* − (2/3). In the next step, *b*(*x*) is divided by *r*\(_{0}\)(*x*) yielding a remainder 1=*r*\(_{1}\)(*x*) = *x*\(^{2}\) + *x* + 2. Finally, dividing *r*\(_{0}\)(*x*) by *r*\(_{1}\)(*x*) yields a zero remainder, indicating that *r*\(_{1}\)(*x*) is the greatest common divisor polynomial of *a*(*x*) and *b*(*x*), consistent with their factorization.

Many of the applications described above for integers carry over to polynomials. The Euclidean algorithm can be used to solve linear Diophantine equations and Chinese remainder problems for polynomials; continued fractions of polynomials can also be defined.

The polynomial Euclidean algorithm has other applications, such as Sturm chains, a method for counting the zeros of a polynomial that lie inside a given real interval. This in turn has applications in several areas, such as the Routh–Hurwitz stability criterion in control theory.

Finally, the coefficients of the polynomials need not be drawn from integers, real numbers or even the complex numbers. For example, the coefficients may be drawn from a general field, such as the finite fields GF(*p*) described above. The corresponding conclusions about the Euclidean algorithm and its applications hold even for such polynomials.

### Gaussian integers

The Gaussian integers are complex numbers of the form 1=*α* = *u* + *vi*, where and are ordinary integers and is the square root of negative one. By defining an analog of the Euclidean algorithm, Gaussian integers can be shown to be uniquely factorizable, by the argument above. This unique factorization is helpful in many applications, such as deriving all Pythagorean triples or proving Fermat's theorem on sums of two squares. In general, the Euclidean algorithm is convenient in such applications, but not essential; for example, the theorems can often be proven by other arguments.

The Euclidean algorithm developed for two Gaussian integers and is nearly the same as that for ordinary integers, but differs in two respects. As before, we set 1=*r*\(_{−2}\) = *α* and 1=*r*\(_{−1}\) = *β*, and the task at each step is to identify a quotient *q*\(_{*k*}\) and a remainder *r*\(_{*k*}\) such that
\(r_k = r_{k-2} - q_k r_{k-1},\)
where every remainder is strictly smaller than its predecessor:. The first difference is that the quotients and remainders are themselves Gaussian integers, and thus are complex numbers. The quotients *q*\(_{*k*}\) are generally found by rounding the real and complex parts of the exact ratio (such as the complex number *α*/*β*) to the nearest integers. The second difference lies in the necessity of defining how one complex remainder can be "smaller" than another. To do this, a norm function 1=*f*(*u* + *vi*) = *u*\(^{2}\) + *v*\(^{2}\) is defined, which converts every Gaussian integer *u* + *vi* into an ordinary integer. After each step of the Euclidean algorithm, the norm of the remainder *f*(*r*\(_{*k*}\)) is smaller than the norm of the preceding remainder, *f*(*r*\(_{*k*−1}\)). Since the norm is a nonnegative integer and decreases with every step, the Euclidean algorithm for Gaussian integers ends in a finite number of steps. The final nonzero remainder is gcd(*α*, *β*), the Gaussian integer of largest norm that divides both and ; it is unique up to multiplication by a unit, ±1 or ±*i*.

Many of the other applications of the Euclidean algorithm carry over to Gaussian integers. For example, it can be used to solve linear Diophantine equations and Chinese remainder problems for Gaussian integers; continued fractions of Gaussian integers can also be defined.

### Euclidean domains
A set of elements under two binary operations, denoted as addition and multiplication, is called a Euclidean domain if it forms a commutative ring and, roughly speaking, if a generalized Euclidean algorithm can be performed on them. The two operations of such a ring need not be the addition and multiplication of ordinary arithmetic; rather, they can be more general, such as the operations of a mathematical group or monoid. Nevertheless, these general operations should respect many of the laws governing ordinary arithmetic, such as commutativity, associativity and distributivity.

The generalized Euclidean algorithm requires a *Euclidean function*, i.e., a mapping from into the set of nonnegative integers such that, for any two nonzero elements and in , there exist and in such that 1=*a* = *qb* + *r* and *f*(*r*) < *f*(*b*). Examples of such mappings are the absolute value for integers, the degree for univariate polynomials, and the norm for Gaussian integers above. The basic principle is that each step of the algorithm reduces *f* inexorably; hence, if can be reduced only a finite number of times, the algorithm must stop in a finite number of steps. This principle relies on the well-ordering property of the non-negative integers, which asserts that every non-empty set of non-negative integers has a smallest member.

The fundamental theorem of arithmetic applies to any Euclidean domain: Any number from a Euclidean domain can be factored uniquely into irreducible elements. Any Euclidean domain is a unique factorization domain (UFD), although the converse is not true. The Euclidean domains and the UFD's are subclasses of the GCD domains, domains in which a greatest common divisor of two numbers always exists. In other words, a greatest common divisor may exist (for all pairs of elements in a domain), although it may not be possible to find it using a Euclidean algorithm. A Euclidean domain is always a principal ideal domain (PID), an integral domain in which every ideal is a principal ideal. Again, the converse is not true: not every PID is a Euclidean domain.

The unique factorization of Euclidean domains is useful in many applications. For example, the unique factorization of the Gaussian integers is convenient in deriving formulae for all Pythagorean triples and in proving Fermat's theorem on sums of two squares. Unique factorization was also a key element in an attempted proof of Fermat's Last Theorem published in 1847 by Gabriel Lamé, the same mathematician who analyzed the efficiency of Euclid's algorithm, based on a suggestion of Joseph Liouville. Lamé's approach required the unique factorization of numbers of the form *x* + *ωy*, where and are integers, and 1=*ω* = *e*\(^{2*iπ*/*n*}\) is an th root of 1, that is, 1=*ω*\(^{*n*}\) = 1. Although this approach succeeds for some values of (such as 1=*n* = 3, the Eisenstein integers), in general such numbers do factor uniquely. This failure of unique factorization in some cyclotomic fields led Ernst Kummer to the concept of ideal numbers and, later, Richard Dedekind to ideals.

#### Unique factorization of quadratic integers

The quadratic integer rings are helpful to illustrate Euclidean domains. Quadratic integers are generalizations of the Gaussian integers in which the imaginary unit *i* is replaced by a number. Thus, they have the form *u* + *vω*, where and are integers and has one of two forms, depending on a parameter. If does not equal a multiple of four plus one, then
\(\omega = \sqrt D. \)

If, however, *D* does equal a multiple of four plus one, then
\(\omega = \frac{1 + \sqrt{D{2}. \)

If the function corresponds to a norm function, such as that used to order the Gaussian integers above, then the domain is known as *norm-Euclidean*. The norm-Euclidean rings of quadratic integers are exactly those where is one of the values −11, −7, −3, −2, −1, 2, 3, 5, 6, 7, 11, 13, 17, 19, 21, 29, 33, 37, 41, 57, or 73. The cases 1=*D* = −1 and 1=*D* = −3 yield the Gaussian integers and Eisenstein integers, respectively.

If is allowed to be any Euclidean function, then the list of possible values of for which the domain is Euclidean is not yet known. The first example of a Euclidean domain that was not norm-Euclidean (with 1=*D* = 69) was published in 1994. In 1973, Weinberger proved that a quadratic integer ring with *D* > 0 is Euclidean if, and only if, it is a principal ideal domain, provided that the generalized Riemann hypothesis holds.

### Noncommutative rings
The Euclidean algorithm may be applied to some noncommutative rings such as the set of Hurwitz quaternions. Let and represent two elements from such a ring. They have a common right divisor if 1=*α* = *ξδ* and 1=*β* = *ηδ* for some choice of and in the ring. Similarly, they have a common left divisor if 1=*α* = *dξ* and 1=*β* = *dη* for some choice of and in the ring. Since multiplication is not commutative, there are two versions of the Euclidean algorithm, one for right divisors and one for left divisors. Choosing the right divisors, the first step in finding the gcd(*α*, *β*) by the Euclidean algorithm can be written
\(\rho_0 = \alpha - \psi_0\beta = (\xi - \psi_0\eta)\delta,\)
where *ψ*\(_{0}\) represents the quotient and *ρ*\(_{0}\) the remainder. Here the quotient and remainder are chosen so that (if nonzero) the remainder has *N*(*ρ*\(_{0}\)) < *N*(*β*) for a "Euclidean function" *N* defined analogously to the Euclidean functions of Euclidean domains in the non-commutative case. This equation shows that any common right divisor of and is likewise a common divisor of the remainder *ρ*\(_{0}\). The analogous equation for the left divisors would be
\(\rho_0 = \alpha - \beta\psi_0 = \delta(\xi - \eta\psi_0).\)

With either choice, the process is repeated as above until the greatest common right or left divisor is identified. As in the Euclidean domain, the "size" of the remainder *ρ*\(_{0}\) (formally, its Euclidean function or "norm") must be strictly smaller than , and there must be only a finite number of possible sizes for *ρ*\(_{0}\), so that the algorithm is guaranteed to terminate.

Many results for the GCD carry over to noncommutative numbers. For example, Bézout's identity states that the right gcd(*α*, *β*) can be expressed as a linear combination of and. In other words, there are numbers and such that
\(\Gamma_\text{right} = \sigma\alpha + \tau\beta.\)

The analogous identity for the left GCD is nearly the same:
\(\Gamma_\text{left} = \alpha\sigma + \beta\tau.\)

Bézout's identity can be used to solve Diophantine equations. For instance, one of the standard proofs of Lagrange's four-square theorem, that every positive integer can be represented as a sum of four squares, is based on quaternion GCDs in this way.
